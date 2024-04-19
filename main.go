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

/* STYLING */
var (
	columnStyle = lipgloss.NewStyle().
			Padding(1, 2)
	focusedStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
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
	err      error
	lists    []list.Model
	loaded   bool
	quitting bool
	focused  status
}

func New() *Model {
	return &Model{}
}

func (m *Model) Next() {
	if m.focused == done {
		m.focused = todo
	} else {
		m.focused++
	}
}

func (m *Model) Prev() {
	if m.focused == todo {
		m.focused = done
	} else {
		m.focused--
	}
}

func (m *Model) initLists(width, height int) {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width/divisor, height-divisor)
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
		Task{status: inProgress, title: "doing this kanban", description: "desc4"},
		Task{status: inProgress, title: "doing this kanban", description: "desc4"},
	})

	// init done
	m.lists[done].Title = "Done"
	m.lists[done].SetItems([]list.Item{
		Task{status: done, title: "first commit", description: "desc5"},
		Task{status: inProgress, title: "doing this kanban", description: "desc4"},
	})
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.loaded {
			columnStyle.Width(msg.Width / divisor)
			focusedStyle.Width(msg.Width / divisor)
			columnStyle.Height(msg.Height - divisor)
			focusedStyle.Height(msg.Height - divisor)
			m.initLists(msg.Width, msg.Height)
			m.loaded = true
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "left", "h":
			m.Prev()
		case "right", "l":
			m.Next()
		}

	}
	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}
	if !m.loaded {
		return "loading..."
	}
	todoView := m.lists[todo].View()
	inProgView := m.lists[inProgress].View()
	doneView := m.lists[done].View()
	switch m.focused {
	case inProgress:
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			columnStyle.Render(todoView),
			focusedStyle.Render(inProgView),
			columnStyle.Render(doneView),
		)
	case done:
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			columnStyle.Render(todoView),
			columnStyle.Render(inProgView),
			focusedStyle.Render(doneView),
		)
	default:
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			focusedStyle.Render(todoView),
			columnStyle.Render(inProgView),
			columnStyle.Render(doneView),
		)
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
