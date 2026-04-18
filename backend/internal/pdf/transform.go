package pdf

import (
	"bytes"
	"io"
	"log"
	"math"
	"strconv"

	"github.com/google/uuid"
	"github.com/pdfcpu/pdfcpu/pkg/api"
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

	pages, err := GenerateTotalPages(*cfg, bytes.NewReader(data))
	if err != nil {
		log.Println("")
		return
	}

	manager := PdfManager{
		MainFile:  bytes.NewBuffer(data),
		MainPages: pages,
		Groups:    map[uuid.UUID]Group{},
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
			manager.MainFile = bytes.NewBuffer(updatedFile.Bytes())
			// NOTE: Update the Pages from our map
			delete(manager.MainPages, *cmd.Page)

			for id, group := range manager.Groups {
				_, pageExists := group.Pages[*cmd.Page]
				if pageExists {
					groupFile, err := Delete(*cfg, bytes.NewReader(manager.Groups[id].File.Bytes()), *cmd.Page)
					if err != nil {
						log.Println("Group File delete failed")
						return
					}
					group.File = groupFile
					// NOTE: Update the Pages from our map
					delete(group.Pages, *cmd.Page)
					manager.Groups[id] = group
				}
			}
		case SPLIT:
			newFile, err := SplitPages(*cfg, bytes.NewReader(manager.MainFile.Bytes()), *cmd.StartCut, *cmd.EndCut)
			if err != nil {
				log.Println("Split failed")
				return
			}
			manager.Groups[uuid.New()] = Group{File: newFile, Pages: GenerateSplitPages(*cmd.StartCut, *cmd.EndCut)}
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

func GenerateSplitPages(start, end int) map[int]struct{} {
	from := math.Min(float64(start), float64(end))
	to := math.Max(float64(start), float64(end))

	pages := make(map[int]struct{}, int(to-from))
	for i := int(from); i < int(to); i++ {
		pages[i] = struct{}{}
	}

	return pages
}

func updateSwapState(pages map[int]struct{}, pageA, pageB int) map[int]struct{} {

	// TODO: check that any groups have the specified page
	// return the new list of pages from the main file
	return nil
}

func updateGroupState(groups map[uuid.UUID]Group) Group {

	return Group{}
}

func GenerateTotalPages(cfg model.Configuration, doc *bytes.Reader) (map[int]struct{}, error) {
	total, err := api.PageCount(doc, &cfg)
	if err != nil {
		return nil, err
	}

	pages := make(map[int]struct{}, total)
	for i := range total {
		pages[i+1] = struct{}{}
	}

	return pages, nil
}
