package main

import (
	"dinky/internal/tui/scrollbar"
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
	log.Println("Starting Scrollbar Demo...")
	verticalScrollbarWidget := scrollbar.NewScrollbar()

	layout := nuview.NewFlex()
	layout.AddItem(nil, 0, 1, false)

	innerLayout := nuview.NewFlex()
	innerLayout.AddItem(nil, 0, 1, false)
	innerLayout.AddItem(verticalScrollbarWidget, 1, 0, true)
	innerLayout.AddItem(nil, 0, 1, false)
	innerLayout.SetDirection(nuview.FlexColumn)

	layout.AddItem(innerLayout, 0, 10, true)

	horizontalScrollbarWidget := scrollbar.NewScrollbar()
	horizontalScrollbarWidget.SetHorizontal(true)

	layout.AddItem(horizontalScrollbarWidget, 1, 0, false)

	layout.AddItem(nil, 0, 1, false)
	layout.SetDirection(nuview.FlexRow)

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
