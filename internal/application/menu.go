package application

import (
	"dinky/internal/tui/menu"
	"strings"

	"github.com/sedwards2009/femto"
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
			{ID: ACTION_FIND, Title: "Find", Callback: handleDinkyAction},
			{Title: "", Callback: nil}, // Separator
			{ID: ACTION_SET_TAB_CHARACTER, Title: "Tab Character…", Callback: handleDinkyAction},
			{ID: ACTION_SET_LINE_ENDINGS, Title: "Line Endings…", Callback: handleDinkyAction},
		}},
		{Title: "Selection", Items: []*menu.MenuItem{
			{ID: femto.ActionSelectAll, Title: "Select All", Callback: handleFemtoAction},
			{Title: "", Callback: nil}, // Separator
			{ID: ACTION_GO_TO_LINE, Title: "Go to Line…", Callback: handleDinkyAction},
			{ID: femto.ActionJumpToMatchingBrace, Title: "Go to Bracket", Callback: handleFemtoAction},
			{Title: "", Callback: nil}, // Separator
			{ID: femto.ActionSpawnMultiCursor, Title: "Add Next Occurrence", Callback: handleFemtoAction},
			{ID: femto.ActionSpawnMultiCursorSelect, Title: "Add Cursors to Selection", Callback: handleFemtoAction},
			{ID: femto.ActionRemoveMultiCursor, Title: "Remove Last Cursor", Callback: handleFemtoAction},
			{ID: femto.ActionRemoveAllMultiCursors, Title: "Remove All Cursors", Callback: handleFemtoAction},
		}},
		{Title: "View", Items: []*menu.MenuItem{
			{ID: femto.ActionToggleRuler, Title: "Line Numbers", Callback: handleFemtoAction},
			{ID: ACTION_TOGGLE_SOFT_WRAP, Title: "Soft Wrap", Callback: handleDinkyAction},
			{ID: ACTION_TOGGLE_MATCH_BRACKET, Title: "Match Brackets", Callback: handleDinkyAction},
			{ID: ACTION_SET_TAB_SIZE, Title: "Tab Size…", Callback: handleDinkyAction},
			{ID: ACTION_SET_SYNTAX_HIGHLIGHTING, Title: "Syntax…", Callback: handleDinkyAction},
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

func syncToggleMenuItem(menus []*menu.Menu, actionID string, label string, on bool) {
	for _, menu := range menus {
		for _, menuItem := range menu.Items {
			if menuItem.ID == actionID {
				if on {
					menuItem.Title = "\u2713 "
				} else {
					menuItem.Title = "  "
				}
				menuItem.Title += label
			}
		}
	}
}

func syncSoftWrap(menus []*menu.Menu, on bool) {
	syncToggleMenuItem(menus, ACTION_TOGGLE_SOFT_WRAP, "Soft Wrap", on)
}

func syncLineNumbers(menus []*menu.Menu, on bool) {
	syncToggleMenuItem(menus, femto.ActionToggleRuler, "Line Numbers", on)
}

func syncMatchBracket(menus []*menu.Menu, on bool) {
	syncToggleMenuItem(menus, ACTION_TOGGLE_MATCH_BRACKET, "Match Brackets", on)
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

func syncTabCharacter(menus []*menu.Menu, buffer *femto.Buffer) {
	for _, menu := range menus {
		for _, menuItem := range menu.Items {
			if menuItem.ID == ACTION_SET_TAB_CHARACTER {
				if buffer == nil {
					menuItem.Title = "Tab Character…"
					return
				}
				expandTab, ok := buffer.Settings["tabstospaces"]
				if !ok || !expandTab.(bool) {
					menuItem.Title = "Tab Character (Tab)…"
				} else {
					menuItem.Title = "Tab Character (Space)…"
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

func syncSyntaxHighlighting(menus []*menu.Menu, buffer *femto.Buffer) {
	for _, menu := range menus {
		for _, menuItem := range menu.Items {
			if menuItem.ID == ACTION_SET_SYNTAX_HIGHLIGHTING {
				if buffer == nil {
					menuItem.Title = "Syntax…"
					return
				}

				currentFiletype := buffer.FileType()
				if currentFiletype == "Unknown" || currentFiletype == "" {
					menuItem.Title = "Syntax (Auto)…"
				} else {
					// Capitalize the first letter for display
					var displayName string
					if len(currentFiletype) > 0 {
						displayName = strings.ToUpper(string(currentFiletype[0])) + currentFiletype[1:]
					} else {
						displayName = currentFiletype
					}
					menuItem.Title = "Syntax (" + displayName + ")…"
				}
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
	syncTabCharacter(menus, buffer)
	syncLineEndings(menus, buffer)
	syncSyntaxHighlighting(menus, buffer)
}
