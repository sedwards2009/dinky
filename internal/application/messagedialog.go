package application

import (
	"dinky/internal/tui/dialog"
	"strings"
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
	OnButtonClick func(button string, index int)) {

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
	modalPages.AddPanel(messageDialogName, messageDialog, true, true)
	messageDialog.OnClose = OnClose
	messageDialog.OnButtonClick = OnButtonClick
	messageDialog.Open(title, message, buttons, width, height)
	messageDialog.FocusButton(0)
}

func CloseMessageDialog() {
	if messageDialog != nil {
		messageDialog.Close()
		modalPages.RemovePanel("messagedialog")
	}
}

func ShowConfirmDialog(message string, onConfirm func(), onCancel func()) {
	ShowMessageDialog("Confirm", message, []string{"OK", "Cancel"},
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

func ShowOkDialog(title string, message string, onClose func()) {
	ShowMessageDialog(title, message, []string{"OK"},
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
