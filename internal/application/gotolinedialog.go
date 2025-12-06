package application

import (
	"strconv"

	"github.com/rivo/tview"
	"github.com/sedwards2009/smidgen/micro/buffer"
)

func handleGoToLine() tview.Primitive {
	return ShowInputDialog("Go to Line", "Enter line number (or line:column):", "",
		func() {
			// On cancel
			CloseInputDialog()
		},
		func(value string, index int) {
			CloseInputDialog()
			if index == 0 || index == -1 { // OK button or Enter key in input field
				parseAndGoToLine(value)
			}
		}, numericInputFilter)
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
