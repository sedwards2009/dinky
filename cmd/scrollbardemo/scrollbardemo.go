package main

import (
	"dinky/internal/tui/scrollbar"
	"dinky/internal/tui/style"
	"log"
	"os"

	"github.com/rivo/tview"
)

func main() {
	logFile := setupLogging()
	defer logFile.Close()

	app := tview.NewApplication()
	app.EnableMouse(true)
	log.Println("Starting Scrollbar Demo...")
	verticalScrollbarWidget := scrollbar.NewScrollbar()
	style.StyleScrollbar(verticalScrollbarWidget)

	layout := tview.NewFlex()
	layout.AddItem(nil, 0, 1, false)

	innerLayout := tview.NewFlex()
	innerLayout.AddItem(nil, 0, 1, false)
	innerLayout.AddItem(verticalScrollbarWidget, 1, 0, true)
	innerLayout.AddItem(nil, 0, 1, false)
	innerLayout.SetDirection(tview.FlexColumn)

	layout.AddItem(innerLayout, 0, 10, true)

	horizontalScrollbarWidget := scrollbar.NewScrollbar()
	style.StyleScrollbar(horizontalScrollbarWidget)
	horizontalScrollbarWidget.SetHorizontal(true)

	layout.AddItem(horizontalScrollbarWidget, 1, 0, false)

	layout.AddItem(nil, 0, 1, false)
	layout.SetDirection(tview.FlexRow)

	app.SetRoot(layout, true)

	if err := app.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}

func setupLogging() *os.File {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic("Failed to open log file: " + err.Error())
	}
	log.SetOutput(logFile)
	return logFile
}
