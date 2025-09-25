package findbar

import (
	"github.com/gdamore/tcell/v2"
	"github.com/sedwards2009/femto"
	"github.com/sedwards2009/nuview"
)

type Findbar struct {
	*nuview.Flex
	app               *nuview.Application
	editor            *femto.View
	searchStringField *nuview.InputField
	OnClose           func()
}

func NewFindbar(app *nuview.Application, editor *femto.View) *Findbar {
	f := &Findbar{
		Flex:   nuview.NewFlex(),
		app:    app,
		editor: editor,
	}
	f.SetDirection(nuview.FlexRow)
	f.SetPadding(0, 0, 0, 0)
	f.SetBackgroundTransparent(false)
	f.SetBorder(false)

	hFlex := nuview.NewFlex()
	hFlex.SetDirection(nuview.FlexColumn)
	hFlex.SetBackgroundTransparent(false)
	hFlex.SetBorder(false)

	searchStringField := nuview.NewInputField()
	searchStringField.SetLabel("Find: ")
	searchStringField.SetLabelWidth(6)
	searchStringField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			if f.OnClose != nil {
				f.OnClose()
			}
			return nil
		case tcell.KeyEnter:
			f.SearchDown()
			return nil
		}
		return event
	})

	hFlex.AddItem(searchStringField, 0, 1, true)
	hFlex.AddItem(nil, 1, 0, false)

	searchUpButton := nuview.NewButton("↑") // U+2191 UPWARDS ARROW
	searchUpButton.SetSelectedFunc(f.SearchUp)
	hFlex.AddItem(searchUpButton, 3, 0, false)

	hFlex.AddItem(nil, 1, 0, false)

	searchDownButton := nuview.NewButton("↓") // U+2193 DOWNWARDS ARROW
	searchDownButton.SetSelectedFunc(f.SearchDown)
	hFlex.AddItem(searchDownButton, 3, 0, false)

	hFlex.AddItem(nil, 1, 0, false)

	closeButton := nuview.NewButton("✕")
	closeButton.SetSelectedFunc(func() {
		if f.OnClose != nil {
			f.OnClose()
		}
	})
	hFlex.AddItem(closeButton, 3, 0, false)

	f.AddItem(hFlex, 1, 0, true)

	f.searchStringField = searchStringField
	return f
}

func (f *Findbar) Focus(delegate func(p nuview.Primitive)) {
	delegate(f.searchStringField)
}

func (f *Findbar) SetSearchText(text string) {
	f.searchStringField.SetText(text)
}

func (f *Findbar) SearchUp() {
	if f.searchStringField.GetText() != "" {
		f.editor.Search(f.searchStringField.GetText(), false, false)
		f.editor.Relocate()
	}
}

func (f *Findbar) SearchDown() {
	if f.searchStringField.GetText() != "" {
		f.editor.Search(f.searchStringField.GetText(), false, true)
		f.editor.Relocate()
	}
}
