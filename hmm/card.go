package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"bytes"

	"github.com/BurntSushi/toml"
)

type card struct {
	Path		string
	Headline	string
	Text		string
        Priority	int		`toml:"priority"`
}

func (c *card) Write() error {
	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(map[string]any{
		"priority": c.Priority,
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

	document, footnote, err := SplitFootnote(f)
	if err != nil {
		return card{}, err
	}

	headline := GetHeadline(document)
	if headline == "" {
		return card{}, fmt.Errorf("Empty card headline")
	}

	var maybeCard struct {
		HammockType	string		`toml:"hammock_type"`
		Priority	int		`toml:"priority"`
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
		document,
		maybeCard.Priority,
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
