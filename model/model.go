package model

type DocType string

const USER DocType = "user"
const ITEM DocType = "item"


type Document interface {
	Identifier() string
}

