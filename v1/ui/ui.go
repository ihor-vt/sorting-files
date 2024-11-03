package ui

import (
	"fmt"
	"v1/analytics"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Cursor        int
	Choices       []string
	Quitting      bool
	FolderPath    string
	AnalyticsData *analytics.Analytics
}

func InitialModel() Model {
	return Model{
		Choices:       []string{"Standard Sorting", "AI Sorting (Coming Soon)"},
		AnalyticsData: analytics.NewAnalytics(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.Quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			if m.Cursor < len(m.Choices)-1 {
				m.Cursor++
			}
		case "enter":
			if m.Cursor == 0 {
				return m, tea.Quit // Proceed to sorting after quitting Bubble Tea UI
			} else {
				m.Quitting = true
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.Quitting {
		return "Thanks for using File Sorter Pro!"
	}

	s := "Welcome to File Sorter Pro!\n\n"
	s += "An innovative file sorting solution for both traditional and AI-based sorting.\n\n"
	s += "Select an option:\n\n"

	for i, choice := range m.Choices {
		cursor := " "
		if m.Cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += "\nPress q to quit.\n"
	return s
}
