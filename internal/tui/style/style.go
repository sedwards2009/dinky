package style

import (
	"dinky/internal/tui/filelist"
	"dinky/internal/tui/menu"
	"dinky/internal/tui/scrollbar"
	"dinky/internal/tui/tabbar"

	"github.com/gdamore/tcell/v2"
	"github.com/sedwards2009/nuview"
)

var white = tcell.NewHexColor(0xffffff).TrueColor() // White foreground
// var blue = tcell.NewHexColor(0x6677aa).TrueColor()  // Blue background
var blue = tcell.NewHexColor(0x007ace).TrueColor()
var black = tcell.NewHexColor(0x000000).TrueColor()
var lightGray = tcell.NewHexColor(0xaaaaaa).TrueColor()

var darkGray = tcell.NewHexColor(0x333333).TrueColor() // Dark gray background

var blackOnGrayStyle = tcell.StyleDefault.Foreground(black).Background(lightGray)
var whiteOnBlueStyle = tcell.StyleDefault.Foreground(white).Background(blue)

func Install() {
	nuview.Styles.ButtonCursorRune = 0
	nuview.Styles.TitleColor = black
	nuview.Styles.BorderColor = black

	nuview.Styles.PrimitiveBackgroundColor = lightGray
	nuview.Styles.PrimaryTextColor = black
	nuview.Styles.SecondaryTextColor = black

	nuview.Styles.ButtonLabelColor = black
	nuview.Styles.ButtonLabelFocusedColor = black
	nuview.Styles.ButtonBackgroundColor = white
	nuview.Styles.ButtonBackgroundFocusedColor = blue
	nuview.Styles.ButtonBackgroundDisabledColor = white
	nuview.Styles.ButtonLabelDisabledColor = lightGray

	nuview.Styles.CheckboxLabelStyle = blackOnGrayStyle
	nuview.Styles.CheckboxUncheckedStyle = blackOnGrayStyle
	nuview.Styles.CheckboxCheckedStyle = blackOnGrayStyle
	nuview.Styles.CheckboxFocusStyle = whiteOnBlueStyle
	nuview.Styles.CheckboxCheckedString = "[✓]"
	nuview.Styles.CheckboxUncheckedString = "[ ]"
	nuview.Styles.CheckboxCursorCheckedString = "[✓]"
	nuview.Styles.CheckboxCursorUncheckedString = "[ ]"

	// nuview.Styles.MoreContrastBackgroundColor = tcell.ColorDarkGray
	// nuview.Styles.ContrastBackgroundColor = black
	// nuview.Styles.PrimaryTextColor = tcell.ColorLightGray
	// nuview.Styles.PrimaryTextColor = white

	nuview.Styles.InputFieldLabelColor = black
	nuview.Styles.InputFieldFieldBackgroundColor = darkGray
	nuview.Styles.InputFieldFieldBackgroundFocusedColor = black
	nuview.Styles.InputFieldFieldTextColor = lightGray
	nuview.Styles.InputFieldFieldTextFocusedColor = white
	nuview.Styles.InputFieldPlaceholderTextColor = lightGray

	nuview.Styles.ListMainTextColor = white
	nuview.Styles.ListSecondaryTextColor = lightGray
	nuview.Styles.ListShortcutColor = lightGray
	nuview.Styles.ListSelectedTextColor = white
	nuview.Styles.ListSelectedBackgroundColor = blue

}

func StyleScrollbarTrack(scrollbarTrack *scrollbar.ScrollbarTrack) {
	scrollbarTrack.SetTrackColor(darkGray)
	scrollbarTrack.SetThumbColor(white)
	scrollbarTrack.SetWidth(1)
}

func StyleScrollbar(scrollbar *scrollbar.Scrollbar) {
	StyleScrollbarTrack(scrollbar.Track)
}

func StyleFileList(fileList *filelist.FileList) {
	fileList.SetTextColor(white)
	fileList.SetBackgroundColor(black)
	fileList.SetSelectedBackgroundColor(blue)
	fileList.SetHeaderLabelColor(black)
	fileList.SetHeaderBackgroundColor(white)

	StyleScrollbar(fileList.VerticalScrollbar)
	StyleScrollbar(fileList.HorizontalScrollbar)
}

func StyleTabBar(tabBar *tabbar.TabBar) {
	tabBar.ActiveTabStyle = tcell.StyleDefault.Foreground(white).Background(black).Bold(true)
	tabBar.InactiveTabStyle = tcell.StyleDefault.Foreground(lightGray).Background(darkGray).Bold(false)
	tabBar.BackgroundStyle = tcell.StyleDefault.Foreground(white).Background(blue)
}

func StyleMenuBar(menuBar *menu.MenuBar) {
	menuBar.MenuBarStyle = tcell.StyleDefault.Foreground(white).Background(blue)
	menuBar.MenuStyle = tcell.StyleDefault.Foreground(black).Background(lightGray)
	menuBar.MenuSelectedStyle = tcell.StyleDefault.Foreground(white).Background(blue)
}

func StyleTable(table *nuview.Table) {
	table.SetBackgroundColor(black)
	table.SetSelectedStyle(tcell.StyleDefault.Background(blue).Foreground(white))
}

func StyleTableCell(cell *nuview.TableCell) {
	cell.SetStyle(tcell.StyleDefault.Foreground(white).Background(black))
}
