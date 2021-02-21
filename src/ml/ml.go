package ml

import (
	"fmt"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"net/http"
	"errors"
	"strconv"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"

	. "main/db"
)

func MLRelated(driver neo4j.Driver, db Firestore, body User) (map[string]interface{}, error) {
	kws, err := MLKeywords(driver, db, body)

	url := "http://egg-network-notes.appspot.com/related"

	postBody, _ := json.Marshal(map[string]interface{}{
		"header_req": os.Getenv("eggs"),
		"prompt": kws,
	})
	responseBody := bytes.NewBuffer(postBody)
	resp, rerr := http.Post(url, "application/json", responseBody)
	if rerr != nil {
		log.Fatalf("An Error Occured %v", rerr)
		return nil, rerr
	}
	defer resp.Body.Close()
	bytes, serr := ioutil.ReadAll(resp.Body)
	if serr != nil {
	   log.Fatalln(serr)
	   return nil, serr
	}

	var dat map[string]interface{}
	err = json.Unmarshal(bytes, &dat)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{}, 0)
	result["generated"] = strings.Split(dat["text"].(string), ", ")
	result["kws"] = kws
	return result, nil
}

func MLKeywords(driver neo4j.Driver, db Firestore, body User) ([]string, error) {
	if !CheckUid(driver, body.Uid) {
		return nil, errors.New("User does not exist")
	}

	if (!CheckDocId(driver, body.Uid, body.Doc.DocId)) {
		return nil, errors.New("User does not own doc with id: "+strconv.Itoa(body.Doc.DocId))
	}

	// Get Document Text 
	texts, terr := ReadDocFS(driver, db, body.Uid, body.Doc.DocId)
	if terr != nil {
		return nil, terr
	}
	docText := texts.(map[string]interface{})["docText"].(string)
	
	fmt.Println("env", os.Getenv("eggs"))

	url := "http://egg-network-notes.appspot.com/keywords"

	postBody, _ := json.Marshal(map[string]string{
		"header_req":  os.Getenv("eggs"),
		"prompt": docText,
	})
	responseBody := bytes.NewBuffer(postBody)
	resp, rerr := http.Post(url, "application/json", responseBody)
	if rerr != nil {
		log.Fatalf("An Error Occured %v", rerr)
		return nil, rerr
	}
	defer resp.Body.Close()
	
	bytes, serr := ioutil.ReadAll(resp.Body)
	if serr != nil {
	   log.Fatalln(serr)
	   return nil, serr
	}

	var dat map[string]interface{}
    if err := json.Unmarshal(bytes, &dat); err != nil {
        panic(err)
    }

	output := dat["text"].(string)
	output = strings.ReplaceAll(output, "'", "")
	output = strings.ReplaceAll(output, "[", "")
	output = strings.ReplaceAll(output, "]", "")
	arr := strings.Split(output, ", ")

	return arr, nil
}