package main

import (
	"dinky/internal/tui/filedialog"
	"dinky/internal/tui/style"
	"log"
	"os"

	nuview "github.com/rivo/tview"
)

func main() {
	logFile := setupLogging()
	defer logFile.Close()

	app := nuview.NewApplication()
	app.EnableMouse(true)

	workspace := nuview.NewBox()

	modalPages := nuview.NewPages()
	modalPages.AddPage("workspace", workspace, true, true)

	fileDialog := filedialog.NewFileDialog(app)
	style.StyleFileDialog(fileDialog)
	fileDialog.SetPath("/home/sbe")
	fileDialog.SetCompletedFunc(func(accepted bool, path string) {
		if accepted {
			log.Printf("File selected: %s\n", path)
		} else {
			log.Println("File selection cancelled.")
		}
		app.Stop()
	})
	modalPages.AddPage("fileDialog", fileDialog, true, true)

	app.SetRoot(modalPages, true)

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
