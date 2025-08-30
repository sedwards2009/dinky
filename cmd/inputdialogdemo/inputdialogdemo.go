package main

import (
	"dinky/internal/tui/dialog"
	"dinky/internal/tui/style"
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/sedwards2009/nuview"
)

func main() {
	logFile := setupLogging()
	defer logFile.Close()

	style.Install()

	app := nuview.NewApplication()
	app.EnableMouse(true)
	log.Println("Starting InputDialog Demo...")

	messageDialog := dialog.NewInputDialog(app)
	// style.StyleMessageDialog(messageDialog)

	topLayout := nuview.NewFlex()
	topLayout.AddItem(nil, 0, 1, false)
	topLayout.AddItem(messageDialog, 80, 0, false)
	topLayout.AddItem(nil, 0, 1, false)
	app.SetRoot(topLayout, true)

	messageDialog.Open(dialog.InputDialogOptions{
		Title:        "Go to Line",
		Message:      "Line: ",
		DefaultValue: "",
		Buttons:      []string{"Ok", "Cancel"},
		Width:        50,
		Height:       7,
		OnCancel: func() {
			log.Println("Input dialog canceled")
			app.Stop()
		},
		OnAccept: func(value string, index int) {
			log.Printf("Input dialog accepted with value: %s (button index: %d)", value, index)
			app.Stop()
		},
		FieldKeyFilter: numberInputFilter,
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

func numberInputFilter(event *tcell.EventKey) bool {
	switch event.Key() {
	case tcell.KeyRune:
		key := event.Rune()
		if key >= '0' && key <= '9' || key == '-' || key == '+' {
			return true
		}
	case tcell.KeyBackspace, tcell.KeyDEL, tcell.KeyDelete, tcell.KeyLeft, tcell.KeyRight, tcell.KeyHome, tcell.KeyEnd:
		return true
	default:
	}
	return false
}
