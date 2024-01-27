package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
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
	Complete  time.Time `toml:"complete"`
	ShowDoneFor string `toml:"show_done_for"`
}

func (p project) Write() error {
	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(p)
	if err != nil {
		return err
	}
	return os.WriteFile(p.Path, []byte(buf.String()), 0644)
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

func longerAgoThan(t time.Time, d time.Duration) bool {
	if t.IsZero() {
		return false
	}
	return time.Now().Sub(t) > d
}

func (p project) BacklogHeadlines() ([]string, error) {
	var headlines []string
	backlog, err := p.Backlog()
	if err != nil {
		return []string{}, err
	}
	for _, card := range backlog {
		d, err := time.ParseDuration(p.ShowDoneFor)
		if err != nil {
			return []string{}, err
		}
		if longerAgoThan(card.Completed, d) {
			continue
		}
		headlines = append(headlines, card.Headline)
	}
	return headlines, nil
}

func (p project) PrintSummary() error {
	summary, err := p.Summary()
	if err != nil {
		return err
	}
	fmt.Println(summary)
	return nil
}

func (p project) Summary() (string, error) {
	s := []string{
		"",
		p.Workspace.Name,
		"",
		p.Name,
		"",
		strings.TrimSpace(p.Goal),
		"",
		p.DeadlineNotice(),
		"",
	}
	headlines, err := p.BacklogHeadlines()
	if err != nil {
		return "", err
	}
	for i, h := range headlines {
		s = append(s, fmt.Sprintf("%v. %v", i+1, h))
	}
	return strings.Join(s, "\n"), nil
}

func WriteConsecutivePriorities(cards []card) error {
	for i, card := range cards {
		card.Priority = i + 1
		err := card.Write()
		if err != nil {
			return err
		}
	}
	return nil
}

// readProject reads a project from a path
// path must be a directory with a valid Project.toml file in it
func readProject(path string) (project, bool, error) {
	tomlPath := filepath.Join(path, "Project.toml")
	if _, err := os.Stat(tomlPath); err != nil {
		return project{}, false, nil
	}
	var p project
	_, err := toml.DecodeFile(tomlPath, &p)
	if err != nil {
		return project{}, false, err
	}
	p.Path = path
	ws, err := inWorkspace(path)
	if err != nil {
		return project{}, false, err
	}
	p.Workspace = ws
	return p, true, nil
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
		p, isProject, err := readProject(projectPath)
		if err != nil {
			return []project{}, err
		}
		if isProject {
			projects = append(projects, p)
		}
	}

	return projects, nil
}
