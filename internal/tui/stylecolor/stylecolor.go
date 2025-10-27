package stylecolor

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var White = tcell.NewHexColor(0xf3f3f3).TrueColor() // White foreground
var Blue = tcell.NewHexColor(0x007ace).TrueColor()
var Black = tcell.NewHexColor(0x000000).TrueColor()
var LightGray = tcell.NewHexColor(0xaaaaaa).TrueColor()
var Green = tcell.NewHexColor(0x0b835c).TrueColor()    // Green foreground
var DarkGray = tcell.NewHexColor(0x333333).TrueColor() // Dark gray background
var Red = tcell.NewHexColor(0xa4090c).TrueColor()      // Red foreground

var blackOnGrayStyle = tcell.StyleDefault.Foreground(Black).Background(LightGray)
var whiteOnBlueStyle = tcell.StyleDefault.Foreground(White).Background(Blue)
var whiteOnGreenStyle = tcell.StyleDefault.Foreground(White).Background(Green)

// var blackOnGreenStyle = tcell.StyleDefault.Foreground(Black).Background(Green)

var ButtonLabelColor tcell.Color
var ButtonLabelFocusedColor tcell.Color
var ButtonBackgroundColor tcell.Color
var ButtonBackgroundFocusedColor tcell.Color
var ButtonBackgroundDisabledColor tcell.Color
var ButtonLabelDisabledColor tcell.Color

var CheckboxLabelStyle tcell.Style
var CheckboxUncheckedStyle tcell.Style
var CheckboxCheckedStyle tcell.Style
var CheckboxFocusStyle tcell.Style
var CheckboxCheckedString string
var CheckboxUncheckedString string
var CheckboxCursorCheckedString string
var CheckboxCursorUncheckedString string

var InputFieldLabelColor tcell.Color
var InputFieldFieldBackgroundColor tcell.Color
var InputFieldFieldBackgroundFocusedColor tcell.Color
var InputFieldFieldTextColor tcell.Color
var InputFieldFieldTextFocusedColor tcell.Color
var InputFieldPlaceholderTextColor tcell.Color
var ListMainTextColor tcell.Color
var ListSecondaryTextColor tcell.Color
var ListShortcutColor tcell.Color
var ListSelectedTextColor tcell.Color
var ListSelectedBackgroundColor tcell.Color

var DropDownTextColor tcell.Color
var DropDownBackgroundColor tcell.Color
var DropDownSelectedTextColor tcell.Color
var DropDownSelectedBackgroundColor tcell.Color

func init() {
	tview.Styles.TitleColor = Black
	tview.Styles.BorderColor = Black

	tview.Styles.PrimitiveBackgroundColor = LightGray
	tview.Styles.PrimaryTextColor = Black
	tview.Styles.SecondaryTextColor = Black

	ButtonLabelColor = Black
	ButtonLabelFocusedColor = Black
	ButtonBackgroundColor = White
	ButtonBackgroundFocusedColor = Blue
	ButtonBackgroundDisabledColor = White
	ButtonLabelDisabledColor = LightGray

	CheckboxLabelStyle = blackOnGrayStyle
	CheckboxUncheckedStyle = blackOnGrayStyle
	CheckboxCheckedStyle = whiteOnGreenStyle
	CheckboxFocusStyle = whiteOnBlueStyle
	CheckboxCheckedString = "[✓]"
	CheckboxUncheckedString = "[ ]"
	CheckboxCursorCheckedString = "[✓]"
	CheckboxCursorUncheckedString = "[ ]"

	// nuview.Styles.MoreContrastBackgroundColor = tcell.ColorDarkGray
	// nuview.Styles.ContrastBackgroundColor = black
	// nuview.Styles.PrimaryTextColor = tcell.ColorLightGray
	// nuview.Styles.PrimaryTextColor = white

	InputFieldLabelColor = Black
	InputFieldFieldBackgroundColor = DarkGray
	InputFieldFieldBackgroundFocusedColor = Black
	InputFieldFieldTextColor = LightGray
	InputFieldFieldTextFocusedColor = White
	InputFieldPlaceholderTextColor = LightGray
	ListMainTextColor = White
	ListSecondaryTextColor = LightGray
	ListShortcutColor = LightGray
	ListSelectedTextColor = White
	ListSelectedBackgroundColor = Blue

	DropDownTextColor = White
	DropDownBackgroundColor = Black
	DropDownSelectedTextColor = White
	DropDownSelectedBackgroundColor = Blue
}
