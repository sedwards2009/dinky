package femtoinputfield

import (
	"github.com/sedwards2009/femto"
	"github.com/sedwards2009/femto/runtime"
)

type FemtoInputField struct {
	*femto.View
}

func NewFemtoInputField() *FemtoInputField {
	buffer := femto.NewBufferFromString("", "")
	editor := femto.NewView(buffer)
	editor.SetRuntimeFiles(runtime.Files)
	buffer.Settings["ruler"] = false
	// editor.SetKeybindings(femtoDefaultKeyBindings)

	return &FemtoInputField{
		View: editor,
	}
}

func (f *FemtoInputField) SetKeybindings(keybindings femto.KeyBindings) {
	f.View.SetKeybindings(keybindings)
}

func (f *FemtoInputField) GetText() string {
	return f.View.Buf.Line(0)
}

func (f *FemtoInputField) SetText(text string) {
	f.View.DeleteLine()
	f.View.Buf.Insert(f.View.Buf.Start(), text)
}
