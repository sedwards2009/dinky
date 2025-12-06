package application

import (
	"dinky/internal/tui/dialog"
	"dinky/internal/tui/style"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var inputDialog *dialog.InputDialog

const inputDialogName = "inputdialog"

func ShowInputDialog(title string, message string, defaultValue string, onCancel func(), onAccept func(value string,
	index int), fieldKeyFilter func(event *tcell.EventKey) bool) tview.Primitive {

	if inputDialog == nil {
		inputDialog = dialog.NewInputDialog(app)
		inputDialog.SetSmidgenKeybindings(smidgenSingleLineKeyBindings)
	}

	width := 50
	height := 7

	// Calculate minimum width based on title and message
	titleWidth := len([]rune(title)) + 4
	messageWidth := len([]rune(message)) + 4
	if titleWidth > width {
		width = titleWidth
	}
	if messageWidth > width {
		width = messageWidth
	}

	modalPages.AddPage(inputDialogName, inputDialog, true, true)

	options := dialog.InputDialogOptions{
		Title:          title,
		Message:        message,
		DefaultValue:   defaultValue,
		Buttons:        []string{"OK", "Cancel"},
		Width:          width,
		Height:         height,
		OnCancel:       onCancel,
		OnAccept:       onAccept,
		FieldKeyFilter: fieldKeyFilter,
	}

	inputDialog.Open(options)
	style.StyleInputDialog(inputDialog)
	return inputDialog
}

func CloseInputDialog() {
	if inputDialog != nil {
		inputDialog.Close()
		modalPages.RemovePage(inputDialogName)
	}
}

func numericInputFilter(event *tcell.EventKey) bool {
	key := event.Key()
	// Allow digits and basic editing keys
	if key == tcell.KeyBackspace || key == tcell.KeyDelete ||
		key == tcell.KeyLeft || key == tcell.KeyRight ||
		key == tcell.KeyHome || key == tcell.KeyEnd ||
		key == tcell.KeyDEL {
		return true
	}
	if event.Rune() >= '0' && event.Rune() <= '9' {
		return true
	}
	return false
}
