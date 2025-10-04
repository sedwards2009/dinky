package main

import (
	"dinky/internal/tui/filelist"
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
	log.Println("Starting Filelist Demo...")

	realFileList := filelist.NewFileList(app)
	style.StyleFileList(realFileList)

	layout := tview.NewFlex()
	layout.AddItem(nil, 0, 1, false)

	innerFlex := tview.NewFlex()
	innerFlex.AddItem(nil, 0, 1, false)
	innerFlex.AddItem(realFileList, 80, 0, true)
	innerFlex.AddItem(nil, 0, 1, false)
	innerFlex.SetDirection(tview.FlexColumn)

	layout.AddItem(innerFlex, 20, 0, true)

	layout.AddItem(nil, 0, 1, false)
	layout.SetDirection(tview.FlexRow)

	app.SetRoot(layout, true)

	realFileList.SetPath(os.Getenv("HOME"))

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
