package application

import (
	"dinky/internal/tui/dialog"
	"dinky/internal/tui/style"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sedwards2009/smidgen/micro/buffer"
)

var inputDialog *dialog.InputDialog

const inputDialogName = "inputdialog"

func ShowGoToLineDialog(title string, message string, defaultValue string, onCancel func(), onAccept func(value string,
	index int)) tview.Primitive {

	if inputDialog == nil {
		inputDialog = dialog.NewInputDialog(app)
		inputDialog.SetFemtoKeybindings(femtoSingleLineKeyBindings)
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
		Title:        title,
		Message:      message,
		DefaultValue: defaultValue,
		Buttons:      []string{"OK", "Cancel"},
		Width:        width,
		Height:       height,
		OnCancel:     onCancel,
		OnAccept:     onAccept,
		FieldKeyFilter: func(event *tcell.EventKey) bool {
			key := event.Key()
			// Allow digits, colon for line:column format, and basic editing keys
			if key == tcell.KeyBackspace || key == tcell.KeyDelete ||
				key == tcell.KeyLeft || key == tcell.KeyRight ||
				key == tcell.KeyHome || key == tcell.KeyEnd ||
				key == tcell.KeyDEL {
				return true
			}
			if event.Rune() >= '0' && event.Rune() <= '9' {
				return true
			}
			if event.Rune() == ':' {
				return true
			}
			return false
		},
	}

	inputDialog.Open(options)
	style.StyleInputDialog(inputDialog)
	return inputDialog
}

func CloseGoToLineDialog() {
	if inputDialog != nil {
		inputDialog.Close()
		modalPages.RemovePage(inputDialogName)
	}
}

func handleGoToLine() tview.Primitive {
	return ShowGoToLineDialog("Go to Line", "Enter line number (or line:column):", "",
		func() {
			// On cancel
			CloseGoToLineDialog()
		},
		func(value string, index int) {
			CloseGoToLineDialog()
			if index == 0 || index == -1 { // OK button or Enter key in input field
				parseAndGoToLine(value)
			}
		})
}

func parseAndGoToLine(input string) {
	if input == "" {
		return
	}

	var lineNum, colNum int
	var err error

	// Check if input contains a colon (line:column format)
	colonIndex := -1
	for i, r := range input {
		if r == ':' {
			colonIndex = i
			break
		}
	}

	if colonIndex != -1 {
		// Parse line:column format
		lineStr := input[:colonIndex]
		colStr := input[colonIndex+1:]

		lineNum, err = strconv.Atoi(lineStr)
		if err != nil {
			statusBar.ShowError("Invalid line number")
			return
		}

		colNum, err = strconv.Atoi(colStr)
		if err != nil {
			statusBar.ShowError("Invalid column number")
			return
		}
	} else {
		// Parse line number only
		lineNum, err = strconv.Atoi(input)
		if err != nil {
			statusBar.ShowError("Invalid line number")
			return
		}
		colNum = 1
	}

	// Convert to 0-based indexing
	lineNum--
	colNum--

	// Validate line number
	if lineNum < 0 || lineNum >= currentFileBuffer.buffer.LinesNum() {
		statusBar.ShowError("Line number out of range")
		return
	}

	// Validate column number
	lineLength := len([]rune(currentFileBuffer.buffer.Line(lineNum)))
	if colNum < 0 {
		colNum = 0
	} else if colNum > lineLength {
		colNum = lineLength
	}

	// Move cursor to the specified location
	currentFileBuffer.editor.GoToLoc(buffer.Loc{X: colNum, Y: lineNum})

	statusBar.ShowMessage("Jumped to line " + strconv.Itoa(lineNum+1))
}
