package application

import (
	"github.com/sedwards2009/femto"
)

var femtoDefaultKeyBindings femto.KeyBindings
var femtoKeyToActionMapping map[string]string
var actionToKeyMapping map[string]string
var dinkyKeyBindings map[femto.KeyDesc]string
var dinkyKeyToActionMapping map[string]string

func initKeyBindings() {
	femtoKeyToActionMapping = map[string]string{
		"Up":    femto.ActionCursorUp,
		"Down":  femto.ActionCursorDown,
		"Right": femto.ActionCursorRight,
		"Left":  femto.ActionCursorLeft,

		"ShiftUp":        femto.ActionSelectUp,
		"ShiftDown":      femto.ActionSelectDown,
		"ShiftLeft":      femto.ActionSelectLeft,
		"ShiftRight":     femto.ActionSelectRight,
		"CtrlLeft":       femto.ActionWordLeft,
		"CtrlRight":      femto.ActionWordRight,
		"AltUp":          femto.ActionMoveLinesUp,
		"AltDown":        femto.ActionMoveLinesDown,
		"CtrlShiftRight": femto.ActionSelectWordRight,
		"CtrlShiftLeft":  femto.ActionSelectWordLeft,
		// "AltLeft":        femto.ActionStartOfTextToggle,
		"AltRight": femto.ActionEndOfLine,
		// "AltShiftLeft":   femto.ActionSelectToStartOfTextToggle,
		// "ShiftHome":      femto.ActionSelectToStartOfTextToggle,
		"AltShiftRight": femto.ActionSelectToEndOfLine,
		"ShiftEnd":      femto.ActionSelectToEndOfLine,
		"CtrlUp":        femto.ActionCursorStart,
		"CtrlDown":      femto.ActionCursorEnd,
		"CtrlShiftUp":   femto.ActionSelectToStart,
		"CtrlShiftDown": femto.ActionSelectToEnd,
		"Alt-{":         femto.ActionParagraphPrevious,
		"Alt-}":         femto.ActionParagraphNext,
		"Enter":         femto.ActionInsertNewline,
		"CtrlH":         femto.ActionBackspace,
		"Backspace":     femto.ActionBackspace,
		"OldBackspace":  femto.ActionBackspace,
		"Alt-CtrlH":     femto.ActionDeleteWordLeft,
		"Alt-Backspace": femto.ActionDeleteWordLeft,
		// "Tab":            "Autocomplete|IndentSelection|InsertTab,
		"Tab":     femto.ActionIndentSelection + "," + femto.ActionInsertTab,
		"Backtab": "CycleAutocompleteBack|OutdentSelection|OutdentLine",
		// "Ctrl-o":  femto.ActionOpenFile,
		// "Ctrl-s":  femto.ActionSave,
		// "Ctrl-f":  femto.ActionFind,
		// "Alt-F":   femto.ActionFindLiteral,
		// "Ctrl-n":  femto.ActionFindNext,
		// "Ctrl-p":  femto.ActionFindPrevious,
		// "Alt-[":          "DiffPrevious|CursorStart",
		// "Alt-]":          "DiffNext|CursorEnd,
		"Ctrl-z": femto.ActionUndo,
		"Ctrl-y": femto.ActionRedo,
		"Ctrl-c": femto.ActionCopy,
		"Ctrl-x": femto.ActionCut,
		"Ctrl-k": femto.ActionCutLine,
		// "Ctrl-d": femto.ActionDuplicateLine,
		"Ctrl-v": femto.ActionPaste,
		"Ctrl-a": femto.ActionSelectAll,

		"Home":     femto.ActionStartOfLine,
		"End":      femto.ActionEndOfLine,
		"CtrlHome": femto.ActionCursorStart,
		"CtrlEnd":  femto.ActionCursorEnd,
		"PageUp":   femto.ActionCursorPageUp,
		"PageDown": femto.ActionCursorPageDown,
		// "CtrlPageUp":     "PreviousTab|LastTab,
		// "CtrlPageDown":   "NextTab|FirstTab,
		"ShiftPageUp":   femto.ActionSelectPageUp,
		"ShiftPageDown": femto.ActionSelectPageDown,
		"Ctrl-r":        femto.ActionToggleRuler,
		"Delete":        femto.ActionDelete,
		// "Ctrl-w":         "NextSplit|FirstSplit,
		// "Ctrl-u":         femto.ActionToggleMacro,
		// "Ctrl-j":         femto.ActionPlayMacro,
		"Insert": femto.ActionToggleOverwriteMode,

		// Emacs-style keybindings
		// "Alt-f": femto.ActionWordRight,
		// "Alt-b": femto.ActionWordLeft,
		// "Alt-a": femto.ActionStartOfText,
		// "Alt-e": femto.ActionEndOfLine,
		// "Alt-p": femto.ActionCursorUp,
		// "Alt-n": femto.ActionCursorDown,

		"Esc": femto.ActionEscape + "," + femto.ActionRemoveAllMultiCursors,

		// "MouseMiddle":      femto.ActionPastePrimary,
		// "Ctrl-MouseLeft":   femto.ActionMouseMultiCursor,

		// "Alt-n": femto.ActionSpawnMultiCursor,
		"Ctrl-d": femto.ActionSpawnMultiCursor,
		"Alt-m":  femto.ActionSpawnMultiCursorSelect,

		//"AltShiftUp":   femto.ActionSpawnMultiCursorUp,
		// "AltShiftDown": femto.ActionSpawnMultiCursorDown,
		"Alt-p":  femto.ActionRemoveMultiCursor,
		"Alt-c":  femto.ActionRemoveAllMultiCursors,
		"Alt-x":  femto.ActionSkipMultiCursor,
		"Ctrl-j": femto.ActionJumpToMatchingBrace,
	}

	dinkyKeyToActionMapping = map[string]string{
		"Ctrl-n": ACTION_NEW,
		"Ctrl-w": ACTION_CLOSE_FILE,
		"Ctrl-o": ACTION_OPEN_FILE,
		"Ctrl-s": ACTION_SAVE_FILE,
		"Ctrl-g": ACTION_GO_TO_LINE,
		"F12":    ACTION_OPEN_MENU,
		"Ctrl-q": ACTION_QUIT,
		"Ctrl-f": ACTION_FIND,
		"F3":     ACTION_FIND_NEXT,
		"F2":     ACTION_FIND_PREVIOUS,
	}

	actionToKeyMapping = make(map[string]string)
	for key, action := range femtoKeyToActionMapping {
		if action == "" {
			continue
		}
		actionToKeyMapping[action] = key
	}

	femtoDefaultKeyBindings = femto.NewKeyBindings(femtoKeyToActionMapping)

	for key, action := range dinkyKeyToActionMapping {
		if action == "" {
			continue
		}
		actionToKeyMapping[action] = key
	}

	dinkyKeyBindings = make(map[femto.KeyDesc]string)
	for key, action := range dinkyKeyToActionMapping {
		if desc, ok := femto.NewKeyDesc(key); ok {
			dinkyKeyBindings[desc] = action
		}
	}
}
