package femtoinputfield

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sedwards2009/smidgen"
)

type FemtoInputField struct {
	*smidgen.View
	done    func(tcell.Key)
	changed func(text string)
}

func NewSmidgenInputField(app *tview.Application) *FemtoInputField {
	buffer := smidgen.NewBufferFromString("", "")
	editor := smidgen.NewView(app, buffer)
	buffer.Settings["ruler"] = false
	buffer.Settings["hidecursoronblur"] = true

	return &FemtoInputField{
		View: editor,
	}
}

func (f *FemtoInputField) SetKeybindings(keybindings smidgen.Keybindings) {
	f.View.SetKeybindings(keybindings)
}

func (f *FemtoInputField) SetTextColor(foreground tcell.Color, background tcell.Color) {
	scheme := make(smidgen.Colorscheme)
	scheme["default"] = tcell.StyleDefault.Foreground(foreground).Background(background)
	f.View.SetColorscheme(scheme)
}

func (f *FemtoInputField) GetText() string {
	return f.View.Buffer().Line(0)
}

func (f *FemtoInputField) SetText(text string) {
	f.View.ActionController().DeleteLine()
	f.View.Buffer().Insert(f.View.Buffer().Start(), text)
}

func (f *FemtoInputField) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
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

func (f *FemtoInputField) SetDoneFunc(handler func(key tcell.Key)) *FemtoInputField {
	f.done = handler
	return f
}

func (f *FemtoInputField) SetChangedFunc(handler func(text string)) *FemtoInputField {
	f.changed = handler
	return f
}
