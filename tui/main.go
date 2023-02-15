package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var app = tview.NewApplication()
var text = tview.NewTextView().
	SetTextColor(tcell.ColorGreen).
	SetText("(q) to quit")

func main() {
	// app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
	// 	if event.Rune() == 113 {
	// 		app.Stop()
	// 	}
	// 	return event
	// })
	app := tview.NewApplication()
	flex := tview.NewFlex().
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Left (1/2 x width of Top)"), 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Top"), 0, 1, false).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Middle (3 x height of Top)"), 0, 3, false).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Bottom (5 rows)"), 5, 1, false), 0, 2, false).
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Right (20 cols)"), 20, 1, false)
	if err := app.SetRoot(flex, true).EnableMouse(true).SetFocus(flex).Run(); err != nil {
		panic(err)
	}
	// if err := app.SetRoot(text, true).EnableMouse(true).Run(); err != nil {
	// 	panic(err)
	// }
}
