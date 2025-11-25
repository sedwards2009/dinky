package findbar

import (
	"dinky/internal/tui/femtoinputfield"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sedwards2009/smidgen"
	"github.com/sedwards2009/smidgen/micro/buffer"
)

type Findbar struct {
	*tview.Flex
	app                   *tview.Application
	editor                *smidgen.View
	SearchStringField     *femtoinputfield.FemtoInputField
	SearchUpButton        *tview.Button
	SearchDownButton      *tview.Button
	CloseButton           *tview.Button
	ReplaceButton         *tview.Button
	ReplaceAllButton      *tview.Button
	ExpanderCheckbox      *tview.Checkbox
	RegexCheckbox         *tview.Checkbox
	CaseSensitiveCheckbox *tview.Checkbox
	OnClose               func()
	OnError               func(err error)
	OnExpand              func(expanded bool)
	OnMessage             func(message string)
	isExpanded            bool

	hFlex2 *tview.Flex
}

func NewFindbar(app *tview.Application, editor *smidgen.View) *Findbar {
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

	expanderCheckbox := tview.NewCheckbox()
	hFlex.AddItem(expanderCheckbox, 3, 0, false)
	expanderCheckbox.SetChangedFunc(f.handleExpandClick)
	f.ExpanderCheckbox = expanderCheckbox

	searchFieldLabel := tview.NewTextView()
	searchFieldLabel.SetText(" Find: ")
	hFlex.AddItem(searchFieldLabel, 7, 0, false)

	searchStringField := femtoinputfield.NewSmidgenInputField(app)
	searchStringField.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEscape:
			if f.OnClose != nil {
				f.OnClose()
			}
		case tcell.KeyEnter:
			f.SearchDown()
		}
	})
	hFlex.AddItem(searchStringField, 0, 1, true)
	f.SearchStringField = searchStringField

	hFlex.AddItem(nil, 1, 0, false)

	// Case Sensitive Checkbox [✓Aa ]
	caseSensitiveCheckbox := tview.NewCheckbox()
	caseSensitiveCheckbox.SetChecked(false)
	hFlex.AddItem(caseSensitiveCheckbox, 7, 0, false)
	f.CaseSensitiveCheckbox = caseSensitiveCheckbox

	// Regex Checkbox [✓Regex ]
	regexCheckbox := tview.NewCheckbox()
	regexCheckbox.SetChecked(false)
	hFlex.AddItem(regexCheckbox, 10, 0, false)
	f.RegexCheckbox = regexCheckbox

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

	hFlex2 := tview.NewFlex()
	f.hFlex2 = hFlex2
	hFlex2.SetDirection(tview.FlexColumn)
	hFlex2.SetBorder(false)

	replaceFieldLabel := tview.NewTextView()
	replaceFieldLabel.SetText(" Replace: ")
	hFlex2.AddItem(replaceFieldLabel, 10, 0, false)

	replaceStringField := femtoinputfield.NewSmidgenInputField(app)
	replaceStringField.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEscape:
			if f.OnClose != nil {
				f.OnClose()
			}
		case tcell.KeyEnter:
			f.SearchDown()
		}
	})
	hFlex2.AddItem(replaceStringField, 0, 1, false)

	hFlex2.AddItem(nil, 1, 0, false)

	replaceButton := tview.NewButton("Replace")
	f.ReplaceButton = replaceButton
	replaceButton.SetSelectedFunc(f.Replace)
	hFlex2.AddItem(replaceButton, 9, 0, false)

	hFlex2.AddItem(nil, 1, 0, false)

	replaceAllButton := tview.NewButton("Replace All")
	f.ReplaceAllButton = replaceAllButton
	replaceAllButton.SetSelectedFunc(f.ReplaceAll)
	hFlex2.AddItem(replaceAllButton, 13, 0, false)

	hFlex2.AddItem(nil, 5, 0, false)

	return f
}

func (f *Findbar) handleExpandClick(checked bool) {
	if checked {
		if f.OnExpand != nil {
			f.OnExpand(checked)
		}
		f.AddItem(f.hFlex2, 1, 0, false)
	} else {
		f.RemoveItem(f.hFlex2)
		if f.OnExpand != nil {
			f.OnExpand(checked)
		}
	}
	f.isExpanded = checked
}

func (f *Findbar) Expand() {
	f.ExpanderCheckbox.SetChecked(true) // This also triggers the handler on this checkbox
}

func (f *Findbar) Focus(delegate func(p tview.Primitive)) {
	delegate(f.SearchStringField)
}

func (f *Findbar) SetOnError(onError func(err error)) {
	f.OnError = onError
}

func (f *Findbar) SetOnMessage(onMessage func(message string)) {
	f.OnMessage = onMessage
}

func (f *Findbar) SetSearchText(text string) {
	f.SearchStringField.SetText(text)
}

func (f *Findbar) search(directionDown bool) bool {
	searchText := f.SearchStringField.GetText()
	if searchText == "" {
		return false
	}
	regex := f.RegexCheckbox.IsChecked()
	caseSensitive := f.CaseSensitiveCheckbox.IsChecked()

	found, err := f.editor.ActionController().Search(searchText, regex, caseSensitive, directionDown)
	if err != nil {
		if f.OnError != nil {
			f.OnError(err)
		}
		return false
	}

	if !found {
		// Wrap the cursor around either to the start or end of the buffer
		prevLoc := f.editor.Cursor().Loc
		prevStartSelection := f.editor.Cursor().CurSelection[0]
		prevEndSelection := f.editor.Cursor().CurSelection[1]

		var wrapLoc buffer.Loc
		if directionDown {
			wrapLoc = f.editor.Buffer().Start()
		} else {
			wrapLoc = f.editor.Buffer().End()
		}
		f.editor.Cursor().Loc = wrapLoc
		f.editor.Cursor().SetSelectionStart(wrapLoc)
		f.editor.Cursor().SetSelectionEnd(wrapLoc)

		found, _ := f.editor.ActionController().Search(searchText, regex, caseSensitive, directionDown)
		if !found {
			// Restore previous cursor position if not found
			f.editor.Cursor().Loc = prevLoc
			f.editor.Cursor().SetSelectionStart(prevStartSelection)
			f.editor.Cursor().SetSelectionEnd(prevEndSelection)
			return false
		}
	}
	f.editor.Relocate()
	return true
}

func (f *Findbar) SearchUp() {
	f.search(false)
}

func (f *Findbar) SearchDown() {
	f.search(true)
}

func (f *Findbar) Replace() {
	// Collapse the selection if there is one
	if f.editor.Cursor().HasSelection() {
		f.editor.Cursor().Loc = f.editor.Cursor().CurSelection[0]
		f.editor.Cursor().ResetSelection()
	}

	found := f.search(true)
	if !found {
		return
	}
	replaceText := f.hFlex2.GetItem(1).(*femtoinputfield.FemtoInputField).GetText()

	if f.editor.Cursor().HasSelection() {
		f.editor.Cursor().DeleteSelection()
		f.editor.Cursor().ResetSelection()
	}
	f.editor.Buffer().Insert(f.editor.Cursor().Loc, replaceText)

	f.search(true)

	f.editor.Relocate()
}

func (f *Findbar) ReplaceAll() {
	regex := f.RegexCheckbox.IsChecked()
	caseSensitive := f.CaseSensitiveCheckbox.IsChecked()
	count, err := f.editor.ActionController().ReplaceAll(f.SearchStringField.GetText(), regex, caseSensitive, f.hFlex2.GetItem(1).(*femtoinputfield.FemtoInputField).GetText())
	if err != nil {
		if f.OnError != nil {
			f.OnError(err)
		}
		return
	}
	if f.OnMessage != nil {
		f.OnMessage(fmt.Sprintf("Replaced %d occurrences", count))
	}
}

func (f *Findbar) SetFemtoKeybindings(keybindings smidgen.Keybindings) {
	f.SearchStringField.SetKeybindings(keybindings)
}

func (f *Findbar) SetOnExpand(onExpand func(expanded bool)) {
	f.OnExpand = onExpand
}
