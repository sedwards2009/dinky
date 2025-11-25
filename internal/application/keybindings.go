package application

import (
	"github.com/gdamore/tcell/v2"
	"github.com/sedwards2009/smidgen"
)

var femtoDefaultKeyBindings smidgen.Keybindings
var femtoSingleLineKeyBindings smidgen.Keybindings
var femtoKeyToActionMapping map[string]string
var femtoSingleLineKeyToActionMapping map[string]string
var actionToKeyMapping map[string]string

type KeyDesc struct {
	KeyCode   tcell.Key
	Modifiers tcell.ModMask
	R         rune
}

var dinkyKeyBindings map[smidgen.KeyDesc]string
var dinkyKeyToActionMapping map[string]string

func initKeyBindings() {
	femtoSingleLineKeyToActionMapping = map[string]string{
		"Right": smidgen.ActionCursorRight,
		"Left":  smidgen.ActionCursorLeft,

		"ShiftLeft":  smidgen.ActionSelectLeft,
		"ShiftRight": smidgen.ActionSelectRight,
		"CtrlLeft":   smidgen.ActionWordLeft,
		"CtrlRight":  smidgen.ActionWordRight,

		"CtrlShiftRight": smidgen.ActionSelectWordRight,
		"CtrlShiftLeft":  smidgen.ActionSelectWordLeft,
		// "AltLeft":        smidgen.ActionStartOfTextToggle,
		"AltRight": smidgen.ActionEndOfLine,
		// "AltShiftLeft":   smidgen.ActionSelectToStartOfTextToggle,
		// "ShiftHome":      smidgen.ActionSelectToStartOfTextToggle,
		"AltShiftRight": smidgen.ActionSelectToEndOfLine,
		"ShiftEnd":      smidgen.ActionSelectToEndOfLine,
		"CtrlUp":        smidgen.ActionCursorStart,
		"CtrlDown":      smidgen.ActionCursorEnd,
		"CtrlShiftUp":   smidgen.ActionSelectToStart,
		"CtrlShiftDown": smidgen.ActionSelectToEnd,
		"Enter":         smidgen.ActionInsertNewline,
		"CtrlH":         smidgen.ActionBackspace,
		"Backspace":     smidgen.ActionBackspace,
		"OldBackspace":  smidgen.ActionBackspace,
		"Alt-CtrlH":     smidgen.ActionDeleteWordLeft,
		"Alt-Backspace": smidgen.ActionDeleteWordLeft,
		// "Tab":            "Autocomplete|IndentSelection|InsertTab,
		// "Ctrl-o":  smidgen.ActionOpenFile,
		// "Ctrl-s":  smidgen.ActionSave,
		// "Ctrl-f":  smidgen.ActionFind,
		// "Alt-F":   smidgen.ActionFindLiteral,
		// "Ctrl-n":  smidgen.ActionFindNext,
		// "Ctrl-p":  smidgen.ActionFindPrevious,
		// "Alt-[":          "DiffPrevious|CursorStart",
		// "Alt-]":          "DiffNext|CursorEnd,
		"Ctrl-z": smidgen.ActionUndo,
		"Ctrl-y": smidgen.ActionRedo,
		"Ctrl-c": smidgen.ActionCopy,
		"Ctrl-x": smidgen.ActionCut,
		"Ctrl-k": smidgen.ActionCutLine,
		// "Ctrl-d": smidgen.ActionDuplicateLine,
		"Ctrl-v": smidgen.ActionPaste,
		"Ctrl-a": smidgen.ActionSelectAll,

		"Home":     smidgen.ActionStartOfLine,
		"End":      smidgen.ActionEndOfLine,
		"CtrlHome": smidgen.ActionCursorStart,
		"CtrlEnd":  smidgen.ActionCursorEnd,
		"Delete":   smidgen.ActionDelete,
	}

	// Copy femtoSimpleKeyToActionMapping into femtoKeyToActionMapping
	femtoKeyToActionMapping = make(map[string]string)
	for k, v := range femtoSingleLineKeyToActionMapping {
		femtoKeyToActionMapping[k] = v
	}
	femtoKeyToActionMapping["Up"] = smidgen.ActionCursorUp
	femtoKeyToActionMapping["Down"] = smidgen.ActionCursorDown
	femtoKeyToActionMapping["ShiftUp"] = smidgen.ActionSelectUp
	femtoKeyToActionMapping["ShiftDown"] = smidgen.ActionSelectDown
	femtoKeyToActionMapping["AltUp"] = smidgen.ActionMoveLinesUp
	femtoKeyToActionMapping["AltDown"] = smidgen.ActionMoveLinesDown
	femtoKeyToActionMapping["Alt-{"] = smidgen.ActionParagraphPrevious
	femtoKeyToActionMapping["Alt-}"] = smidgen.ActionParagraphNext
	femtoKeyToActionMapping["Tab"] = smidgen.ActionIndentSelection + "," + smidgen.ActionInsertTab
	femtoKeyToActionMapping["Backtab"] = "CycleAutocompleteBack|OutdentSelection|OutdentLine"
	femtoKeyToActionMapping["PageUp"] = smidgen.ActionCursorPageUp
	femtoKeyToActionMapping["PageDown"] = smidgen.ActionCursorPageDown
	// "CtrlPageUp":     "PreviousTab|LastTab,
	// "CtrlPageDown":   "NextTab|FirstTab,
	femtoKeyToActionMapping["ShiftPageUp"] = smidgen.ActionSelectPageUp
	femtoKeyToActionMapping["ShiftPageDown"] = smidgen.ActionSelectPageDown
	// "Ctrl-r" smidgen.ActionToggleRuler
	// "Ctrl-w":         "NextSplit|FirstSplit,
	// "Ctrl-u":         smidgen.ActionToggleMacro,
	// "Ctrl-j":         smidgen.ActionPlayMacro,
	femtoKeyToActionMapping["Insert"] = smidgen.ActionToggleOverwriteMode
	femtoKeyToActionMapping["Esc"] = smidgen.ActionEscape + "," + smidgen.ActionRemoveAllMultiCursors
	// "MouseMiddle":      smidgen.ActionPastePrimary,
	// "Ctrl-MouseLeft":   smidgen.ActionMouseMultiCursor,
	// "Alt-n": smidgen.ActionSpawnMultiCursor,
	femtoKeyToActionMapping["Ctrl-d"] = smidgen.ActionSpawnMultiCursor
	femtoKeyToActionMapping["Alt-m"] = smidgen.ActionSpawnMultiCursorSelect
	// "AltShiftUp":   smidgen.ActionSpawnMultiCursorUp,
	// "AltShiftDown": smidgen.ActionSpawnMultiCursorDown,
	femtoKeyToActionMapping["Alt-p"] = smidgen.ActionRemoveMultiCursor
	femtoKeyToActionMapping["Alt-c"] = smidgen.ActionRemoveAllMultiCursors
	femtoKeyToActionMapping["Alt-x"] = smidgen.ActionSkipMultiCursor
	femtoKeyToActionMapping["Ctrl-j"] = smidgen.ActionJumpToMatchingBrace

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
		"F9":     ACTION_NEXT_EDITOR,
		"F10":    ACTION_PREVIOUS_EDITOR,
		"Ctrl-r": ACTION_FIND_AND_REPLACE,
	}

	actionToKeyMapping = make(map[string]string)
	for key, action := range femtoKeyToActionMapping {
		if action == "" {
			continue
		}
		actionToKeyMapping[action] = key
	}

	femtoDefaultKeyBindings = smidgen.ParseKeybindings(femtoKeyToActionMapping)
	femtoSingleLineKeyBindings = smidgen.ParseKeybindings(femtoSingleLineKeyToActionMapping)

	for key, action := range dinkyKeyToActionMapping {
		if action == "" {
			continue
		}
		actionToKeyMapping[action] = key
	}

	dinkyKeyBindings = make(map[smidgen.KeyDesc]string)
	for key, action := range dinkyKeyToActionMapping {
		if desc, ok := smidgen.ParseKeySequence(key); ok {
			dinkyKeyBindings[desc] = action
		}
	}
}
