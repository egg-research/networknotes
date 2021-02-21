package db

import (
	"errors"
	"strconv"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func CheckDocId(driver neo4j.Driver, uid int, docId int) (bool) {
	// null if user does not exist
	session := driver.NewSession(neo4j.SessionConfig{AccessMode:neo4j.AccessModeRead})
	defer session.Close()

	retDocId, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			`
			MATCH (u:User) WHERE id(u) = $uid
			MATCH (d:Document) <-- (u) WHERE id(d) = $docId 
			RETURN id(d)
			`,
			map[string]interface{}{"uid":uid, "docId":docId})

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


func GetDoc(driver neo4j.Driver, uid int, docId int) (interface{}, error) {
	if (!CheckDocId(driver, uid, docId)) {
		return nil, errors.New("User does not own doc with id: "+strconv.Itoa(docId))
	}

	// null if user does not exist
	session := driver.NewSession(neo4j.SessionConfig{AccessMode:neo4j.AccessModeRead})
	defer session.Close()

	kws, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			`
			MATCH (d:Document) WHERE id(d) = $docId
			MATCH (d) -[r]-> (kw)
			RETURN collect(kw.kw), collect(r.kwText)
			`,
			map[string]interface{}{
				"docId":docId,
			})

		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().Values, nil
		}

		return nil, result.Err()
	})
	return kws, err
}

//get all docs

func ReadDocFS(driver neo4j.Driver, db Firestore, uid int, docId int) (interface{}, error) {
	if (!CheckDocId(driver, uid, docId)) {
		return nil, errors.New("User does not own doc with id: "+strconv.Itoa(docId))
	}

	res, err := db.Client.Collection("docs").Doc(strconv.Itoa(docId)).Get(db.Ctx)
	if err != nil {
		return nil, err
	}
	return res.Data(), err
}

func AddDoc(driver neo4j.Driver, uid int) (interface{}, error) {
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	
	docId, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			`
			MATCH (u:User) WHERE id(u) = $uid
			CREATE (d:Document)
			WITH d,u
			CREATE (u)-[:DOCUMENT]->(d) 
			RETURN id(d)
			`,
			map[string]interface{}{
				"uid":uid,
			})

		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().Values[0], nil
		}

		return nil, result.Err()
	})
	return docId, err
}

func WriteDocFS(driver neo4j.Driver, db Firestore, uid int, docId int, docName string, docText string, rawDocText interface{}) (interface{}, error) {		
	if (!CheckDocId(driver, uid, docId)) {
		return nil, errors.New("User does not own doc with id: "+strconv.Itoa(docId))
	}

	_, err := db.Client.Collection("docs").Doc(strconv.Itoa(docId)).Set(db.Ctx, map[string]interface{}{
		"docName": docName,
		"docText": docText,
		"rawDocText": rawDocText,
	})
	if err != nil {
		return nil, err
	}
	return nil, err
}

func WriteDoc(driver neo4j.Driver, db Firestore, body User) (interface{}, error) {
	if !CheckUid(driver, body.Uid) {
		return nil, errors.New("User does not exist")
	}

	if (body.Doc.DocId < 0) {
		docId, err := AddDoc(driver, body.Uid)
		_, err = WriteDocFS(driver, db, body.Uid, toInt(docId), body.Doc.DocName, body.Doc.DocText, body.Doc.RawDocText)
		if err != nil {
			return nil, err
		}
		return 	docId, err
	} else {
		if !CheckDocId(driver, body.Uid, body.Doc.DocId) {
			return nil, errors.New("User does not own document with id:"+strconv.Itoa(body.Doc.DocId))
		}
		_, err := WriteDocFS(driver, db, body.Uid, body.Doc.DocId, body.Doc.DocName, body.Doc.DocText, body.Doc.RawDocText)
		if err != nil {
			return nil, err
		}
		return body.Doc.DocId, err
	}
}

func ReadDoc(driver neo4j.Driver, db Firestore, body User) (interface{}, error) {
	if !CheckUid(driver, body.Uid) {
		return nil, errors.New("User does not exist")
	}

	kws, kerr := GetDoc(driver, body.Uid, body.Doc.DocId)
	if kerr != nil {
		return nil, kerr
	}
	texts, terr := ReadDocFS(driver, db, body.Uid, body.Doc.DocId) 
	if terr != nil {
		return nil, terr
	}

	v := []interface{} {
		kws,
		texts,
	}

	return v, nil
}


func GetAllDocs(driver neo4j.Driver, db Firestore, uid int) (map[int]interface{}, error) {
	if !CheckUid(driver, uid) {
		return nil, errors.New("User does not exist")
	}

	// null if user does not exist
	session := driver.NewSession(neo4j.SessionConfig{AccessMode:neo4j.AccessModeRead})
	defer session.Close()

	docIds, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			`
			MATCH (u:User) WHERE id(u) = $uid
			MATCH (d:Document) <-- (u)  
			RETURN collect(id(d))
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

	res := make(map[int]interface{}, len(docIds.([]interface{})))

	for _, docId := range docIds.([]interface{}) {
		id := toInt(docId)
		texts, terr := ReadDocFS(driver, db, uid, id) 
		if terr != nil {
			return nil, terr
		}
		res[id] = texts
	}
	return res, nil
}