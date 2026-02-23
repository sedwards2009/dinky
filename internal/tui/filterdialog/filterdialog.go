package filterdialog

import (
	"dinky/internal/tui/dialog"
	"dinky/internal/tui/smidgeninputfield"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sedwards2009/smidgen"
)

type FilterDialog struct {
	*tview.Flex
	app *tview.Application

	verticalContentsFlex *tview.Flex
	buttonsFlex          *tview.Flex
	CommandInputField    *smidgeninputfield.SmidgenInputField
	DirectoryInputField  *smidgeninputfield.SmidgenInputField
	innerFlex            *tview.Flex
	inputFieldFlex       *tview.Flex
	directoryFieldFlex   *tview.Flex

	Buttons        []*tview.Button
	options        FilterDialogOptions
	InputLabel     *tview.TextView
	DirectoryLabel *tview.TextView
}

type FilterDialogOptions struct {
	OnCancel func()
	OnAccept func(command string, directory string, buttonIndex int)
}

const FilterDialogWidth = 60
const FilterDialogHeight = 11

func NewFilterDialog(app *tview.Application) *FilterDialog {
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
	verticalContentsFlex.SetTitle("Filter via Shell")

	inputField := smidgeninputfield.NewSmidgenInputField(app)

	inputFieldFlex := tview.NewFlex()
	inputFieldFlex.SetDirection(tview.FlexColumn)
	inputFieldFlex.SetBorder(false)

	inputLabel := tview.NewTextView()
	inputLabel.SetText("Shell command: ")
	inputFieldFlex.AddItem(inputLabel, 15, 0, false)
	inputFieldFlex.AddItem(inputField, 0, 1, true)
	verticalContentsFlex.AddItem(inputFieldFlex, 1, 0, false)

	directoryField := smidgeninputfield.NewSmidgenInputField(app)
	dir, err := os.Getwd()
	if err == nil {
		directoryField.SetText(dir)
	}

	directoryFieldFlex := tview.NewFlex()
	directoryFieldFlex.SetDirection(tview.FlexColumn)
	directoryFieldFlex.SetBorder(false)

	directoryLabel := tview.NewTextView()
	directoryLabel.SetText("Directory: ")
	directoryFieldFlex.AddItem(directoryLabel, 15, 0, false)
	directoryFieldFlex.AddItem(directoryField, 0, 1, true)
	verticalContentsFlex.AddItem(directoryFieldFlex, 1, 0, false)

	verticalContentsFlex.AddItem(nil, 1, 0, false)

	explanationLabel := tview.NewTextView()
	explanationLabel.SetText("The command is run via `sh`. The selection\nis piped in and replaced with the output.")
	verticalContentsFlex.AddItem(explanationLabel, 2, 0, false)
	verticalContentsFlex.AddItem(nil, 1, 0, false)

	buttonsFlex := tview.NewFlex()
	buttonsFlex.SetDirection(tview.FlexColumn)
	buttonsFlex.SetBorder(false)
	verticalContentsFlex.AddItem(buttonsFlex, 1, 0, false)

	innerFlex.AddItem(verticalContentsFlex, FilterDialogWidth, 0, true)
	innerFlex.AddItem(nil, 0, 1, false)
	innerFlex.SetDirection(tview.FlexColumn)

	topLayout.AddItem(innerFlex, FilterDialogHeight, 0, true)
	topLayout.AddItem(nil, 0, 1, false)
	topLayout.SetDirection(tview.FlexRow)

	result := &FilterDialog{
		Flex:                 topLayout,
		app:                  app,
		verticalContentsFlex: verticalContentsFlex,
		innerFlex:            innerFlex,
		buttonsFlex:          buttonsFlex,
		inputFieldFlex:       inputFieldFlex,
		directoryFieldFlex:   directoryFieldFlex,
		CommandInputField:    inputField,
		DirectoryInputField:  directoryField,
		InputLabel:           inputLabel,
		DirectoryLabel:       directoryLabel,
	}
	return result
}

func (d *FilterDialog) Open(options FilterDialogOptions) {
	d.options = options

	onButtonClick := func(button string, index int) {
		d.options.OnAccept(d.CommandInputField.GetText(), d.DirectoryInputField.GetText(), index)
	}

	d.Buttons = dialog.CreateButtonsRow(d.buttonsFlex, []string{"OK", "Cancel"}, onButtonClick)
	for _, btn := range d.Buttons {
		btn.SetInputCapture(d.inputFilter)
	}
	d.CommandInputField.SetInputCapture(d.inputFilter)
	d.DirectoryInputField.SetInputCapture(d.inputFilter)

	d.app.SetFocus(d.CommandInputField)
}

func (d *FilterDialog) Close() {
}

func (d *FilterDialog) inputFilter(event *tcell.EventKey) *tcell.EventKey {
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
		if d.CommandInputField.HasFocus() || d.DirectoryInputField.HasFocus() {
			if d.options.OnAccept != nil {
				d.options.OnAccept(d.CommandInputField.GetText(), d.DirectoryInputField.GetText(), -1)
			}
		}
		return nil
	}

	// if d.InputField.HasFocus() && d.options.FieldKeyFilter != nil {
	// 	if !d.options.FieldKeyFilter(event) {
	// 		return nil
	// 	}
	// }

	return event
}

func (d *FilterDialog) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return d.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		d.verticalContentsFlex.MouseHandler()(action, event, setFocus)
		return true, nil
	})
}

// Focus is called when this primitive receives focus.
func (d *FilterDialog) Focus(delegate func(p tview.Primitive)) {
	delegate(d.CommandInputField)
}

func (d *FilterDialog) handleTabKey(direction int) {
	widgets := []tview.Primitive{}
	for _, btn := range d.Buttons {
		widgets = append(widgets, btn)
	}
	widgets = append(widgets, d.CommandInputField)
	widgets = append(widgets, d.DirectoryInputField)

	for i := 0; i < len(widgets); i++ {
		if widgets[i].HasFocus() {
			d.app.SetFocus(widgets[(i+direction)%len(widgets)])
			return
		}
	}
}

func (d *FilterDialog) SetSmidgenKeybindings(keybindings smidgen.Keybindings) {
	d.CommandInputField.SetKeybindings(keybindings)
	d.DirectoryInputField.SetKeybindings(keybindings)
}
