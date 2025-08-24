package main

import (
	"dinky/internal/tui/messagedialog"
	"dinky/internal/tui/style"
	"log"
	"os"

	"github.com/sedwards2009/nuview"
)

func main() {
	logFile := setupLogging()
	defer logFile.Close()

	style.Install()

	app := nuview.NewApplication()
	app.EnableMouse(true)
	log.Println("Starting MessageDialog Demo...")

	messageDialog := messagedialog.NewMessageDialog(app)
	messageDialog.OnButtonClick = func(button string, index int) {
		log.Printf("Button clicked: %s (index: %d)", button, index)
		app.Stop()
	}
	messageDialog.OnClose = func() {
		log.Println("Message dialog closed")
		app.Stop()
	}

	// style.StyleMessageDialog(messageDialog)

	topLayout := nuview.NewFlex()
	topLayout.AddItem(nil, 0, 1, false)
	topLayout.AddItem(messageDialog, 80, 0, false)
	topLayout.AddItem(nil, 0, 1, false)
	app.SetRoot(topLayout, true)

	messageDialog.Open("Question", "Do you want to proceed?\n\nIt will be way cool.",
		[]string{"Yes", "No", "Cancel"}, 50, 9)
	messageDialog.FocusButton(0)

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
