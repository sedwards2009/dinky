package dialog

import (
	"dinky/internal/tui/scrollbar"
	"dinky/internal/tui/table2"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ListDialog struct {
	*tview.Flex
	app *tview.Application

	messageView                 *tview.TextView
	verticalContentsFlex        *tview.Flex
	buttonsFlex                 *tview.Flex
	TableField                  *table2.Table
	tableFlex                   *tview.Flex
	VerticalScrollbar           *scrollbar.Scrollbar
	innerFlex                   *tview.Flex
	Buttons                     []*tview.Button
	options                     ListDialogOptions
	itemTextColor               tcell.Color
	itemBackgroundColor         tcell.Color
	selectedItemBackgroundColor tcell.Color
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

func NewListDialog(app *tview.Application) *ListDialog {
	topLayout := tview.NewFlex()
	topLayout.AddItem(nil, 0, 1, false)

	innerFlex := tview.NewFlex()
	innerFlex.AddItem(nil, 0, 1, false)

	verticalContentsFlex := tview.NewFlex()

	verticalContentsFlex.Box = tview.NewBox() // Nasty hack to clear the `dontClear` flag inside Box.
	verticalContentsFlex.Box.Primitive = verticalContentsFlex

	verticalContentsFlex.SetDirection(tview.FlexRow)
	verticalContentsFlex.SetBorderPadding(1, 1, 1, 1)
	// verticalContentsFlex.SetBackgroundTransparent(false)
	verticalContentsFlex.SetBorder(true)
	verticalContentsFlex.SetTitleAlign(tview.AlignLeft)

	messageView := tview.NewTextView()
	verticalContentsFlex.AddItem(messageView, 1, 0, false)
	verticalContentsFlex.AddItem(nil, 1, 0, false)

	tableField := table2.NewTable()

	tableFlex := tview.NewFlex()
	tableFlex.SetDirection(tview.FlexColumn)
	tableFlex.SetBorder(false)
	tableFlex.AddItem(tableField, 0, 1, false)

	verticalScrollbar := scrollbar.NewScrollbar()
	tableFlex.AddItem(verticalScrollbar, 1, 0, false)

	verticalContentsFlex.AddItem(tableFlex, 0, 1, false)
	verticalContentsFlex.AddItem(nil, 1, 0, false)

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

	d := &ListDialog{
		Flex:                 topLayout,
		app:                  app,
		messageView:          messageView,
		verticalContentsFlex: verticalContentsFlex,
		innerFlex:            innerFlex,
		buttonsFlex:          buttonsFlex,
		TableField:           tableField,
		tableFlex:            tableFlex,
		VerticalScrollbar:    verticalScrollbar,
	}

	// Set up vertical scrollbar
	verticalScrollbar.Track.SetBeforeDrawFunc(func(_ tcell.Screen) {
		row, _ := tableField.GetOffset()
		verticalScrollbar.Track.SetMax(tableField.GetRowCount() - 1)
		_, _, _, height := d.TableField.GetInnerRect()
		verticalScrollbar.Track.SetThumbSize(height)
		verticalScrollbar.Track.SetPosition(row)
	})
	verticalScrollbar.SetChangedFunc(func(position int) {
		_, column := d.TableField.GetOffset()
		d.TableField.SetOffset(position, column)
	})

	d.TableField.SetDoubleClickFunc(func(row int, _ int) {
		if d.options.OnAccept == nil {
			return
		}
		selection, _ := d.TableField.GetSelection()
		d.options.OnAccept(d.options.Items[selection].Value, -1)
	})

	return d
}

func (d *ListDialog) Open(options ListDialogOptions) {
	d.options = options
	d.verticalContentsFlex.SetTitle(options.Title)
	d.messageView.SetText(options.Message)

	onButtonClick := func(button string, index int) {
		selection, _ := d.TableField.GetSelection()
		d.options.OnAccept(d.options.Items[selection].Value, index)
	}

	// Fill in the table with items
	d.TableField.Clear()
	for rowIndex, item := range options.Items {
		cell := &table2.TableCell{
			Text:  item.Text,
			Style: tcell.StyleDefault.Foreground(d.itemTextColor).Background(d.itemBackgroundColor),
		}
		d.TableField.SetCell(rowIndex, 0, cell)
	}
	d.TableField.SetSelectable(true, false)

	d.TableField.Select(0, 0)
	for rowIndex, item := range options.Items {
		if item.Value == options.DefaultSelected {
			d.TableField.Select(rowIndex, 0)
			break
		}
	}

	d.Buttons = createButtonsRow(d.buttonsFlex, options.Buttons, onButtonClick)
	d.ResizeItem(d.innerFlex, options.Height, 0)
	d.innerFlex.ResizeItem(d.verticalContentsFlex, options.Width, 0)

	for _, btn := range d.Buttons {
		btn.SetInputCapture(d.inputFilter)
	}
	d.TableField.SetInputCapture(d.inputFilter)

	d.app.SetFocus(d.TableField)
}

func (d *ListDialog) Close() {
}

func (d *ListDialog) inputFilter(event *tcell.EventKey) *tcell.EventKey {
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
		if d.TableField.HasFocus() {
			if d.options.OnAccept != nil {
				selection, _ := d.TableField.GetSelection()
				d.options.OnAccept(d.options.Items[selection].Value, -1)
			}
		}
		return nil
	}

	return event
}

func (d *ListDialog) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return d.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		d.verticalContentsFlex.MouseHandler()(action, event, setFocus)
		return true, nil
	})
}

// Focus is called when this primitive receives focus.
func (d *ListDialog) Focus(delegate func(p tview.Primitive)) {
	delegate(d.TableField)
}

func (d *ListDialog) handleTabKey(direction int) {
	widgets := []tview.Primitive{}
	for _, btn := range d.Buttons {
		widgets = append(widgets, btn)
	}
	widgets = append(widgets, d.TableField)

	for i := 0; i < len(widgets); i++ {
		if widgets[i].HasFocus() {
			d.app.SetFocus(widgets[(i+direction)%len(widgets)])
			return
		}
	}
}

func (d *ListDialog) SetItemTextColor(color tcell.Color) {
	d.itemTextColor = color
	d.TableField.SetSelectedStyle(tcell.StyleDefault.Background(color).Foreground(d.itemTextColor))
}

func (d *ListDialog) SetItemBackgroundColor(color tcell.Color) {
	d.itemBackgroundColor = color
}

func (d *ListDialog) SetSelectedItemBackgroundColor(color tcell.Color) {
	d.selectedItemBackgroundColor = color
	d.TableField.SetSelectedStyle(tcell.StyleDefault.Background(color).Foreground(d.itemTextColor))
}
