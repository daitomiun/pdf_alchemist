package internal

import (
	"bytes"
	"mime/multipart"
)

type Document struct {
	File     *bytes.Reader
	Commands []Command
}

type ActionType string

const (
	DELETE ActionType = "DELETE"
	SWAP   ActionType = "SWAP"
	SPLIT  ActionType = "SPLIT"
	RESET  ActionType = "RESET"
)

type Body struct {
	File     *multipart.FileHeader `form:"File" binding:"required"`
	Commands string                `form:"commands" binding:"required"`
}

type Command struct {
	Type ActionType `json:"type"`
	// DELETE
	Page *Page `json:"page,omitempty"`

	// SWAP
	PageA *Page `json:"pageA,omitempty"`
	PageB *Page `json:"pageB,omitempty"`

	// SPLIT
	StartCut *int `json:"startCut,omitempty"`
	EndCut   *int `json:"endCut,omitempty"`
}

type Page struct {
	Id      string `json:"id"`
	PageNum int    `json:"pageNum"`
}
