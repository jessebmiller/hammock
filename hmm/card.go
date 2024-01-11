package main

import (
	"fmt"
	"os"
	"bufio"
	"path/filepath"
	"strings"
	"bytes"
	"time"

	"github.com/BurntSushi/toml"
)

type card struct {
	Path		string
	Headline	string
	Text		string
        Priority	int		`toml:"priority"`
	CreatedAt	time.Time	`toml:"created_at"`
}

func (c *card) Write() error {
	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(map[string]any{
		"priority": c.Priority,
		"created_at": c.CreatedAt,
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
	err = os.WriteFile(c.Path, []byte(content), 0644)
	if err != nil {
		return err
	}
	
	return nil
}

// readCard tries to read a path into a card struct
func readCard(path string) (card, error) {
	f, err := os.Open(path)
	if err != nil {
		return card{}, err
	}
	defer f.Close()
	var lines []string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() && scanner.Text() != "+++" {
		lines = append(lines, scanner.Text())
	}

	if len(lines) < 1 {
		return card{}, fmt.Errorf("No card text found")
	}

	headline := lines[0]
	headlinePrefix := "# "
	if strings.HasPrefix(lines[0], headlinePrefix) {
		headline = strings.TrimPrefix(lines[0], headlinePrefix)
	}

	var footnote string
	for scanner.Scan() && scanner.Text() != "+++" {
		footnote = footnote + "\n" + scanner.Text()
	}

	var maybeCard struct {
		HammockType	string		`toml:"hammock_type"`
		Priority	int		`toml:"priority"`
		CreatedAt	time.Time	`toml:"created_at"`
	}

	_, err = toml.Decode(footnote, &maybeCard)
	if err != nil {
		return card{}, err
	}

	if maybeCard.HammockType != "Card" {
		return card{}, fmt.Errorf("Not a card (%s)", path)
	}

	return card{
		path,
		headline,
		strings.Join(lines, "\n"),
		maybeCard.Priority,
		maybeCard.CreatedAt,
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
