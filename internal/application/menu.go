package application

import (
	"dinky/internal/tui/menu"

	"github.com/pgavlin/femto"
)

func createMenus() []*menu.Menu {
	return []*menu.Menu{
		{Title: "File", Items: []*menu.MenuItem{
			{ID: ACTION_NEW, Title: "New", Callback: handleDinkyAction},
			{ID: ACTION_OPEN_FILE, Title: "Open", Callback: handleDinkyAction},
			{ID: ACTION_SAVE_FILE, Title: "Save", Callback: handleDinkyAction},
			{ID: ACTION_SAVE_FILE_AS, Title: "Save As…", Callback: handleDinkyAction},
			{ID: ACTION_CLOSE_FILE, Title: "Close", Callback: handleDinkyAction},
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
			{Title: "", Callback: nil}, // Separator
			{ID: ACTION_SET_LINE_ENDINGS, Title: "Line Endings…", Callback: handleDinkyAction},
		}},
		{Title: "View", Items: []*menu.MenuItem{
			{ID: ACTION_TOGGLE_LINE_NUMBERS, Title: "Line Numbers", Callback: handleDinkyAction},
			{ID: ACTION_TOGGLE_SOFT_WRAP, Title: "Soft Wrap", Callback: handleDinkyAction},
			{ID: ACTION_TOGGLE_MATCH_BRACKET, Title: "Match Brackets", Callback: handleDinkyAction},
			{ID: ACTION_SET_TAB_SIZE, Title: "Tab Size…", Callback: handleDinkyAction},
			{Title: "", Callback: nil}, // Separator
			{ID: ACTION_GO_TO_LINE, Title: "Go to Line…", Callback: handleDinkyAction},
		}},
		{Title: "Help", Items: []*menu.MenuItem{
			{ID: ACTION_ABOUT, Title: "About", Callback: handleDinkyAction},
		}},
	}
}

func syncMenuKeyBindings(menus []*menu.Menu, femtoActionToKeyMapping map[string]string) {
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
				menuItem.Title += "Line Numbers"
			}
		}
	}
}

func syncMatchBracket(menus []*menu.Menu, on bool) {
	for _, menu := range menus {
		for _, menuItem := range menu.Items {
			if menuItem.ID == ACTION_TOGGLE_MATCH_BRACKET {
				if on {
					menuItem.Title = "\u2713 "
				} else {
					menuItem.Title = "  "
				}
				menuItem.Title += "Match Brackets"
			}
		}
	}
}

func syncLineEndings(menus []*menu.Menu, buffer *femto.Buffer) {
	for _, menu := range menus {
		for _, menuItem := range menu.Items {
			if menuItem.ID == ACTION_SET_LINE_ENDINGS {
				if isBufferCRLF(buffer) {
					menuItem.Title = "Line Endings (CRLF)…"
				} else {
					menuItem.Title = "Line Endings (LF)…"
				}
			}
		}
	}
}

func syncTabSize(menus []*menu.Menu, size int) {
	for _, menu := range menus {
		for _, menuItem := range menu.Items {
			if menuItem.ID == ACTION_SET_TAB_SIZE {
				menuItem.Title = "Tab Size (" + string(rune('0'+size)) + ")…"
			}
		}
	}
}

func syncMenuFromBuffer(buffer *femto.Buffer) {
	softwrap := buffer.Settings["softwrap"].(bool)
	syncSoftWrap(menus, softwrap)
	lineNumbers := buffer.Settings["ruler"].(bool)
	syncLineNumbers(menus, lineNumbers)
	matchBracket := buffer.Settings["matchbrace"].(bool)
	syncMatchBracket(menus, matchBracket)
	syncTabSize(menus, int(buffer.Settings["tabsize"].(float64)))
	syncLineEndings(menus, buffer)
}
