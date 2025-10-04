package dialog

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type InputDialog struct {
	*tview.Flex
	app *tview.Application

	verticalContentsFlex *tview.Flex
	buttonsFlex          *tview.Flex
	InputField           *tview.InputField
	innerFlex            *tview.Flex
	Buttons              []*tview.Button
	options              InputDialogOptions
}

type InputDialogOptions struct {
	Title          string
	Message        string
	DefaultValue   string
	Buttons        []string
	Width          int
	Height         int
	OnCancel       func()
	OnAccept       func(value string, index int)
	FieldKeyFilter func(event *tcell.EventKey) bool // Returns true if the key permitted
}

func NewInputDialog(app *tview.Application) *InputDialog {
	topLayout := tview.NewFlex()

	topLayout.AddItem(nil, 0, 1, false)

	innerFlex := tview.NewFlex()
	innerFlex.AddItem(nil, 0, 1, false)

	verticalContentsFlex := tview.NewFlex()

	verticalContentsFlex.Box = tview.NewBox() // Nasty hack to clear the `dontClear` flag inside Box.
	verticalContentsFlex.Box.Primitive = topLayout

	verticalContentsFlex.SetDirection(tview.FlexRow)
	verticalContentsFlex.SetBorderPadding(1, 1, 1, 1)
	verticalContentsFlex.SetBorder(true)
	verticalContentsFlex.SetTitleAlign(tview.AlignLeft)

	inputField := tview.NewInputField()
	verticalContentsFlex.AddItem(inputField, 0, 1, false)

	buttonsFlex := tview.NewFlex()
	buttonsFlex.SetDirection(tview.FlexColumn)
	// buttonsFlex.SetBackgroundTransparent(false)
	buttonsFlex.SetBorder(false)
	verticalContentsFlex.AddItem(buttonsFlex, 1, 0, false)

	innerFlex.AddItem(verticalContentsFlex, 80, 0, true)
	innerFlex.AddItem(nil, 0, 1, false)
	innerFlex.SetDirection(tview.FlexColumn)

	topLayout.AddItem(innerFlex, 20, 0, true)
	topLayout.AddItem(nil, 0, 1, false)
	topLayout.SetDirection(tview.FlexRow)

	result := &InputDialog{
		Flex:                 topLayout,
		app:                  app,
		verticalContentsFlex: verticalContentsFlex,
		innerFlex:            innerFlex,
		buttonsFlex:          buttonsFlex,
		InputField:           inputField,
	}
	return result
}

func (d *InputDialog) Open(options InputDialogOptions) {
	d.options = options
	d.verticalContentsFlex.SetTitle(options.Title)
	d.InputField.SetLabel(options.Message)
	d.InputField.SetText(options.DefaultValue)

	onButtonClick := func(button string, index int) {
		d.options.OnAccept(d.InputField.GetText(), index)
	}

	d.Buttons = createButtonsRow(d.buttonsFlex, options.Buttons, onButtonClick)
	d.ResizeItem(d.innerFlex, options.Height, 0)
	d.innerFlex.ResizeItem(d.verticalContentsFlex, options.Width, 0)

	for _, btn := range d.Buttons {
		btn.SetInputCapture(d.inputFilter)
	}
	d.InputField.SetInputCapture(d.inputFilter)

	d.app.SetFocus(d.InputField)
}

func (d *InputDialog) Close() {
}

func (d *InputDialog) inputFilter(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEscape:
		if d.options.OnCancel != nil {
			d.options.OnCancel()
		}
		return nil

	case tcell.KeyLeft:
		for i := 1; i < len(d.Buttons); i++ {
			if d.Buttons[i].HasFocus() {
				d.app.SetFocus(d.Buttons[i-1])
				return nil
			}
		}

	case tcell.KeyRight:
		for i := 0; i < len(d.Buttons)-1; i++ {
			if d.Buttons[i].HasFocus() {
				d.app.SetFocus(d.Buttons[i+1])
				return nil
			}
		}

	case tcell.KeyTab:
		if event.Modifiers() == tcell.ModNone {
			d.handleTabKey(1)
		} else if event.Modifiers() == tcell.ModShift {
			d.handleTabKey(-1)
		}

	case tcell.KeyEnter:
		if d.InputField.HasFocus() {
			if d.options.OnAccept != nil {
				d.options.OnAccept(d.InputField.GetText(), -1)
			}
		}
		return nil
	}

	if d.InputField.HasFocus() && d.options.FieldKeyFilter != nil {
		if !d.options.FieldKeyFilter(event) {
			return nil
		}
	}

	return event
}

func (d *InputDialog) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return d.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		d.verticalContentsFlex.MouseHandler()(action, event, setFocus)
		return true, nil
	})
}

// Focus is called when this primitive receives focus.
func (d *InputDialog) Focus(delegate func(p tview.Primitive)) {
	delegate(d.InputField)
}

func (d *InputDialog) handleTabKey(direction int) {
	widgets := []tview.Primitive{}
	for _, btn := range d.Buttons {
		widgets = append(widgets, btn)
	}
	widgets = append(widgets, d.InputField)

	for i := 0; i < len(widgets); i++ {
		if widgets[i].HasFocus() {
			d.app.SetFocus(widgets[(i+direction)%len(widgets)])
			return
		}
	}
}
