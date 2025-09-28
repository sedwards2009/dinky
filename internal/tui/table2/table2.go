package table2

import (
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	nuview "github.com/rivo/tview"
)

// TableCell represents one cell inside a Table. You can instantiate this type
// directly but all colors (background and text) will be set to their default
// which is black.
type TableCell struct {
	// The reference object.
	Reference interface{}

	// The text to be displayed in the table cell.
	Text string

	// The alignment of the cell text. One of AlignLeft (default), AlignCenter,
	// or AlignRight.
	Align int

	// The maximum width of the cell in screen space. This is used to give a
	// column a maximum width. Any cell text whose screen width exceeds this width
	// is cut off. Set to 0 if there is no maximum width.
	MaxWidth int

	// If the total table width is less than the available width, this value is
	// used to add extra width to a column. See SetExpansion() for details.
	Expansion int

	// The color of the cell text. You should not use this anymore, it is only
	// here for backwards compatibility. Use the Style field instead.
	Color tcell.Color

	// The background color of the cell. You should not use this anymore, it is
	// only here for backwards compatibility. Use the Style field instead.
	BackgroundColor tcell.Color

	// The style attributes of the cell. You should not use this anymore, it is
	// only here for backwards compatibility. Use the Style field instead.
	Attributes tcell.AttrMask

	// The style of the cell. If this is uninitialized (tcell.StyleDefault), the
	// Color and BackgroundColor fields are used instead.
	Style tcell.Style

	// The style of the cell when it is selected. If this is uninitialized
	// (tcell.StyleDefault), the table's selected style is used instead. If that
	// is uninitialized as well, the cell's background and text color are
	// swapped.
	SelectedStyle tcell.Style

	// If set to true, the BackgroundColor is not used and the cell will have
	// the background color of the table.
	Transparent bool

	// If set to true, this cell cannot be selected.
	NotSelectable bool

	// An optional handler for mouse clicks. This also fires if the cell is not
	// selectable. If true is returned, no additional "selected" event is fired
	// on selectable cells.
	Clicked func() bool

	// The position and width of the cell the last time table was drawn.
	x int
	y int

	// The effective width of the cell. This takes into account the text size and MaxWidth
	width int

	sync.RWMutex
}

const StandardDoubleClick = 500

// NewTableCell returns a new table cell with sensible defaults. That is, left
// aligned text with the primary text color (see Styles) and a transparent
// background (using the background of the Table).
func NewTableCell(text string) *TableCell {
	cell := &TableCell{
		Text:        text,
		Align:       nuview.AlignLeft,
		Style:       tcell.StyleDefault.Foreground(nuview.Styles.PrimaryTextColor).Background(nuview.Styles.PrimitiveBackgroundColor),
		Transparent: true,
	}
	cell.updateWidth()
	return cell
}

// updateWidth calculates and sets the effective width of the cell based on
// the text content and maximum width constraint.
func (c *TableCell) updateWidth() {
	textWidth := nuview.TaggedStringWidth(c.Text)
	if c.MaxWidth > 0 && textWidth > c.MaxWidth {
		c.width = c.MaxWidth
	} else {
		c.width = textWidth
	}
}

// SetText sets the cell's text.
func (c *TableCell) SetText(text string) {
	c.Lock()
	defer c.Unlock()
	c.Text = text
	c.updateWidth()
}

// SetAlign sets the cell's text alignment, one of AlignLeft, AlignCenter, or
// AlignRight.
func (c *TableCell) SetAlign(align int) {
	c.Lock()
	defer c.Unlock()
	c.Align = align
}

// SetMaxWidth sets maximum width of the cell in screen space. This is used to
// give a column a maximum width. Any cell text whose screen width exceeds this
// width is cut off. Set to 0 if there is no maximum width.
func (c *TableCell) SetMaxWidth(maxWidth int) {
	c.Lock()
	defer c.Unlock()
	c.MaxWidth = maxWidth
	c.updateWidth()
}

// SetExpansion sets the value by which the column of this cell expands if the
// available width for the table is more than the table width (prior to applying
// this expansion value). This is a proportional value. The amount of unused
// horizontal space is divided into widths to be added to each column. How much
// extra width a column receives depends on the expansion value: A value of 0
// (the default) will not cause the column to increase in width. Other values
// are proportional, e.g. a value of 2 will cause a column to grow by twice
// the amount of a column with a value of 1.
//
// Since this value affects an entire column, the maximum over all visible cells
// in that column is used.
//
// This function panics if a negative value is provided.
func (c *TableCell) SetExpansion(expansion int) {
	c.Lock()
	defer c.Unlock()
	if expansion < 0 {
		panic("Table cell expansion values may not be negative")
	}
	c.Expansion = expansion
}

// SetTextColor sets the cell's text color.
func (c *TableCell) SetTextColor(color tcell.Color) {
	c.Lock()
	defer c.Unlock()
	if c.Style == tcell.StyleDefault {
		c.Color = color
	} else {
		c.Style = c.Style.Foreground(color)
	}
}

// SetBackgroundColor sets the cell's background color. This will also cause the
// cell's Transparent flag to be set to "false".
func (c *TableCell) SetBackgroundColor(color tcell.Color) {
	c.Lock()
	defer c.Unlock()
	if c.Style == tcell.StyleDefault {
		c.BackgroundColor = color
	} else {
		c.Style = c.Style.Background(color)
	}
	c.Transparent = false
}

// SetTransparency sets the background transparency of this cell. A value of
// "true" will cause the cell to use the table's background color. A value of
// "false" will cause it to use its own background color.
func (c *TableCell) SetTransparency(transparent bool) {
	c.Lock()
	defer c.Unlock()
	c.Transparent = transparent
}

// SetAttributes sets the cell's text attributes. You can combine different
// attributes using bitmask operations:
//
//	cell.SetAttributes(tcell.AttrUnderline | tcell.AttrBold)
func (c *TableCell) SetAttributes(attr tcell.AttrMask) {
	c.Lock()
	defer c.Unlock()
	if c.Style == tcell.StyleDefault {
		c.Attributes = attr
	} else {
		c.Style = c.Style.Attributes(attr)
	}
}

// SetStyle sets the cell's style (foreground color, background color, and
// attributes) all at once.
func (c *TableCell) SetStyle(style tcell.Style) {
	c.Lock()
	defer c.Unlock()
	c.Style = style
}

// SetSelectedStyle sets the cell's style when it is selected. If this is
// uninitialized (tcell.StyleDefault), the table's selected style is used
// instead. If that is uninitialized as well, the cell's background and text
// color are swapped.
func (c *TableCell) SetSelectedStyle(style tcell.Style) {
	c.Lock()
	defer c.Unlock()
	c.SelectedStyle = style
}

// SetSelectable sets whether or not this cell can be selected by the user.
func (c *TableCell) SetSelectable(selectable bool) {
	c.Lock()
	defer c.Unlock()
	c.NotSelectable = !selectable
}

// SetReference allows you to store a reference of any type in this cell. This
// will allow you to establish a mapping between the cell and your
// actual data.
func (c *TableCell) SetReference(reference interface{}) {
	c.Lock()
	defer c.Unlock()
	c.Reference = reference
}

// GetReference returns this cell's reference object.
func (c *TableCell) GetReference() interface{} {
	c.RLock()
	defer c.RUnlock()
	return c.Reference
}

// GetLastPosition returns the position of the table cell the last time it was
// drawn on screen. If the cell is not on screen, the return values are
// undefined.
//
// Because the Table class will attempt to keep selected cells on screen, this
// function is most useful in response to a "selected" event (see
// SetSelectedFunc()) or a "selectionChanged" event (see
// SetSelectionChangedFunc()).
func (c *TableCell) GetLastPosition() (x, y, width int) {
	c.RLock()
	defer c.RUnlock()
	return c.x, c.y, c.width
}

// SetClickedFunc sets a handler which fires when this cell is clicked. This is
// independent of whether the cell is selectable or not. But for selectable
// cells, if the function returns "true", the "selected" event is not fired.
func (c *TableCell) SetClickedFunc(clicked func() bool) {
	c.Lock()
	defer c.Unlock()
	c.Clicked = clicked
}

type tableContent interface {
	// Return the cell at the given position or nil if there is no cell. The
	// row and column arguments start at 0 and end at what GetRowCount() and
	// GetColumnCount() return, minus 1.
	GetCell(row, column int) *TableCell

	// Return the total number of rows in the table.
	GetRowCount() int

	// Return the total number of columns in the table.
	GetColumnCount() int

	// The following functions are provided for completeness reasons as the
	// original Table implementation was not read-only. If you do not wish to
	// forward modifying operations to your data, you may opt to leave these
	// functions empty. To make this easier, you can include the
	// TableContentReadOnly type in your struct. See also the
	// demos/table/virtualtable example.

	// Set the cell at the given position to the provided cell.
	SetCell(row, column int, cell *TableCell)

	// Remove the row at the given position by shifting all following rows up
	// by one. Out of range positions may be ignored.
	RemoveRow(row int)

	// Remove the column at the given position by shifting all following columns
	// left by one. Out of range positions may be ignored.
	RemoveColumn(column int)

	// Insert a new empty row at the given position by shifting all rows at that
	// position and below down by one. Implementers may decide what to do with
	// out of range positions.
	InsertRow(row int)

	// Insert a new empty column at the given position by shifting all columns
	// at that position and to the right by one to the right. Implementers may
	// decide what to do with out of range positions.
	InsertColumn(column int)

	// Remove all table data.
	Clear()
}

// tableDefaultContent implements the default TableContent interface for the
// Table class.
type tableDefaultContent struct {
	// The cells of the table. Rows first, then columns.
	cells [][]*TableCell

	// The rightmost column in the data set.
	lastColumn int
}

// Clear clears all data.
func (t *tableDefaultContent) Clear() {
	t.cells = nil
	t.lastColumn = -1
}

// SetCell sets a cell's content.
func (t *tableDefaultContent) SetCell(row, column int, cell *TableCell) {
	if row >= len(t.cells) {
		t.cells = append(t.cells, make([][]*TableCell, row-len(t.cells)+1)...)
	}
	rowLen := len(t.cells[row])
	if column >= rowLen {
		t.cells[row] = append(t.cells[row], make([]*TableCell, column-rowLen+1)...)
		for c := rowLen; c < column; c++ {
			t.cells[row][c] = &TableCell{}
		}
	}
	t.cells[row][column] = cell
	cell.updateWidth()
	if column > t.lastColumn {
		t.lastColumn = column
	}
}

// RemoveRow removes a row from the data.
func (t *tableDefaultContent) RemoveRow(row int) {
	if row < 0 || row >= len(t.cells) {
		return
	}
	t.cells = append(t.cells[:row], t.cells[row+1:]...)
}

// RemoveColumn removes a column from the data.
func (t *tableDefaultContent) RemoveColumn(column int) {
	for row := range t.cells {
		if column < 0 || column >= len(t.cells[row]) {
			continue
		}
		t.cells[row] = append(t.cells[row][:column], t.cells[row][column+1:]...)
	}
	if column >= 0 && column <= t.lastColumn {
		t.lastColumn--
	}
}

// InsertRow inserts a new row at the given position.
func (t *tableDefaultContent) InsertRow(row int) {
	if row >= len(t.cells) {
		return
	}
	t.cells = append(t.cells, nil)       // Extend by one.
	copy(t.cells[row+1:], t.cells[row:]) // Shift down.
	t.cells[row] = nil                   // New row is uninitialized.
}

// InsertColumn inserts a new column at the given position.
func (t *tableDefaultContent) InsertColumn(column int) {
	for row := range t.cells {
		if column >= len(t.cells[row]) {
			continue
		}
		t.cells[row] = append(t.cells[row], nil)             // Extend by one.
		copy(t.cells[row][column+1:], t.cells[row][column:]) // Shift to the right.
		t.cells[row][column] = &TableCell{}                  // New element is an uninitialized table cell.
	}
}

// GetCell returns the cell at the given position.
func (t *tableDefaultContent) GetCell(row, column int) *TableCell {
	if row < 0 || column < 0 || row >= len(t.cells) || column >= len(t.cells[row]) {
		return nil
	}
	return t.cells[row][column]
}

// GetRowCount returns the number of rows in the data set.
func (t *tableDefaultContent) GetRowCount() int {
	return len(t.cells)
}

// GetColumnCount returns the number of columns in the data set.
func (t *tableDefaultContent) GetColumnCount() int {
	if len(t.cells) == 0 {
		return 0
	}
	return t.lastColumn + 1
}

// Table visualizes two-dimensional data consisting of rows and columns. Each
// Table cell is defined via [Table.SetCell] by the [TableCell] type. They can
// be added dynamically to the table and changed any time.
//
// The most compact display of a table is without borders. Each row will then
// occupy one row on screen and columns are separated by the rune defined via
// [Table.SetSeparator] (a space character by default).
//
// When borders are turned on (via [Table.SetBorders]), each table cell is
// surrounded by lines. Therefore one table row will require two rows on screen.
//
// Columns will use as much horizontal space as they need. You can constrain
// their size with the [TableCell.MaxWidth] parameter of the [TableCell] type.
//
// # Fixed Columns
//
// You can define fixed rows and rolumns via [Table.SetFixed]. They will always
// stay in their place, even when the table is scrolled. Fixed rows are always
// the top rows. Fixed columns are always the leftmost columns.
//
// # Selections
//
// You can call [Table.SetSelectable] to set columns and/or rows to
// "selectable". If the flag is set only for columns, entire columns can be
// selected by the user. If it is set only for rows, entire rows can be
// selected. If both flags are set, individual cells can be selected. The
// "selected" handler set via [Table.SetSelectedFunc] is invoked when the user
// presses Enter on a selection.
//
// # Navigation
//
// If the table extends beyond the available space, it can be navigated with
// key bindings similar to Vim:
//
//   - h, left arrow: Move left by one column.
//   - l, right arrow: Move right by one column.
//   - j, down arrow: Move down by one row.
//   - k, up arrow: Move up by one row.
//   - g, home: Move to the top.
//   - G, end: Move to the bottom.
//   - Ctrl-F, page down: Move down by one page.
//   - Ctrl-B, page up: Move up by one page.
//
// When there is no selection, this affects the entire table (except for fixed
// rows and columns). When there is a selection, the user moves the selection.
// The class will attempt to keep the selection from moving out of the screen.
//
// Use [Box.SetInputCapture] to override or modify keyboard input.
//
// See https://github.com/rivo/tview/wiki/Table for an example.
type Table struct {
	*nuview.Box

	// Whether or not this table has borders around each cell.
	borders bool

	// The color of the borders or the separator.
	bordersColor tcell.Color

	// If there are no borders, the column separator.
	separator rune

	// The table's data structure.
	content tableContent

	// The number of fixed rows / columns.
	fixedRows, fixedColumns int

	// Whether or not rows or columns can be selected. If both are set to true,
	// cells can be selected.
	rowsSelectable, columnsSelectable bool

	// The currently selected row and column.
	selectedRow, selectedColumn int

	// A temporary flag which causes the next call to Draw() to force the
	// current selection to remain visible. It is set to false afterwards.
	clampToSelection bool

	// If set to true, moving the selection will wrap around horizontally (last
	// to first column and vice versa) or vertically (last to first row and vice
	// versa).
	wrapHorizontally, wrapVertically bool

	// The number of rows/columns by which the table is scrolled down/to the
	// right.
	rowOffset    int
	columnOffset int // -1 means that xScroll is being used
	xScroll      int

	// If set to true, the table's last row will always be visible.
	trackEnd bool

	// The number of visible rows the last time the table was drawn.
	visibleRows int

	// The indices of the visible columns as of the last time the table was
	// drawn.
	visibleColumnIndices []int

	// The style of the selected rows. If this value is the empty struct,
	// selected rows are simply inverted.
	selectedStyle tcell.Style

	// An optional function which gets called when the user presses Enter on a
	// selected cell. If entire rows selected, the column value is undefined.
	// Likewise for entire columns.
	selected func(row, column int)

	// An optional function which gets called when the user changes the selection.
	// If entire rows selected, the column value is undefined.
	// Likewise for entire columns.
	selectionChanged func(row, column int)

	// An optional function which gets called when the user double-clicks a
	// cell. If entire rows are selected, the column index is undefined.
	doubleClick func(row, column int)

	// An optional function which gets called when the user presses Escape, Tab,
	// or Backtab. Also when the user presses Enter if nothing is selectable.
	done func(key tcell.Key)

	lastMouseDown       time.Time
	doubleClickDuration time.Duration
	sync.RWMutex
}

// NewTable returns a new table.
func NewTable() *Table {
	t := &Table{
		Box:                 nuview.NewBox(),
		bordersColor:        nuview.Styles.GraphicsColor,
		separator:           ' ',
		doubleClickDuration: StandardDoubleClick,
		content: &tableDefaultContent{
			lastColumn: -1,
		},
	}
	return t
}

// Clear removes all table data.
func (t *Table) Clear() {
	t.Lock()
	defer t.Unlock()
	t.content.Clear()
}

// SetBorders sets whether or not each cell in the table is surrounded by a
// border.
func (t *Table) SetBorders(show bool) {
	t.Lock()
	defer t.Unlock()
	t.borders = show
}

// SetBordersColor sets the color of the cell borders.
func (t *Table) SetBordersColor(color tcell.Color) {
	t.Lock()
	defer t.Unlock()
	t.bordersColor = color
}

// SetSelectedStyle sets a specific style for selected cells. If no such style
// is set, the cell's background and text color are swapped. If a cell defines
// its own selected style, that will be used instead.
//
// To reset a previous setting to its default, make the following call:
//
//	table.SetSelectedStyle(tcell.StyleDefault)
func (t *Table) SetSelectedStyle(style tcell.Style) {
	t.Lock()
	defer t.Unlock()
	t.selectedStyle = style
}

// SetSeparator sets the character used to fill the space between two
// neighboring cells. This is a space character ' ' per default but you may
// want to set it to Borders.Vertical (or any other rune) if the column
// separation should be more visible. If cell borders are activated, this is
// ignored.
//
// Separators have the same color as borders.
func (t *Table) SetSeparator(separator rune) {
	t.Lock()
	defer t.Unlock()
	t.separator = separator
}

// SetFixed sets the number of fixed rows and columns which are always visible
// even when the rest of the cells are scrolled out of view. Rows are always the
// top-most ones. Columns are always the left-most ones.
func (t *Table) SetFixed(rows int, columns int) {
	t.Lock()
	defer t.Unlock()
	t.fixedRows = rows
	t.fixedColumns = columns
	t.columnOffset = 0
}

// SetSelectable sets the flags which determine what can be selected in a table.
// There are three selection modi:
//
//   - rows = false, columns = false: Nothing can be selected.
//   - rows = true, columns = false: Rows can be selected.
//   - rows = false, columns = true: Columns can be selected.
//   - rows = true, columns = true: Individual cells can be selected.
func (t *Table) SetSelectable(rows bool, columns bool) {
	t.Lock()
	defer t.Unlock()
	t.rowsSelectable = rows
	t.columnsSelectable = columns
}

// GetSelectable returns what can be selected in a table. Refer to
// SetSelectable() for details.
func (t *Table) GetSelectable() (rows bool, columns bool) {
	t.RLock()
	defer t.RUnlock()
	return t.rowsSelectable, t.columnsSelectable
}

// GetSelection returns the position of the current selection.
// If entire rows are selected, the column index is undefined.
// Likewise for entire columns.
func (t *Table) GetSelection() (row int, column int) {
	t.RLock()
	defer t.RUnlock()
	return t.selectedRow, t.selectedColumn
}

// Select sets the selected cell. Depending on the selection settings
// specified via SetSelectable(), this may be an entire row or column, or even
// ignored completely. The "selection changed" event is fired if such a callback
// is available (even if the selection ends up being the same as before and even
// if cells are not selectable).
func (t *Table) Select(row int, column int) {
	t.Lock()
	defer t.Unlock()
	t.selectedRow = row
	t.selectedColumn = column
	t.clampToSelection = true
	if t.selectionChanged != nil {
		t.selectionChanged(row, column)
	}
}

// SetOffset sets how many rows and columns should be skipped when drawing the
// table. This is useful for large tables that do not fit on the screen.
// Navigating a selection can change these values.
//
// Fixed rows and columns are never skipped.
func (t *Table) SetOffset(row int, column int) {
	t.Lock()
	defer t.Unlock()
	t.rowOffset = row
	t.columnOffset = column
	t.trackEnd = false
}

// GetOffset returns the current row and column offset. This indicates how many
// rows and columns the table is scrolled down and to the right.
func (t *Table) GetOffset() (row int, column int) {
	t.RLock()
	defer t.RUnlock()
	return t.rowOffset, t.columnOffset
}

// SetXScroll sets the horizontal scroll position of the table. This overrides columnOffset.
func (t *Table) SetXScroll(x int) {
	t.Lock()
	defer t.Unlock()
	t.xScroll = x
	t.columnOffset = -1 // This indicates that xScroll is being used.
	t.trackEnd = false
}

// GetXScroll returns the current horizontal scroll position of the table.
func (t *Table) GetXScroll() int {
	t.RLock()
	defer t.RUnlock()
	return t.xScroll
}

// SetSelectedFunc sets a handler which is called whenever the user presses the
// Enter key on a selected cell/row/column. The handler receives the position of
// the selection and its cell contents. If entire rows are selected, the column
// index is undefined. Likewise for entire columns.
func (t *Table) SetSelectedFunc(handler func(row int, column int)) {
	t.Lock()
	defer t.Unlock()
	t.selected = handler
}

// SetSelectionChangedFunc sets a handler which is called whenever the current
// selection changes. The handler receives the position of the new selection.
// If entire rows are selected, the column index is undefined. Likewise for
// entire columns.
func (t *Table) SetSelectionChangedFunc(handler func(row int, column int)) {
	t.Lock()
	defer t.Unlock()
	t.selectionChanged = handler
}

// SetDoneFunc sets a handler which is called whenever the user presses the
// Escape, Tab, or Backtab key. If nothing is selected, it is also called when
// user presses the Enter key (because pressing Enter on a selection triggers
// the "selected" handler set via SetSelectedFunc()).
func (t *Table) SetDoneFunc(handler func(key tcell.Key)) {
	t.Lock()
	defer t.Unlock()
	t.done = handler
}

// SetDoubleClickFunc sets a handler which is called whenever the user
// double-clicks a cell. The handler receives the position of the double-clicked
// cell. If entire rows are selected, the column index is undefined. Likewise
// for entire columns.
func (t *Table) SetDoubleClickFunc(handler func(row int, column int)) {
	t.Lock()
	defer t.Unlock()
	t.doubleClick = handler
}

// SetCell sets the content of a cell the specified position. It is ok to
// directly instantiate a TableCell object. If the cell has content, at least
// the Text and Color fields should be set.
//
// Note that setting cells in previously unknown rows and columns will
// automatically extend the internal table representation with empty TableCell
// objects, e.g. starting with a row of 100,000 will immediately create 100,000
// empty rows.
//
// To avoid unnecessary garbage collection, fill columns from left to right.
func (t *Table) SetCell(row int, column int, cell *TableCell) {
	t.Lock()
	defer t.Unlock()
	t.content.SetCell(row, column, cell)
}

// SetCellSimple calls SetCell() with the given text, left-aligned, in white.
func (t *Table) SetCellSimple(row int, column int, text string) {
	t.SetCell(row, column, NewTableCell(text))
}

// GetCell returns the contents of the cell at the specified position. A valid
// TableCell object is always returned but it will be uninitialized if the cell
// was not previously set. Such an uninitialized object will not automatically
// be inserted. Therefore, repeated calls to this function may return different
// pointers for uninitialized cells.
func (t *Table) GetCell(row int, column int) *TableCell {
	t.RLock()
	defer t.RUnlock()
	cell := t.content.GetCell(row, column)
	if cell == nil {
		cell = &TableCell{}
	}
	return cell
}

// RemoveRow removes the row at the given position from the table. If there is
// no such row, this has no effect.
func (t *Table) RemoveRow(row int) {
	t.Lock()
	defer t.Unlock()
	t.content.RemoveRow(row)
}

// RemoveColumn removes the column at the given position from the table. If
// there is no such column, this has no effect.
func (t *Table) RemoveColumn(column int) {
	t.Lock()
	defer t.Unlock()
	t.content.RemoveColumn(column)
}

// InsertRow inserts a row before the row with the given index. Cells on the
// given row and below will be shifted to the bottom by one row. If "row" is
// equal or larger than the current number of rows, this function has no effect.
func (t *Table) InsertRow(row int) {
	t.Lock()
	defer t.Unlock()
	t.content.InsertRow(row)
}

// InsertColumn inserts a column before the column with the given index. Cells
// in the given column and to its right will be shifted to the right by one
// column. Rows that have fewer initialized cells than "column" will remain
// unchanged.
func (t *Table) InsertColumn(column int) {
	t.Lock()
	defer t.Unlock()
	t.content.InsertColumn(column)
}

// GetRowCount returns the number of rows in the table.
func (t *Table) GetRowCount() int {
	t.RLock()
	defer t.RUnlock()
	return t.content.GetRowCount()
}

// GetColumnCount returns the (maximum) number of columns in the table.
func (t *Table) GetColumnCount() int {
	t.RLock()
	defer t.RUnlock()
	return t.content.GetColumnCount()
}

// CellAt returns the row and column located at the given screen coordinates.
// Each returned value may be negative if there is no row and/or cell. This
// function will also process coordinates outside the table's inner rectangle so
// callers will need to check for bounds themselves.
//
// The layout of the table when it was last drawn is used so if anything has
// changed in the meantime, the results may not be reliable.
func (t *Table) CellAt(x, y int) (row int, column int) {
	t.RLock()
	defer t.RUnlock()
	rectX, rectY, _, _ := t.GetInnerRect()

	// Determine row as seen on screen.
	if t.borders {
		row = (y - rectY - 1) / 2
	} else {
		row = y - rectY
	}

	// Respect fixed rows and row offset.
	if row >= 0 {
		if row >= t.fixedRows {
			row += t.rowOffset
		}
		if row >= t.content.GetRowCount() {
			row = -1
		}
	}

	column = -1
	columnWidths := t.calculateColumnWidths()
	relX := x - rectX
	posX := 0
	for i := 0; i < t.fixedColumns; i++ {
		posX += columnWidths[i]
		if t.borders {
			posX++ // Add space for the borders.
		}
		if relX < posX {
			column = i
			return row, column
		}
	}
	if t.borders {
		posX++ // Add space for the borders.
	}

	relX += t.effectiveXOffset(columnWidths)
	for i := t.fixedColumns; i < len(columnWidths); i++ {
		posX += columnWidths[i]
		if t.borders {
			posX++ // Add space for the borders.
		}
		if relX < posX {
			column = i
			return row, column
		}
	}

	return row, column
}

// ScrollToBeginning scrolls the table to the beginning to that the top left
// corner of the table is shown. Note that this position may be corrected if
// there is a selection.
func (t *Table) ScrollToBeginning() {
	t.Lock()
	defer t.Unlock()
	t.trackEnd = false
	if t.columnOffset != -1 {
		t.columnOffset = 0
	}
	t.xScroll = 0
	t.rowOffset = 0
}

// ScrollToEnd scrolls the table to the beginning to that the bottom left corner
// of the table is shown. Adding more rows to the table will cause it to
// automatically scroll with the new data. Note that this position may be
// corrected if there is a selection.
func (t *Table) ScrollToEnd() {
	t.Lock()
	defer t.Unlock()
	t.trackEnd = true
	t.columnOffset = 0
	t.rowOffset = t.content.GetRowCount()
}

func (t *Table) GetVisibleColumnRange() (first int, last int) {
	totalVisibleColumns := len(t.visibleColumnIndices)
	return t.visibleColumnIndices[0], t.visibleColumnIndices[totalVisibleColumns-1]
}

// SetWrapSelection determines whether a selection wraps vertically or
// horizontally when moved. Vertically wrapping selections will jump from the
// last selectable row to the first selectable row and vice versa. Horizontally
// wrapping selections will jump from the last selectable column to the first
// selectable column (on the next selectable row) or from the first selectable
// column to the last selectable column (on the previous selectable row). If set
// to false, the selection is not moved when it is already on the first/last
// selectable row/column.
//
// The default is for both values to be false.
func (t *Table) SetWrapSelection(vertical, horizontal bool) {
	t.Lock()
	defer t.Unlock()
	t.wrapHorizontally = horizontal
	t.wrapVertically = vertical
}

func (t *Table) Draw(screen tcell.Screen) {
	t.Box.Draw(screen)
	// What's our available screen space?
	// _, totalHeight := screen.Size()
	x, y, width, height := t.GetInnerRect()
	netWidth := width
	if t.borders {
		t.visibleRows = height / 2
		netWidth -= 2
	} else {
		t.visibleRows = height
	}

	screenAdapter := NewTranslateScreenWriterAdapter(screen)
	screenWriter := NewClippingScreenWriter(screenAdapter, x, y, width, height)

	// Setup selection and get table dimensions
	rowCount := t.content.GetRowCount()
	columnCount := t.content.GetColumnCount()

	t.ensureValidSelection(rowCount, columnCount)
	t.clampOffsets(height, width, rowCount, columnCount)

	// Determine visible rows
	rows, _ := t.calculateVisibleRows(height, rowCount)
	columnWidths := t.calculateColumnWidths()

	normalColumnCount := columnCount - t.fixedColumns

	xOffset := t.effectiveXOffset(columnWidths)
	fixedColumnsWidth := t.effectiveColumnsWidth(columnWidths[0:t.fixedColumns])

	t.drawCellColumnRange(screenWriter.NewClipXY(fixedColumnsWidth, 0).NewTranslate(-xOffset, 0), rows, t.fixedColumns,
		normalColumnCount, columnWidths)
	if t.fixedColumns > 0 {
		t.drawCellColumnRange(screenWriter, rows, 0, t.fixedColumns, columnWidths)
	}

	t.drawCellBackgroundColumnRange(screenWriter.NewClipXY(fixedColumnsWidth, 0).NewTranslate(-xOffset, 0), rows, t.fixedColumns,
		normalColumnCount, columnWidths)
	if t.fixedColumns > 0 {
		t.drawCellBackgroundColumnRange(screenWriter, rows, 0, t.fixedColumns, columnWidths)
	}
}

func (t *Table) effectiveXOffset(columnWidths []int) int {
	xOffset := t.xScroll
	if t.columnOffset != -1 {
		xOffset = 0
		for i := 0; i < t.columnOffset; i++ {
			xOffset += columnWidths[i+t.fixedColumns]
		}
		xOffset += t.columnOffset // Add space for the borders.
	}
	return xOffset
}

func (t *Table) MaximumXOffset() int {
	_, _, width, _ := t.GetInnerRect()
	columnWidths := t.calculateColumnWidths()
	fixedColumnsWidth := t.effectiveColumnsWidth(columnWidths[0:t.fixedColumns])
	effectiveWidth := width - fixedColumnsWidth
	normalColumnsWidth := t.effectiveColumnsWidth(columnWidths[t.fixedColumns:])
	return max(0, normalColumnsWidth-effectiveWidth+1)
}

func (t *Table) ScrollableWidth() int {
	columnWidths := t.calculateColumnWidths()
	normalColumnsWidth := t.effectiveColumnsWidth(columnWidths[t.fixedColumns:])
	return normalColumnsWidth
}

func (t *Table) ScrollableViewportWidth() int {
	columnWidths := t.calculateColumnWidths()
	fixedColumnsWidth := t.effectiveColumnsWidth(columnWidths[0:t.fixedColumns])
	_, _, width, _ := t.GetInnerRect()
	return max(0, width-fixedColumnsWidth)
}

func (t *Table) effectiveColumnsWidth(widths []int) int {
	columnsWidth := 0
	for _, width := range widths {
		columnsWidth += width
	}
	if t.borders {
		columnsWidth += len(widths) // Add space for the borders.
	}
	return columnsWidth
}

func (t *Table) drawCellColumnRange(screenWriter TranslateScreenWriter, rows []int, startColumn int, columnCount int,
	columnWidths []int) int {

	posX := 0
	if t.borders {
		for columnIndex := startColumn; columnIndex < startColumn+columnCount; columnIndex++ {
			columnWidth := columnWidths[columnIndex]

			t.drawCellColumn(screenWriter.NewTranslate(posX+1, 0), rows, columnIndex, columnWidth, 1)
			isLastColumn := columnIndex == startColumn+columnCount-1
			t.drawColumnBorders(screenWriter.NewTranslate(posX, 0), rows, columnIndex, columnWidth, isLastColumn)
			posX += columnWidth
			posX++
		}
	} else {
		for columnIndex := startColumn; columnIndex < startColumn+columnCount; columnIndex++ {
			columnWidth := columnWidths[columnIndex]
			t.drawCellColumn(screenWriter.NewTranslate(posX, 0), rows, columnIndex, columnWidth, 0)
			posX += columnWidth
			posX++
		}
	}
	return posX
}

func (t *Table) drawCellColumn(screenWriter TranslateScreenWriter, rows []int, column int,
	columnWidth int, verticalSpacing int) {

	for rowIndex, row := range rows {
		// Get the cell.
		cell := t.content.GetCell(row, column)
		if cell == nil {
			continue
		}

		rowY := verticalSpacing + ((1 + verticalSpacing) * rowIndex)

		// Draw text.
		// finalWidth := min(columnWidth, width)
		cell.x, cell.y = screenWriter.AbsolutePosition(0, rowY)

		style := cell.Style
		if style == tcell.StyleDefault {
			style = tcell.StyleDefault.Background(cell.BackgroundColor).Foreground(cell.Color).Attributes(cell.Attributes)
		}

		t.drawRectangleColorScreenWriter(screenWriter, 0, rowY, columnWidth+1, 1, style)

		start, end := printStyle(screenWriter, []byte(cell.Text), 0, rowY, columnWidth, cell.Align, style)
		printed := end - start
		if nuview.TaggedStringWidth(cell.Text)-printed > 0 && printed > 0 {
			_, _, style, _ := screenWriter.GetContent(cell.width-1, rowY)
			printStyle(screenWriter, []byte(string(nuview.SemigraphicsHorizontalEllipsis)), cell.width-1, rowY, 1, nuview.AlignLeft, style)
		}
	}
}

func (t *Table) drawColumnBorders(screenWriter ScreenWriter, rows []int, columnIndex int, columnWidth int,
	drawRightEdge bool) {

	borderStyle := tcell.StyleDefault.Background(t.GetBackgroundColor()).Foreground(t.bordersColor)

	leftJointRune := nuview.Borders.Cross
	if columnIndex == 0 {
		leftJointRune = nuview.Borders.LeftT
	}

	_, height := screenWriter.Size()
	for rowIndex, row := range rows {
		rowY := 2 * rowIndex

		topLeftJointRune := leftJointRune
		if row == 0 {
			if columnIndex == 0 {
				topLeftJointRune = nuview.Borders.TopLeft
			} else {
				topLeftJointRune = nuview.Borders.TopT
			}
		}

		screenWriter.SetContent(0, rowY, topLeftJointRune, nil, borderStyle)
		for i := range columnWidth {
			screenWriter.SetContent(i+1, rowY, nuview.Borders.Horizontal, nil, borderStyle)
		}
		if rowY+1 < height {
			screenWriter.SetContent(0, rowY+1, nuview.Borders.Vertical, nil, borderStyle)
		}

		if drawRightEdge {
			rightJoint := nuview.Borders.RightT
			if row == 0 {
				rightJoint = nuview.Borders.TopRight
			}
			screenWriter.SetContent(columnWidth+1, rowY, rightJoint, nil, borderStyle)
			screenWriter.SetContent(columnWidth+1, rowY+1, nuview.Borders.Vertical, nil, borderStyle)
		}
	}
}

func (t *Table) drawCellBackgroundColumnRange(screenWriter ScreenWriter, rows []int, startColumn int,
	columnCount int, columnWidths []int) {

	verticalSpacing := 0
	if t.borders {
		verticalSpacing = 1
	}

	if t.rowsSelectable && t.columnsSelectable {
		for _, rowIndex := range rows {
			rowSelected := rowIndex == t.selectedRow
			if rowSelected {
				columnStartX := 0
				for columnIndex := startColumn; columnIndex < startColumn+columnCount; columnIndex++ {
					columnWidth := columnWidths[columnIndex]
					if t.selectedColumn == columnIndex {
						rowY := verticalSpacing + ((1 + verticalSpacing) * (rowIndex - t.rowOffset))
						selectStyle := t.getSelectStyleForCell(rowIndex, columnIndex)
						if t.borders {
							t.drawRectangleColorScreenWriter(screenWriter, columnStartX, rowY-1, columnWidth+2, 3, selectStyle)
						} else {
							t.drawRectangleColorScreenWriter(screenWriter, columnStartX, rowY, columnWidth+1, 1, selectStyle)
						}
					}
					columnStartX += columnWidth + 1
				}
			}
		}
	} else if t.rowsSelectable {
		for _, rowIndex := range rows {
			rowSelected := rowIndex == t.selectedRow
			if rowSelected {
				rowY := verticalSpacing + ((1 + verticalSpacing) * (rowIndex - t.rowOffset))
				columnStartX := 0
				for columnIndex := startColumn; columnIndex < startColumn+columnCount; columnIndex++ {
					columnWidth := columnWidths[columnIndex]
					selectStyle := t.getSelectStyleForCell(rowIndex, columnIndex)
					if t.borders {
						t.drawRectangleColorScreenWriter(screenWriter, columnStartX, rowY-1, columnWidth+2, 3, selectStyle)
					} else {
						t.drawRectangleColorScreenWriter(screenWriter, columnStartX, rowY, columnWidth+1, 1, selectStyle)
					}
					columnStartX += columnWidth + 1
				}
			}
		}
	} else if t.columnsSelectable {
		columnStartX := 0
		for columnIndex := startColumn; columnIndex < startColumn+columnCount; columnIndex++ {
			columnWidth := columnWidths[columnIndex]
			if t.selectedColumn == columnIndex {
				for _, rowIndex := range rows {
					rowY := verticalSpacing + ((1 + verticalSpacing) * (rowIndex - t.rowOffset))
					selectStyle := t.getSelectStyleForCell(rowIndex, columnIndex)
					if t.borders {
						t.drawRectangleColorScreenWriter(screenWriter, columnStartX, rowY-1, columnWidth+2, 3, selectStyle)
					} else {
						t.drawRectangleColorScreenWriter(screenWriter, columnStartX, rowY, columnWidth+1, 1, selectStyle)
					}
				}
			}
			columnStartX += columnWidth + 1
		}
	}
}

func (t *Table) getSelectStyleForCell(rowIndex int, columnIndex int) tcell.Style {
	cell := t.content.GetCell(rowIndex, columnIndex)
	var selectStyle tcell.Style
	if cell.SelectedStyle != tcell.StyleDefault {
		selectStyle = cell.SelectedStyle
	} else if t.selectedStyle != tcell.StyleDefault {
		selectStyle = t.selectedStyle
	} else {
		textColor := cell.Color
		backgroundColor := cell.BackgroundColor
		if cell.Style != tcell.StyleDefault {
			textColor, backgroundColor, _ = cell.Style.Decompose()
		}
		// _, _, style, _ := screen.GetContent(columnStartX+1, rowY)
		// fg, bg, _ := style.Decompose()
		selectStyle = tcell.StyleDefault.Background(textColor).Foreground(backgroundColor)
	}
	return selectStyle
}

func (t *Table) drawRectangleColorScreenWriter(screenWriter ScreenWriter, x int, y int, width int, height int, style tcell.Style) {
	for row := 0; row < height; row++ {
		for col := 0; col < width; col++ {
			primaryRune, comboRune, _, _ := screenWriter.GetContent(x+col, y+row)
			screenWriter.SetContent(x+col, y+row, primaryRune, comboRune, style)
		}
	}
}

// ensureValidSelection ensures that the current selection is valid and selectable.
func (t *Table) ensureValidSelection(rowCount int, columnCount int) {
	if t.rowsSelectable || t.columnsSelectable {
		if t.selectedRow < 0 {
			t.selectedRow = 0
		}
		for t.selectedRow < rowCount {
			cell := t.content.GetCell(t.selectedRow, t.selectedColumn)
			if cell != nil && !cell.NotSelectable {
				break
			}
			t.selectedColumn++
			if t.selectedColumn > columnCount-1 {
				t.selectedColumn = 0
				t.selectedRow++
			}
		}

		if t.selectedColumn < 0 {
			t.selectedColumn = 0
		}
		if t.selectedColumn >= columnCount {
			t.selectedColumn = columnCount - 1
		}
	}
}

// clampOffsets calculates and adjusts row and column offsets based on selection and constraints.
func (t *Table) clampOffsets(height int, width int, rowCount int, columnCount int) {
	screenHeightRows := height
	if t.borders {
		screenHeightRows = height / 2 // With borders, every table row takes two screen rows.
	}

	// Clamp row offsets if requested.
	if t.clampToSelection && t.rowsSelectable {
		if t.selectedRow >= t.fixedRows && t.selectedRow < t.fixedRows+t.rowOffset {
			t.rowOffset = t.selectedRow - t.fixedRows
			t.trackEnd = false
		}
		if t.selectedRow+1-t.rowOffset >= screenHeightRows {
			t.rowOffset = t.selectedRow + 1 - screenHeightRows
			t.trackEnd = false
		}
	}
	if t.rowOffset < 0 {
		t.rowOffset = 0
	}
	if rowCount-t.rowOffset < screenHeightRows {
		t.trackEnd = true
	}

	if t.trackEnd {
		t.rowOffset = rowCount - screenHeightRows
	}
	if t.rowOffset < 0 {
		t.rowOffset = 0
	}

	if t.clampToSelection && t.columnsSelectable {
		if t.columnOffset != -1 {
			if t.selectedColumn >= t.fixedColumns && t.selectedColumn < t.fixedColumns+t.columnOffset {
				t.columnOffset = t.selectedColumn - t.fixedColumns
			}

			if t.selectedColumn >= t.fixedColumns {
				columnWidths := t.calculateColumnWidths()
				effectiveWidth := width - t.effectiveColumnsWidth(columnWidths[0:t.fixedColumns])

				maxColumnOffset := columnCount - t.fixedColumns - 1
				for {
					selectionRightEdge := t.effectiveColumnsWidth(columnWidths[t.fixedColumns+t.columnOffset : t.selectedColumn+1])
					if t.columnOffset >= maxColumnOffset || selectionRightEdge > effectiveWidth {
						if t.columnOffset >= maxColumnOffset {
							break
						}
						t.columnOffset++
					} else {
						break
					}
				}
			}
		} else {
			// If columnOffset is -1, we use xScroll.
			if t.selectedColumn >= t.fixedColumns {
				columnWidths := t.calculateColumnWidths()
				effectiveWidth := width - t.effectiveColumnsWidth(columnWidths[0:t.fixedColumns])

				left, right := t.normalColumnLeftRightPositions(columnWidths, t.selectedColumn)
				if left-t.xScroll < 0 {
					t.xScroll = left
				} else if right-t.xScroll > effectiveWidth {
					t.xScroll = -(effectiveWidth - right)
				}
			}
		}
	}

	if t.columnOffset == -1 {
		t.xScroll = min(t.MaximumXOffset(), t.xScroll)
		t.xScroll = max(0, t.xScroll)
	} else {
		// Avoid invalid column offsets.
		if t.columnOffset >= columnCount-t.fixedColumns {
			t.columnOffset = columnCount - t.fixedColumns - 1
		}
		if t.columnOffset < 0 {
			t.columnOffset = 0
		}
	}

	t.clampToSelection = false // Only once.
}

func (t *Table) normalColumnLeftRightPositions(columnWidths []int, columnIndex int) (left int, right int) {
	left = t.effectiveColumnsWidth(columnWidths[t.fixedColumns:columnIndex])
	right = t.effectiveColumnsWidth(columnWidths[t.fixedColumns : columnIndex+1])
	if t.borders {
		right++
	}
	return
}

// calculateVisibleRows determines which rows should be visible on screen.
func (t *Table) calculateVisibleRows(height int, rowCount int) (rows []int, allRows []int) {

	rowStep := 1
	if t.borders {
		rowStep = 2 // With borders, every table row takes two screen rows.
	}

	allRows = make([]int, rowCount)
	for row := 0; row < rowCount; row++ {
		allRows[row] = row
	}

	tableHeight := 0
	for row := 0; row < t.fixedRows && row < rowCount && tableHeight < height; row++ { // Do the fixed rows first.
		rows = append(rows, row)
		tableHeight += rowStep
	}

	for row := t.fixedRows + t.rowOffset; row < rowCount && tableHeight < height; row++ { // Then the remaining rows.
		rows = append(rows, row)
		tableHeight += rowStep
	}

	return rows, allRows
}

// calculateVisibleColumns determines which columns should be visible and their widths.
func (t *Table) calculateColumnWidths() []int {
	rowCount := t.content.GetRowCount()
	columnCount := t.content.GetColumnCount()

	columnWidths := make([]int, columnCount)
	for i := range columnCount {
		maxWidth := 0
		for j := range rowCount {
			if cell := t.content.GetCell(j, i); cell != nil {
				maxWidth = max(maxWidth, cell.width)
			}
		}
		columnWidths[i] = maxWidth
	}
	return columnWidths
}

// moveSelectionForward moves the selection forward, don't go beyond final cell, return
// true if a selection was found.
func (t *Table) moveSelectionForward(finalRow int, finalColumn int) bool {
	lastColumn := t.content.GetColumnCount() - 1
	rowCount := t.content.GetRowCount()
	row := t.selectedRow
	column := t.selectedColumn
	for {
		// Stop if the current selection is fine.
		cell := t.content.GetCell(row, column)
		if cell != nil && !cell.NotSelectable {
			t.selectedRow, t.selectedColumn = row, column
			return true
		}

		// If we reached the final cell, stop.
		if row == finalRow && column == finalColumn {
			return false
		}

		// Move forward.
		column++
		if column > lastColumn {
			column = 0
			row++
			if row >= rowCount {
				row = 0
			}
		}
	}
}

// moveSelectionBackwards moves the selection backwards, don't go beyond final cell, return
// true if a selection was found.
func (t *Table) moveSelectionBackwards(finalRow int, finalColumn int) bool {
	lastColumn := t.content.GetColumnCount() - 1
	rowCount := t.content.GetRowCount()
	row := t.selectedRow
	column := t.selectedColumn
	for {
		// Stop if the current selection is fine.
		cell := t.content.GetCell(row, column)
		if cell != nil && !cell.NotSelectable {
			t.selectedRow = row
			t.selectedColumn = column
			return true
		}

		// If we reached the final cell, stop.
		if row == finalRow && column == finalColumn {
			return false
		}

		// Move backwards.
		column--
		if column < 0 {
			column = lastColumn
			row--
			if row < 0 {
				row = rowCount - 1
			}
		}
	}
}

// navigateHome moves the selection to the beginning of the table.
func (t *Table) navigateHome() {
	if t.rowsSelectable {
		lastColumn := t.content.GetColumnCount() - 1
		t.selectedRow = 0
		t.selectedColumn = 0
		rowCount := t.content.GetRowCount()
		t.moveSelectionForward(rowCount-1, lastColumn)
		t.clampToSelection = true
	} else {
		t.trackEnd = false
		t.rowOffset = 0
		t.columnOffset = 0
	}
}

// navigateEnd moves the selection to the end of the table.
func (t *Table) navigateEnd() {
	if t.rowsSelectable {
		lastColumn := t.content.GetColumnCount() - 1
		rowCount := t.content.GetRowCount()
		t.selectedRow = rowCount - 1
		t.selectedColumn = lastColumn
		t.moveSelectionBackwards(0, 0)
		t.clampToSelection = true
	} else {
		t.trackEnd = true
		t.columnOffset = 0
	}
}

// navigateDown moves the selection down by one row.
func (t *Table) navigateDown() {
	if t.rowsSelectable {
		rowCount := t.content.GetRowCount()
		t.selectedRow++
		if t.selectedRow >= rowCount {
			if t.wrapVertically {
				t.selectedRow = 0
			} else {
				t.selectedRow = rowCount - 1
			}
		}
		row := t.selectedRow
		column := t.selectedColumn
		lastColumn := t.content.GetColumnCount() - 1
		finalRow, finalColumn := rowCount-1, lastColumn
		if t.wrapVertically {
			finalRow = row
			finalColumn = column
		}
		if !t.moveSelectionForward(finalRow, finalColumn) {
			t.moveSelectionBackwards(row, column)
		}
		t.clampToSelection = true
	} else {
		t.rowOffset++
	}
}

// navigateUp moves the selection up by one row.
func (t *Table) navigateUp() {
	if t.rowsSelectable {
		rowCount := t.content.GetRowCount()
		t.selectedRow--
		if t.selectedRow < 0 {
			if t.wrapVertically {
				t.selectedRow = rowCount - 1
			} else {
				t.selectedRow = 0
			}
		}
		row := t.selectedRow
		column := t.selectedColumn
		finalRow, finalColumn := 0, 0
		if t.wrapVertically {
			finalRow = row
			finalColumn = column
		}
		if !t.moveSelectionBackwards(finalRow, finalColumn) {
			t.moveSelectionForward(row, column)
		}
		t.clampToSelection = true
	} else {
		t.trackEnd = false
		t.rowOffset--
	}
}

// navigateLeft moves the selection left by one column.
func (t *Table) navigateLeft() {
	if t.columnsSelectable {
		row := t.selectedRow
		column := t.selectedColumn
		t.selectedColumn--
		if t.selectedColumn < 0 {
			if t.wrapHorizontally {
				lastColumn := t.content.GetColumnCount() - 1
				rowCount := t.content.GetRowCount()
				t.selectedColumn = lastColumn
				t.selectedRow--
				if t.selectedRow < 0 {
					if t.wrapVertically {
						t.selectedRow = rowCount - 1
					} else {
						t.selectedColumn = 0
						t.selectedRow = 0
					}
				}
			} else {
				t.selectedColumn = 0
			}
		}
		finalRow, finalColumn := row, column
		if !t.wrapHorizontally {
			finalColumn = 0
		} else if !t.wrapVertically {
			finalRow = 0
			finalColumn = 0
		}
		if !t.moveSelectionBackwards(finalRow, finalColumn) {
			t.moveSelectionForward(row, column)
		}
		t.clampToSelection = true
	} else {
		t.columnOffset = max(0, t.columnOffset-1)
	}
}

// navigateRight moves the selection right by one column.
func (t *Table) navigateRight() {
	if t.columnsSelectable {
		row := t.selectedRow
		column := t.selectedColumn
		t.selectedColumn++
		lastColumn := t.content.GetColumnCount() - 1
		rowCount := t.content.GetRowCount()
		if t.selectedColumn > lastColumn {
			if t.wrapHorizontally {
				t.selectedColumn = 0
				t.selectedRow++
				if t.selectedRow >= rowCount {
					if t.wrapVertically {
						t.selectedRow = 0
					} else {
						t.selectedColumn = lastColumn
						t.selectedRow = rowCount - 1
					}
				}
			} else {
				t.selectedColumn = lastColumn
			}
		}
		finalRow := row
		finalColumn := column
		if !t.wrapHorizontally {
			finalColumn = lastColumn
		} else if !t.wrapVertically {
			finalRow = rowCount - 1
			finalColumn = lastColumn
		}
		if !t.moveSelectionForward(finalRow, finalColumn) {
			t.moveSelectionBackwards(row, column)
		}
		t.clampToSelection = true
	} else {
		maxColumnOffset := t.content.GetColumnCount() - t.fixedColumns - 1
		t.columnOffset = min(t.columnOffset+1, maxColumnOffset)
	}
}

// navigatePageDown moves the selection down by one page.
func (t *Table) navigatePageDown() {
	offsetAmount := t.visibleRows - t.fixedRows
	if offsetAmount < 0 {
		offsetAmount = 0
	}
	if t.rowsSelectable {
		row := t.selectedRow
		column := t.selectedColumn
		t.selectedRow += offsetAmount
		rowCount := t.content.GetRowCount()
		if t.selectedRow >= rowCount {
			t.selectedRow = rowCount - 1
		}
		lastColumn := t.content.GetColumnCount() - 1
		finalRow := rowCount - 1
		finalColumn := lastColumn
		if !t.moveSelectionForward(finalRow, finalColumn) {
			t.moveSelectionBackwards(row, column)
		}
		t.clampToSelection = true
	} else {
		t.rowOffset += offsetAmount
	}
}

// navigatePageUp moves the selection up by one page.
func (t *Table) navigatePageUp() {
	offsetAmount := t.visibleRows - t.fixedRows
	if offsetAmount < 0 {
		offsetAmount = 0
	}
	if t.rowsSelectable {
		row := t.selectedRow
		column := t.selectedColumn
		t.selectedRow -= offsetAmount
		if t.selectedRow < 0 {
			t.selectedRow = 0
		}
		finalRow := 0
		finalColumn := 0
		if !t.moveSelectionBackwards(finalRow, finalColumn) {
			t.moveSelectionForward(row, column)
		}
		t.clampToSelection = true
	} else {
		t.trackEnd = false
		t.rowOffset -= offsetAmount
	}
}

// InputHandler returns the handler for this primitive.
func (t *Table) InputHandler() func(event *tcell.EventKey, setFocus func(p nuview.Primitive)) {
	return t.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p nuview.Primitive)) {
		key := event.Key()

		if (!t.rowsSelectable && !t.columnsSelectable && key == tcell.KeyEnter) ||
			key == tcell.KeyEscape ||
			key == tcell.KeyTab ||
			key == tcell.KeyBacktab {
			if t.done != nil {
				t.done(key)
			}
			return
		}

		// Movement functions.
		previouslySelectedRow, previouslySelectedColumn := t.selectedRow, t.selectedColumn
		if t.content.GetRowCount() == 0 {
			return // No movement on empty tables.
		}

		switch key {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'g':
				t.navigateHome()
			case 'G':
				t.navigateEnd()
			case 'j':
				t.navigateDown()
			case 'k':
				t.navigateUp()
			case 'h':
				t.navigateLeft()
			case 'l':
				t.navigateRight()
			}
		case tcell.KeyHome:
			t.navigateHome()
		case tcell.KeyEnd:
			t.navigateEnd()
		case tcell.KeyUp:
			t.navigateUp()
		case tcell.KeyDown:
			t.navigateDown()
		case tcell.KeyLeft:
			t.navigateLeft()
		case tcell.KeyRight:
			t.navigateRight()
		case tcell.KeyPgDn, tcell.KeyCtrlF:
			t.navigatePageDown()
		case tcell.KeyPgUp, tcell.KeyCtrlB:
			t.navigatePageUp()
		case tcell.KeyEnter:
			if (t.rowsSelectable || t.columnsSelectable) && t.selected != nil {
				t.selected(t.selectedRow, t.selectedColumn)
			}
		}

		// If the selection has changed, notify the handler.
		if t.selectionChanged != nil &&
			(t.rowsSelectable && previouslySelectedRow != t.selectedRow ||
				t.columnsSelectable && previouslySelectedColumn != t.selectedColumn) {
			t.selectionChanged(t.selectedRow, t.selectedColumn)
		}
	})
}

// MouseHandler returns the mouse handler for this primitive.
func (t *Table) MouseHandler() func(action nuview.MouseAction, event *tcell.EventMouse, setFocus func(p nuview.Primitive)) (consumed bool, capture nuview.Primitive) {
	return t.WrapMouseHandler(func(action nuview.MouseAction, event *tcell.EventMouse, setFocus func(p nuview.Primitive)) (consumed bool, capture nuview.Primitive) {
		x, y := event.Position()
		if !t.InRect(x, y) {
			return false, nil
		}

		switch action {
		case nuview.MouseLeftDown:
			setFocus(t)

			selectEvent := true
			row, column := t.CellAt(x, y)
			cell := t.content.GetCell(row, column)
			if cell != nil && cell.Clicked != nil {
				if noSelect := cell.Clicked(); noSelect {
					selectEvent = false
				}
			}
			isAlreadySelected := t.selectedRow == row && t.selectedColumn == column
			if !isAlreadySelected && selectEvent && (t.rowsSelectable || t.columnsSelectable) {
				t.Select(row, column)
			}

			if isAlreadySelected {
				now := time.Now()
				if !t.lastMouseDown.IsZero() && (now.Sub(t.lastMouseDown) < t.doubleClickDuration) {
					// Double-click: Notify the handler.
					if t.doubleClick != nil {
						t.doubleClick(row, column)
					}
				}
			}
			t.lastMouseDown = time.Now()
			consumed = true

		case nuview.MouseScrollUp:
			t.trackEnd = false
			t.rowOffset--
			consumed = true

		case nuview.MouseScrollDown:
			t.rowOffset++
			consumed = true
		}

		return
	})
}

// PrintStyle works like Print() but it takes a style instead of just a
// foreground color.
func printStyle(screen ScreenWriter, text []byte, x, y, maxWidth, align int, style tcell.Style) (int, int) {
	if maxWidth <= 0 || len(text) == 0 {
		return 0, 0
	}

	// Decompose the text.
	colorIndices, colors, _, _, escapeIndices, strippedText, strippedWidth := decomposeText(text, true, false)

	// We want to reduce all alignments to AlignLeft.
	if align == nuview.AlignRight {
		if strippedWidth <= maxWidth {
			// There's enough space for the entire text.
			return printStyle(screen, text, x+maxWidth-strippedWidth, y, maxWidth, nuview.AlignLeft, style)
		}
		// Trim characters off the beginning.
		var (
			bytes, width, colorPos, escapePos, tagOffset int
			foregroundColor, backgroundColor, attributes string
		)
		_, originalBackground, _ := style.Decompose()
		iterateString(string(strippedText), func(main rune, comb []rune, textPos, textWidth, screenPos, screenWidth int) bool {
			// Update color/escape tag offset and style.
			if colorPos < len(colorIndices) && textPos+tagOffset >= colorIndices[colorPos][0] && textPos+tagOffset < colorIndices[colorPos][1] {
				foregroundColor, backgroundColor, attributes = styleFromTag(foregroundColor, backgroundColor, attributes, colors[colorPos])
				style = overlayStyle(originalBackground, style, foregroundColor, backgroundColor, attributes)
				tagOffset += colorIndices[colorPos][1] - colorIndices[colorPos][0]
				colorPos++
			}
			if escapePos < len(escapeIndices) && textPos+tagOffset >= escapeIndices[escapePos][0] && textPos+tagOffset < escapeIndices[escapePos][1] {
				tagOffset++
				escapePos++
			}
			if strippedWidth-screenPos < maxWidth {
				// We chopped off enough.
				if escapePos > 0 && textPos+tagOffset-1 >= escapeIndices[escapePos-1][0] && textPos+tagOffset-1 < escapeIndices[escapePos-1][1] {
					// Unescape open escape sequences.
					escapeCharPos := escapeIndices[escapePos-1][1] - 2
					text = append(text[:escapeCharPos], text[escapeCharPos+1:]...)
				}
				// Print and return.
				bytes, width = printStyle(screen, text[textPos+tagOffset:], x, y, maxWidth, nuview.AlignLeft, style)
				return true
			}
			return false
		})
		return bytes, width
	} else if align == nuview.AlignCenter {
		if strippedWidth == maxWidth {
			// Use the exact space.
			return printStyle(screen, text, x, y, maxWidth, nuview.AlignLeft, style)
		} else if strippedWidth < maxWidth {
			// We have more space than we need.
			half := (maxWidth - strippedWidth) / 2
			return printStyle(screen, text, x+half, y, maxWidth-half, nuview.AlignLeft, style)
		} else {
			// Chop off runes until we have a perfect fit.
			var choppedLeft, choppedRight, leftIndex, rightIndex int
			rightIndex = len(strippedText)
			for rightIndex-1 > leftIndex && strippedWidth-choppedLeft-choppedRight > maxWidth {
				if choppedLeft < choppedRight {
					// Iterate on the left by one character.
					iterateString(string(strippedText[leftIndex:]), func(main rune, comb []rune, textPos, textWidth, screenPos, screenWidth int) bool {
						choppedLeft += screenWidth
						leftIndex += textWidth
						return true
					})
				} else {
					// Iterate on the right by one character.
					iterateStringReverse(string(strippedText[leftIndex:rightIndex]), func(main rune, comb []rune, textPos, textWidth, screenPos, screenWidth int) bool {
						choppedRight += screenWidth
						rightIndex -= textWidth
						return true
					})
				}
			}

			// Add tag offsets and determine start style.
			var (
				colorPos, escapePos, tagOffset               int
				foregroundColor, backgroundColor, attributes string
			)
			_, originalBackground, _ := style.Decompose()
			for index := range strippedText {
				// We only need the offset of the left index.
				if index > leftIndex {
					// We're done.
					if escapePos > 0 && leftIndex+tagOffset-1 >= escapeIndices[escapePos-1][0] && leftIndex+tagOffset-1 < escapeIndices[escapePos-1][1] {
						// Unescape open escape sequences.
						escapeCharPos := escapeIndices[escapePos-1][1] - 2
						text = append(text[:escapeCharPos], text[escapeCharPos+1:]...)
					}
					break
				}

				// Update color/escape tag offset.
				if colorPos < len(colorIndices) && index+tagOffset >= colorIndices[colorPos][0] && index+tagOffset < colorIndices[colorPos][1] {
					if index <= leftIndex {
						foregroundColor, backgroundColor, attributes = styleFromTag(foregroundColor, backgroundColor, attributes, colors[colorPos])
						style = overlayStyle(originalBackground, style, foregroundColor, backgroundColor, attributes)
					}
					tagOffset += colorIndices[colorPos][1] - colorIndices[colorPos][0]
					colorPos++
				}
				if escapePos < len(escapeIndices) && index+tagOffset >= escapeIndices[escapePos][0] && index+tagOffset < escapeIndices[escapePos][1] {
					tagOffset++
					escapePos++
				}
			}
			return printStyle(screen, text[leftIndex+tagOffset:], x, y, maxWidth, nuview.AlignLeft, style)
		}
	}

	// Draw text.
	var (
		drawn, drawnWidth, colorPos, escapePos, tagOffset int
		foregroundColor, backgroundColor, attributes      string
	)
	iterateString(string(strippedText), func(main rune, comb []rune, textPos, length, screenPos, screenWidth int) bool {
		// Only continue if there is still space.
		if drawnWidth+screenWidth > maxWidth {
			return true
		}

		// Handle color tags.
		for colorPos < len(colorIndices) && textPos+tagOffset >= colorIndices[colorPos][0] && textPos+tagOffset < colorIndices[colorPos][1] {
			foregroundColor, backgroundColor, attributes = styleFromTag(foregroundColor, backgroundColor, attributes, colors[colorPos])
			tagOffset += colorIndices[colorPos][1] - colorIndices[colorPos][0]
			colorPos++
		}

		// Handle scape tags.
		if escapePos < len(escapeIndices) && textPos+tagOffset >= escapeIndices[escapePos][0] && textPos+tagOffset < escapeIndices[escapePos][1] {
			if textPos+tagOffset == escapeIndices[escapePos][1]-2 {
				tagOffset++
				escapePos++
			}
		}

		// Print the rune sequence.
		finalX := x + drawnWidth
		_, _, finalStyle, _ := screen.GetContent(finalX, y)
		_, background, _ := finalStyle.Decompose()
		finalStyle = overlayStyle(background, style, foregroundColor, backgroundColor, attributes)
		for offset := screenWidth - 1; offset >= 0; offset-- {
			// To avoid undesired effects, we populate all cells.
			if offset == 0 {
				screen.SetContent(finalX+offset, y, main, comb, finalStyle)
			} else {
				screen.SetContent(finalX+offset, y, ' ', nil, finalStyle)
			}
		}

		// Advance.
		drawn += length
		drawnWidth += screenWidth

		return false
	})

	return drawn + tagOffset + len(escapeIndices), drawnWidth
}
