package pdf

import (
	"bytes"
	"errors"
	"io"
	"log"
	"math"
	"os"
	"slices"
	"strconv"

	"github.com/google/uuid"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func ProcessCommands(d Document) error {
	// TODO:
	// 1. After transforming the document, write output to the file server via gin
	pdfManager, err := Transform(d)
	if err != nil {
		return err
	}
	filepath := ""
	err = copyToPath(pdfManager.MainFile, filepath)
	if err != nil {
		return err
	}
	for _, group := range pdfManager.Groups {
		err = copyToPath(group, filepath)
		if err != nil {
			return err
		}
	}
	return nil
	// 2. Save metadata to tables (also the json commands!)
	// 3. return json response or error response
}

func copyToPath(buf *bytes.Buffer, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, buf)
	if err != nil {
		return err
	}

	return nil
}

func Transform(d Document) (PdfManager, error) {
	cfg := model.NewDefaultConfiguration()
	data, err := io.ReadAll(d.File)
	if err != nil {
	}
	d.File.Seek(0, io.SeekStart)

	pages, err := GenerateTotalPages(*cfg, bytes.NewReader(data))
	if err != nil {
		log.Println("")
		return PdfManager{}, err
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
				return PdfManager{}, errors.New("Delete failed, cannot find page")
			}

			updatedFile, err := Delete(*cfg, bytes.NewReader(manager.MainFile.Bytes()), *cmd.Page)
			if err != nil {
				log.Println("Delete failed")
				return PdfManager{}, errors.New("Could not delete Page")
			}

			idx := slices.IndexFunc(manager.Pages, func(p Page) bool { return p.PageNum == *cmd.Page })
			if idx == -1 {
				log.Println("Delete Failed, couldn't find page index")
				return PdfManager{}, errors.New("Delete Failed, couldn't find page index")
			}
			groups := manager.Pages[idx].GroupIds

			for _, id := range groups {
				groupFile, err := Delete(*cfg, bytes.NewReader(manager.Groups[id].Bytes()), *cmd.Page)
				if err != nil {
					log.Println("Group File delete failed")
					return PdfManager{}, err
				}
				manager.Groups[id] = groupFile
			}
			manager.Pages = slices.Delete(manager.Pages, idx, idx+1)
			manager.MainFile = bytes.NewBuffer(updatedFile.Bytes())
		case SPLIT:
			newFile, err := SplitPages(*cfg, bytes.NewReader(manager.MainFile.Bytes()), *cmd.StartCut, *cmd.EndCut)
			if err != nil {
				log.Println("Split failed")
				return PdfManager{}, err
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
				return PdfManager{}, err
			}
			manager.MainFile = bytes.NewBuffer(newFile.Bytes())

			pages := manager.Pages
			idxA := slices.IndexFunc(manager.Pages, func(p Page) bool { return p.PageNum == *cmd.PageA })
			idxB := slices.IndexFunc(manager.Pages, func(p Page) bool { return p.PageNum == *cmd.PageB })
			if idxA == -1 || idxB == -1 {
				log.Println("The pages do not exist")
				return PdfManager{}, errors.New("The Page does not exists")
			}
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
				newFile, err := TrimPages(*cfg, bytes.NewReader(manager.MainFile.Bytes()), p)
				if err != nil {
					log.Println("The swap failed")
					return PdfManager{}, err
				}
				manager.Groups[id] = newFile
			}
		default:
			continue
		}
	}
	return manager, nil
}

func GenerateTotalPages(cfg model.Configuration, doc *bytes.Reader) ([]Page, error) {
	total, err := api.PageCount(doc, &cfg)
	if err != nil {
		return nil, err
	}

	pages := make([]Page, 0, total)
	for i := range total {
		pages = append(pages, Page{PageNum: i + 1})
	}

	return pages, nil
}
