package db

import (
	"fmt"
	"errors"
	"github.com/fatih/structs"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func mapKeywords(kws []Keyword) []map[string]interface{} {
	var result = make([]map[string]interface{}, len(kws))

	for index, item := range kws {
		result[index] = structs.Map(item)
	}
	fmt.Println("kws", result)
	return result
}

func AddKw(driver neo4j.Driver, uid string, docId string, kws []Keyword) (interface{}, error) {
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	resp, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		
		result, err := transaction.Run(
			`
			UNWIND $kws as kw
			MATCH (u:User) WHERE id(u) = $uid
			MATCH (d:Document) WHERE id(d) = $docId  
			MERGE (k:Keyword {kw:kw.Kw})
			WITH u,d,k,kw
			MERGE (u)-[:DOCUMENT]->(d)
			MERGE (k)-[:DOCUMENT]->(d)
			MERGE (u)-[:KEYWORD {kwText:kw.KwText}]->(k)
			MERGE (d)-[:KEYWORD {kwText:kw.KwText}]->(k)
			RETURN id(k)
			`,
			map[string]interface{}{
				"uid":uid,
				"docId":docId,
				"kws":mapKeywords(kws),
			})
		fmt.Println("error", err)

		if err != nil {
			return nil, err
		}
		fmt.Println("Keyword Result", result)

		if result.Next() {
			return result.Record().Values[0], nil
		}

		return nil, result.Err()
	})
	return resp, err
}


func WriteKw(driver neo4j.Driver, db Firestore, body User) (interface{}, error) {
	if !CheckUid(driver, body.Uid) {
		return nil, errors.New("User does not exist")
	}

	if !CheckDocId(driver, body.Doc.DocId) {
		return nil, errors.New("Document does not exist")
	}

	resp, err := AddKw(driver, body.Uid, body.Doc.DocId, body.Kws)
	if (err != nil) {
		return nil, err
	}
	return resp, err
}
