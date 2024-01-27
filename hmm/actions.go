package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

func summarize() error {
	hmm, err := readHammock()
	if err != nil {
		return err
	}

	projects, err := hmm.PriorityProjects()
	if err != nil {
		return err
	}

	width := 60
	var rows [][]string
	for _, p := range projects {
		col := []string{
			withMaxWidth(width)(p.Workspace.Name),
			withMaxWidth(width)(p.Name),
			" ",
			withJustWidth(width)(p.Goal),
			" ",
			withMaxWidth(width)(p.DeadlineNotice()),
			" ",
		}
		headlines, err := p.CardHeadlines()
		if err != nil {
			return err
		}
		for i, h := range headlines {
			headlines[i] = fmt.Sprintf("%v. %v", i+1, h)
		}
		col = append(col, each(headlines, withMaxWidth(width))...)
		rows = pushCol(rows, col)
	}

	lines := rowsToLines(rows)

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 0, 8, ' ', 0)
	defer w.Flush()

	fmt.Println()
	for _, line := range lines {
		fmt.Fprintln(w, strings.Join(line, "\t"))
	}

	return nil
}

// done marks a card done by rank
func done(args []string, ws workspace) error {
	var rank int
	if len(args) > 0 {
		input, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		rank = input - 1
	}

	activeProjects, err := ws.ActiveProjects()
	if err != nil {
		return err
	}

	if len(activeProjects) == 0 {
		fmt.Println("No active projects")
		return nil
	}

	if len(activeProjects) > 1 {
		fmt.Println("Multiple active projects! only showing card from one of them")
		fmt.Println("Hangling this is future work")
		return nil
	}

	backlog, err := activeProjects[0].Backlog()

	rank = min(len(backlog)-1, rank)
	rank = max(0, rank)

	return backlog[rank].ToggleComplete()
}

// list the cards in active projects of a workspace
func list(args []string, ws workspace) error {
	projects, err := ws.ActiveProjects()
	if err != nil {
		return err
	}

	for _, project := range projects {
		project.PrintSummary()
	}

	return nil
}

func initWorkspace(args []string) error {
	var name string
	var err error
	if len(args) == 0 {
		name, err = prompt("Workspace Name:")
	} else {
		name = strings.Join(args, " ")
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	path := filepath.Join(wd, "Workspace.toml")
	ws := workspace{
		name,
		path,
		"projects",
	}
	return ws.Write()
}

func cdToWorkspace(args []string) error {
	hmm, err := readHammock()
	if err != nil {
		return err
	}
	workspaces, err := hmm.Workspaces()
	if err != nil {
		return err
	}
	if len(args) < 1 {
		// TODO don't print the message, print an echo statement this is
		// evaluated by a shell script to change directories
		fmt.Println("usage: hmm go <workspaceName>")
		fmt.Println("Workspaces:")
		for _, ws := range workspaces {
			fmt.Println(ws.Name)
		}
		return fmt.Errorf("Missing argument")
	}

	for _, ws := range workspaces {
		if ws.Name == args[0] {
			fmt.Println("cd", ws.Path)
			return nil
		}
	}
	fmt.Println("Workspace", args[0], "not found.")
	fmt.Println("Workspaces:")
	for _, ws := range workspaces {
		fmt.Println(ws.Name)
	}
	return fmt.Errorf("workspace not found")
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
	}

	activeProjects, err := ws.ActiveProjects()
	if err != nil {
		return err
	}

	if len(activeProjects) == 0 {
		fmt.Println("No active projects, no card to show")
		return nil
	}

	if len(activeProjects) > 1 {
		fmt.Println("Multiple active projects! only showing card from one of them")
		fmt.Println("Hangling this is future work")
		return nil
	}

	project := activeProjects[0]
	backlog, err := project.Backlog()

	rank = min(len(backlog)-1, rank)
	rank = max(0, rank)

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
		time.Time{},
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

	if rank < 0 || rank >= len(backlog) {
		fmt.Println("No card at rank", rank)
	}

	card := backlog[rank]

	if card.Headline == args[1] {
		os.Remove(card.Path)
		fmt.Println("removed", card.Headline)
	} else {
		fmt.Println("The headline you entered did not match the headline of the card")
		fmt.Println("Card rank", rank+1, "has headline:", card.Headline)
		fmt.Println("You entered:             ", args[1])
	}

	return nil
}
