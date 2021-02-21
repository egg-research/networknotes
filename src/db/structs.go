package db

import (
	"cloud.google.com/go/firestore" // https://godoc.org/cloud.google.com/go/firestore"
    "context"                       // https://blog.golang.org/context
)

type Document struct {
	DocName string  `json:"docName,omitempty"`
	DocId int `json:"docId,omitempty"`
	RawDocText interface{} `json:"rawDocText,omitempty"`
	DocText string `json:"docText,omitempty"`
} 

type Keyword struct {
	Kw string `json:"kw,omitempty"`
	KwId int `json:"kwId,omitempty"`
	KwText string `json:"kwText,omitempty"`
}

type User struct {
	AuthUid string `json:"authUid,omitempty"`
	Uid int `json:"uid,omitempty"`
	Doc Document `json:"doc,omitempty"`
	Kws []Keyword `json:"kws,omitempty"`
}

type Firestore struct {
	Client *firestore.Client
	Ctx    context.Context
}