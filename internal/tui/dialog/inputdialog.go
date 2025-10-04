package dialog

import (
	"github.com/gdamore/tcell/v2"
	nuview "github.com/rivo/tview"
)

type InputDialog struct {
	*nuview.Flex
	app *nuview.Application

	verticalContentsFlex *nuview.Flex
	buttonsFlex          *nuview.Flex
	InputField           *nuview.InputField
	innerFlex            *nuview.Flex
	Buttons              []*nuview.Button
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

	verticalContentsFlex.Box = nuview.NewBox() // Nasty hack to clear the `dontClear` flag inside Box.
	verticalContentsFlex.Box.Primitive = topLayout

	verticalContentsFlex.SetDirection(nuview.FlexRow)
	verticalContentsFlex.SetBorderPadding(1, 1, 1, 1)
	verticalContentsFlex.SetBorder(true)
	verticalContentsFlex.SetTitleAlign(nuview.AlignLeft)

	inputField := nuview.NewInputField()
	verticalContentsFlex.AddItem(inputField, 0, 1, false)

	buttonsFlex := nuview.NewFlex()
	buttonsFlex.SetDirection(nuview.FlexColumn)
	// buttonsFlex.SetBackgroundTransparent(false)
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

	d.app.SetInputCapture(d.inputFilter)
	d.app.SetFocus(d.InputField)
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

func (d *InputDialog) MouseHandler() func(action nuview.MouseAction, event *tcell.EventMouse, setFocus func(p nuview.Primitive)) (consumed bool, capture nuview.Primitive) {
	return d.WrapMouseHandler(func(action nuview.MouseAction, event *tcell.EventMouse, setFocus func(p nuview.Primitive)) (consumed bool, capture nuview.Primitive) {
		d.verticalContentsFlex.MouseHandler()(action, event, setFocus)
		return true, nil
	})
}

// Focus is called when this primitive receives focus.
func (d *InputDialog) Focus(delegate func(p nuview.Primitive)) {
	delegate(d.InputField)
}

func (d *InputDialog) handleTabKey(direction int) {
	widgets := []nuview.Primitive{}
	for _, btn := range d.Buttons {
		widgets = append(widgets, btn)
	}
	widgets = append(widgets, d.InputField)

	for i := 0; i < len(widgets); i++ {
		if widgets[i].HasFocus() {
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
