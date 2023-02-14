package display

// not too familiar with bubbletea yet, so for now im seperating my display code out into a package. this might be a bad idea
// this file defines the main menu view

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type sessionState int

const (
	menuView sessionState = iota
	highlightsView
	kindleView
)

// the MainModel for this file
type MainModel struct {
	state    sessionState
	choices  []string
	cursorY  int
	selected map[int]struct{}
}

func InitTUI() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func initialModel() MainModel {
	return MainModel{
		// Our to-do list is a grocery list
		choices: []string{"View your clippings", "Update clippings from your Kindle", "Exit"},

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
	}
}

// This function should not be called externally
// Init handles initial IO
func (m MainModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursorY > 0 {
				m.cursorY--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursorY < len(m.choices)-1 {
				m.cursorY++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.selected[m.cursorY]
			if ok {
				delete(m.selected, m.cursorY)
			} else {
				m.selected[m.cursorY] = struct{}{}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m MainModel) View() string {
	// The header
	s := "Welcome to kindlenotes\n\n"

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursorY == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}
