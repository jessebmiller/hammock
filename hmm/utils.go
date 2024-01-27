package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func check(err error) {
	if os.Getenv("HAMMOCK_DEBUG") == "true" && err != nil {
		panic(err)
	}
	if err != nil {
		os.Exit(1)
	}
}

func PPstr(i any) string {
	s, _ := json.MarshalIndent(i, "", "  ")
	return string(s)
}

func PPrint(i interface{}) {
	fmt.Println(PPstr(i))
}

func prompt(p string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(p + " ")
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(text), nil
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

func splitAtWidth(line string, width int) string {
	if len(line) <= width {
		return line
	}
	scanner := bufio.NewScanner(strings.NewReader(line))
	scanner.Split(bufio.ScanWords)
	var l []string
	var llen int
	var ls []string
	for scanner.Scan() {
		scanLen := len(scanner.Text()) + 1
		if llen+scanLen <= width {
			l = append(l, scanner.Text())
			llen += scanLen
			continue
		}
		ls = append(ls, strings.Join(l, " "))
		l = []string{}
		llen = 0
	}
	return strings.Join(ls, "\n")
}

func withMaxWidth(width int) func(string) string {
	return func(text string) string {
		if len(text) <= width {
			return text
		}
		return text[:width-3] + "..."
	}
}

func each[X any, Y any](xs []X, f func(x X) Y) []Y {
	ys := make([]Y, len(xs))
	for i, x := range xs {
		ys[i] = f(x)
	}
	return ys
}

func withJustWidth(width int) func(string) string {
	return func(text string) string {
		if len(text) <= width {
			return text
		}
		var lines []string
		for _, line := range strings.Split(text, "\n") {
			lines = append(lines, splitAtWidth(line, width))
		}
		return strings.Join(lines, "\n")

	}
}

// rowsToLines adds blank lines for rows with multiline cells
// pushing down later rows as needed
func rowsToLines(rows [][]string) [][]string {
	var lines [][]string
	for _, row := range rows {
		lines = append(lines, rowToLines(row)...)
	}
	return lines
}

// rowToLines splits a row into lines and pads cells to match the max
func rowToLines(row []string) [][]string {
	var lines [][]string
	for _, cell := range row {
		cell = strings.TrimSpace(cell)
		cellLines := strings.Split(cell, "\n")
		lines = pushCol(lines, cellLines)
	}
	return lines
}

// pushCol pushes a column into a 2d string array
// adding blank rows as needed. This function expects
// the rows to all be the same length
func pushCol(rows [][]string, col []string) [][]string {
	if len(rows) == 0 {
		for _, cell := range col {
			rows = append(rows, []string{cell})
		}
		return rows
	}

	// append a cell from the column or a blank
	// cell to the end of each row
	var pushedRows [][]string
	var pushedTo int
	for i, row := range rows {
		cell := " "
		if i < len(col) {
			cell = col[i]
		}
		r := append(row, cell)
		pushedRows = append(pushedRows, r)
		pushedTo = i
	}

	// if there were more rows than the column needed
	// we're done
	if len(rows) >= len(col) {
		return pushedRows
	}

	// there may be extra cells in the column
	// if so, add a row and fill it with blanks before
	// the new column cells
	padding := make([]string, len(pushedRows[0])-1)
	for _, cell := range col[pushedTo+1:] {
		r := append(padding, cell)
		pushedRows = append(pushedRows, r)
	}

	return pushedRows
}
