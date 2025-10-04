package findbar

import (
	"github.com/gdamore/tcell/v2"
	nuview "github.com/rivo/tview"
	"github.com/sedwards2009/femto"
)

type Findbar struct {
	*nuview.Flex
	app               *nuview.Application
	editor            *femto.View
	SearchStringField *nuview.InputField
	SearchUpButton    *nuview.Button
	SearchDownButton  *nuview.Button
	CloseButton       *nuview.Button
	OnClose           func()
}

func NewFindbar(app *nuview.Application, editor *femto.View) *Findbar {
	f := &Findbar{
		Flex:   nuview.NewFlex(),
		app:    app,
		editor: editor,
	}
	f.SetDirection(nuview.FlexRow)
	f.SetBorderPadding(0, 0, 0, 0)
	f.SetBorder(false)

	hFlex := nuview.NewFlex()
	hFlex.SetDirection(nuview.FlexColumn)
	// hFlex.SetBackgroundTransparent(false)
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
	f.SearchUpButton = searchUpButton
	hFlex.AddItem(searchUpButton, 3, 0, false)

	hFlex.AddItem(nil, 1, 0, false)

	searchDownButton := nuview.NewButton("↓") // U+2193 DOWNWARDS ARROW
	f.SearchDownButton = searchDownButton
	searchDownButton.SetSelectedFunc(f.SearchDown)
	hFlex.AddItem(searchDownButton, 3, 0, false)

	hFlex.AddItem(nil, 1, 0, false)

	closeButton := nuview.NewButton("✕")
	f.CloseButton = closeButton
	closeButton.SetSelectedFunc(func() {
		if f.OnClose != nil {
			f.OnClose()
		}
	})
	hFlex.AddItem(closeButton, 3, 0, false)

	f.AddItem(hFlex, 1, 0, true)

	f.SearchStringField = searchStringField
	return f
}

func (f *Findbar) Focus(delegate func(p nuview.Primitive)) {
	delegate(f.SearchStringField)
}

func (f *Findbar) SetSearchText(text string) {
	f.SearchStringField.SetText(text)
}

func (f *Findbar) SearchUp() {
	if f.SearchStringField.GetText() != "" {
		f.editor.Search(f.SearchStringField.GetText(), false, false)
		f.editor.Relocate()
	}
}

func (f *Findbar) SearchDown() {
	if f.SearchStringField.GetText() != "" {
		f.editor.Search(f.SearchStringField.GetText(), false, true)
		f.editor.Relocate()
	}
}
