package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func check(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

// SplitFootnote splits text into a document and a footnote
func SplitFootnote(input io.Reader) (string, string, error) {
	data, err := io.ReadAll(input)
	if err != nil {
		return "", "", err
	}
	pair := strings.Split(string(data), "\n+++\n")
	if len(pair) != 2 {
		err := fmt.Errorf("Found more than one footnote separator")
		return "", "", err
	}
	return pair[0], pair[1], nil
}

// GetHeadline infers a headline from a document string
// returns empty string if it can't find a headline
func GetHeadline(document string) string {
	lines := strings.Split(document, "\n")
	if len(lines) == 0 {
		return ""
	}
	headline := strings.TrimSpace(strings.TrimLeft(lines[0], "#  \n\t\r"))
	return headline
}

func PopI(s []card, i int) (card, []card) {
	ret := make([]card, 0)
	ret = append(ret, s[:i]...)
	return s[i], append(ret, s[i+1:]...)
}

func InsertI(s []card, elem card, i int) []card {
	ret := make([]card, 0)
	ret = append(ret, s[:i]...)
	ret = append(ret, elem)
	return append(ret, s[i:]...)
}

func Swap(f int, t int, s []card) []card {
	a, y := PopI(s, f)
	return InsertI(y, a, t)
}

func openEditor(path string) error {
	editor := os.Getenv("VISUAL")
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}

	editor, err := exec.LookPath(editor)
	if err != nil {
		return err
	}

	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	err = cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()

	return err
}

func textFromEditor(presetText string) (string, error) {
	f, err := os.CreateTemp("", "new-card")
	if err != nil {
		return "", err
	}

	path := f.Name()
	defer os.Remove(path)

	f.Write([]byte(presetText))
	f.Close()
	openEditor(path)

	data, err := os.ReadFile(path)
	
	return string(data), nil
}

