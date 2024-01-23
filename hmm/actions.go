package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/tabwriter"
)

type ProjectColumn struct {
	WsName string
	Name string
	Goal string
	DeadlineNotice string
	CardHeadlines []string
}

/*
func projectColumn(wsName string, p project) (ProjectColumn, error) {
	col := ProjectColumn{
		wsName,
		p.Name,
		p.Goal,
		p.DeadlineNotice(),
		[]string{},
	}
	backlog, err := p.Backlog()
	if err != nil {
		return []string{}, err
	}
	for i, c := range backlog {
		col.CardHeadlines = append(
			col.CardHeadlines,
			fmt.Sprintf("%v. %s", i, c.Headline),
		)
	}
	return col, nil
}

func maxHeight(m int, s string) int {
	return max(m, len(strings.Split(s)))
}

// maxHeights returns the max number of lines for each property
func maxHeights(cols []ProjectColumn) map[string]int {
	var maxes map[string]int
	for _, col := range cols {
		maxes["WsName"] = maxHeight(maxes["WsName"], col.WsName)
		maxes["Name"] = maxHeight(maxes["Name"], col.Name)
		maxes["Goal"] = maxHeight(maxes["Goal"], col.Goal)
		maxes["DeadlineNotice"] = maxHeight(
			maxes["DeadlineNotice"],
			col.DeadlineNotice,
		)
		maxes["CardHeadlines"] = max(
			maxes["CardHeadlines"],
			len(col.CardHeadlines),
		)
	}
	return maxes
}

func pad(propName string, propVal string) string {
	maxes = maxHeights(cols)
	lines := strings.Split(propVal, "\n")
	dif := maxes[propName] - len(lines)
	return strings.Join(append(lines, make([]string, dif)...), "\n")
}
//padFields pads all fields with empty lines to match the max
func padFields(cols []ProjectColumn) {
	for _, col := range cols {
		col.WsName = pad("WsName", col.WsName)
		col.Name = pad("Name", col.Name)
		col.Goal = pad("Goal", col.Goal)
		col.DeadlineNotice = pad("DeadlineNotice", col.DeadlineNotice)
		dif := maxes["CardHeadlines"] - len(col.CardHeadlines)
		col.CardHeadlines = append(
			col.CardHeadlines,
			make([]string, dif)...,
		)
	}
}
*/

func PPstr(i any) string {
	s, _ := json.MarshalIndent(i, "", "  ")
	return string(s)
}

func PPrint(i interface{}) {
      fmt.Println(PPstr(i))
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
	for i, row := range rows {
		cell := " "
		if i < len(col) {
			cell = col[i]
		}
		row = append(row, cell)
		pushedRows = append(pushedRows, row)
	}

	// if there were more rows than the column needed
	// we're done
	if len(rows) >= len(col) {
		return pushedRows
	}

	// there may be extra cells in the column
	// if so, add a row and fill it with blanks before
	// the new column cells
	padding := make([]string, len(pushedRows[0]) - 1)
	for _, cell := range col[len(padding):] { 
		row := append(padding, cell)
		pushedRows = append(pushedRows, row)
	}

	return pushedRows
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

// rowsToLines adds blank lines for rows with multiline cells
// pushing down later rows as needed
func rowsToLines(rows [][]string) [][]string {
	var lines [][]string
	for _, row := range rows {
		lines = append(lines, rowToLines(row)...)
	}
	return lines
}

func splitAtWidth(line string, width int) string {
	scanner := bufio.NewScanner(strings.NewReader(line))
	scanner.Split(bufio.ScanWords)
	var l string
	var ls []string
	for scanner.Scan() {
		l = l + " " + scanner.Text()
		if len(l) - 1 + len(scanner.Text()) < width {
			continue
		}
		ls = append(ls, strings.Trim(l, " "))
		l = ""

	}
	return strings.Join(ls, "\n")
}

func withJustWidth(width int) func(string)string {
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

func withMaxWidth(width int) func(string)string {
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

func summarize() error {
	hmm, err := readHammock()
	if err != nil {
		return err
	}

	projects, err := hmm.PriorityProjects()
	if err != nil {
		return err
	}
	
	width := 32
	var rows [][]string
	for _, p := range projects {
		col := []string{
			withMaxWidth(width)(p.Workspace.Name),
			withMaxWidth(width)(p.Name),
			withJustWidth(width)(p.Goal),
			withMaxWidth(width)(p.DeadlineNotice()),
		}
		headlines, err := p.CardHeadlines()
		if err != nil {
			return err
		}
		col = append(col, each(headlines, withMaxWidth(width))...)
		rows = pushCol(rows, col)
	}

	lines := rowsToLines(rows)

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 0, 4, ' ', 0)
	defer w.Flush()

	for _, line := range lines {
		fmt.Println(strings.Join(line, "\\t"))
		fmt.Fprintln(w, strings.Join(line, "\t"))
	}

	return nil
}

// list the cards in active projects of a workspace
func list(args []string, ws workspace) error {
	fmt.Println()
	fmt.Println("  ", ws.Name)
	projects, err := ws.ActiveProjects()
	if err != nil {
		return err
	}

	for _, project := range projects {
		project.PrintSummary()
	}

	return nil
}

// show a card or few
// By default, show the card with the highest priority rank
// Optional arg specifying the rank of the card to show
func show(args []string, ws workspace) error {
	var rank int
	if len(args) > 0 {
		input, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		rank = input - 1
		fmt.Println(rank)
	}
	fmt.Println(rank)
	
	activeProjects, err := ws.ActiveProjects()
	if err != nil {
		return err
	}
	
	if len(activeProjects) == 0 {
		fmt.Println("No active projects, no card to show")
	}

	if len(activeProjects) > 1 {
		fmt.Println("Multiple active projects! only showing card from one of them")
		fmt.Println("Hangling this is future work")
	}

	project := activeProjects[0]
	backlog, err := project.Backlog()

	fmt.Println(rank)
	rank = min(len(backlog) - 1, rank)
	fmt.Println(rank)
	rank = max(0, rank)
	fmt.Println(rank)

	fmt.Println(project.Name)
	for i, card := range backlog {
		if i == rank {
			fmt.Println(fmt.Sprintf("%v) %s\n", i+1, card.Text))
		} else {
			fmt.Println(fmt.Sprintf("%v) %s", i+1, card.Headline))
		}
	}

	return nil
}

func rank(args []string, ws workspace) error {
	usage := "Usage: hmm rank FROM TO"
	if len(args) < 2 {
		fmt.Println(usage)
		return nil
	}

	from, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println(usage)
		fmt.Println("FROM and TO must both be integers")
		return nil
	}
	to, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println(usage)
		fmt.Println("FROM and TO must both be integers")
		return nil
	}
	ps, err := ws.ActiveProjects()
	if len(ps) != 1 {
		fmt.Println("Rank only supports a single active project")
		fmt.Println("There are", len(ps))
		for _, p := range ps {
			fmt.Println("  ", p.Name)
		}
		return nil
	}

	// translate from ranks to backlog array index
	from = from - 1
	to = to - 1

	activeProjects := ps[0]
	backlog, err := activeProjects.Backlog()
	if err != nil {
		return err
	}

	if from >= len(backlog) || from < 0 {
		fmt.Println(usage)
		fmt.Println("FROM must be 1 -", len(backlog))
		return nil
	}

	backlog = Swap(from, to, backlog)
	err = WriteConsecutivePriorities(backlog)
	if err != nil {
		return err
	}
	return nil
}

func create(args []string, ws workspace) error {
	presetText := "# "
	cardText, err := textFromEditor(presetText)
	if err != nil {
		return err
	}
	if cardText == presetText {
		fmt.Println("Aborting new card due to empty card text")
		return nil
	}

	// extract headline for filename
	scanner := bufio.NewScanner(strings.NewReader(cardText))
	scanner.Scan()

	headline := strings.Trim(
		strings.TrimSpace(scanner.Text()),
		"# ",
	)
	if headline == "" {
		fmt.Println("Aborting new card due to missing headline on first line")
		return nil
	}

	textLines := []string{scanner.Text()}
	for scanner.Scan() {
		textLines = append(textLines, scanner.Text())
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}

	// find project directory path
	projects, err := ws.ActiveProjects()
	if err != nil {
		return err
	}
	if len(projects) != 1 {
		fmt.Println("Aborting new card")
		fmt.Println("Not implemented for other than exactly one active project")
		fmt.Println("Active Projects:", len(projects))
		for i, project := range projects {
			fmt.Println(i+1, ")", project.Name)
		}
	}

	priorityProject := projects[0]

	path := priorityProject.Path

	// create file and write text into it
	reg, err := regexp.Compile("[^A-Za-z0-9_]+")
	if err != nil {
		return err
	}
	fileRoot := strings.Trim(reg.ReplaceAllString(headline, "-"), "-")
	fileName := fmt.Sprintf("%s.md", fileRoot)
	path = filepath.Join(path, fileName)
	newCard := card{
		path,
		headline,
		strings.Join(textLines, "\n"),
		-1, // Priority, Future work, take this in a flag default 0
	}
	err = newCard.Write()
	if err != nil {
		return err
	}

	err = priorityProject.NormalizePriorities()
	if err != nil {
		return err
	}

	return nil
}

func remove(args []string, ws workspace) error {
	if len(args) != 2 {
		fmt.Println("Usage: hmm remove RANK \"QUOTED CARD HEADLINE\"")
		return nil
	}

	rank, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("RANK must be an integer, got", args[0])
		return nil
	}
	rank = rank - 1

	activeProjects, err := ws.ActiveProjects()
	if err != nil {
		return err
	}

	if len(activeProjects) != 1 {
		fmt.Println("Not implemented for more than 1 active project")
		fmt.Println("Active projects:")
		for i, p := range activeProjects {
			fmt.Println(i, p.Name)
		}
		return nil
	}

	backlog, err := activeProjects[0].Backlog()
	if err != nil {
		return err
	}

	if (rank < 0 || rank >= len(backlog)) {
		fmt.Println("No card at rank", rank)
	}

	card := backlog[rank]

	if card.Headline == args[1] {
		os.Remove(card.Path)
		fmt.Println("removed", card.Headline)
	} else {
		fmt.Println("The headline you entered did not match the headline of the card")
		fmt.Println("Card rank", rank + 1, "has headline:", card.Headline)
		fmt.Println("You entered:             ", args[1])
	}
	
	return nil
}
