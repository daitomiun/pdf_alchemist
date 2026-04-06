package transform

import (
	"bytes"
	"io"
	"log"

	"github.com/daitonium/pdf_alchemist/backend/internal"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

/* TODO:
// 1. Linearly go through all the commands given of the pdf file and transform them
// > Similar to the transformations in the UI
// 2. For each one of the actions it should
//		- know when to swap the previously deleted file and it's position of the group
//		- create separate files for the split functionality
//	  - for the main file add the separate file without the groups but with the swaps and deletions
//		- NEVER modify the original file, add new ones and make copies of the original
// 3. Once the transformations are complete return the files and send them to store in the file server separately
// 4. the metadata will be saved onto the pg database with
//		- the location of the files, the UUID link to the files, the command outputs
*/

func TransformFile(d internal.Document) {
	cfg := model.NewDefaultConfiguration()
	data, _ := io.ReadAll(d.File)
	d.File.Seek(0, io.SeekStart)

	var mainFile bytes.Buffer
	for i, cmd := range d.Commands {
		log.Printf("Step %d. %s \n", i, cmd.Type)

		switch cmd.Type {
		case internal.DELETE:
			updatedFile, err := internal.Delete(*cfg, bytes.NewReader(data), *cmd.Page)
			if err != nil {
				log.Println("Delete failed")
				return
			}
			mainFile = *updatedFile
		case internal.SPLIT:

		case internal.SWAP:

		default:
			continue
		}

	}

}
