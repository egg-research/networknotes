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
	"github.com/fatih/structs"
												
    "context"                       // https://blog.golang.org/context
	firebase "firebase.google.com/go"
    "google.golang.org/api/option"
	
	. "main/db"
)

func kw(driver neo4j.Driver, db Firestore, body User) (interface{}, error) {
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	resp, err := addKeywords(session, body)
	if (err != nil) {
		return nil, err
	}
	return resp, err
}

func mapKeywords(kws []Keyword) []map[string]interface{} {
	var result = make([]map[string]interface{}, len(kws))

	for index, item := range kws {
		result[index] = structs.Map(item)
	}
	fmt.Println("kws", result)
	return result
}

func addKeywords(session neo4j.Session, body User) (interface{}, error) {
	resp, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		fmt.Println("addKeywords:", body)
		fmt.Printf("%T\n %T\n %T\n", body.Doc, body.Kws, mapKeywords(body.Kws))
		// MERGE (k:Keyword {kw:kw.Kw})
		// MERGE (d:Document {doc: $doc}) 
		// MERGE (u:User {user: $uid})
		// WITH u,d,k,kw
		// MERGE (u)-[:DOCUMENT]->(d)
		// MERGE (k)-[:DOCUMENT]->(d)
		// MERGE (u)-[:KEYWORD {kwText:kw.KwText}]->(k)
		// MERGE (d)-[:KEYWORD {kwText:kw.KwText}]->(k)
		result, err := transaction.Run(
			`
			UNWIND $kws as kw
			RETURN id(kw)
			`,
			map[string]interface{}{
				"uid":body.Uid,
				"doc":body.Uid+"$"+body.Doc.DocId,
				"kws":mapKeywords(body.Kws),
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

func clear(driver neo4j.Driver, db Firestore, body User) (interface{}, error) {
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
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

func requestHandler(driver neo4j.Driver, db Firestore, request string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		req.Body = http.MaxBytesReader(w, req.Body, 1048576)
		
		body, err := decodeRequestBody(w, req)
		if err != nil {
			return
		}

		var resp interface{}		
		switch request {
		case "signup":
			resp, err = Signup(driver, db, body)
		case "login":
			resp, err = Login(driver, db, body)
		case "getDocId":
			resp, err = GetDocId(driver, body.Uid, body.Doc.DocName)
		case "writeDoc":
			resp, err = WriteDoc(driver, db, body)
		case "writeDocFS":
			resp, err = WriteDocFS(db, body.Doc.DocId, body.Doc.DocText, body.Doc.RawDocText)
		case "readDoc":
			resp, err = ReadDoc(driver, db, body)
		case "kw":
			resp, err = kw(driver, db, body)
		case "clear":
			resp, err = clear(driver, db, body)
		default:
			http.Error(w, fmt.Sprintf("Error unknown request type: %s", request), 500)
			log.Println("error unknown request type:", request)
		}

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error writing response: %s", request), 500)
			return
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

	ctx := context.Background()
	// conf := &firebase.Config{ProjectID: "networknotes-305405"}
	// app, err := firebase.NewApp(ctx, conf)
	sa := option.WithCredentialsFile("networknotes-305405-107f89927bbb.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	db := Firestore{Client: client, Ctx: ctx}

	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/clear", requestHandler(driver, db, "clear"))
	serveMux.HandleFunc("/signup", requestHandler(driver, db, "signup"))
	serveMux.HandleFunc("/login", requestHandler(driver, db, "login"))
	serveMux.HandleFunc("/getDocId", requestHandler(driver, db, "getDocId"))
	serveMux.HandleFunc("/writeDoc", requestHandler(driver, db, "writeDoc"))
	serveMux.HandleFunc("/writeDocFS", requestHandler(driver, db, "writeDocFS"))
	serveMux.HandleFunc("/readDoc", requestHandler(driver, db, "readDos"))
	serveMux.HandleFunc("/kw", requestHandler(driver, db, "kw"))

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