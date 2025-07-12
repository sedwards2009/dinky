package statusbar

import (
	"dinky/internal/tui/utils"
	"fmt"
	"time"

	runewidth "github.com/mattn/go-runewidth"

	"github.com/gdamore/tcell/v2"
	"github.com/sedwards2009/nuview"
)

const (
	MESSAGE_TIMEOUT = 5 * time.Second
	ERROR_TIMEOUT   = 15 * time.Second
)

type StatusBar struct {
	*nuview.Box
	app          *nuview.Application
	Style        tcell.Style
	MessageStyle tcell.Style
	ErrorStyle   tcell.Style
	Filename     string
	IsModified   bool
	Line         int
	Col          int
	LineEndings  string
	TabSize      int
	message      string
	errorMessage string
	UpdateHook   func(statusBar *StatusBar) // Hook for updating the status bar
}

func NewStatusBar(app *nuview.Application) *StatusBar {
	fg := tcell.NewHexColor(0xf3f3f3)
	bg := tcell.NewHexColor(0x007ace)
	messageBg := tcell.NewHexColor(0x0b835c)
	errorBg := tcell.NewHexColor(0xa4090c)
	return &StatusBar{
		Box:          nuview.NewBox(),
		app:          app,
		Style:        tcell.StyleDefault.Foreground(fg).Background(bg),
		MessageStyle: tcell.StyleDefault.Foreground(fg).Background(messageBg),
		ErrorStyle:   tcell.StyleDefault.Foreground(fg).Background(errorBg),
	}
}

func (statusBar *StatusBar) ShowMessage(message string) {
	statusBar.message = message
	statusBar.scheduleMessageReset(MESSAGE_TIMEOUT)
}

func (statusBar *StatusBar) ShowError(errorMessage string) {
	statusBar.errorMessage = errorMessage
	statusBar.scheduleMessageReset(ERROR_TIMEOUT)
}

func (statusBar *StatusBar) scheduleMessageReset(timeOut time.Duration) {
	time.AfterFunc(timeOut, func() {
		statusBar.app.QueueUpdateDraw(func() {
			statusBar.message = ""
		})
	})
}

func (statusBar *StatusBar) Draw(screen tcell.Screen) {
	if statusBar.UpdateHook != nil {
		statusBar.UpdateHook(statusBar)
	}

	x, y, width, _ := statusBar.GetInnerRect()

	style := statusBar.Style
	leftMessage := statusBar.Filename
	if statusBar.message != "" {
		leftMessage = statusBar.message
		style = statusBar.MessageStyle
	}
	if statusBar.errorMessage != "" {
		leftMessage = statusBar.errorMessage
		style = statusBar.ErrorStyle
	}

	utils.DrawHorizontalLine(screen, x, y, width, style, ' ')

	if statusBar.IsModified {
		utils.DrawText(screen, x, y, "M", style)
	}

	utils.DrawText(screen, x+2, y, leftMessage, style)

	cursorMessage := fmt.Sprintf("Ln %d, Col %d", statusBar.Line, statusBar.Col)
	padding := 19 - runewidth.StringWidth(cursorMessage)
	for i := 0; i < padding; i++ {
		cursorMessage += " "
	}
	rightMessage := fmt.Sprintf("%s  Tab Size: %d  %s  F12: Menu", cursorMessage, statusBar.TabSize, statusBar.LineEndings)
	rx := width - runewidth.StringWidth(rightMessage) - 1
	utils.DrawText(screen, rx, y, rightMessage, style)
}
