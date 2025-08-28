package dialog

import (
	"github.com/gdamore/tcell/v2"
	"github.com/sedwards2009/nuview"
)

type MessageDialog struct {
	*nuview.Flex
	app *nuview.Application

	messageView          *nuview.TextView
	verticalContentsFlex *nuview.Flex
	buttonsFlex          *nuview.Flex
	innerFlex            *nuview.Flex
	buttons              []*nuview.Button

	OnClose       func()
	OnButtonClick func(button string, index int)
}

func NewMessageDialog(app *nuview.Application) *MessageDialog {
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

	messageView := nuview.NewTextView()
	verticalContentsFlex.AddItem(messageView, 0, 1, false)

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

	d.buttons = createButtonsRow(d.buttonsFlex, buttons, d.OnButtonClick)
	d.ResizeItem(d.innerFlex, height, 0)
	d.innerFlex.ResizeItem(d.verticalContentsFlex, width, 0)

	d.app.SetInputCapture(d.inputFilter)
	d.FocusButton(0)
}

func (d *MessageDialog) Close() {
	d.app.SetInputCapture(nil)
}

func (d *MessageDialog) inputFilter(event *tcell.EventKey) *tcell.EventKey {
	key := event.Key()
	if key == tcell.KeyTab {
		if event.Modifiers() == tcell.ModNone {
			key = tcell.KeyRight
		} else if event.Modifiers() == tcell.ModShift {
			key = tcell.KeyLeft
		}
	}

	switch key {
	case tcell.KeyEscape:
		if d.OnClose != nil {
			d.OnClose()
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
	case tcell.KeyTab:
		for i := 0; i < len(d.buttons)-1; i++ {
			if d.buttons[i].HasFocus() {
				d.app.SetFocus(d.buttons[i+1])
				return nil
			}
		}
	}
	return event
}

func (d *MessageDialog) FocusButton(index int) {
	d.app.SetFocus(d.buttons[index])
}

func (d *MessageDialog) MouseHandler() func(action nuview.MouseAction, event *tcell.EventMouse, setFocus func(p nuview.Primitive)) (consumed bool, capture nuview.Primitive) {
	return d.WrapMouseHandler(func(action nuview.MouseAction, event *tcell.EventMouse, setFocus func(p nuview.Primitive)) (consumed bool, capture nuview.Primitive) {
		d.verticalContentsFlex.MouseHandler()(action, event, setFocus)
		return true, nil
	})
}
