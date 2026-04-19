package pdf

import (
	"bytes"
	"io"
	"log"
	"math"
	"slices"
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
		MainFile: bytes.NewBuffer(data),
		Pages:    pages,
		Groups:   make(map[uuid.UUID]*bytes.Buffer),
	}

	for i, cmd := range d.Commands {
		log.Printf("Step %d. %s \n", i, cmd.Type)

		switch cmd.Type {
		case DELETE:
			pageExists := slices.ContainsFunc(manager.Pages, func(p Page) bool { return p.PageNum == *cmd.Page })
			if !pageExists {
				log.Println("Delete failed, cannot find page")
				return
			}

			updatedFile, err := Delete(*cfg, bytes.NewReader(manager.MainFile.Bytes()), *cmd.Page)
			if err != nil {
				log.Println("Delete failed")
				return
			}

			idx := slices.IndexFunc(manager.Pages, func(p Page) bool { return p.PageNum == *cmd.Page })
			groups := manager.Pages[idx].GroupIds

			for _, id := range groups {
				groupFile, err := Delete(*cfg, bytes.NewReader(manager.Groups[id].Bytes()), *cmd.Page)
				if err != nil {
					log.Println("Group File delete failed")
					return
				}
				manager.Groups[id] = groupFile
			}
			manager.Pages = slices.Delete(manager.Pages, idx, idx+1)
			manager.MainFile = bytes.NewBuffer(updatedFile.Bytes())
		case SPLIT:
			newFile, err := SplitPages(*cfg, bytes.NewReader(manager.MainFile.Bytes()), *cmd.StartCut, *cmd.EndCut)
			if err != nil {
				log.Println("Split failed")
				return
			}
			newId := uuid.New()
			manager.Groups[newId] = newFile

			from := math.Min(float64(*cmd.StartCut), float64(*cmd.EndCut))
			to := math.Max(float64(*cmd.StartCut), float64(*cmd.EndCut))

			// INFO: Update state of pages
			for i, p := range manager.Pages {
				if i >= int(from) && i < int(to) {
					if !slices.Contains(p.GroupIds, newId) {
						manager.Pages[i].GroupIds = append(manager.Pages[i].GroupIds, newId)
					}
				}
			}
		case SWAP:
			newFile, err := SwapPages(*cfg, bytes.NewReader(manager.MainFile.Bytes()), *cmd.PageA, *cmd.PageB)
			if err != nil {
				log.Println("Split failed")
				return
			}
			manager.MainFile = bytes.NewBuffer(newFile.Bytes())

			pages := manager.Pages
			idxA := slices.IndexFunc(manager.Pages, func(p Page) bool { return p.PageNum == *cmd.PageA })
			idxB := slices.IndexFunc(manager.Pages, func(p Page) bool { return p.PageNum == *cmd.PageB })
			groupsA := pages[idxA].GroupIds
			groupsB := pages[idxB].GroupIds

			pages[idxA].GroupIds = groupsB
			pages[idxB].GroupIds = groupsA
			pages[idxA], pages[idxB] = pages[idxB], pages[idxA]
			manager.Pages = pages

			affectedSet := make(map[uuid.UUID]struct{})

			for _, id := range groupsA {
				affectedSet[id] = struct{}{}
			}
			for _, id := range groupsB {
				affectedSet[id] = struct{}{}
			}

			groupedPages := make(map[uuid.UUID][]string)
			for _, p := range manager.Pages {
				for id := range affectedSet {
					if slices.Contains(p.GroupIds, id) {
						groupedPages[id] = append(groupedPages[id], strconv.Itoa(p.PageNum))
					}
				}
			}

			for id, p := range groupedPages {
				newFile := TrimPages(*cfg, bytes.NewReader(manager.MainFile.Bytes()), p)
				manager.Groups[id] = newFile
			}
		default:
			continue
		}

	}

}

func GenerateTotalPages(cfg model.Configuration, doc *bytes.Reader) ([]Page, error) {
	total, err := api.PageCount(doc, &cfg)
	if err != nil {
		return nil, err
	}

	pages := make([]Page, total)
	for i := range total {
		pages = append(pages, Page{PageNum: i + 1})
	}

	return pages, nil
}
