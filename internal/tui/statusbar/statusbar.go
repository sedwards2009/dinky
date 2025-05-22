package statusbar

import (
	"dinky/internal/tui/utils"
	"fmt"

	runewidth "github.com/mattn/go-runewidth"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type StatusBar struct {
	*tview.Box
	Style       tcell.Style
	Filename    string
	Line        int
	Col         int
	LineEndings string
	TabSize     int
}

func NewStatusBar() *StatusBar {
	fg := tcell.NewHexColor(0xf3f3f3)
	bg := tcell.NewHexColor(0x007ace)
	return &StatusBar{
		Box:   tview.NewBox(),
		Style: tcell.StyleDefault.Foreground(fg).Background(bg),
	}
}

func (statusBar *StatusBar) Draw(screen tcell.Screen) {
	statusBar.Box.DrawForSubclass(screen, statusBar)
	x, y, width, _ := statusBar.GetInnerRect()

	utils.DrawHorizontalLine(screen, x, y, width, statusBar.Style, ' ')

	utils.DrawText(screen, x+1, y, statusBar.Filename, statusBar.Style)

	cursorMessage := fmt.Sprintf("Ln %d, Col %d", statusBar.Line, statusBar.Col)
	padding := 19 - runewidth.StringWidth(cursorMessage)
	for i := 0; i < padding; i++ {
		cursorMessage += " "
	}
	rightMessage := fmt.Sprintf("%s  Tab Size: %d  %s  F12: Menu", cursorMessage, statusBar.TabSize, statusBar.LineEndings)
	rx := width - runewidth.StringWidth(rightMessage) - 1
	utils.DrawText(screen, rx, y, rightMessage, statusBar.Style)
}
