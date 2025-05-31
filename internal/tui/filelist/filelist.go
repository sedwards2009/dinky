package filelist

import (
	"dinky/internal/tui/scrollbar"
	"os"
	"slices"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type FileList struct {
	*tview.Flex
	app               *tview.Application
	table             *tview.Table
	scrollbar         *scrollbar.Scrollbar
	path              string
	entries           []os.DirEntry
	dirRequestsChan   chan string
	columnDescriptors []columnDescriptor

	sortColumn    int // Index of the currently sorted column
	sortDirection int // 1 for ascending, -1 for descending
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
				if entry.IsDir() {
					return "📁 " + entry.Name() + "/"
				}
				return entry.Name()
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

	table := tview.NewTable()
	table.SetSelectable(true, false)
	table.SetBorder(false)
	table.SetFixed(1, 0)
	table.SetEvaluateAllRows(true)

	topFlex.AddItem(table, 0, 1, false)

	fileListScrollbar := scrollbar.NewScrollbar()
	topFlex.AddItem(fileListScrollbar, 1, 0, false)

	dirRequestsChan := make(chan string, 10)

	fileList := &FileList{
		app:               app,
		Flex:              topFlex,
		table:             table,
		scrollbar:         fileListScrollbar,
		path:              "/home/sbe",
		columnDescriptors: columnDescriptors,
		dirRequestsChan:   dirRequestsChan,
		sortColumn:        0,
		sortDirection:     1,
	}

	fileListScrollbar.Track.SetBeforeDrawFunc(func(_ tcell.Screen) {
		row, _ := table.GetOffset()
		fileListScrollbar.Track.SetMax(table.GetRowCount() - 1)
		_, _, _, height := fileList.table.GetInnerRect()
		fileListScrollbar.Track.SetThumbSize(height)
		fileListScrollbar.Track.SetPosition(row)
	})
	fileListScrollbar.SetChangedFunc(func(position int) {
		_, column := fileList.table.GetOffset()
		fileList.table.SetOffset(position, column)
	})

	for i, desc := range columnDescriptors {
		cell := &tview.TableCell{
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
		table.SetCell(0, i, cell)
	}
	fileList.updateColumnHeaders()

	go fileList.runDirectoryLister(dirRequestsChan)

	return fileList
}

func (fileList *FileList) SetPath(path string) {
	fileList.dirRequestsChan <- path
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
			fileList.app.QueueUpdateDraw(func() {
				fileList.setEntries(entries, dirPath)
			})
		}
	}
}

func (fileList *FileList) setEntries(entries []os.DirEntry, dirPath string) {
	fileList.entries = entries
	fileList.sortEntries(entries, fileList.sortColumn, fileList.sortDirection)
	fileList.loadEntries(entries)
}

func (fileList *FileList) SetSortColumn(sortColumn int, sortDirection int) {
	fileList.sortColumn = sortColumn
	fileList.sortDirection = sortDirection
	fileList.sortEntries(fileList.entries, sortColumn, sortDirection)
	fileList.loadEntries(fileList.entries)
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
	for i, desc := range fileList.columnDescriptors {
		cell := fileList.table.GetCell(0, i)
		cell.SetText(desc.name)
	}

	cell := fileList.table.GetCell(0, fileList.sortColumn)

	if fileList.sortDirection == 1 {
		cell.SetText(cell.Text + " ▼")
	} else {
		cell.SetText(cell.Text + " ▲")
	}
}

func (fileList *FileList) loadEntries(entries []os.DirEntry) {
	for fileList.table.GetRowCount() > 1 {
		fileList.table.RemoveRow(fileList.table.GetRowCount() - 1)
	}

	for i, entry := range entries {
		for j, desc := range fileList.columnDescriptors {
			cell := &tview.TableCell{
				Text:  desc.formatFunc(entry),
				Color: tview.Styles.PrimaryTextColor,
				Align: desc.align,
			}
			fileList.table.SetCell(i+1, j, cell)
		}
	}
}
