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

	style.Init()

	app := tview.NewApplication()
	app.EnableMouse(true)
	log.Println("Starting Scrollbar Demo...")
	scrollbarWidget := scrollbar.NewScrollbar()

	layout := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(scrollbarWidget, 80, 0, true).
			AddItem(nil, 0, 1, false).
			SetDirection(tview.FlexColumn),
			20, 0, true).
		AddItem(nil, 0, 1, false).
		SetDirection(tview.FlexRow)

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
