package main

import (
	"net/http"
	"encoding/json"
	"fmt"
	"log"
	"io"
	"io/ioutil"
	"strings"
	"errors"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/daaku/go.httpgzip"
	"github.com/golang/gddo/httputil/header"
												
    "context"                       // https://blog.golang.org/context
	firebase "firebase.google.com/go"
    // "google.golang.org/api/option"
	
	. "main/db"
	. "main/ml"
)

func Max(x, y int) int {
    if x < y {
        return y
    }
    return x
}

func GetAll(driver neo4j.Driver, db Firestore, uid int, req *http.Request) (interface{}, error) {
	if !CheckUid(driver, uid) {
		return nil, errors.New("User does not exist")
	}

	session := driver.NewSession(neo4j.SessionConfig{AccessMode:neo4j.AccessModeRead})
	defer session.Close()

	nodes := make([]interface{}, 0)
	linkMap := make(map[int]string, 0)
	nodeIds := make(map[int]int, 0)
	links := make([]interface{}, 0)

	kwMap := make(map[int]string, 0)
	max := 1

	_, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		var query string
		
		docquery := 
		`
		MATCH (u:User) WHERE id(u) = $uid
		MATCH (d:Document) <-- (u)  
		OPTIONAL MATCH (d) -[r:KEYWORD]-> (k:Keyword) <-[s:KEYWORD]- (e:Document) <-- (u) WHERE id(d) < id(e)
		RETURN id(d), id(e), id(k), k.kw
		`

		kwquery := 
		`
		MATCH (u:User) WHERE id(u) = $uid
		MATCH (k:Keyword) <-- (u)  
		OPTIONAL MATCH (k) <-[r:KEYWORD]- (d:Document) -[s:KEYWORD]-> (l:Keyword) <-- (u) WHERE id(k) < id(l)
		RETURN id(k), id(l), id(d), k.kw, l.kw
		`

		if req.URL.Query()["q"][0] == "kw" {
			query = kwquery
		} else if req.URL.Query()["q"][0] == "doc" {
			query = docquery
		} else {
			return nil, errors.New("query param string not recognized")
		}

		result, err := transaction.Run(
			query,
			map[string]interface{}{"uid":uid})

			if err != nil {
				return nil, err
			}

			collection, err := result.Collect()
			if err != nil {
				return nil, err
			}

			for _, record := range collection {
				if (req.URL.Query()["q"][0] == "doc") {
					fmt.Println(record.Values)
					_,ok := nodeIds[ToInt(record.Values[0])]
					if !ok {
						nodeIds[ToInt(record.Values[0])] = 0
					}

					if (record.Values[0] == record.Values[1]) {
						// same document
						continue
					}

					if record.Values[1] == nil {
						// no links
						continue
					}

					_,ok = nodeIds[ToInt(record.Values[1])]
					if !ok {
						nodeIds[ToInt(record.Values[1])] = 0
					}

					sourceId := ToInt(record.Values[0])
					destId := ToInt(record.Values[1])
					kwId := ToInt(record.Values[2])
					kw := ToString(record.Values[3])

					// add source node
					nodeIds[sourceId] += 1
					max = Max(max, nodeIds[sourceId])

					// add dest node
					nodeIds[destId] += 1
					max = Max(max, nodeIds[sourceId])

					linkMap[kwId] = kw
					
					link := make(map[string]interface{},0)
					link["id"] = kwId
					link["name"] = kw
					link["source"] = sourceId
					link["target"] = destId

					links = append(links, link)
				} else {
					fmt.Println(record.Values)

					// id(k), id(l), id(d), k.kw, l.kw

					_,ok := nodeIds[ToInt(record.Values[0])]
					if !ok {
						nodeIds[ToInt(record.Values[0])] = 0
					}

					kwMap[ToInt(record.Values[0])] = ToString(record.Values[3])

					if (record.Values[0] == record.Values[1]) {
						// same document
						continue
					}

					if record.Values[1] == nil {
						// no links
						continue
					}

					_,ok = nodeIds[ToInt(record.Values[1])]
					if !ok {
						nodeIds[ToInt(record.Values[1])] = 0
					}

					sourceId := ToInt(record.Values[0])
					destId := ToInt(record.Values[1])
					docId := ToInt(record.Values[2])
					kw1 := ToString(record.Values[3])
					kw2 := ToString(record.Values[4])

					kwMap[sourceId] = kw1 
					kwMap[destId] = kw2 

					// add source node
					nodeIds[sourceId] += 1
					max = Max(max, nodeIds[sourceId])

					// add dest node
					nodeIds[destId] += 1
					max = Max(max, nodeIds[sourceId])

					texts, terr := ReadDocFS(driver, db, uid, docId) 
					if terr != nil {
						return nil, terr
					}

					link := make(map[string]interface{},0)
					link["id"] = docId
					link["name"] = texts.(map[string]interface{})["docName"].(string)
					link["source"] = sourceId
					link["target"] = destId

					links = append(links, link)
				}
			}
			return nil, nil
		})
	
	if err != nil {
		return nil, err
	}

	if (req.URL.Query()["q"][0] == "doc") {
		for docId, weight := range nodeIds {
			texts, terr := ReadDocFS(driver, db, uid, docId) 
			if terr != nil {
				return nil, terr
			}
			node := make(map[string]interface{},0)
			node["id"] = docId
			node["name"] = texts.(map[string]interface{})["docName"].(string)
			node["weight"] = weight / max
			nodes = append(nodes, node)
		}
	} else {
		for kwId, weight := range nodeIds {
			node := make(map[string]interface{},0)
			node["id"] = kwId
			node["name"] = kwMap[kwId]
			node["weight"] = weight / max
			nodes = append(nodes, node)
		}
	}
	
	resMap := make(map[string][]interface{}, 0)
	resMap["nodes"] = nodes
	resMap["links"] = links

	return resMap, err
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
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-Auth-Token")
		if req.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
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
		case "writeDoc":
			resp, err = WriteDoc(driver, db, body)
		case "writeDocFS":
			resp, err = WriteDocFS(driver, db, body.Uid, body.Doc.DocId, body.Doc.DocName, body.Doc.DocText, body.Doc.RawDocText)
		case "readDoc":
			resp, err = ReadDoc(driver, db, body)
		case "readDocFS":
			resp, err = ReadDocFS(driver, db, body.Uid, body.Doc.DocId)
		case "getAllDocs":
			resp, err = GetAllDocs(driver, db, body.Uid)
		case "writeKw":
			resp, err = WriteKw(driver, db, body)
		case "getAllKws":
			resp, err = GetAllKws(driver, db, body.Uid)
		case "getAll":
			resp, err = GetAll(driver, db, body.Uid, req)
		case "delKw":
			resp, err = DelKw(driver, db, body)
		case "related":
			resp, err = MLRelated(driver, db, body)
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
	conf := &firebase.Config{ProjectID: "networknotes2"}
	app, err := firebase.NewApp(ctx, conf)
	// sa := option.WithCredentialsFile("../credentials2.json")
	// app, err := firebase.NewApp(ctx, nil, sa)
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
	serveMux.HandleFunc("/writeDoc", requestHandler(driver, db, "writeDoc"))
	serveMux.HandleFunc("/writeDocFS", requestHandler(driver, db, "writeDocFS"))
	serveMux.HandleFunc("/readDoc", requestHandler(driver, db, "readDoc"))
	serveMux.HandleFunc("/readDocFS", requestHandler(driver, db, "readDocFS"))
	serveMux.HandleFunc("/getAllDocs", requestHandler(driver, db, "getAllDocs"))
	serveMux.HandleFunc("/getAllKws", requestHandler(driver, db, "getAllKws"))
	serveMux.HandleFunc("/getAll", requestHandler(driver, db, "getAll"))
	serveMux.HandleFunc("/writeKw", requestHandler(driver, db, "writeKw"))
	serveMux.HandleFunc("/delKw", requestHandler(driver, db, "delKw"))

	serveMux.HandleFunc("/related", requestHandler(driver, db, "related"))

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

	doc := Document {
		DocId: -1,
	}

	u = User {
		Doc: doc,
	}

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "unable to parse request body"
		http.Error(w, msg, http.StatusBadRequest)
		return u, err
	}

    err = json.Unmarshal(bytes, &u)
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