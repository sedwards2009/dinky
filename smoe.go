package main

import (
	"log"
	"os"
	"smoe/internal/tui/menu"

	"github.com/rivo/tview"
)

func setupLogging() *os.File {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic("Failed to open log file: " + err.Error())
	}
	log.SetOutput(logFile)
	return logFile
}

func main() {
	logFile := setupLogging()
	defer logFile.Close()

	app := tview.NewApplication()
	app.EnableMouse(true)

	menuBar := menu.NewMenuBar()

	menuBar.SetMenus([]*menu.Menu{
		{Title: "File", Items: []*menu.MenuItem{
			{Title: "New", Shortcut: "Ctrl+N", Callback: func() { log.Println("New file") }},
			{Title: "Open", Shortcut: "Ctrl+O", Callback: func() { log.Println("Open file") }},
			{Title: "Save", Shortcut: "Ctrl+S", Callback: func() { log.Println("Save file") }},
			{Title: "", Shortcut: "", Callback: nil}, // Separator
			{Title: "Quit", Shortcut: "Ctrl+Q", Callback: func() { app.Stop() }},
		}},
		{Title: "Edit", Items: []*menu.MenuItem{
			{Title: "Undo", Shortcut: "Ctrl+Z", Callback: func() { log.Println("Undo") }},
			{Title: "Redo", Shortcut: "Ctrl+Y", Callback: func() { log.Println("Redo") }},
			{Title: "", Shortcut: "", Callback: nil}, // Separator
			{Title: "Cut", Shortcut: "Ctrl+X", Callback: func() { log.Println("Cut") }},
			{Title: "Copy", Shortcut: "Ctrl+C", Callback: func() { log.Println("Copy") }},
			{Title: "Paste", Shortcut: "Ctrl+V", Callback: func() { log.Println("Paste") }},
			{Title: "", Shortcut: "", Callback: nil}, // Separator
			{Title: "Select All", Shortcut: "Ctrl+A", Callback: func() { log.Println("Select All") }},
		}},
		{Title: "Help", Items: []*menu.MenuItem{}},
	})

	log.Println("Starting")
	app.SetRoot(menuBar, true)
	app.SetAfterDrawFunc(menuBar.AfterDraw())

	menuBar.Focus(nil)
	if err := app.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
