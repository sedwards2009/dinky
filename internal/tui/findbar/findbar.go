package findbar

import (
	"dinky/internal/tui/femtoinputfield"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sedwards2009/femto"
)

type Findbar struct {
	*tview.Flex
	app               *tview.Application
	editor            *femto.View
	SearchStringField *femtoinputfield.FemtoInputField
	SearchUpButton    *tview.Button
	SearchDownButton  *tview.Button
	CloseButton       *tview.Button
	OnClose           func()
}

func NewFindbar(app *tview.Application, editor *femto.View) *Findbar {
	f := &Findbar{
		Flex:   tview.NewFlex(),
		app:    app,
		editor: editor,
	}
	f.SetDirection(tview.FlexRow)
	f.SetBorderPadding(0, 0, 0, 0)
	f.SetBorder(false)

	hFlex := tview.NewFlex()
	hFlex.SetDirection(tview.FlexColumn)
	hFlex.SetBorder(false)

	searchStringField := femtoinputfield.NewFemtoInputField()
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

	searchFieldLabel := tview.NewTextView()
	searchFieldLabel.SetText("Find: ")
	hFlex.AddItem(searchFieldLabel, 6, 0, false)

	hFlex.AddItem(searchStringField, 0, 1, true)
	hFlex.AddItem(nil, 1, 0, false)

	searchUpButton := tview.NewButton("↑") // U+2191 UPWARDS ARROW
	searchUpButton.SetSelectedFunc(f.SearchUp)
	f.SearchUpButton = searchUpButton
	hFlex.AddItem(searchUpButton, 3, 0, false)

	hFlex.AddItem(nil, 1, 0, false)

	searchDownButton := tview.NewButton("↓") // U+2193 DOWNWARDS ARROW
	f.SearchDownButton = searchDownButton
	searchDownButton.SetSelectedFunc(f.SearchDown)
	hFlex.AddItem(searchDownButton, 3, 0, false)

	hFlex.AddItem(nil, 1, 0, false)

	closeButton := tview.NewButton("✕")
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

func (f *Findbar) Focus(delegate func(p tview.Primitive)) {
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

func (f *Findbar) SetFemtoKeybindings(keybindings femto.KeyBindings) {
	f.SearchStringField.SetKeybindings(keybindings)
}
