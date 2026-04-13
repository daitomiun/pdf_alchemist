package pdf

import (
	"bytes"
	"io"
	"log"

	"github.com/google/uuid"
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

func Transform(d Document) {
	cfg := model.NewDefaultConfiguration()
	data, _ := io.ReadAll(d.File)
	d.File.Seek(0, io.SeekStart)

	manager := PdfManager{
		MainFile: bytes.NewBuffer(data),
		Groups:   map[uuid.UUID]Group{},
	}

	for i, cmd := range d.Commands {
		log.Printf("Step %d. %s \n", i, cmd.Type)

		switch cmd.Type {
		case DELETE:
			updatedFile, err := Delete(*cfg, bytes.NewReader(manager.MainFile.Bytes()), *cmd.Page)
			if err != nil {
				log.Println("Delete failed")
				return
			}
			for _, group := range manager.Groups {
				_, pageExists := group.Pages[cmd.Page.PageNum]
				if pageExists {
					// TODO: delete the page from the specified group
				}

			}
			// TODO: check that the is is inside the group if it does, delete the page from the split
			manager.MainFile = bytes.NewBuffer(updatedFile.Bytes())
		case SPLIT:
			newFile, err := SplitPages(*cfg, bytes.NewReader(manager.MainFile.Bytes()), *cmd.StartCut, *cmd.EndCut)
			if err != nil {
				log.Println("Split failed")
				return
			}

			// manager.Groups = append(manager.Groups, Group{File: *newFile, GroupId: uuid.New()})

			manager.Groups[uuid.New()] = Group{File: *newFile, Pages: }
		case SWAP:
			newFile, err := SwapPages(*cfg, bytes.NewReader(manager.MainFile.Bytes()), *cmd.PageA, *cmd.PageB)
			if err != nil {
				log.Println("Split failed")
				return
			}
			// TODO: check that if page is inside group Id, modify and swap them (similar to UI)
			manager.MainFile = bytes.NewBuffer(newFile.Bytes())
		default:
			continue
		}

	}

}
