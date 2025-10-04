package dialog

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type MessageDialog struct {
	*tview.Flex
	app *tview.Application

	messageView          *tview.TextView
	verticalContentsFlex *tview.Flex
	buttonsFlex          *tview.Flex
	innerFlex            *tview.Flex
	Buttons              []*tview.Button

	OnClose       func()
	OnButtonClick func(button string, index int)
}

func NewMessageDialog(app *tview.Application) *MessageDialog {
	topLayout := tview.NewFlex()

	topLayout.AddItem(nil, 0, 1, false)

	innerFlex := tview.NewFlex()
	innerFlex.AddItem(nil, 0, 1, false)

	verticalContentsFlex := tview.NewFlex()
	verticalContentsFlex.SetDirection(tview.FlexRow)

	verticalContentsFlex.Box = tview.NewBox() // Nasty hack to clear the `dontClear` flag inside Box.
	verticalContentsFlex.Box.Primitive = verticalContentsFlex

	verticalContentsFlex.SetBorderPadding(1, 1, 1, 1)
	// verticalContentsFlex.SetBackgroundTransparent(false)
	verticalContentsFlex.SetBorder(true)
	verticalContentsFlex.SetTitleAlign(tview.AlignLeft)

	messageView := tview.NewTextView()
	verticalContentsFlex.AddItem(messageView, 0, 1, false)

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

	result := &MessageDialog{
		Flex:                 topLayout,
		app:                  app,
		messageView:          messageView,
		verticalContentsFlex: verticalContentsFlex,
		innerFlex:            innerFlex,
		buttonsFlex:          buttonsFlex,
	}
	return result
}

// Open displays the message dialog with the specified title, message, and buttons.
// The width and height parameters control the size of the dialog.
// The buttons parameter is a slice of strings representing the button labels.
// The first button will be focused by default.
func (d *MessageDialog) Open(title string, message string, buttons []string, width int, height int) {
	d.verticalContentsFlex.SetTitle(title)
	d.messageView.SetText(message)

	d.Buttons = createButtonsRow(d.buttonsFlex, buttons, d.OnButtonClick)
	for _, btn := range d.Buttons {
		btn.SetInputCapture(d.inputFilter)
	}
	d.ResizeItem(d.innerFlex, height, 0)
	d.innerFlex.ResizeItem(d.verticalContentsFlex, width, 0)
}

func (d *MessageDialog) Close() {
}

func (d *MessageDialog) inputFilter(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEscape:
		if d.OnClose != nil {
			d.OnClose()
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
		for i := 0; i < len(d.Buttons); i++ {
			if d.Buttons[i].HasFocus() {
				d.app.SetFocus(d.Buttons[(i+1)%len(d.Buttons)])
				return nil
			}
		}
	default:
	}
	return event
}

// Focus is called when this primitive receives focus.
func (d *MessageDialog) Focus(delegate func(p tview.Primitive)) {
	if len(d.Buttons) == 0 {
		return
	}
	delegate(d.Buttons[0])
}

func (d *MessageDialog) FocusButton(index int) {
	d.app.SetFocus(d.Buttons[index])
}

func (d *MessageDialog) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return d.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		d.verticalContentsFlex.MouseHandler()(action, event, setFocus)
		return true, nil
	})
}
