package statusbar

import (
	"dinky/internal/tui/stylecolor"
	"dinky/internal/tui/utils"
	"fmt"
	"time"

	runewidth "github.com/mattn/go-runewidth"
	"github.com/rivo/tview"

	"github.com/gdamore/tcell/v2"
)

const (
	MESSAGE_TIMEOUT = 5 * time.Second
	ERROR_TIMEOUT   = 15 * time.Second
)

type StatusBar struct {
	*tview.Box
	app             *tview.Application
	Style           tcell.Style
	MessageStyle    tcell.Style
	ErrorStyle      tcell.Style
	Filename        string
	IsModified      bool
	Line            int
	Col             int
	LineEndings     string
	TabSize         int
	message         string
	errorMessage    string
	IsOverwriteMode bool
	UpdateHook      func(statusBar *StatusBar) // Hook for updating the status bar
}

func NewStatusBar(app *tview.Application) *StatusBar {
	fg := stylecolor.White
	bg := stylecolor.Blue

	return &StatusBar{
		Box:          tview.NewBox(),
		app:          app,
		Style:        tcell.StyleDefault.Foreground(fg).Background(bg),
		MessageStyle: tcell.StyleDefault.Foreground(fg).Background(stylecolor.Green),
		ErrorStyle:   tcell.StyleDefault.Foreground(fg).Background(stylecolor.Red),
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
	if statusBar.IsModified {
		leftMessage = "M " + leftMessage
	}
	if statusBar.message != "" {
		leftMessage = statusBar.message
		style = statusBar.MessageStyle
	}
	if statusBar.errorMessage != "" {
		leftMessage = statusBar.errorMessage
		style = statusBar.ErrorStyle
	}

	utils.DrawHorizontalLine(screen, x, y, width, style, ' ')
	utils.DrawText(screen, x, y, leftMessage, style)

	overwrite := "    "
	if statusBar.IsOverwriteMode {
		overwrite = "OVR "
	}
	cursorMessage := fmt.Sprintf("%sLn %d, Col %d", overwrite, statusBar.Line, statusBar.Col)
	padding := 15 - runewidth.StringWidth(cursorMessage)
	for i := 0; i < padding; i++ {
		cursorMessage += " "
	}
	rightMessage := fmt.Sprintf("%s  Tab Size: %d  %s  F12: Menu", cursorMessage, statusBar.TabSize, statusBar.LineEndings)
	rx := width - runewidth.StringWidth(rightMessage) - 1
	utils.DrawText(screen, rx, y, rightMessage, style)
}
