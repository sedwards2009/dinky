package main

import (
	"dinky/internal/tui/menu"
	"log"
	"os"

	"github.com/pgavlin/femto"
	"github.com/pgavlin/femto/runtime"
	"github.com/rivo/tview"
)

const ID_SOFT_WRAP = "softwrap"
const ID_LINE_NUMBERS = "linenumbers"

func setupLogging() *os.File {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic("Failed to open log file: " + err.Error())
	}
	log.SetOutput(logFile)
	return logFile
}

var app *tview.Application
var menus []*menu.Menu
var buffer *femto.Buffer
var editor *femto.View

func syncSoftWrap(menus []*menu.Menu, on bool) {
	for _, menu := range menus {
		for _, menuItem := range menu.Items {
			if menuItem.ID == ID_SOFT_WRAP {
				if on {
					menuItem.Title = "\u2713 "
				} else {
					menuItem.Title = "  "
				}
				menuItem.Title += "Soft Wrap"
			}
		}
	}
}

func syncLineNumbers(menus []*menu.Menu, on bool) {
	for _, menu := range menus {
		for _, menuItem := range menu.Items {
			if menuItem.ID == ID_LINE_NUMBERS {
				if on {
					menuItem.Title = "\u2713 "
				} else {
					menuItem.Title = "  "
				}
				menuItem.Title += "Line numbers"
			}
		}
	}
}

func handleSelectAll(_ string) {
	editor.SelectAll()
}

func handleLineNumbers(_ string) {
	on := buffer.Settings["ruler"].(bool)
	buffer.Settings["ruler"] = !on
	syncLineNumbers(menus, !on)
}

func handleSoftWrap(_ string) {
	on := buffer.Settings["softwrap"].(bool)
	buffer.Settings["softwrap"] = !on
	syncSoftWrap(menus, !on)
}

func syncMenuFromBuffer(buffer *femto.Buffer) {
	softwrap := buffer.Settings["softwrap"].(bool)
	syncSoftWrap(menus, softwrap)
	lineNumbers := buffer.Settings["ruler"].(bool)
	syncLineNumbers(menus, lineNumbers)
}

func handleUndo(_ string) {
	editor.Undo()
}

func handleRedo(_ string) {
	editor.Redo()
}

func handleQuit(_ string) {
	app.Stop()
}

func main() {
	logFile := setupLogging()
	defer logFile.Close()

	app = tview.NewApplication()
	app.EnableMouse(true)

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)

	menuBar := menu.NewMenuBar()

	menus = []*menu.Menu{
		{Title: "File", Items: []*menu.MenuItem{
			{Title: "New", Shortcut: "Ctrl+N", Callback: func(id string) { log.Println("New file") }},
			{Title: "Open", Shortcut: "Ctrl+O", Callback: func(id string) { log.Println("Open file") }},
			{Title: "Save", Shortcut: "Ctrl+S", Callback: func(id string) { log.Println("Save file") }},
			{Title: "", Shortcut: "", Callback: nil}, // Separator
			{Title: "Quit", Shortcut: "Ctrl+Q", Callback: handleQuit},
		}},
		{Title: "Edit", Items: []*menu.MenuItem{
			{Title: "Undo", Shortcut: "Ctrl+Z", Callback: handleUndo},
			{Title: "Redo", Shortcut: "Ctrl+Y", Callback: handleRedo},
			{Title: "", Shortcut: "", Callback: nil}, // Separator
			{Title: "Cut", Shortcut: "Ctrl+X", Callback: func(id string) { log.Println("Cut") }},
			{Title: "Copy", Shortcut: "Ctrl+C", Callback: func(id string) { log.Println("Copy") }},
			{Title: "Paste", Shortcut: "Ctrl+V", Callback: func(id string) { log.Println("Paste") }},
			{Title: "", Shortcut: "", Callback: nil}, // Separator
			{Title: "Select All", Shortcut: "Ctrl+A", Callback: handleSelectAll},
		}},
		{Title: "View", Items: []*menu.MenuItem{
			{ID: ID_LINE_NUMBERS, Title: "Line Numbers", Callback: handleLineNumbers},
			{ID: ID_SOFT_WRAP, Title: "Soft Wrap", Callback: handleSoftWrap},
		}},
		{Title: "Help", Items: []*menu.MenuItem{
			{Title: "About", Shortcut: "F1", Callback: func(id string) { log.Println("About") }},
		}},
	}
	menuBar.SetMenus(menus)

	flex.AddItem(menuBar, 1, 0, false)

	var colorscheme femto.Colorscheme
	if monokai := runtime.Files.FindFile(femto.RTColorscheme, "monokai"); monokai != nil {
		if data, err := monokai.Data(); err == nil {
			colorscheme = femto.ParseColorscheme(string(data))
		}
	}

	buffer = femto.NewBufferFromString("Hello Smoe\nSome words to click on\n", "/home/sbe/smoe.txt")

	editor = femto.NewView(buffer)

	syncMenuFromBuffer(buffer)

	editor.SetRuntimeFiles(runtime.Files)
	editor.SetColorscheme(colorscheme)
	// editor.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
	// 	switch event.Key() {
	// 	case tcell.KeyCtrlS:
	// 		return nil
	// 	case tcell.KeyCtrlQ:
	// 		app.Stop()
	// 		return nil
	// 	}
	// 	return event
	// })
	flex.AddItem(editor, 0, 1, true)

	app.SetRoot(flex, true)
	app.SetAfterDrawFunc(menuBar.AfterDraw())

	menuBar.SetOnClose(func() {
		app.SetFocus(editor)
	})

	if err := app.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
