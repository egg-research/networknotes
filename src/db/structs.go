package db

import (
	"cloud.google.com/go/firestore" // https://godoc.org/cloud.google.com/go/firestore"
    "context"                       // https://blog.golang.org/context
)

type Document struct {
	DocId string  `json:"docId,omitempty"`
	DocName string  `json:"docName,omitempty"`
	RawDocText interface{} `json:"rawDocText,omitempty"`
	DocText string `json:"docText,omitempty"`
} 

type Keyword struct {
	Kw string `json:"kw"`
	KwText string `json:"kwText"`
}

type User struct {
	Uid string `json:"uid,omitempty"`
	Doc Document `json:"doc,omitempty"`
	Kws []Keyword `json:"kws,omitempty"`
}

type Firestore struct {
	Client *firestore.Client
	Ctx    context.Context
}