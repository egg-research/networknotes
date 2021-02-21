package db

import (
	"fmt"
	"errors"
	"strconv"
	"github.com/fatih/structs"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func CheckKwId(driver neo4j.Driver, uid int, kwId int) (bool) {
	// null if Keyword does not exist
	session := driver.NewSession(neo4j.SessionConfig{AccessMode:neo4j.AccessModeRead})
	defer session.Close()

	retDocId, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			`
			MATCH (u:User) WHERE id(u) = $uid
			MATCH (k:Keyword) <-- (u) WHERE id(k) = $kwId
			RETURN id(k)
			`,
			map[string]interface{}{"uid":uid, "kwId":kwId})

		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().Values[0], nil
		}

		return nil, result.Err()
	})

	if err != nil {
		return false
	}

	if retDocId == nil {
		return false
	}
	return true
}

func mapKeywords(kws []Keyword) []map[string]interface{} {
	var result = make([]map[string]interface{}, len(kws))

	for index, item := range kws {
		result[index] = structs.Map(item)
	}
	return result
}

func AddKw(driver neo4j.Driver, db Firestore, uid int, docId int, kws []Keyword) (interface{}, error) {
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	kwIds, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			`
			UNWIND $kws as kw
			MATCH (u:User) WHERE id(u) = $uid 
			MATCH (d:Document) WHERE id(d) = $docId 
			MERGE (k:Keyword {kw:kw.Kw})
			MERGE (u)-[:DOCUMENT]->(d)
			MERGE (k)-[:DOCUMENT]->(d)
			MERGE (u)-[:KEYWORD {kwText:kw.KwText}]->(k)
			MERGE (d)-[:KEYWORD {kwText:kw.KwText}]->(k)
			RETURN collect(id(k))
			`,
			map[string]interface{}{
				"uid":uid,
				"docId":docId,
				"kws":mapKeywords(kws),
			})
		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().Values[0], nil
		}

		return nil, result.Err()
	})

	fmt.Println(kwIds)
	for index, kwId := range kwIds.([]interface{}) {
		kw := kws[index]
		WriteKwFS(driver, db, kw, toInt(kwId))
	}
	return kwIds, err
}

func WriteKw(driver neo4j.Driver, db Firestore, body User) (interface{}, error) {
	if !CheckUid(driver, body.Uid) {
		return nil, errors.New("User does not exist")
	}

	if !CheckDocId(driver, body.Uid, body.Doc.DocId) {
		return nil, errors.New("Document does not exist for this user")
	}

	resp, err := AddKw(driver, db, body.Uid, body.Doc.DocId, body.Kws)
	if (err != nil) {
		return nil, err
	}
	return resp, err
}


func DelKw(driver neo4j.Driver, db Firestore, body User) (interface{}, error) {
	if !CheckUid(driver, body.Uid) {
		return nil, errors.New("User does not exist")
	}

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			`
			UNWIND $kws as kw
			MATCH (u:User) WHERE id(u) = $uid
			MATCH (k:Keyword {kw:kw.Kw}) <-[e]- (u)
			MATCH (k:Keyword {kw:kw.Kw}) -[f]-> (d:Document) <-- (u)
			MATCH (k:Keyword {kw:kw.Kw}) <-[g]- (d:Document) <-- (u)
			DELETE e
			DELETE f
			DELETE g
			`,
			map[string]interface{}{
				"uid":body.Uid,
				"kws":mapKeywords(body.Kws),
			})
		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().Values[0], nil
		}
		return nil, result.Err()
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func WriteKwFS(driver neo4j.Driver, db Firestore, kw Keyword, kwId int) (interface{}, error) {	
	_, err := db.Client.Collection("kws").Doc(strconv.Itoa(kwId)).Set(db.Ctx, map[string]interface{}{
		"kw": kw.Kw,
		"kwText": kw.KwText,
	})
	if err != nil {
		return nil, err
	}
	return nil, err
}

func ReadKwFS(driver neo4j.Driver, db Firestore, uid int, kwId int) (interface{}, error) {
	if (!CheckKwId(driver, uid, kwId)) {
		return nil, errors.New("User does not own keyword with id: "+strconv.Itoa(kwId))
	}

	res, err := db.Client.Collection("kws").Doc(strconv.Itoa(kwId)).Get(db.Ctx)
	if err != nil {
		return nil, err
	}
	return res.Data(), err
}

func GetAllKws(driver neo4j.Driver, db Firestore, uid int) (interface{}, error) {
	if !CheckUid(driver, uid) {
		return nil, errors.New("User does not exist")
	}

	session := driver.NewSession(neo4j.SessionConfig{AccessMode:neo4j.AccessModeRead})
	defer session.Close()

	kws, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			`
			MATCH (u:User) WHERE id(u) = $uid
			MATCH (k:Keyword) <-- (u)  
			RETURN collect(id(k))
			`,
			map[string]interface{}{"uid":uid})

		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().Values[0], nil
		}

		return nil, result.Err()
	})
    if err != nil {
		return nil, err
	}

	res := make(map[int]interface{}, len(kws.([]interface{})))

	for _, kwId := range kws.([]interface{}) {
		id := toInt(kwId)
		kw, terr := ReadKwFS(driver, db, uid, id) 
		if terr != nil {
			return nil, terr
		}
		res[id] = kw
	}
	return res, nil
}