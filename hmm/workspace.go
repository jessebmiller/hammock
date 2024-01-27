package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
)

type workspace struct {
	Name        string `toml:"name"`
	Path        string
	ProjectsDir string `toml:"projects_dir"`
}

// workspace.Write writes the workspace to the file
func (ws workspace) Write() error {
	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(ws)
	if err != nil {
		return err
	}
	return os.WriteFile(ws.Path, []byte(buf.String()), 0644)
}

// workspace.Projects gets all the projects in the workspace
func (ws workspace) Projects() ([]project, error) {
	return readProjects(filepath.Join(ws.Path, ws.ProjectsDir))
}

func hasPassed(t time.Time) bool {
	if t.IsZero() {
		return false
	}
	return t.Before(time.Now())
}

// workspace.ActiveProjects gets all the active projects in the workspace
// all projects with starts in the past and are not complete
func (ws workspace) ActiveProjects() ([]project, error) {
	allProjects, err := ws.Projects()
	if err != nil {
		return []project{}, err
	}
	var activeProjects []project
	for _, p := range allProjects {
		if hasPassed(p.Start) && !hasPassed(p.Completed) {
			activeProjects = append(activeProjects, p)
		}
	}

	return activeProjects, nil
}

// readWorkspace reads a workspace from a path
// path must be a directory with a Workspace.toml file in it
func readWorkspace(path string) (workspace, error) {
	var w workspace
	_, err := toml.DecodeFile(filepath.Join(path, "Workspace.toml"), &w)
	if err != nil {
		return workspace{}, err
	}
	w.Path = path
	return w, nil
}

type NotInWorkspace struct {
	WorkingDir string
}

func (e NotInWorkspace) Error() string {
	return fmt.Sprintf("Working directory (%s) not in any workspace", e.WorkingDir)
}

// currentWorkspace gets the workspace the command is run in
// A workspace is the closest parent with a valid Workspace.toml file
func currentWorkspace() (workspace, error) {
	dir, err := os.Getwd()
	if err != nil {
		return workspace{}, err
	}
	return inWorkspace(dir)
}

func inWorkspace(dir string) (workspace, error) {
	for dir != "/" {
		maybeTOMLpath := filepath.Join(dir, "Workspace.toml")
		_, err := os.Stat(maybeTOMLpath)
		if err == nil {
			return readWorkspace(dir)
		}
		if errors.Is(err, os.ErrNotExist) {
			dir = filepath.Dir(dir)
		}
	}

	wd, err := os.Getwd()
	if err != nil {
		return workspace{}, err
	}
	return workspace{}, NotInWorkspace{wd}
}
