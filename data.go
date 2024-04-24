package main

import (
	"strconv"

	"github.com/charmbracelet/bubbles/list"
)

func (b *Board) initLists(totalStages int) {
	b.cols = []column{}

	for i := 0; i < totalStages; i++ {
		b.cols = append(b.cols, newColumn(i))
		b.cols[i].list.Title = "Stage " + strconv.Itoa(i+1)
		b.cols[i].list.SetItems([]list.Item{
			// Step{status: i, title: "Exampel Stage", description: "delete it"},
		})
	}
}
