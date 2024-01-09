package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"
)

const usage = `
Usage: hmm [COMMAND]

Commands:
  list
  show [RANK]
  move [FROM_RANK] [TO_RANK]
`

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
		rank, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		rank = rank - 1
	}
	
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

	rank = min(len(backlog) - 1, rank)
	rank = max(0, rank)

	fmt.Println(project.Name)
	for i, card := range backlog {
		if i == rank {
			fmt.Println(i, ") ", card.Text)
		} else {
			fmt.Println(i, ") ", card.Headline)
		}
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
	presetTextTempl := `


+++
created_at = %v
hammock_type = "Card"
priority_rank = %v
+++
`

	presetText := fmt.Sprintf(presetTextTempl, time.Now().Format(time.RFC3339), 1)
	cardText, err := textFromEditor(presetText)
	if err != nil {
		return err
	}
	if cardText == presetText {
		fmt.Println("No card text, did you save?")
		return nil
	}

	// extract headline for filename
	fmt.Println(cardText)

	// find project directory path

	// create file and write text into it
	return nil
}

func main() {
	fmt.Println("got args", os.Args)
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
	case "move":
		err = move(os.Args[2:], ws)
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
