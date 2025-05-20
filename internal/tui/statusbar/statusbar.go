package statusbar

import (
	"dinky/internal/tui/utils"

	runewidth "github.com/mattn/go-runewidth"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type StatusBar struct {
	*tview.Box
	Style    tcell.Style
	Filename string
}

func NewStatusBar() *StatusBar {
	fg := tcell.NewHexColor(0xf3f3f3)
	bg := tcell.NewHexColor(0x007ace)
	return &StatusBar{
		Box:   tview.NewBox(),
		Style: tcell.StyleDefault.Foreground(fg).Background(bg).Bold(true),
	}
}

func (statusBar *StatusBar) Draw(screen tcell.Screen) {
	statusBar.Box.DrawForSubclass(screen, statusBar)
	x, y, width, _ := statusBar.GetInnerRect()

	utils.DrawHorizontalLine(screen, x, y, width, statusBar.Style, ' ')

	utils.DrawText(screen, x+1, y, statusBar.Filename, statusBar.Style)

	rightMessage := "F12: Menu"
	rx := width - runewidth.StringWidth(rightMessage)
	utils.DrawText(screen, rx, y, rightMessage, statusBar.Style)
}
