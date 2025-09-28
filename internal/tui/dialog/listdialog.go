package dialog

import (
	"dinky/internal/tui/scrollbar"
	"dinky/internal/tui/style"
	"dinky/internal/tui/table2"

	"github.com/gdamore/tcell/v2"
	nuview "github.com/rivo/tview"
)

type ListDialog struct {
	*nuview.Flex
	app *nuview.Application

	messageView          *nuview.TextView
	verticalContentsFlex *nuview.Flex
	buttonsFlex          *nuview.Flex
	tableField           *table2.Table
	tableFlex            *nuview.Flex
	VerticalScrollbar    *scrollbar.Scrollbar
	innerFlex            *nuview.Flex
	buttons              []*nuview.Button
	options              ListDialogOptions
}

type ListDialogOptions struct {
	Title           string
	Message         string
	Buttons         []string
	Width           int
	Height          int
	DefaultSelected string
	Items           []ListItem
	OnCancel        func()
	OnAccept        func(value string, index int)
}

type ListItem struct {
	Text  string
	Value string
}

func NewListDialog(app *nuview.Application) *ListDialog {
	topLayout := nuview.NewFlex()
	topLayout.AddItem(nil, 0, 1, false)

	innerFlex := nuview.NewFlex()
	innerFlex.AddItem(nil, 0, 1, false)

	verticalContentsFlex := nuview.NewFlex()
	verticalContentsFlex.SetDirection(nuview.FlexRow)
	verticalContentsFlex.SetBorderPadding(1, 1, 1, 1)
	// verticalContentsFlex.SetBackgroundTransparent(false)
	verticalContentsFlex.SetBorder(true)
	verticalContentsFlex.SetTitleAlign(nuview.AlignLeft)

	messageView := nuview.NewTextView()
	verticalContentsFlex.AddItem(messageView, 1, 0, false)
	verticalContentsFlex.AddItem(nil, 1, 0, false)

	tableField := table2.NewTable()
	style.StyleTable(tableField)

	tableFlex := nuview.NewFlex()
	tableFlex.SetDirection(nuview.FlexColumn)
	tableFlex.SetBorder(false)
	tableFlex.AddItem(tableField, 0, 1, false)

	verticalScrollbar := scrollbar.NewScrollbar()
	style.StyleScrollbar(verticalScrollbar)
	tableFlex.AddItem(verticalScrollbar, 1, 0, false)

	verticalContentsFlex.AddItem(tableFlex, 0, 1, false)
	verticalContentsFlex.AddItem(nil, 1, 0, false)

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

	d := &ListDialog{
		Flex:                 topLayout,
		app:                  app,
		messageView:          messageView,
		verticalContentsFlex: verticalContentsFlex,
		innerFlex:            innerFlex,
		buttonsFlex:          buttonsFlex,
		tableField:           tableField,
		tableFlex:            tableFlex,
		VerticalScrollbar:    verticalScrollbar,
	}

	// Set up vertical scrollbar
	verticalScrollbar.Track.SetBeforeDrawFunc(func(_ tcell.Screen) {
		row, _ := tableField.GetOffset()
		verticalScrollbar.Track.SetMax(tableField.GetRowCount() - 1)
		_, _, _, height := d.tableField.GetInnerRect()
		verticalScrollbar.Track.SetThumbSize(height)
		verticalScrollbar.Track.SetPosition(row)
	})
	verticalScrollbar.SetChangedFunc(func(position int) {
		_, column := d.tableField.GetOffset()
		d.tableField.SetOffset(position, column)
	})

	d.tableField.SetDoubleClickFunc(func(row int, _ int) {
		if d.options.OnAccept == nil {
			return
		}
		selection, _ := d.tableField.GetSelection()
		d.options.OnAccept(d.options.Items[selection].Value, -1)
	})

	return d
}

func (d *ListDialog) Open(options ListDialogOptions) {
	d.options = options
	d.verticalContentsFlex.SetTitle(options.Title)
	d.messageView.SetText(options.Message)

	onButtonClick := func(button string, index int) {
		selection, _ := d.tableField.GetSelection()
		d.options.OnAccept(d.options.Items[selection].Value, index)
	}

	// Fill in the table with items
	d.tableField.Clear()
	for rowIndex, item := range options.Items {
		cell := &table2.TableCell{
			Text: item.Text,
		}
		style.StyleTableCell(cell)
		d.tableField.SetCell(rowIndex, 0, cell)
	}
	d.tableField.SetSelectable(true, false)

	d.tableField.Select(0, 0)
	for rowIndex, item := range options.Items {
		if item.Value == options.DefaultSelected {
			d.tableField.Select(rowIndex, 0)
			break
		}
	}

	d.buttons = createButtonsRow(d.buttonsFlex, options.Buttons, onButtonClick)
	d.ResizeItem(d.innerFlex, options.Height, 0)
	d.innerFlex.ResizeItem(d.verticalContentsFlex, options.Width, 0)

	d.app.SetInputCapture(d.inputFilter)
	d.app.SetFocus(d.tableField)
}

func (d *ListDialog) Close() {
	d.app.SetInputCapture(nil)
}

func (d *ListDialog) inputFilter(event *tcell.EventKey) *tcell.EventKey {
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
		if d.tableField.HasFocus() {
			if d.options.OnAccept != nil {
				selection, _ := d.tableField.GetSelection()
				d.options.OnAccept(d.options.Items[selection].Value, -1)
			}
		}
		return nil
	}

	return event
}

func (d *ListDialog) MouseHandler() func(action nuview.MouseAction, event *tcell.EventMouse, setFocus func(p nuview.Primitive)) (consumed bool, capture nuview.Primitive) {
	return d.WrapMouseHandler(func(action nuview.MouseAction, event *tcell.EventMouse, setFocus func(p nuview.Primitive)) (consumed bool, capture nuview.Primitive) {
		d.verticalContentsFlex.MouseHandler()(action, event, setFocus)
		return true, nil
	})
}

// Focus is called when this primitive receives focus.
func (d *ListDialog) Focus(delegate func(p nuview.Primitive)) {
	delegate(d.tableField)
}

func (d *ListDialog) handleTab(direction int) {
	widgets := []nuview.Primitive{}
	for _, btn := range d.buttons {
		widgets = append(widgets, btn)
	}
	widgets = append(widgets, d.tableField)

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
