package application

import (
	"dinky/internal/tui/menu"
	"log"

	"github.com/pgavlin/femto"
)

func createMenus() []*menu.Menu {
	return []*menu.Menu{
		{Title: "File", Items: []*menu.MenuItem{
			{Title: "New", Callback: func(id string) { log.Println("New file") }},
			{Title: "Open", Callback: func(id string) { log.Println("Open file") }},
			{Title: "Save", Callback: func(id string) { log.Println("Save file") }},
			{Title: "", Callback: nil}, // Separator
			{ID: ACTION_QUIT, Title: "Quit", Callback: handleDinkyAction},
		}},
		{Title: "Edit", Items: []*menu.MenuItem{
			{ID: femto.ActionUndo, Title: "Undo", Callback: handleFemtoAction},
			{ID: femto.ActionRedo, Title: "Redo", Callback: handleFemtoAction},
			{Title: "", Callback: nil}, // Separator
			{ID: femto.ActionCut, Title: "Cut", Callback: handleFemtoAction},
			{ID: femto.ActionCopy, Title: "Copy", Callback: handleFemtoAction},
			{ID: femto.ActionPaste, Title: "Paste", Callback: handleFemtoAction},
			{Title: "", Callback: nil}, // Separator
			{ID: femto.ActionSelectAll, Title: "Select All", Callback: handleFemtoAction},
		}},
		{Title: "View", Items: []*menu.MenuItem{
			{ID: ACTION_TOGGLE_LINE_NUMBERS, Title: "Line Numbers", Callback: handleDinkyAction},
			{ID: ACTION_TOGGLE_SOFT_WRAP, Title: "Soft Wrap", Callback: handleDinkyAction},
		}},
		{Title: "Help", Items: []*menu.MenuItem{
			{Title: "About", Callback: func(id string) { log.Println("About") }},
		}},
	}
}

func syncMenuKeyBindings(menu []*menu.Menu, keyBindings femto.KeyBindings) {
	for _, menu := range menus {
		for _, menuItem := range menu.Items {
			if key, ok := femtoActionToKeyMapping[menuItem.ID]; ok {
				menuItem.Shortcut = key
			} else {
				menuItem.Shortcut = ""
			}
		}
	}
}

func syncSoftWrap(menus []*menu.Menu, on bool) {
	for _, menu := range menus {
		for _, menuItem := range menu.Items {
			if menuItem.ID == ACTION_TOGGLE_SOFT_WRAP {
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
			if menuItem.ID == ACTION_TOGGLE_LINE_NUMBERS {
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
