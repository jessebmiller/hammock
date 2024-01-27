package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

type card struct {
	Path      string
	Headline  string
	Text      string
	Priority  int		`toml:"priority"`
	Completed time.Time	`toml:"completed"`
}

type CardOpt func(*card)

func WithPriority(n int) CardOpt {
	return func(c *card) {
		c.Priority = n
	}
}

func WithCompleted(t time.Time) CardOpt {
	return func(c *card) {
		c.Completed = t
	}
}

func getHeadline(s string) string {	
	scanner := bufio.NewScanner(strings.NewReader(s))
	scanner.Scan()

	headline := strings.Trim(
		strings.TrimSpace(scanner.Text()),
		"# ",
	)
	return headline
}

func fileNameFromHeadline(h string) string {
	reg, err := regexp.Compile("[^A-Za-z0-9_]+")
	if err != nil {
		panic(err)
	}
	fileRoot := strings.Trim(reg.ReplaceAllString(h, "-"), "-")
	return fmt.Sprintf("%s.md", fileRoot)
}

func CreateCard(
	projectPath string,
	text string,
	opts []CardOpt,
) (card, error) {
	headline := getHeadline(text)
	if headline == "" {
		e := "Aborting new card, missing headline"
		return card{}, fmt.Errorf(e)
	}
	fileName := fileNameFromHeadline(headline)
	path := filepath.Join(projectPath, fileName)
	c := card{
		path,
		headline,
		text,
		0,
		time.Time{},
	}
	return c, c.Write()
}

func (c *card) Write() error {
	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(map[string]any{
		"priority":     c.Priority,
		"completed":    c.Completed,
		"hammock_type": "Card",
	})
		if err != nil {
		return err
	}
	content := strings.Join([]string{
		c.Text,
		"+++",
		buf.String(),
	}, "\n")

	return os.WriteFile(c.Path, []byte(content), 0644)
}

func (c card) IsComplete() bool {
	return !c.Completed.IsZero()
}

func (c *card) MarkComplete() error {
	c.Completed = time.Now()
	return c.Write()
}

func (c *card) MarkNotComplete() error {
	c.Completed = time.Time{}
	return c.Write()
}

func (c *card) ToggleComplete() error {
	if c.Completed.IsZero() {
		return c.MarkComplete()
	}
	return c.MarkNotComplete()
}

// readCard tries to read a path into a card struct
func readCard(path string) (card, error) {
	f, err := os.Open(path)
	if err != nil {
		return card{}, err
	}
	defer f.Close()

	document, footnote, err := SplitFootnote(f)
	if err != nil {
		return card{}, err
	}

	headline := GetHeadline(document)
	if headline == "" {
		return card{}, fmt.Errorf("Empty card headline")
	}

	var maybeCard struct {
		HammockType string	`toml:"hammock_type"`
		Priority    int		`toml:"priority"`
		Completed   time.Time   `toml:"completed"`
	}

	_, err = toml.Decode(footnote, &maybeCard)
	if err != nil {
		return card{}, err
	}

	if maybeCard.HammockType != "Card" {
		return card{}, fmt.Errorf("Not a card (%s)", path)
	}

	if !maybeCard.Completed.IsZero() {
		headline = fmt.Sprintf("✔ %s", headline)
	} else {
		headline = fmt.Sprintf("• %s", headline)
	}

	return card{
		path,
		headline,
		document,
		maybeCard.Priority,
		maybeCard.Completed,
	}, nil
}

// readCards reads the cards it can find in the given directory
func readCards(path string) ([]card, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return []card{}, err
	}
	if !fileInfo.IsDir() {
		return []card{}, fmt.Errorf("Expected a directory, got %s", path)
	}

	f, err := os.Open(path)
	if err != nil {
		return []card{}, err
	}
	defer f.Close()

	dirNames, err := f.Readdirnames(0)
	if err != nil {
		return []card{}, err
	}

	var cardPaths []string
	for _, dirName := range dirNames {
		cardPaths = append(cardPaths, filepath.Join(path, dirName))
	}

	var cards []card
	for _, path := range cardPaths {
		c, err := readCard(path)
		if err == nil {
			cards = append(cards, c)
		}
	}

	return cards, nil
}
