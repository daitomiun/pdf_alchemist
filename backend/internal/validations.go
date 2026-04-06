package internal

import "fmt"

func Validate(cmds []Command) []string {
	var errors []string

	for i, cmd := range cmds {
		step := fmt.Sprintf("Step %d. The %s command needs", i, cmd.Type)
		fmt.Printf("cmd %s \n", cmd.Type)
		switch cmd.Type {
		case DELETE:
			fmt.Printf("contents -> %v \n", cmd)
			if cmd.Page == nil {
				fmt.Println("Nil page contents")
				errors = append(errors, fmt.Sprintf("%s the page parameter", step))
			}
		case SWAP:
			if cmd.PageA == nil && cmd.PageB == nil {
				errors = append(errors, fmt.Sprintf("%s the PageA and PageB parameters", step))
			}
		case SPLIT:
			if cmd.StartCut == nil && cmd.EndCut == nil {
				errors = append(errors, fmt.Sprintf("%s the start and EndCut parameters", step))
			}
		default:
			continue
		}
	}
	fmt.Printf("errs -> %v \n", errors)
	return errors
}
