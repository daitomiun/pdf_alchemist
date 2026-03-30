package internal

import (
	"bytes"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

type Conf struct {
	pdfConf model.Configuration
}

type Document struct {
	File *bytes.Reader
}

// TODO: create body document struct
