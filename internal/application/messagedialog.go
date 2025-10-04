package application

import (
	"dinky/internal/tui/dialog"
	"dinky/internal/tui/style"
	"strings"

	"github.com/rivo/tview"
)

var messageDialog *dialog.MessageDialog

const messageDialogName = "messagedialog"

// MeasureStringDimensions calculates the width and height in characters of a string containing CRs
func MeasureStringDimensions(text string) (width int, height int) {
	if text == "" {
		return 0, 0
	}

	lines := strings.Split(text, "\n")
	height = len(lines)
	width = 0

	for _, line := range lines {
		lineWidth := len([]rune(line)) // Use runes to handle Unicode characters properly
		if lineWidth > width {
			width = lineWidth
		}
	}

	return width, height
}

func MessageButtonsSize(buttons []string) int {
	if len(buttons) == 0 {
		return 0
	}

	width := 0
	for _, button := range buttons {
		width += len([]rune(button)) // Use runes to handle Unicode characters properly
	}

	// Add 2 for each button for padding
	width += len(buttons) * 2
	// Add 1 for each space between buttons
	width += len(buttons) - 1

	width *= 2
	return width
}

func ShowMessageDialog(title string, message string, buttons []string, OnClose func(),
	OnButtonClick func(button string, index int)) tview.Primitive {

	width, height := MeasureStringDimensions(message)
	height += 6

	minimumButtonsWidth := MessageButtonsSize(buttons)
	if width < minimumButtonsWidth {
		width = minimumButtonsWidth
	}
	width += 4

	if messageDialog == nil {
		messageDialog = dialog.NewMessageDialog(app)
	}
	modalPages.AddPage(messageDialogName, messageDialog, true, true)
	messageDialog.OnClose = OnClose
	messageDialog.OnButtonClick = OnButtonClick
	messageDialog.Open(title, message, buttons, width, height)
	style.StyleMessageDialog(messageDialog)

	return messageDialog
}

func CloseMessageDialog() {
	if messageDialog != nil {
		messageDialog.Close()
		modalPages.RemovePage("messagedialog")
	}
}

func ShowConfirmDialog(message string, onConfirm func(), onCancel func()) tview.Primitive {
	return ShowMessageDialog("Confirm", message, []string{"OK", "Cancel"},
		func() {
			CloseMessageDialog()
			onCancel()
		},
		func(button string, index int) {
			CloseMessageDialog()
			if index == 0 {
				onConfirm()
			} else {
				onCancel()
			}
		})
}

func ShowOkDialog(title string, message string, onClose func()) tview.Primitive {
	return ShowMessageDialog(title, message, []string{"OK"},
		func() {
			CloseMessageDialog()
			if onClose != nil {
				onClose()
			}
		},
		func(button string, index int) {
			CloseMessageDialog()
			if onClose != nil {
				onClose()
			}
		})
}
