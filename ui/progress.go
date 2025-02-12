package ui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/progress"
)

type ProgressModel struct {
	progress progress.Model
	Complete bool
}

func NewProgressModel() ProgressModel {
	return ProgressModel{
		progress: progress.New(progress.WithDefaultGradient()),
	}
}

func (m ProgressModel) Init() tea.Cmd {
	return tickCmd()
}

func (m ProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	case tickMsg:
		if m.Complete {
			return m, tea.Quit
		}

		cmd := m.progress.IncrPercent(0.01)
		return m, tea.Batch(tickCmd(), cmd)

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
}

func (m ProgressModel) View() string {
	return fmt.Sprintf("\n%s\n\nPress any key to quit", m.progress.View())
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m *ProgressModel) ProgressUpdate(percent float64) {
	m.progress.SetPercent(percent)
}
