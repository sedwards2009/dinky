package application

import (
	"dinky/internal/tui/menu"
	"log"

	"github.com/pgavlin/femto"
)

const ID_SOFT_WRAP = "softwrap"
const ID_LINE_NUMBERS = "linenumbers"

func createMenus() []*menu.Menu {
	return []*menu.Menu{
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
}

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

func syncMenuFromBuffer(buffer *femto.Buffer) {
	softwrap := buffer.Settings["softwrap"].(bool)
	syncSoftWrap(menus, softwrap)
	lineNumbers := buffer.Settings["ruler"].(bool)
	syncLineNumbers(menus, lineNumbers)
}
