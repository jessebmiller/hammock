package main

import (
	"fmt"
	"os"
	"bufio"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

type card struct {
	Headline string
	Text     string
	Priority int
}

// readCard tries to read a path into a card struct
func readCard(path string) (card, error) {
	f, err := os.Open(path)
	if err != nil {
		return card{}, err
	}
	defer f.Close()

	headline := ""
	text := ""

	scanner := bufio.NewScanner(f)
	for scanner.Scan() && scanner.Text() != "+++" {
		headlinePrefix := "# "
		if headline == "" && strings.HasPrefix(scanner.Text(), headlinePrefix) {
			headline = strings.TrimPrefix(scanner.Text(), headlinePrefix)
		}
		text = text + "\n" + scanner.Text()
	}

	var metaTOML string
	for scanner.Scan() && scanner.Text() != "+++" {
		metaTOML = metaTOML + "\n" + scanner.Text()
	}

	var maybeCard struct {
		HammockType  string `toml:"hammock_type"`
		PriorityRank int    `toml:"priority_rank"`
	}

	_, err = toml.Decode(metaTOML, &maybeCard)
	if err != nil {
		return card{}, err
	}

	if maybeCard.HammockType != "Card" {
		return card{}, fmt.Errorf("Not a card (%s)", path)
	}

	return card{headline, text, maybeCard.PriorityRank}, nil
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
