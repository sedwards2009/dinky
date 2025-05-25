package style

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var white = tcell.NewHexColor(0xf3f3f3) // White foreground
var blue = tcell.NewHexColor(0x007ace)  // Blue background
var black = tcell.NewHexColor(0x1e1e1e)

func Init() {

	tview.Styles.PrimitiveBackgroundColor = blue
	tview.Styles.PrimaryTextColor = white
	tview.Styles.SecondaryTextColor = white
}

var TextAreaBackgroundColor = black

func StyleInputField(inputField *tview.InputField) {
	// inputField.SetBackgroundColor(TextAreaBackgroundColor)
	// inputField.SetLabelColor(tview.Styles.PrimaryTextColor)
	inputField.SetFieldBackgroundColor(TextAreaBackgroundColor)
	// inputField.SetTextColor(tview.Styles.PrimaryTextColor)
}

func StyleList(fileList *tview.List) {
	fileList.SetBackgroundColor(TextAreaBackgroundColor)
	// fileList.SetMainTextColor(tview.Styles.PrimaryTextColor)
	fileList.SetMainTextStyle(tcell.StyleDefault.Foreground(white).Background(black))
	fileList.SetSelectedStyle(tcell.StyleDefault.Foreground(black).Background(white))
}
