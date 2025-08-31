package dialog

import (
	"github.com/gdamore/tcell/v2"
	"github.com/sedwards2009/nuview"
)

type InputDialog struct {
	*nuview.Flex
	app *nuview.Application

	verticalContentsFlex *nuview.Flex
	buttonsFlex          *nuview.Flex
	inputField           *nuview.InputField
	innerFlex            *nuview.Flex
	buttons              []*nuview.Button
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

func NewInputDialog(app *nuview.Application) *InputDialog {
	topLayout := nuview.NewFlex()
	topLayout.AddItem(nil, 0, 1, false)

	innerFlex := nuview.NewFlex()
	innerFlex.AddItem(nil, 0, 1, false)

	verticalContentsFlex := nuview.NewFlex()
	verticalContentsFlex.SetDirection(nuview.FlexRow)
	verticalContentsFlex.SetPadding(1, 1, 1, 1)
	verticalContentsFlex.SetBackgroundTransparent(false)
	verticalContentsFlex.SetBorder(true)
	verticalContentsFlex.SetTitleAlign(nuview.AlignLeft)

	inputField := nuview.NewInputField()
	verticalContentsFlex.AddItem(inputField, 0, 1, false)

	buttonsFlex := nuview.NewFlex()
	buttonsFlex.SetDirection(nuview.FlexColumn)
	buttonsFlex.SetBackgroundTransparent(false)
	buttonsFlex.SetBorder(false)
	verticalContentsFlex.AddItem(buttonsFlex, 1, 0, false)

	innerFlex.AddItem(verticalContentsFlex, 80, 0, true)
	innerFlex.AddItem(nil, 0, 1, false)
	innerFlex.SetDirection(nuview.FlexColumn)

	topLayout.AddItem(innerFlex, 20, 0, true)
	topLayout.AddItem(nil, 0, 1, false)
	topLayout.SetDirection(nuview.FlexRow)

	result := &InputDialog{
		Flex:                 topLayout,
		app:                  app,
		verticalContentsFlex: verticalContentsFlex,
		innerFlex:            innerFlex,
		buttonsFlex:          buttonsFlex,
		inputField:           inputField,
	}
	return result
}

func (d *InputDialog) Open(options InputDialogOptions) {
	d.options = options
	d.verticalContentsFlex.SetTitle(options.Title)
	d.inputField.SetLabel(options.Message)
	d.inputField.SetText(options.DefaultValue)

	onButtonClick := func(button string, index int) {
		d.options.OnAccept(d.inputField.GetText(), index)
	}

	d.buttons = createButtonsRow(d.buttonsFlex, options.Buttons, onButtonClick)
	d.ResizeItem(d.innerFlex, options.Height, 0)
	d.innerFlex.ResizeItem(d.verticalContentsFlex, options.Width, 0)

	d.app.SetInputCapture(d.inputFilter)
	d.app.SetFocus(d.inputField)
}

func (d *InputDialog) Close() {
	d.app.SetInputCapture(nil)
}

func (d *InputDialog) inputFilter(event *tcell.EventKey) *tcell.EventKey {
	key := event.Key()

	switch key {
	case tcell.KeyEscape:
		if d.options.OnCancel != nil {
			d.options.OnCancel()
		}
		return nil

	case tcell.KeyLeft:
		for i := 1; i < len(d.buttons); i++ {
			if d.buttons[i].HasFocus() {
				d.app.SetFocus(d.buttons[i-1])
				return nil
			}
		}

	case tcell.KeyRight:
		for i := 0; i < len(d.buttons)-1; i++ {
			if d.buttons[i].HasFocus() {
				d.app.SetFocus(d.buttons[i+1])
				return nil
			}
		}

	case tcell.KeyTab:
		if event.Modifiers() == tcell.ModNone {
			d.handleTab(1)
		} else if event.Modifiers() == tcell.ModShift {
			d.handleTab(-1)
		}

	case tcell.KeyEnter:
		if d.inputField.HasFocus() {
			if d.options.OnAccept != nil {
				d.options.OnAccept(d.inputField.GetText(), -1)
			}
		}
		return nil
	}

	if d.inputField.HasFocus() && d.options.FieldKeyFilter != nil {
		if !d.options.FieldKeyFilter(event) {
			return nil
		}
	}

	return event
}

func (d *InputDialog) MouseHandler() func(action nuview.MouseAction, event *tcell.EventMouse, setFocus func(p nuview.Primitive)) (consumed bool, capture nuview.Primitive) {
	return d.WrapMouseHandler(func(action nuview.MouseAction, event *tcell.EventMouse, setFocus func(p nuview.Primitive)) (consumed bool, capture nuview.Primitive) {
		d.verticalContentsFlex.MouseHandler()(action, event, setFocus)
		return true, nil
	})
}

// Focus is called when this primitive receives focus.
func (d *InputDialog) Focus(delegate func(p nuview.Primitive)) {
	delegate(d.inputField)
}

func (d *InputDialog) handleTab(direction int) {
	widgets := []nuview.Primitive{}
	for _, btn := range d.buttons {
		widgets = append(widgets, btn)
	}
	widgets = append(widgets, d.inputField)

	for i := 0; i < len(widgets); i++ {
		if widgets[i].GetFocusable().HasFocus() {
			x := i + direction
			if x < 0 {
				x = len(widgets) - 1
			} else if x >= len(widgets) {
				x = 0
			} else {
			}
			d.app.SetFocus(widgets[x])
			return
		}
	}
}
