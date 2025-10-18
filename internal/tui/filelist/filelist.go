package filelist

import (
	"dinky/internal/tui/scrollbar"
	"dinky/internal/tui/stylecolor"
	"dinky/internal/tui/table2"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type FileList struct {
	*tview.Flex
	app                 *tview.Application
	table               *table2.Table
	VerticalScrollbar   *scrollbar.Scrollbar
	HorizontalScrollbar *scrollbar.Scrollbar
	path                string

	allEntries     []os.DirEntry
	visibleEntries []os.DirEntry

	dirRequestsChan   chan string
	columnDescriptors []columnDescriptor
	showHidden        bool

	sortColumn    int // Index of the currently sorted column
	sortDirection int // 1 for ascending, -1 for descending

	changedFunc  func(path string, entry os.DirEntry)
	selectedFunc func(path string, entry os.DirEntry)

	textColor               tcell.Color
	backgroundColor         tcell.Color
	selectedBackgroundColor tcell.Color
	headerLabelColor        tcell.Color
	headerBackgroundColor   tcell.Color
}

type columnDescriptor struct {
	name       string
	align      int
	formatFunc func(entry os.DirEntry) string
	sortFunc   func(a os.DirEntry, b os.DirEntry) int
}

func NewFileList(app *tview.Application) *FileList {
	columnDescriptors := []columnDescriptor{
		{
			name:  "Name",
			align: tview.AlignLeft,
			formatFunc: func(entry os.DirEntry) string {
				return emojiForFileType(entry) + " " + entry.Name()
			},
			sortFunc: sortNameFunc,
		},
		{
			name:  "Size",
			align: tview.AlignRight,
			formatFunc: func(entry os.DirEntry) string {
				info, err := entry.Info()
				if err != nil {
					return "?"
				}
				return formatSize(info.Size())
			},
			sortFunc: sortSizeFunc,
		},
		{
			name:  "Modified",
			align: tview.AlignLeft,
			formatFunc: func(entry os.DirEntry) string {
				info, err := entry.Info()
				if err != nil {
					return "?"
				}
				return info.ModTime().Format("2006-01-02 15:04:05")
			},
			sortFunc: sortModifiedFunc,
		},
		{
			name:       "Permissions",
			align:      tview.AlignLeft,
			formatFunc: permissions,
			sortFunc:   sortPermissionsFunc,
		},
		{
			name:       "Owner",
			align:      tview.AlignLeft,
			formatFunc: ownerName,
			sortFunc:   sortOwnerFunc,
		},
		{
			name:       "Group",
			align:      tview.AlignLeft,
			formatFunc: groupName,
			sortFunc:   sortGroupFunc,
		},
	}

	topFlex := tview.NewFlex()
	topFlex.SetDirection(tview.FlexColumn)
	topFlex.SetBorder(false)

	table := table2.NewTable()
	table.SetSelectable(true, false)
	table.SetBorder(false)
	table.SetFixed(1, 0)

	middleFlex := tview.NewFlex()
	middleFlex.SetDirection(tview.FlexRow)
	middleFlex.SetBorder(false)
	middleFlex.AddItem(table, 0, 1, true)
	horizontalScrollbar := scrollbar.NewScrollbar()
	horizontalScrollbar.SetHorizontal(true)
	middleFlex.AddItem(horizontalScrollbar, 1, 0, false)

	topFlex.AddItem(middleFlex, 0, 1, false)

	verticalScrollbar := scrollbar.NewScrollbar()
	topFlex.AddItem(verticalScrollbar, 1, 0, false)

	dirRequestsChan := make(chan string, 10)

	fileList := &FileList{
		app:                 app,
		Flex:                topFlex,
		table:               table,
		VerticalScrollbar:   verticalScrollbar,
		HorizontalScrollbar: horizontalScrollbar,
		path:                "/",
		columnDescriptors:   columnDescriptors,
		dirRequestsChan:     dirRequestsChan,
		sortColumn:          0,
		sortDirection:       1,

		textColor:               tview.Styles.PrimaryTextColor,
		backgroundColor:         tview.Styles.PrimitiveBackgroundColor,
		selectedBackgroundColor: tview.Styles.ContrastBackgroundColor,
		headerLabelColor:        stylecolor.ButtonLabelColor,
		headerBackgroundColor:   stylecolor.ButtonBackgroundColor,
	}
	fileList.table.SetXScroll(0)

	verticalScrollbar.Track.SetBeforeDrawFunc(func(_ tcell.Screen) {
		row, _ := table.GetOffset()
		verticalScrollbar.Track.SetMax(table.GetRowCount() - 1)
		_, _, _, height := fileList.table.GetInnerRect()
		verticalScrollbar.Track.SetThumbSize(height)
		verticalScrollbar.Track.SetPosition(row)
	})
	verticalScrollbar.SetChangedFunc(func(position int) {
		_, column := fileList.table.GetOffset()
		fileList.table.SetOffset(position, column)
	})

	horizontalScrollbar.Track.SetBeforeDrawFunc(func(_ tcell.Screen) {
		horizontalScrollbar.Track.SetMax(table.ScrollableWidth())
		horizontalScrollbar.Track.SetThumbSize(table.ScrollableViewportWidth())
		horizontalScrollbar.Track.SetPosition(table.GetXScroll())
	})
	horizontalScrollbar.SetChangedFunc(func(position int) {
		fileList.table.SetXScroll(position)
	})

	notifyFunc := func(consumerFunc func(path string, entry os.DirEntry), row int) {
		if consumerFunc != nil && row > 0 {
			entryPath := filepath.Join(fileList.path, fileList.visibleEntries[row-1].Name())
			consumerFunc(entryPath, fileList.visibleEntries[row-1])
		}
	}
	table.SetSelectionChangedFunc(func(row int, _ int) {
		notifyFunc(fileList.changedFunc, row)
	})
	table.SetSelectedFunc(func(row int, _ int) {
		notifyFunc(fileList.selectedFunc, row)
	})
	table.SetDoubleClickFunc(func(row int, _ int) {
		notifyFunc(fileList.selectedFunc, row)
	})

	fileList.loadColumnHeaders()
	go fileList.runDirectoryLister(dirRequestsChan)

	return fileList
}

func (fileList *FileList) loadColumnHeaders() {
	for i, desc := range fileList.columnDescriptors {
		cell := &table2.TableCell{
			Text:          desc.name,
			NotSelectable: true,
			Clicked: func() bool {
				newSortDirection := fileList.sortDirection
				if i == fileList.sortColumn {
					newSortDirection *= -1 // Toggle sort direction
				}
				fileList.SetSortColumn(i, newSortDirection)
				return true
			},
		}
		fileList.table.SetCell(0, i, cell)
	}
	fileList.updateColumnHeaders()
}

func (fileList *FileList) SetPath(path string) {
	fileList.dirRequestsChan <- path
}

func (fileList *FileList) Path() string {
	return fileList.path
}

func (fileList *FileList) runDirectoryLister(dirRequests chan string) {
	for {
		dirPath, ok := <-dirRequests
		if !ok || dirPath == "" {
			return
		}
		entries, err := os.ReadDir(dirPath)
		if err != nil {

		} else {
			if dirPath != "/" {
				parentEntry, err := parentDirEntry(dirPath)
				if err == nil {
					entries = append([]os.DirEntry{aliasDirEntry(parentEntry, "..")}, entries...)
				}
			}
			fileList.app.QueueUpdateDraw(func() {
				fileList.setEntries(entries, dirPath)
			})
		}
	}
}

func (fileList *FileList) setEntries(entries []os.DirEntry, dirPath string) {
	fileList.path = dirPath
	fileList.allEntries = entries
	fileList.visibleEntries = entries

	if !fileList.showHidden {
		fileList.visibleEntries = fileList.filterVisible(entries)
	}

	fileList.sortEntries(fileList.visibleEntries, fileList.sortColumn, fileList.sortDirection)
	fileList.loadEntries(fileList.visibleEntries)
}

func (fileList *FileList) filterVisible(entries []os.DirEntry) []os.DirEntry {
	var visibleEntries []os.DirEntry
	for _, entry := range entries {
		name := entry.Name()
		if name == ".." || !strings.HasPrefix(name, ".") {
			visibleEntries = append(visibleEntries, entry)
		}
	}
	return visibleEntries
}

func (fileList *FileList) SetSortColumn(sortColumn int, sortDirection int) {
	fileList.sortColumn = sortColumn
	fileList.sortDirection = sortDirection
	fileList.sortEntries(fileList.visibleEntries, sortColumn, sortDirection)
	fileList.loadEntries(fileList.visibleEntries)
	fileList.updateColumnHeaders()
}

func (fileList *FileList) sortEntries(entries []os.DirEntry, sortColumn int, sortDirection int) {
	reverse := func(sortFunc func(a os.DirEntry, b os.DirEntry) int) func(a os.DirEntry, b os.DirEntry) int {
		return sortFunc
	}

	if sortDirection == -1 {
		reverse = func(sortFunc func(a os.DirEntry, b os.DirEntry) int) func(a os.DirEntry, b os.DirEntry) int {
			return func(a os.DirEntry, b os.DirEntry) int {
				return -sortFunc(a, b)
			}
		}
	}

	slices.SortFunc(entries, reverse(fileList.columnDescriptors[sortColumn].sortFunc))
}

func (fileList *FileList) updateColumnHeaders() {
	cellStyle := tcell.StyleDefault.Foreground(fileList.headerLabelColor).Background(fileList.headerBackgroundColor)
	for i, desc := range fileList.columnDescriptors {
		cell := fileList.table.GetCell(0, i)
		cell.SetText(desc.name)
		cell.SetStyle(cellStyle)
	}

	cell := fileList.table.GetCell(0, fileList.sortColumn)

	if fileList.sortDirection == 1 {
		cell.SetText(string(cell.Text) + " ▼")
	} else {
		cell.SetText(string(cell.Text) + " ▲")
	}
}

func (fileList *FileList) loadEntries(entries []os.DirEntry) {
	for fileList.table.GetRowCount() > 1 {
		fileList.table.RemoveRow(fileList.table.GetRowCount() - 1)
	}

	for i, entry := range entries {
		for j, desc := range fileList.columnDescriptors {
			cell := &table2.TableCell{
				Text:  desc.formatFunc(entry),
				Style: tcell.StyleDefault.Foreground(fileList.textColor).Background(fileList.backgroundColor),
				Align: desc.align,
			}
			fileList.table.SetCell(i+1, j, cell)
		}
	}
	if len(entries) > 0 {
		fileList.table.Select(1, 0)
	}
	fileList.table.SetOffset(0, 0)
}

func (fileList *FileList) SetTextColor(color tcell.Color) {
	fileList.textColor = color
}

func (fileList *FileList) SetBackgroundColor(color tcell.Color) {
	fileList.backgroundColor = color
	fileList.table.SetBackgroundColor(color)
	fileList.Box.SetBackgroundColor(color)
}

func (fileList *FileList) SetSelectedBackgroundColor(color tcell.Color) {
	fileList.selectedBackgroundColor = color
	fileList.table.SetSelectedStyle(tcell.StyleDefault.Background(color).Foreground(fileList.textColor))
}

func (fileList *FileList) SetHeaderLabelColor(color tcell.Color) {
	fileList.headerLabelColor = color
	fileList.updateColumnHeaders()
}

func (fileList *FileList) SetHeaderBackgroundColor(color tcell.Color) {
	fileList.headerBackgroundColor = color
	fileList.updateColumnHeaders()
}

func (fileList *FileList) SetChangedFunc(changedFunc func(path string, entry os.DirEntry)) {
	fileList.changedFunc = changedFunc
}

func (fileList *FileList) SetSelectedFunc(selectedFunc func(path string, entry os.DirEntry)) {
	fileList.selectedFunc = selectedFunc
}

func (fileList *FileList) SetShowHidden(showHidden bool) {
	if fileList.showHidden == showHidden {
		return
	}

	fileList.showHidden = showHidden
	fileList.setEntries(fileList.allEntries, fileList.path)
}

func (fileList *FileList) Focus(delegate func(p tview.Primitive)) {
	delegate(fileList.table)
}
