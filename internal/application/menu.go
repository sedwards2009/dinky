package application

import (
	"dinky/internal/tui/menu"
	"strconv"
	"strings"

	"github.com/sedwards2009/smidgen"
	"github.com/sedwards2009/smidgen/micro/buffer"
)

func createMenus() []*menu.Menu {
	return []*menu.Menu{
		{Title: "[::u]F[::U]ile", Shortcut: 'f', Items: []*menu.MenuItem{
			{ID: ACTION_NEW, Title: "New", Callback: handleDinkyAction},
			{ID: ACTION_OPEN_FILE, Title: "Open", Callback: handleDinkyAction},
			{ID: ACTION_SAVE_FILE, Title: "Save", Callback: handleDinkyAction},
			{ID: ACTION_SAVE_FILE_AS, Title: "Save As…", Callback: handleDinkyAction},
			{ID: ACTION_CLOSE_FILE, Title: "Close", Callback: handleDinkyAction},
			{Title: "", Callback: nil}, // Separator
			{ID: ACTION_SUSPEND, Title: "Suspend", Callback: handleDinkyAction},
			{Title: "", Callback: nil}, // Separator
			{ID: ACTION_SETTINGS, Title: "Settings", Callback: handleDinkyAction},
			{Title: "", Callback: nil}, // Separator
			{ID: ACTION_QUIT, Title: "Quit", Callback: handleDinkyAction},
		}},
		{Title: "[::u]E[::U]dit", Shortcut: 'e', Items: []*menu.MenuItem{
			{ID: smidgen.ActionUndo, Title: "Undo", Callback: handleSmidgenAction},
			{ID: smidgen.ActionRedo, Title: "Redo", Callback: handleSmidgenAction},
			{Title: "", Callback: nil}, // Separator
			{ID: smidgen.ActionCut, Title: "Cut", Callback: handleSmidgenAction},
			{ID: smidgen.ActionCopy, Title: "Copy", Callback: handleSmidgenAction},
			{ID: smidgen.ActionPaste, Title: "Paste", Callback: handleSmidgenAction},
			{Title: "", Callback: nil}, // Separator
			{ID: ACTION_FIND, Title: "Find", Callback: handleDinkyAction},
			{ID: ACTION_FIND_NEXT, Title: "Find Next", Callback: handleDinkyAction},
			{ID: ACTION_FIND_PREVIOUS, Title: "Find Previous", Callback: handleDinkyAction},
			{ID: ACTION_FIND_AND_REPLACE, Title: "Find & Replace", Callback: handleDinkyAction},
			{Title: "", Callback: nil}, // Separator
			{ID: ACTION_SET_TAB_CHARACTER, Title: "Tab Character…", Callback: handleDinkyAction},
			{ID: ACTION_SET_LINE_ENDINGS, Title: "Line Endings…", Callback: handleDinkyAction},
			{ID: ACTION_CONVERT_TAB_SPACES, Title: "Convert All Tabs to Spaces", Callback: handleDinkyAction},
		}},
		{Title: "[::u]S[::U]election", Shortcut: 's', Items: []*menu.MenuItem{
			{ID: smidgen.ActionSelectAll, Title: "Select All", Callback: handleSmidgenAction},
			{ID: smidgen.ActionSetManualSelectionStart, Title: "Set Selection Start", Callback: handleSmidgenAction},
			{ID: smidgen.ActionSetManualSelectionEnd, Title: "Set Selection End", Callback: handleSmidgenAction},
			{Title: "", Callback: nil}, // Separator
			{ID: smidgen.ActionSpawnMultiCursor, Title: "Add Next Occurrence", Callback: handleSmidgenAction},
			{ID: smidgen.ActionSpawnMultiCursorSelect, Title: "Add Cursors to Selection", Callback: handleSmidgenAction},
			{ID: smidgen.ActionRemoveMultiCursor, Title: "Remove Last Cursor", Callback: handleSmidgenAction},
			{ID: smidgen.ActionRemoveAllMultiCursors, Title: "Remove All Cursors", Callback: handleSmidgenAction},
		}},
		{Title: "[::u]G[::U]o", Shortcut: 'g', Items: []*menu.MenuItem{
			{ID: ACTION_GO_TO_LINE, Title: "Go to Line…", Callback: handleDinkyAction},
			{Title: "", Callback: nil}, // Separator
			{ID: smidgen.ActionJumpToMatchingBrace, Title: "Go to Bracket", Callback: handleSmidgenAction},
			{ID: smidgen.ActionParagraphPrevious, Title: "Previous Paragraph", Callback: handleSmidgenAction},
			{ID: smidgen.ActionParagraphNext, Title: "Next Paragraph", Callback: handleSmidgenAction},
			{ID: ACTION_NEXT_EDITOR, Title: "Next Editor", Callback: handleDinkyAction},
			{ID: ACTION_PREVIOUS_EDITOR, Title: "Previous Editor", Callback: handleDinkyAction},
			{Title: "", Callback: nil}, // Separator
			{ID: smidgen.ActionToggleBookmark, Title: "Toggle Bookmark", Callback: handleSmidgenAction},
			{ID: ACTION_NEXT_BOOKMARK, Title: "Next Bookmark", Callback: handleDinkyAction},
			{ID: ACTION_PREVIOUS_BOOKMARK, Title: "Previous Bookmark", Callback: handleDinkyAction},
		}},
		{Title: "[::u]T[::U]ransform", Shortcut: 't', Items: []*menu.MenuItem{
			{ID: ACTION_TO_UPPERCASE, Title: "To Uppercase", Callback: handleDinkyAction},
			{ID: ACTION_TO_LOWERCASE, Title: "To Lowercase", Callback: handleDinkyAction},
			{ID: ACTION_URL_ENCODE, Title: "URL Encode", Callback: handleDinkyAction},
			{ID: ACTION_URL_DECODE, Title: "URL Decode", Callback: handleDinkyAction},
			{Title: "", Callback: nil}, // Separator
			{ID: ACTION_SORT_LINES, Title: "Sort Lines", Callback: handleDinkyAction},
			{ID: ACTION_REVERSE_LINES, Title: "Reverse Lines", Callback: handleDinkyAction},
			{Title: "", Callback: nil}, // Separator
			{ID: ACTION_HARD_WORD_WRAP, Title: "Hard Word Wrap", Callback: handleDinkyAction},
			{Title: "", Callback: nil}, // Separator
			{ID: ACTION_FORMAT_JSON, Title: "Format JSON", Callback: handleDinkyAction},
		}},
		{Title: "[::u]V[::U]iew", Shortcut: 'v', Items: []*menu.MenuItem{
			{ID: smidgen.ActionToggleRuler, Title: "Line Numbers", Callback: handleSmidgenAction},
			{ID: ACTION_TOGGLE_WHITESPACE, Title: "Show Whitespace", Callback: handleDinkyAction},
			{ID: ACTION_TOGGLE_TRAILING_WHITESPACE, Title: "Show Trailing Whitespace", Callback: handleDinkyAction},
			{ID: ACTION_TOGGLE_SOFT_WRAP, Title: "Soft Wrap", Callback: handleDinkyAction},
			{ID: ACTION_TOGGLE_MATCH_BRACKET, Title: "Match Brackets", Callback: handleDinkyAction},
			{ID: ACTION_SET_TAB_SIZE, Title: "Tab Size…", Callback: handleDinkyAction},
			{ID: ACTION_SET_VERTICAL_RULER, Title: "Vertical Ruler…", Callback: handleDinkyAction},
			{ID: ACTION_SET_SYNTAX_HIGHLIGHTING, Title: "Syntax…", Callback: handleDinkyAction},
		}},
		{Title: "[::u]H[::U]elp", Shortcut: 'h', Items: []*menu.MenuItem{
			{ID: ACTION_ABOUT, Title: "About", Callback: handleDinkyAction},
		}},
	}
}

func syncMenuKeyBindings(menus []*menu.Menu, smidgenActionToKeyMapping map[string]string) {
	for _, menu := range menus {
		for _, menuItem := range menu.Items {
			if key, ok := smidgenActionToKeyMapping[menuItem.ID]; ok {
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
	syncToggleMenuItem(menus, smidgen.ActionToggleRuler, "Line Numbers", on)
}

func syncMatchBracket(menus []*menu.Menu, on bool) {
	syncToggleMenuItem(menus, ACTION_TOGGLE_MATCH_BRACKET, "Match Brackets", on)
}

func syncShowWhitespace(menus []*menu.Menu, on bool) {
	syncToggleMenuItem(menus, ACTION_TOGGLE_WHITESPACE, "Show Whitespace", on)
}

func syncShowTrailingWhitespace(menus []*menu.Menu, on bool) {
	syncToggleMenuItem(menus, ACTION_TOGGLE_TRAILING_WHITESPACE, "Show Trailing Whitespace", on)
}

func syncLineEndings(menus []*menu.Menu, buffer *buffer.Buffer) {
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

func syncTabCharacter(menus []*menu.Menu, buffer *buffer.Buffer) {
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
			if menuItem.ID == ACTION_CONVERT_TAB_SPACES {
				expandTab, ok := buffer.Settings["tabstospaces"]
				if !ok || !expandTab.(bool) {
					menuItem.Title = "Convert All Spaces to Tabs"
				} else {
					menuItem.Title = "Convert All Tabs to Spaces"
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

func syncSyntaxHighlighting(menus []*menu.Menu, buffer *buffer.Buffer) {
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

func syncVerticalRuler(menus []*menu.Menu, colorColumn int) {
	for _, menu := range menus {
		for _, menuItem := range menu.Items {
			if menuItem.ID == ACTION_SET_VERTICAL_RULER {
				if colorColumn == 0 {
					menuItem.Title = "Vertical Ruler (Off)…"
				} else {
					menuItem.Title = "Vertical Ruler (" + strconv.Itoa(colorColumn) + ")…"
				}
			}
		}
	}
}

func syncMenuFromBuffer(buffer *buffer.Buffer) {
	softwrap := buffer.Settings["softwrap"].(bool)
	syncSoftWrap(menus, softwrap)
	lineNumbers := buffer.Settings["ruler"].(bool)
	syncLineNumbers(menus, lineNumbers)
	syncMatchBracket(menus, buffer.Settings["matchbrace"].(bool))
	syncShowWhitespace(menus, buffer.Settings["showwhitespace"].(bool))
	syncShowTrailingWhitespace(menus, buffer.Settings["hltrailingws"].(bool))
	syncTabSize(menus, int(buffer.Settings["tabsize"].(float64)))
	syncTabCharacter(menus, buffer)
	syncLineEndings(menus, buffer)
	syncSyntaxHighlighting(menus, buffer)
	syncVerticalRuler(menus, int(buffer.Settings["colorcolumn"].(float64)))
}
