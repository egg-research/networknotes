package db

import (
	"fmt"
	"errors"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func GetDocId(driver neo4j.Driver, uid string, docName string) (interface{}, error) {
	// null if user does not exist
	session := driver.NewSession(neo4j.SessionConfig{AccessMode:neo4j.AccessModeRead})
	defer session.Close()

	docId, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			`
			MATCH (d:Document {docName: $docName}) <-- (:User {user: $uid})
			RETURN id(d)
			`,
			map[string]interface{}{"uid":uid, "docName":docName})

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

func GetDoc(driver neo4j.Driver, docId string) (interface{}, error) {
	// null if user does not exist
	session := driver.NewSession(neo4j.SessionConfig{AccessMode:neo4j.AccessModeRead})
	defer session.Close()

	kws, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			`
			MATCH (n) WHERE id(n) = $docId
			MATCH (n) -[r]-> (kw)
			RETURN kw.kw, r.kwText
			`,
			map[string]interface{}{"docId":docId})

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

func ReadDocFS(db Firestore, docId string) (interface{}, error) {
	res, err := db.Client.Collection("docs").Doc(docId).Get(db.Ctx)
	if err != nil {
		return nil, err
	}
	return res.Data(), err
}

func AddDoc(driver neo4j.Driver, uid string, docName string) (interface{}, error) {
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	
	docId, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			`
			CREATE (d:Document {docName: $docName})
			MERGE (u:User {user: $uid})
			WITH d,u
			CREATE (u)-[:DOCUMENT]->(d) 
			RETURN id(d)
			`,
			map[string]interface{}{
				"uid":uid,
				"docName":docName,
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

func WriteDocFS(db Firestore, docId string, docText string, rawDocText interface{}) (interface{}, error) {
	res, err := db.Client.Collection("docs").Doc(docId).Set(db.Ctx, map[string]interface{}{
		"docText": docText,
		"rawDocText": rawDocText,
	})
	if err != nil {
		return nil, err
	}
	return res, err
}

func WriteDoc(driver neo4j.Driver, db Firestore, body User) (interface{}, error) {
	uid, uerr := GetUid(driver, body.Uid)
	if uerr != nil {
		return nil, uerr
	}

	if uid == nil {
		return nil, errors.New("User does not have an account")
	}

	// if doc already exists, return doc id
	docId, err := GetDocId(driver, body.Uid, body.Doc.DocName)
	if err != nil {
		return nil, errors.New("Error getting document")
	}

	if docId != nil {
		return docId, err
	}

	docId, err = AddDoc(driver, body.Uid, body.Doc.DocName)

	WriteDocFS(db, toString(docId.(int64)), body.Doc.DocText, body.Doc.RawDocText)

	return 	docId, err
}

func ReadDoc(driver neo4j.Driver, db Firestore, body User) (interface{}, error) {
	if body.Doc.DocId == "" {
		return  nil, errors.New("Reading document requires document ID")
	}

	kws, kerr := GetDoc(driver, body.Doc.DocId)
	if kerr != nil {
		return nil, kerr
	}
	texts, terr := ReadDocFS(db, body.Doc.DocId) 
	if terr != nil {
		return nil, terr
	}

	v := []interface{} {
		kws,
		texts,
	}

	fmt.Println("data", v)
	return v, nil
}
