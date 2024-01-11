package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const usage = `
Usage: hmm [COMMAND]

Commands:
  list
  show [RANK]
  rank FROM TO
  new
`

// TODO move helper functions in main.go to their own files

// SplitFootnote splits text into a document and a footnote
func SplitFootnote(input io.Reader) (string, string, error) {
	data, err := io.ReadAll(input)
	if err != nil {
		return "", "", err
	}
	pair := strings.Split(string(data), "\n+++\n")
	if len(pair) != 2 {
		return "", "", fmt.Errorf("Found more than one footnote separator")
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

func list(args []string, ws workspace) error {
	fmt.Println(ws.Name)
	projects, err := ws.ActiveProjects()
	if err != nil {
		return err
	}

	for _, project := range projects {
		fmt.Println(project.Name, "due", project.Deadline)
		backlog, err := project.Backlog()
		if err != nil {
			return err
		}
		for i, card := range backlog {
			fmt.Println(fmt.Sprintf("%v) %s", i+1, card.Headline))
		}
		fmt.Println()
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

func move(args []string, ws workspace) error {
	fmt.Println("Move not implemented", args)
	return nil
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
		time.Now(),
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

func main() {
	if (len(os.Args) < 2) {
		fmt.Print(usage)
		return
	}

	ws, err := currentWorkspace()
	if err != nil {
		panic(err)
	}
	action := os.Args[1]
	switch action {
	case "list":
		err = list(os.Args[2:], ws)
	case "show":
		err = show(os.Args[2:], ws)
	case "rank":
		err = rank(os.Args[2:], ws)
		_ = list([]string{}, ws)
	case "new":
		err = create(os.Args[2:], ws)
	default:
		fmt.Println("ERROR: Unknown command", action)
		fmt.Print(usage)
		os.Exit(1)
	}
	
	if err != nil {
		panic(err)
	}
}
