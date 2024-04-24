package main

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const APPEND = -1

type column struct {
	list   list.Model
	height int
	focus  bool
	status int
	width  int
}

func (c *column) Focus() {
	c.focus = true
}

func (c *column) Blur() {
	c.focus = false
}

func (c *column) Focused() bool {
	return c.focus
}

func newColumn(status int) column {
	var focus bool
	if status == 0 {
		focus = true
	}
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	defaultList.SetShowHelp(false)
	return column{focus: focus, status: status, list: defaultList}
}

// Init does initial setup for the column.
func (c column) Init() tea.Cmd {
	return nil
}

// Update handles all the I/O for columns.
func (c column) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.setSize(msg.Width, msg.Height)
		c.list.SetSize(msg.Width/TotalStages, msg.Height/2)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Edit):
			if len(c.list.VisibleItems()) != 0 {
				task := c.list.SelectedItem().(Step)
				f := NewForm(task.title, task.description)
				f.index = c.list.Index()
				f.col = c
				return f.Update(nil)
			}
		case key.Matches(msg, keys.New):
			f := newDefaultForm()
			f.index = APPEND
			f.col = c
			return f.Update(nil)
		case key.Matches(msg, keys.Delete):
			return c, c.DeleteCurrent()
		case key.Matches(msg, keys.Enter):
			// TODO: instead of move next make it EDIT
			return c, c.MoveToNext()
		}
	}
	c.list, cmd = c.list.Update(msg)
	return c, cmd
}

func (c column) View() string {
	return c.getStyle().Render(c.list.View())
}

func (c *column) DeleteCurrent() tea.Cmd {
	if len(c.list.VisibleItems()) > 0 {
		c.list.RemoveItem(c.list.Index())
	}

	var cmd tea.Cmd
	c.list, cmd = c.list.Update(nil)
	return cmd
}

func (c *column) Set(i int, t Step) tea.Cmd {
	if i != APPEND {
		return c.list.SetItem(i, t)
	}
	return c.list.InsertItem(APPEND, t)
}

func (c *column) setSize(width, height int) {
	c.width = width / TotalStages
}

func (c *column) getStyle() lipgloss.Style {
	if c.Focused() {
		return lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Height(c.height).
			Width(c.width)
	}
	return lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.HiddenBorder()).
		Height(c.height).
		Width(c.width)
}

type moveMsg struct {
	Step
}

func (c *column) MoveToNext() tea.Cmd {
	var step Step
	var ok bool
	// If nothing is selected, the SelectedItem will return Nil.
	if step, ok = c.list.SelectedItem().(Step); !ok {
		return nil
	}
	// move item
	c.list.RemoveItem(c.list.Index())
	step.status = getNext(c.status)

	// refresh list
	var cmd tea.Cmd
	c.list, cmd = c.list.Update(nil)

	return tea.Sequence(cmd, func() tea.Msg { return moveMsg{step} })
}
