package smidgeninputfield

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sedwards2009/smidgen"
)

type SmidgenInputField struct {
	*smidgen.View
	done    func(tcell.Key)
	changed func(text string)
}

func NewSmidgenInputField(app *tview.Application) *SmidgenInputField {
	buffer := smidgen.NewBufferFromString("", "")
	editor := smidgen.NewView(app, buffer)
	buffer.Settings["ruler"] = false
	buffer.Settings["hidecursoronblur"] = true

	return &SmidgenInputField{
		View: editor,
	}
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
		}

		f.View.Buffer().ClearModified()
		f.View.InputHandler()(event, setFocus)
		if f.View.Buffer().Modified() {
			f.View.Buffer().ClearModified()
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
