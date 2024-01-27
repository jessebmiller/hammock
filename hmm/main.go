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
  done RANK
  new [project]
  remove RANK "QUOTED CARD HEADLINE"
  init
`

type Hammock struct {
	Path                   string
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
		e := fmt.Errorf(
			"Error decoding HAMMOCK_PATH=\"%v\": %v",
			hammock_path,
			err,
		)
		return Hammock{}, e
	}
	hmm.Path = hammock_path
	return hmm, nil
}

func (hmm Hammock) Workspaces() ([]workspace, error) {
	hmmRoot := filepath.Dir(hmm.Path)
	files, err := os.ReadDir(hmmRoot)
	if err != nil {
		return []workspace{}, err
	}
	var workspaces []workspace
	for _, f := range files {
		ws, err := readWorkspace(filepath.Join(hmmRoot, f.Name()))
		if err != nil {
			continue
		}
		workspaces = append(workspaces, ws)
	}
	return workspaces, nil
}

func main() {
	if len(os.Args) == 1 {
		check(summarize())
		return
	}

	action := os.Args[1]

	// actions that don't need a current workspace
	switch action {
	case "go":
		check(cdToWorkspace(os.Args[2:]))
		return
	case "init":
		check(initWorkspace(os.Args[2:]))
		return
	}

	var actionFunc func([]string, workspace) error
	var showList bool
	switch action {
	case "list":
		actionFunc = list
	case "show":
		actionFunc = show
	case "rank":
		actionFunc = rank
		showList = true
	case "new":
		actionFunc = create
		showList = true
	case "remove":
		actionFunc = remove
		showList = true
	case "done":
		actionFunc = done
		showList = true
	default:
		fmt.Println("Unknown action", action)
		fmt.Print(usage)
		return
	}

	ws, err := currentWorkspace()
	check(err)
	check(actionFunc(os.Args[2:], ws))
	if showList {
		check(list([]string{}, ws))
	}
}
