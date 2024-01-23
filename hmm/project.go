package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/BurntSushi/toml"
)

type project struct {
	Path      string
	Workspace workspace
	Name      string    `toml:"name"`
	Goal      string    `toml:"goal"`
	Start     time.Time `toml:"start"`
	Deadline  time.Time `toml:"deadline"`
	Complete  bool      `toml:"complete"`
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

func (p project) CardHeadlines() ([]string, error) {
	backlog, err := p.Backlog()
	if err != nil {
		return []string{}, err
	}
	var headlines []string
	for _, c := range backlog {
		headlines = append(headlines, c.Headline)
	}
	return headlines, nil
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

// DeadlineNotice text notice of the deadline
func (p project) DeadlineNotice() string {
	duration := p.Deadline.Sub(time.Now())
	days := int(duration.Hours() / 24)
	var deadline string
	switch {
	case p.Deadline.IsZero():
		deadline = "No deadline"
	case days == 0:
		deadline = "Due today"
	case days > 0:
		deadline = fmt.Sprintf(
			"Due in %v days",
			days,
		)
	case days < 0:
		deadline = fmt.Sprintf(
			"Overdue by %v days",
			days,
		)
	}
	return deadline
}

// push elements into a slice
func push[x any](elements []x, elem x) {
	elements = append(elements, elem)
}

func (p project) BacklogHeadlines() ([]string, error) {
	var headlines []string
	backlog, err := p.Backlog()
	if err != nil {
		return []string{}, err
	}
	for _, card := range backlog {
		push(headlines, card.Headline)
	}
	return headlines, nil
}

func (p project) PrintSummary() error {
	fmt.Println()
	fmt.Println(p.Name)
	fmt.Println(p.DeadlineNotice())
	headlines, err := p.BacklogHeadlines()
	if err != nil {
		return err
	}
	for i, h := range headlines {
		fmt.Println(fmt.Sprintf("%v. %v", i, h))
	}
	fmt.Println()
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
	ws, err := inWorkspace(path)
	if err != nil {
		return project{}, err
	}
	p.Workspace = ws
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
