package main

import (
	"fmt"
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

var TotalStages int

func getNext(s int) int {
	if s != TotalStages-1 {
		return s + 1
	}
	return s
}

func getPrev(s int) int {
	if s != 0 {
		return s - 1
	}
	return s
}

var board *Board

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	levelName := os.Args[1]
	TotalStages, _ = strconv.Atoi(os.Args[2])
	fmt.Println(levelName, TotalStages)

	board = NewBoard()
	board.initLists(TotalStages)
	p := tea.NewProgram(board, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
