package main

import (
	"time"
	"path/filepath"
	"sort"
	"os"
	"fmt"

	"github.com/BurntSushi/toml"
)

type project struct {
	Path     string
	Name     string	   `toml:"name"`
	Goal     string    `toml:"goal"`
	Start    time.Time `toml:"start"`
	Deadline time.Time `toml:"deadline"`
	Complete bool      `toml:"complete"`
}

// project.Backlog gets the backlog of the project
// all cards in the project folder are considered to be the project's backlog
func (p project) Backlog() ([]card, error) {
	cards, err := readCards(p.Path)
	if err != nil {
		return []card{}, err
	}
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Priority < cards[j].Priority
	})
	return cards, nil
}

// project.NormalizePriorityRanks makes priority ranks of each card sequential
// starting at 1 and incrementing by 1
// Resolves conflicts by whim
func (p project) NormalizePriorities() error {
	backlog, err := p.Backlog()
	if err != nil {
		return err
	}
	err = WriteConsecutivePriorities(backlog)
	if err != nil {
		return err
	}
	return nil
}

func WriteConsecutivePriorities(cards []card) error {
	for i, card := range cards {
		card.Priority = i+1
		err := card.Write()
		if err != nil {
			return err 
		}
	}
	return nil
}

// readProject reads a project from a path
// path must be a directory with a valid Project.toml file in it
func readProject(path string) (project, error) {
	var p project
	_, err := toml.DecodeFile(filepath.Join(path, "Project.toml"), &p)
	if err != nil {
		return project{}, err
	}
	p.Path = path
	return p, nil
}

// readProjects reads the projects it can find in the given directory
// all directories in the project directory with a Project.toml are considered
// projects
func readProjects(path string) ([]project, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return []project{}, err
	}
	if !fileInfo.IsDir() {
		return []project{}, fmt.Errorf("projects_dir %s is not a directory", path)
	}

	f, err := os.Open(path)
	if err != nil {
		return []project{}, err
	}
	defer f.Close()

	dirNames, err := f.Readdirnames(0)
	if err != nil {
		return []project{}, err
	}

	var projectPaths []string
	for _, dirName := range dirNames {
		projectPaths = append(projectPaths, filepath.Join(path, dirName))
	}

	var projects []project
	for _, projectPath := range projectPaths {
		p, err := readProject(projectPath)
		if err == nil {
			projects = append(projects, p)
		}
	}

	return projects, nil
}
