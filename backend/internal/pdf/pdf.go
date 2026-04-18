package pdf

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/pdfcpu/pdfcpu/pkg/api"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func Delete(cfg model.Configuration, doc *bytes.Reader, page int) (*bytes.Buffer, error) {
	var newFile bytes.Buffer

	pages := []string{strconv.Itoa(page)}

	if err := api.RemovePages(doc, &newFile, pages, &cfg); err != nil {
		return nil, err
	}
	return &newFile, nil
}

func SplitPages(cfg model.Configuration, doc *bytes.Reader, start, end int) (*bytes.Buffer, error) {
	var newFile bytes.Buffer

	var pages []string

	from := math.Min(float64(start), float64(end))
	to := math.Max(float64(start), float64(end))

	for i := from; i < to; i++ {
		fmt.Printf("cycle -> %v \n", i)
		pages = append(pages, strconv.Itoa(int(i)))
	}

	if err := api.Trim(doc, &newFile, pages, &cfg); err != nil {
		return nil, err
	}
	return &newFile, nil
}

func SwapPages(cfg model.Configuration, doc *bytes.Reader, pageA, pageB int) (*bytes.Buffer, error) {
	totalPages, err := api.PageCount(doc, &cfg)
	if err != nil {
		return nil, errors.New("Cannot get page count")
	}
	pages := swapOrder(totalPages, pageA, pageB)
	var newFile bytes.Buffer

	api.Trim(doc, &newFile, pages, &cfg)

	return &newFile, nil
}

func swapOrder(totalPages, pageA, pageB int) []string {
	newOrder := make([]string, 0, totalPages)
	for i := 1; i < totalPages; i++ {
		switch i {
		case pageA:
			newOrder = append(newOrder, strconv.Itoa(pageB))
		case pageB:
			newOrder = append(newOrder, strconv.Itoa(pageA))
		default:
			newOrder = append(newOrder, strconv.Itoa(i))
		}
	}
	return newOrder
}
