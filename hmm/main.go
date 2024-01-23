package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const usage = `
Usage: hmm [COMMAND]

Commands:
  list
  show [RANK]
  rank FROM TO
  new
  remove RANK "QUOTED CARD HEADLINE"
`

type Hammock struct {
	Path string
	PriorityWorkspaceNames []string `toml:"priority_workspaces"`
}

// PriorityProjects joins active projects from priority workspaces
func (hmm Hammock) PriorityProjects() ([]project, error) {
	var projects []project
	for _, name := range hmm.PriorityWorkspaceNames {
		hmmRoot := filepath.Dir(hmm.Path)
		wsPath := filepath.Join(hmmRoot, name)
		ws, err := readWorkspace(wsPath)
		if err != nil {
			return []project{}, err
		}
		activeProjects, err := ws.ActiveProjects()
		if err != nil {
			return []project{}, err
		}
		projects = append(projects, activeProjects...)
	}
	return projects, nil
} 

func readHammock() (Hammock, error) {
	hammock_path := filepath.Join(os.Getenv("HAMMOCK_PATH"), "Hammock.toml")
	var hmm Hammock
	_, err := toml.DecodeFile(hammock_path, &hmm)
	if err != nil {
		e :=  fmt.Errorf(
			"Error decoding HAMMOCK_PATH=\"%v\": %v",
			hammock_path,
			err,
		)
		return Hammock{}, e
	}
	hmm.Path = hammock_path
	return hmm, nil
}

func main() {
	if len(os.Args) == 1 {
		check(summarize())
		return
	}

	ws, err := currentWorkspace()
	if err != nil {
		panic(err)
	}

	action := os.Args[1]
	switch action {
	case "list":
		check(list(os.Args[2:], ws))
	case "show":
		check(show(os.Args[2:], ws))
	case "rank":
		check(rank(os.Args[2:], ws))
		check(list([]string{}, ws))
	case "new":
		check(create(os.Args[2:], ws))
	case "remove":
		check(remove(os.Args[2:], ws))
	default:
		fmt.Println("Unknown action", action)
		fmt.Print(usage)
	}
}
