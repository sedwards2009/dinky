package main

import (
	"dinky/internal/tui/dialog"
	"log"
	"os"

	nuview "github.com/rivo/tview"
)

func main() {
	logFile := setupLogging()
	defer logFile.Close()

	app := nuview.NewApplication()
	app.EnableMouse(true)
	log.Println("Starting ListDialog Demo...")

	listDialog := dialog.NewListDialog(app)
	// style.StyleListDialog(listDialog)

	topLayout := nuview.NewFlex()
	topLayout.AddItem(nil, 0, 1, false)
	topLayout.AddItem(listDialog, 80, 0, false)
	topLayout.AddItem(nil, 0, 1, false)
	app.SetRoot(topLayout, true)

	listDialog.Open(dialog.ListDialogOptions{
		Title:   "Select an Item",
		Message: "Please select an item from the list:",
		Buttons: []string{"Ok", "Cancel"},
		Width:   50,
		Height:  20,
		Items: []dialog.ListItem{
			{Text: "Item 1", Value: "value1"},
			{Text: "Item 2", Value: "value2"},
			{Text: "Item 3", Value: "value3"},
			{Text: "Item 4", Value: "value4"},
			{Text: "Item 5", Value: "value5"},
		},
		DefaultSelected: "value3",
		OnCancel: func() {
			log.Println("List dialog canceled")
			app.Stop()
		},
		OnAccept: func(value string, index int) {
			log.Printf("List dialog accepted with value: %s (button index: %d)", value, index)
			app.Stop()
		},
	})

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
