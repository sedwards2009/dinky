package smidgeninputfield

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sedwards2009/smidgen"
)

type SmidgenInputField struct {
	*smidgen.View
	userText       string
	history        []string
	historyPointer int

	done    func(tcell.Key)
	changed func(text string)
}

func NewSmidgenInputField(app *tview.Application) *SmidgenInputField {
	buffer := smidgen.NewBufferFromString("", "")
	editor := smidgen.NewView(app, buffer)
	buffer.Settings["ruler"] = false
	buffer.Settings["hidecursoronblur"] = true

	return &SmidgenInputField{
		View:           editor,
		historyPointer: -1,
	}
}

func (f *SmidgenInputField) SetHistory(historyText []string) {
	f.history = historyText
	f.historyPointer = len(historyText)
}

func (f *SmidgenInputField) SetKeybindings(keybindings smidgen.Keybindings) {
	f.View.SetKeybindings(keybindings)
}

func (f *SmidgenInputField) SetTextColor(foreground tcell.Color, background tcell.Color) {
	scheme := make(smidgen.Colorscheme)
	scheme["default"] = tcell.StyleDefault.Foreground(foreground).Background(background)
	f.View.SetColorscheme(scheme)
}

func (f *SmidgenInputField) GetText() string {
	return f.View.Buffer().Line(0)
}

func (f *SmidgenInputField) SetText(text string) {
	f.userText = text
	f.internalSetText(text)
}

func (f *SmidgenInputField) internalSetText(text string) {
	f.View.ActionController().DeleteLine()
	f.View.Buffer().Insert(f.View.Buffer().Start(), text)
	f.View.ActionController().StartOfLine()
	f.View.Relocate()
}

func (f *SmidgenInputField) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return f.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		switch event.Key() {
		case tcell.KeyEnter, tcell.KeyEscape, tcell.KeyTab, tcell.KeyBacktab:
			if f.done != nil {
				f.done(event.Key())
			}
			return
		case tcell.KeyUp:
			if len(f.history) != 0 && f.historyPointer > 0 {
				f.historyPointer--
				f.internalSetText(f.history[f.historyPointer])
				if f.changed != nil {
					f.changed(f.GetText())
				}
				return
			}
		case tcell.KeyDown:
			if len(f.history) != 0 {
				if f.historyPointer < len(f.history)-1 {
					f.historyPointer++
					f.internalSetText(f.history[f.historyPointer])
				} else {
					f.internalSetText(f.userText)
					f.historyPointer = len(f.history)
					// Set to one past the end of history, so pressing up goes to the most recent entry
				}
				if f.changed != nil {
					f.changed(f.GetText())
				}
				return
			}
		}

		f.View.Buffer().ClearModified()
		f.View.InputHandler()(event, setFocus)
		if f.View.Buffer().Modified() {
			f.View.Buffer().ClearModified()
			f.userText = f.GetText()
			if len(f.history) != 0 {
				f.historyPointer = len(f.history)
				// Reset history pointer to end of history, so pressing up goes to the most recent entry
			}
			if f.changed != nil {
				f.changed(f.GetText())
			}
		}
	})
}

func (f *SmidgenInputField) SetDoneFunc(handler func(key tcell.Key)) *SmidgenInputField {
	f.done = handler
	return f
}

func (f *SmidgenInputField) SetChangedFunc(handler func(text string)) *SmidgenInputField {
	f.changed = handler
	return f
}
