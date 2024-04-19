package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type status int

const divisor int = 4

const (
	todo status = iota
	inProgress
	done
)

type Task struct {
	title       string
	description string
	status      status
}

// list.Item interface

func (t Task) FilterValue() string {
	return t.title
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}

/* Main Model */

type Model struct {
	err     error
	lists   []list.Model
	loaded  bool
	focused status
}

func New() *Model {
	return &Model{}
}

func (m *Model) initLists(width, height int) {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width/divisor, height)
	defaultList.SetShowHelp(false)
	m.lists = []list.Model{defaultList, defaultList, defaultList}

	// init To Do
	m.lists[todo].Title = "To Do"
	m.lists[todo].SetItems([]list.Item{
		Task{status: todo, title: "title1", description: "desc1"},
		Task{status: todo, title: "title2", description: "desc2"},
		Task{status: todo, title: "title3", description: "desc3"},
	})

	// init progress
	m.lists[inProgress].Title = "Progress"
	m.lists[inProgress].SetItems([]list.Item{
		Task{status: inProgress, title: "doing this kanban", description: "desc4"},
	})

	// init done
	m.lists[done].Title = "Done"
	m.lists[done].SetItems([]list.Item{
		Task{status: done, title: "first commit", description: "desc5"},
	})
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.loaded {
			m.initLists(msg.Width, msg.Height)
			m.loaded = true
		}
	}
	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.loaded {
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			m.lists[todo].View(),
			m.lists[inProgress].View(),
			m.lists[done].View(),
		)
	} else {
		return "loading..."
	}
}

func main() {
	m := New()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
