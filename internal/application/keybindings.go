package application

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/sedwards2009/smidgen"
)

var smidgenDefaultKeyBindings smidgen.Keybindings
var smidgenSingleLineKeyBindings smidgen.Keybindings
var smidgenKeyToActionMapping map[string]string
var smidgenSingleLineKeyToActionMapping map[string]string
var actionToKeyMapping map[string]string

type KeyDesc struct {
	KeyCode   tcell.Key
	Modifiers tcell.ModMask
	R         rune
}

var dinkyKeyBindings map[smidgen.KeyDesc]string
var dinkyKeyToActionMapping map[string]string

func initKeyBindings() {
	smidgenSingleLineKeyToActionMapping = map[string]string{
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
		"Ctrl-v":   smidgen.ActionPaste,
		"Ctrl-a":   smidgen.ActionSelectAll,
		"Home":     smidgen.ActionStartOfTextToggle,
		"End":      smidgen.ActionEndOfLine,
		"CtrlHome": smidgen.ActionCursorStart,
		"CtrlEnd":  smidgen.ActionCursorEnd,
		"Delete":   smidgen.ActionDelete,

		"MouseWheelUp":     smidgen.ActionScrollUp,
		"MouseWheelDown":   smidgen.ActionScrollDown,
		"MouseLeft":        smidgen.ActionMousePress,
		"MouseLeftDrag":    smidgen.ActionMouseDrag,
		"MouseLeftRelease": smidgen.ActionMouseRelease,
	}

	// Copy smidgenSimpleKeyToActionMapping into smidgenKeyToActionMapping
	smidgenKeyToActionMapping = make(map[string]string)
	for k, v := range smidgenSingleLineKeyToActionMapping {
		smidgenKeyToActionMapping[k] = v
	}
	smidgenKeyToActionMapping["Up"] = smidgen.ActionCursorUp
	smidgenKeyToActionMapping["Down"] = smidgen.ActionCursorDown
	smidgenKeyToActionMapping["ShiftUp"] = smidgen.ActionSelectUp
	smidgenKeyToActionMapping["ShiftDown"] = smidgen.ActionSelectDown
	smidgenKeyToActionMapping["AltUp"] = smidgen.ActionMoveLinesUp
	smidgenKeyToActionMapping["AltDown"] = smidgen.ActionMoveLinesDown
	smidgenKeyToActionMapping["Alt-{"] = smidgen.ActionParagraphPrevious
	smidgenKeyToActionMapping["Alt-}"] = smidgen.ActionParagraphNext
	smidgenKeyToActionMapping["Tab"] = smidgen.ActionIndentSelection + "," + smidgen.ActionInsertTab
	smidgenKeyToActionMapping["Backtab"] = smidgen.ActionOutdentSelection + "|" + smidgen.ActionOutdentLine
	smidgenKeyToActionMapping["PageUp"] = smidgen.ActionCursorPageUp
	smidgenKeyToActionMapping["PageDown"] = smidgen.ActionCursorPageDown
	// "CtrlPageUp":     "PreviousTab|LastTab,
	// "CtrlPageDown":   "NextTab|FirstTab,
	smidgenKeyToActionMapping["ShiftPageUp"] = smidgen.ActionSelectPageUp
	smidgenKeyToActionMapping["ShiftPageDown"] = smidgen.ActionSelectPageDown
	// "Ctrl-r" smidgen.ActionToggleRuler
	// "Ctrl-w":         "NextSplit|FirstSplit,
	// "Ctrl-u":         smidgen.ActionToggleMacro,
	// "Ctrl-j":         smidgen.ActionPlayMacro,
	smidgenKeyToActionMapping["Insert"] = smidgen.ActionToggleOverwriteMode
	smidgenKeyToActionMapping["Esc"] = smidgen.ActionEscape + "," + smidgen.ActionRemoveAllMultiCursors
	// "MouseMiddle":      smidgen.ActionPastePrimary,
	// "Ctrl-MouseLeft":   smidgen.ActionMouseMultiCursor,
	// "Alt-n": smidgen.ActionSpawnMultiCursor,
	smidgenKeyToActionMapping["Ctrl-d"] = smidgen.ActionSpawnMultiCursor
	smidgenKeyToActionMapping["Alt-m"] = smidgen.ActionSpawnMultiCursorSelect
	// "AltShiftUp":   smidgen.ActionSpawnMultiCursorUp,
	// "AltShiftDown": smidgen.ActionSpawnMultiCursorDown,
	smidgenKeyToActionMapping["Alt-p"] = smidgen.ActionRemoveMultiCursor
	smidgenKeyToActionMapping["Alt-c"] = smidgen.ActionRemoveAllMultiCursors
	smidgenKeyToActionMapping["Alt-x"] = smidgen.ActionSkipMultiCursor
	smidgenKeyToActionMapping["Ctrl-j"] = smidgen.ActionJumpToMatchingBrace

	dinkyKeyToActionMapping = map[string]string{
		"Ctrl-n": ACTION_NEW,
		"Ctrl-w": ACTION_CLOSE_FILE,
		"Ctrl-o": ACTION_OPEN_FILE,
		"Ctrl-s": ACTION_SAVE_FILE,
		"Ctrl-g": ACTION_GO_TO_LINE,
		"F12":    ACTION_OPEN_FILE_MENU,
		"Alt-f":  ACTION_OPEN_FILE_MENU,
		"Alt-e":  ACTION_OPEN_EDIT_MENU,
		"Alt-s":  ACTION_OPEN_SELECTION_MENU,
		"Alt-t":  ACTION_OPEN_TRANSFORM_MENU,
		"Alt-v":  ACTION_OPEN_VIEW_MENU,
		"Alt-h":  ACTION_OPEN_HELP_MENU,
		"Ctrl-q": ACTION_QUIT,
		"Ctrl-f": ACTION_FIND,
		"F3":     ACTION_FIND_NEXT,
		"F2":     ACTION_FIND_PREVIOUS,
		"F9":     ACTION_NEXT_EDITOR,
		"F10":    ACTION_PREVIOUS_EDITOR,
		"Ctrl-r": ACTION_FIND_AND_REPLACE,
	}

	actionToKeyMapping = make(map[string]string)
	for key, action := range smidgenKeyToActionMapping {
		if action == "" {
			continue
		}
		actionToKeyMapping[action] = key
	}

	smidgenDefaultKeyBindings = smidgen.ParseKeybindings(smidgenKeyToActionMapping)
	smidgenSingleLineKeyBindings = smidgen.ParseKeybindings(smidgenSingleLineKeyToActionMapping)

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
		} else {
			log.Printf("Failed to parse key sequence: %s", key)
		}
	}
}
