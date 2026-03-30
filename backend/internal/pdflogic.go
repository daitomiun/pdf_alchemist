package internal

import (
	"bytes"
	"fmt"
	"math"
	"strconv"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func (c *Conf) Delete(d Document, pages []string) *bytes.Buffer {
	var newFile bytes.Buffer

	api.RemovePages(d.File, &newFile, pages, &c.pdfConf)
	return &newFile
}

func (c *Conf) Split(d Document, start, end int) *bytes.Buffer {
	var newFile bytes.Buffer

	var pages []string

	from := math.Min(float64(start), float64(end))
	to := math.Max(float64(start), float64(end))

	for i := from; i < to; i++ {
		fmt.Printf("cycle -> %v \n", i)
		pages = append(pages, strconv.Itoa(int(i)))
	}

	api.Trim(d.File, &newFile, pages, &c.pdfConf)
	return &newFile
}

func (c *Conf) Swap(d Document, pageA, pageB int) {

}
