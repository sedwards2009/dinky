package application

import (
	"github.com/sedwards2009/femto"
)

var femtoDefaultKeyBindings femto.KeyBindings
var femtoSingleLineKeyBindings femto.KeyBindings
var femtoKeyToActionMapping map[string]string
var femtoSingleLineKeyToActionMapping map[string]string
var actionToKeyMapping map[string]string
var dinkyKeyBindings map[femto.KeyDesc]string
var dinkyKeyToActionMapping map[string]string

func initKeyBindings() {
	femtoSingleLineKeyToActionMapping = map[string]string{
		"Right": femto.ActionCursorRight,
		"Left":  femto.ActionCursorLeft,

		"ShiftLeft":  femto.ActionSelectLeft,
		"ShiftRight": femto.ActionSelectRight,
		"CtrlLeft":   femto.ActionWordLeft,
		"CtrlRight":  femto.ActionWordRight,

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
		"Enter":         femto.ActionInsertNewline,
		"CtrlH":         femto.ActionBackspace,
		"Backspace":     femto.ActionBackspace,
		"OldBackspace":  femto.ActionBackspace,
		"Alt-CtrlH":     femto.ActionDeleteWordLeft,
		"Alt-Backspace": femto.ActionDeleteWordLeft,
		// "Tab":            "Autocomplete|IndentSelection|InsertTab,
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
		"Delete":   femto.ActionDelete,
	}

	// Copy femtoSimpleKeyToActionMapping into femtoKeyToActionMapping
	femtoKeyToActionMapping = make(map[string]string)
	for k, v := range femtoSingleLineKeyToActionMapping {
		femtoKeyToActionMapping[k] = v
	}
	femtoKeyToActionMapping["Up"] = femto.ActionCursorUp
	femtoKeyToActionMapping["Down"] = femto.ActionCursorDown
	femtoKeyToActionMapping["ShiftUp"] = femto.ActionSelectUp
	femtoKeyToActionMapping["ShiftDown"] = femto.ActionSelectDown
	femtoKeyToActionMapping["AltUp"] = femto.ActionMoveLinesUp
	femtoKeyToActionMapping["AltDown"] = femto.ActionMoveLinesDown
	femtoKeyToActionMapping["Alt-{"] = femto.ActionParagraphPrevious
	femtoKeyToActionMapping["Alt-}"] = femto.ActionParagraphNext
	femtoKeyToActionMapping["Tab"] = femto.ActionIndentSelection + "," + femto.ActionInsertTab
	femtoKeyToActionMapping["Backtab"] = "CycleAutocompleteBack|OutdentSelection|OutdentLine"
	femtoKeyToActionMapping["PageUp"] = femto.ActionCursorPageUp
	femtoKeyToActionMapping["PageDown"] = femto.ActionCursorPageDown
	// "CtrlPageUp":     "PreviousTab|LastTab,
	// "CtrlPageDown":   "NextTab|FirstTab,
	femtoKeyToActionMapping["ShiftPageUp"] = femto.ActionSelectPageUp
	femtoKeyToActionMapping["ShiftPageDown"] = femto.ActionSelectPageDown
	femtoKeyToActionMapping["Ctrl-r"] = femto.ActionToggleRuler
	// "Ctrl-w":         "NextSplit|FirstSplit,
	// "Ctrl-u":         femto.ActionToggleMacro,
	// "Ctrl-j":         femto.ActionPlayMacro,
	femtoKeyToActionMapping["Insert"] = femto.ActionToggleOverwriteMode
	femtoKeyToActionMapping["Esc"] = femto.ActionEscape + "," + femto.ActionRemoveAllMultiCursors
	// "MouseMiddle":      femto.ActionPastePrimary,
	// "Ctrl-MouseLeft":   femto.ActionMouseMultiCursor,
	// "Alt-n": femto.ActionSpawnMultiCursor,
	femtoKeyToActionMapping["Ctrl-d"] = femto.ActionSpawnMultiCursor
	femtoKeyToActionMapping["Alt-m"] = femto.ActionSpawnMultiCursorSelect
	// "AltShiftUp":   femto.ActionSpawnMultiCursorUp,
	// "AltShiftDown": femto.ActionSpawnMultiCursorDown,
	femtoKeyToActionMapping["Alt-p"] = femto.ActionRemoveMultiCursor
	femtoKeyToActionMapping["Alt-c"] = femto.ActionRemoveAllMultiCursors
	femtoKeyToActionMapping["Alt-x"] = femto.ActionSkipMultiCursor
	femtoKeyToActionMapping["Ctrl-j"] = femto.ActionJumpToMatchingBrace

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
	femtoSingleLineKeyBindings = femto.NewKeyBindings(femtoSingleLineKeyToActionMapping)

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
