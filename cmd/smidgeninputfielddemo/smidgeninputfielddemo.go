package main

import (
	"dinky/internal/tui/smidgeninputfield"
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	logFile := setupLogging()
	defer logFile.Close()

	app := tview.NewApplication()
	app.EnableMouse(true)
	log.Println("Starting smidgeninputfield Demo...")

	smidgenInputField := smidgeninputfield.NewSmidgenInputField(app)
	smidgenInputField.SetText("Type something here...")
	smidgenInputField.SetBackgroundColor(tcell.NewHexColor(0x007ace).TrueColor())
	smidgenInputField.SetTextColor(tcell.NewHexColor(0xffffff).TrueColor(), tcell.NewHexColor(0x007ace).TrueColor())
	smidgenInputField.SetHistory([]string{"First command", "Second command", "Third command"})

	topLayout := tview.NewFlex()
	topLayout.AddItem(nil, 0, 1, false)
	topLayout.AddItem(smidgenInputField, 80, 0, false)
	topLayout.AddItem(nil, 0, 1, false)
	app.SetRoot(topLayout, true)

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
