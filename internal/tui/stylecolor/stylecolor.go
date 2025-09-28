package stylecolor

import (
	"github.com/gdamore/tcell/v2"
	nuview "github.com/rivo/tview"
)

var White = tcell.NewHexColor(0xffffff).TrueColor() // White foreground
// var Blue = tcell.NewHexColor(0x6677aa).TrueColor()  // Blue background
var Blue = tcell.NewHexColor(0x007ace).TrueColor()
var Black = tcell.NewHexColor(0x000000).TrueColor()
var LightGray = tcell.NewHexColor(0xaaaaaa).TrueColor()

var DarkGray = tcell.NewHexColor(0x333333).TrueColor() // Dark gray background

var blackOnGrayStyle = tcell.StyleDefault.Foreground(Black).Background(LightGray)
var whiteOnBlueStyle = tcell.StyleDefault.Foreground(White).Background(Blue)

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

func init() {
	nuview.Styles.TitleColor = Black
	nuview.Styles.BorderColor = Black

	nuview.Styles.PrimitiveBackgroundColor = LightGray
	nuview.Styles.PrimaryTextColor = Black
	nuview.Styles.SecondaryTextColor = Black

	ButtonLabelColor = Black
	ButtonLabelFocusedColor = Black
	ButtonBackgroundColor = White
	ButtonBackgroundFocusedColor = Blue
	ButtonBackgroundDisabledColor = White
	ButtonLabelDisabledColor = LightGray

	CheckboxLabelStyle = blackOnGrayStyle
	CheckboxUncheckedStyle = blackOnGrayStyle
	CheckboxCheckedStyle = blackOnGrayStyle
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
}
