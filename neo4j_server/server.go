package main

import (
	"net/http"
	"encoding/json"
	"fmt"
	"log"
	"io"
	"strings"
	"errors"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/daaku/go.httpgzip"
	"github.com/golang/gddo/httputil/header"
)

type User struct {
	Uid string `json:"uid"`
	Doc string `json:"doc,omitempty"`
	Kws []string `json:"kws,omitempty"`
}

func addUser(session neo4j.Session, body User) (interface{}, error) {
	resp, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		fmt.Println("addUser:", body)
		result, err := transaction.Run(
			`
			MERGE (u:User {user: $uid})
			`,
			map[string]interface{}{"uid":body.Uid})

		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().Values[0], nil
		}

		return nil, result.Err()
	})
	return resp, err
}

func addDoc(session neo4j.Session, body User) (interface{}, error) {
	resp, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		fmt.Println("addDoc:", body)
		result, err := transaction.Run(
			`
			MERGE (d:Document {doc: $doc}) 
			MERGE (u:User {user: $uid})
			WITH u,d
			MERGE (u)-[:DOCUMENT]->(d)
			`,
			map[string]interface{}{
				"uid":body.Uid,
				"doc":body.Doc,
			})

		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().Values[0], nil
		}

		return nil, result.Err()
	})
	return resp, err
}

func addKeywords(session neo4j.Session, body User) (interface{}, error) {
	resp, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		fmt.Println("addKeywords:", body)
		result, err := transaction.Run(
			`
			UNWIND $kws as kw
			MERGE (k:Keyword {kw: kw}) 
			MERGE (d:Document {doc: $doc}) 
			MERGE (u:User {user: $uid})
			WITH u,d,k
			MERGE (u)-[:DOCUMENT]->(d)
			MERGE (k)-[:DOCUMENT]->(d)
			MERGE (u)-[:KEYWORD]->(k)
			MERGE (d)-[:KEYWORD]->(k)
			`,
			map[string]interface{}{
				"uid":body.Uid,
				"doc":body.Doc,
				"kws":body.Kws,
			})

		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().Values[0], nil
		}

		return nil, result.Err()
	})
	return resp, err
}



func clear(session neo4j.Session, body User) (interface{}, error) {
	resp, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		fmt.Println("Deleting database!")
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
	return resp, err
}

func requestHandler(driver neo4j.Driver, request string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		req.Body = http.MaxBytesReader(w, req.Body, 1048576)
		
		
		body, err := decodeRequestBody(w, req)
		if err != nil {
			return
		}


		session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
		defer session.Close()

		var resp interface{}
		
		switch request {
		case "user":
			resp, err = addUser(session, body)
		case "doc":
			resp, err = addDoc(session, body)
		case "kw":
			resp, err = addKeywords(session, body)
		case "clear":
			resp, err = clear(session, body)
		default:
			log.Println("error unknown request type:", request)
		}

		if err != nil {
			log.Println("error querying:", request, err)
			return
		}
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			log.Println("error writing response:", request, err)
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
	serveMux.HandleFunc("/clear", requestHandler(driver, "clear"))
	serveMux.HandleFunc("/addUser", requestHandler(driver, "user"))
	serveMux.HandleFunc("/addDoc", requestHandler(driver, "doc"))
	serveMux.HandleFunc("/addKeywords", requestHandler(driver, "kw"))




	fmt.Println("Running on localhost:8080")
	panic(http.ListenAndServe(":8080", httpgzip.NewHandler(serveMux)))
}




func decodeRequestBody(w http.ResponseWriter, r *http.Request) (User, error) {
	var u User

    // If the Content-Type header is present, check that it has the value
    // application/json. Note that we are using the gddo/httputil/header
    // package to parse and extract the value here, so the check works
    // even if the client includes additional charset or boundary
    // information in the header.
    if r.Header.Get("Content-Type") != "" {
        value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
        if value != "application/json" {
            msg := "Content-Type header is not application/json"
            http.Error(w, msg, http.StatusUnsupportedMediaType)
            return u, errors.New("error")
        }
    }

    // Use http.MaxBytesReader to enforce a maximum read of 1MB from the
    // response body. A request body larger than that will now result in
    // Decode() returning a "http: request body too large" error.
    r.Body = http.MaxBytesReader(w, r.Body, 1048576)

    // Setup the decoder and call the DisallowUnknownFields() method on it.
    // This will cause Decode() to return a "json: unknown field ..." error
    // if it encounters any extra unexpected fields in the JSON. Strictly
    // speaking, it returns an error for "keys which do not match any
    // non-ignored, exported fields in the destination".
    dec := json.NewDecoder(r.Body)
    dec.DisallowUnknownFields()

    err := dec.Decode(&u)
    if err != nil {
        var syntaxError *json.SyntaxError
        var unmarshalTypeError *json.UnmarshalTypeError

        switch {
        // Catch any syntax errors in the JSON and send an error message
        // which interpolates the location of the problem to make it
        // easier for the client to fix.
        case errors.As(err, &syntaxError):
            msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
            http.Error(w, msg, http.StatusBadRequest)

        // In some circumstances Decode() may also return an
        // io.ErrUnexpectedEOF error for syntax errors in the JSON. There
        // is an open issue regarding this at
        // https://github.com/golang/go/issues/25956.
        case errors.Is(err, io.ErrUnexpectedEOF):
            msg := fmt.Sprintf("Request body contains badly-formed JSON")
            http.Error(w, msg, http.StatusBadRequest)

        // Catch any type errors, like trying to assign a string in the
        // JSON request body to a int field in our Person struct. We can
        // interpolate the relevant field name and position into the error
        // message to make it easier for the client to fix.
        case errors.As(err, &unmarshalTypeError):
            msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
            http.Error(w, msg, http.StatusBadRequest)

        // Catch the error caused by extra unexpected fields in the request
        // body. We extract the field name from the error message and
        // interpolate it in our custom error message. There is an open
        // issue at https://github.com/golang/go/issues/29035 regarding
        // turning this into a sentinel error.
        case strings.HasPrefix(err.Error(), "json: unknown field "):
            fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
            msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
            http.Error(w, msg, http.StatusBadRequest)

        // An io.EOF error is returned by Decode() if the request body is
        // empty.
        case errors.Is(err, io.EOF):
            msg := "Request body must not be empty"
            http.Error(w, msg, http.StatusBadRequest)

        // Catch the error caused by the request body being too large. Again
        // there is an open issue regarding turning this into a sentinel
        // error at https://github.com/golang/go/issues/30715.
        case err.Error() == "http: request body too large":
            msg := "Request body must not be larger than 1MB"
            http.Error(w, msg, http.StatusRequestEntityTooLarge)

        // Otherwise default to logging the error and sending a 500 Internal
        // Server Error response.
        default:
            log.Println(err.Error())
            http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
        }
        return u, errors.New("error")
    }

    // Call decode again, using a pointer to an empty anonymous struct as 
    // the destination. If the request body only contained a single JSON 
    // object this will return an io.EOF error. So if we get anything else, 
    // we know that there is additional data in the request body.
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
        msg := "Request body must only contain a single JSON object"
        http.Error(w, msg, http.StatusBadRequest)
        return u, errors.New("error")
    }

    return u, nil
}