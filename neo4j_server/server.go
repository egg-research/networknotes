package main

import (
	"net/http"
	"encoding/json"
	"fmt"
	"log"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/daaku/go.httpgzip"
)

func helloWorldHandler(driver neo4j.Driver) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
		defer session.Close()

		resp, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
			result, err := transaction.Run(
				"CREATE (a:Greeting) SET a.message = $message RETURN a.message + ', from node ' + id(a)",
				map[string]interface{}{"message": "hello, world"})
			if err != nil {
				return nil, err
			}

			if result.Next() {
				return result.Record().Values[0], nil
			}

			return nil, result.Err()
		})

		if err != nil {
			log.Println("error querying graph:", err)
			return
		}
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			log.Println("error writing graph response:", err)
		}
	}
}

func addUserHandler(driver neo4j.Driver) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
		defer session.Close()

		resp, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
			fmt.Println("addUser:", req.URL.Query()["uid"][0])
			result, err := transaction.Run(
				"CREATE (a:User) SET a.uid = $uid RETURN a.uid + ', from node ' + id(a)",
				map[string]interface{}{"uid":req.URL.Query()["uid"][0]})

			if err != nil {
				return nil, err
			}

			if result.Next() {
				return result.Record().Values[0], nil
			}

			return nil, result.Err()
		})

		if err != nil {
			log.Println("error querying graph:", err)
			return
		}
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			log.Println("error writing graph response:", err)
		}
	}
}

func addDocumentHandler(driver neo4j.Driver) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
		defer session.Close()

		resp, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
			fmt.Println("addUser:", req.URL.Query()["uid"][0])
			result, err := transaction.Run(
				`
				CREATE (a:Document) SET a.doc = $doc 
				CREATE ($uid)-[:DOCUMENT]->()
				RETURN a.uid + ', from node ' + id(a)
				`,
				map[string]interface{}{
					"uid":req.URL.Query()["uid"][0],
					"doc": req.URL.Query()["doc"][0],
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
			log.Println("error querying graph:", err)
			return
		}
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			log.Println("error writing graph response:", err)
		}
	}
}



func clearHandler(driver neo4j.Driver) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
		defer session.Close()

		resp, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
			result, err := transaction.Run(
				"MATCH (n) DETACH DELETE n",
				map[string]interface{}{})
			if err != nil {
				return nil, err
			}

			if result.Next() {
				return result.Record().Values[0], nil
			}

			return nil, result.Err()
		})


		if err != nil {
			log.Println("error querying graph:", err)
			return
		}
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			log.Println("error writing graph response:", err)
		}
	}
}



func main() {
	// sandbox neo4j, only will last 3 days
	url := "bolt://3.92.216.55:7687"
	user := "neo4j"
	pass := "tempers-dyes-orders"

	driver, err := neo4j.NewDriver(url, neo4j.BasicAuth(user, pass, "" /*database*/))
	defer driver.Close()
	if err != nil {
		log.Println("error setting up driver:", err)
		return
	}

	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/hello", helloWorldHandler(driver))
	serveMux.HandleFunc("/clear", clearHandler(driver))

	serveMux.HandleFunc("/addUser", addUserHandler(driver))
	serveMux.HandleFunc("/addDoc", addDocumentHandler(driver))




	fmt.Printf("Running on localhost:8080")
	panic(http.ListenAndServe(":8080", httpgzip.NewHandler(serveMux)))
}

