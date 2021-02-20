package db

import (
	"strconv"
	"errors"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func GetUid(driver neo4j.Driver, uid string) (interface{}, error) {
	// null if user does not exist
	session := driver.NewSession(neo4j.SessionConfig{AccessMode:neo4j.AccessModeRead})
	defer session.Close()

	userId, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			`
			MATCH (u:User {user: $authUid}) 
			RETURN id(u)
			`,
			map[string]interface{}{"authUid":uid})

		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().Values[0], nil
		}

		return nil, result.Err()
	})
	return userId, err
}

func AddUser(driver neo4j.Driver, uid string) (interface{}, error) {
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	
	userId, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			`
			CREATE (u:User {user: $authUid}) 
			RETURN id(u)
			`,
			map[string]interface{}{"authUid":uid})

		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().Values[0], nil
		}

		return nil, result.Err()
	})
	return userId, err
}

func Signup(driver neo4j.Driver, db Firestore, body User) (interface{}, error) {
	uid, err := GetUid(driver, body.AuthUid)
	if err != nil {
		return nil, errors.New("Error in signup on matching on uid")
	}

	if uid != nil {
		return nil, errors.New("User already has an account")
	}

	uid, err = AddUser(driver, body.AuthUid)

	db.Client.Collection("users").Doc(toString(uid.(int64))).Set(db.Ctx, map[string]interface{}{
		"uid": toString(uid.(int64)),
	})
	if err != nil {
		return nil, err
	}

	return uid, err 
}

func Login(driver neo4j.Driver,  db Firestore, body User) (interface{}, error) {
	uid, err := GetUid(driver, body.AuthUid)
	if err != nil {
		return nil, errors.New("Error on login")
	}

	if uid == nil {
		return nil, errors.New("User does not exist")
	}

	return uid, err
}

func CheckUid(driver neo4j.Driver, uid string) (bool) {
	session := driver.NewSession(neo4j.SessionConfig{AccessMode:neo4j.AccessModeRead})
	defer session.Close()

	userId, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			`
			MATCH (u:User) WHERE id(u) == $uid 
			RETURN id(u)
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
		return false
	}

	if userId != nil {
		return true
	}
	return false
}

func toString(i int64) (string) {
	return strconv.FormatInt(i, 10)
}