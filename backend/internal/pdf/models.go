package pdf

import (
	"bytes"
	"mime/multipart"

	"github.com/google/uuid"
)

type Document struct {
	File     *bytes.Reader
	Commands []Command
}

// The manager will handle the group split logic and the original file changes
type PdfManager struct {
	MainFile *bytes.Buffer
	Pages    []Page
	Groups   map[uuid.UUID]*bytes.Buffer
}

type Page struct {
	PageNum  int
	GroupIds []uuid.UUID
}

type ActionType string

const (
	DELETE ActionType = "DELETE"
	SWAP   ActionType = "SWAP"
	SPLIT  ActionType = "SPLIT"
	RESET  ActionType = "RESET"
)

type Body struct {
	FileHeader *multipart.FileHeader `form:"File" binding:"required"`
	Commands   string                `form:"commands" binding:"required"`
}

type Command struct {
	Type ActionType `json:"type"`
	// DELETE
	Page *int `json:"page,omitempty"`

	// SWAP
	PageA *int `json:"pageA,omitempty"`
	PageB *int `json:"pageB,omitempty"`

	// SPLIT
	StartCut *int `json:"startCut,omitempty"`
	EndCut   *int `json:"endCut,omitempty"`
}
