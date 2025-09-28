package style

import (
	"dinky/internal/tui/filelist"
	"dinky/internal/tui/menu"
	"dinky/internal/tui/scrollbar"
	"dinky/internal/tui/stylecolor"
	"dinky/internal/tui/tabbar"
	"dinky/internal/tui/table2"

	"github.com/gdamore/tcell/v2"
)

func StyleScrollbarTrack(scrollbarTrack *scrollbar.ScrollbarTrack) {
	scrollbarTrack.SetTrackColor(stylecolor.DarkGray)
	scrollbarTrack.SetThumbColor(stylecolor.White)
	scrollbarTrack.SetWidth(1)
}

func StyleScrollbar(scrollbar *scrollbar.Scrollbar) {
	StyleScrollbarTrack(scrollbar.Track)
}

func StyleFileList(fileList *filelist.FileList) {
	fileList.SetTextColor(stylecolor.White)
	fileList.SetBackgroundColor(stylecolor.Black)
	fileList.SetSelectedBackgroundColor(stylecolor.Blue)
	fileList.SetHeaderLabelColor(stylecolor.Black)
	fileList.SetHeaderBackgroundColor(stylecolor.White)

	StyleScrollbar(fileList.VerticalScrollbar)
	StyleScrollbar(fileList.HorizontalScrollbar)
}

func StyleTabBar(tabBar *tabbar.TabBar) {
	tabBar.ActiveTabStyle = tcell.StyleDefault.Foreground(stylecolor.White).Background(stylecolor.Black).Bold(true)
	tabBar.InactiveTabStyle = tcell.StyleDefault.Foreground(stylecolor.LightGray).Background(stylecolor.DarkGray).Bold(false)
	tabBar.BackgroundStyle = tcell.StyleDefault.Foreground(stylecolor.White).Background(stylecolor.Blue)
}

func StyleMenuBar(menuBar *menu.MenuBar) {
	menuBar.MenuBarStyle = tcell.StyleDefault.Foreground(stylecolor.White).Background(stylecolor.Blue)
	menuBar.MenuStyle = tcell.StyleDefault.Foreground(stylecolor.Black).Background(stylecolor.LightGray)
	menuBar.MenuSelectedStyle = tcell.StyleDefault.Foreground(stylecolor.White).Background(stylecolor.Blue)
}

func StyleTable(table *table2.Table) {
	table.SetBackgroundColor(stylecolor.Black)
	table.SetSelectedStyle(tcell.StyleDefault.Background(stylecolor.Blue).Foreground(stylecolor.White))
}

func StyleTableCell(cell *table2.TableCell) {
	cell.SetStyle(tcell.StyleDefault.Foreground(stylecolor.White).Background(stylecolor.Black))
}
