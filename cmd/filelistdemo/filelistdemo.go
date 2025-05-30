package main

import (
	"dinky/internal/tui/style"
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type fileListItem struct {
	Name string
	Size int
	Type string
}

var fakeFiles []fileListItem

var sortColumn = 0
var sortDirection = 1 // 1 for ascending, -1 for descending

func sortFiles(columnIndex int, sortDirection int) {
	reverse := func(sortFunc func(a fileListItem, b fileListItem) int) func(a fileListItem, b fileListItem) int {
		return sortFunc
	}

	if sortDirection == -1 {
		reverse = func(sortFunc func(a fileListItem, b fileListItem) int) func(a fileListItem, b fileListItem) int {
			return func(a fileListItem, b fileListItem) int {
				return -sortFunc(a, b)
			}
		}
	}

	switch columnIndex {
	case 0: // Name
		slices.SortFunc(fakeFiles, reverse(func(a fileListItem, b fileListItem) int {
			if a.Name < b.Name {
				return -1
			} else if a.Name > b.Name {
				return 1
			}
			return 0
		}))
	case 1: // Size
		slices.SortFunc(fakeFiles, reverse(func(a fileListItem, b fileListItem) int {
			if a.Size < b.Size {
				return -1
			} else if a.Size > b.Size {
				return 1
			}
			return 0
		}))
	case 2: // Type
		slices.SortFunc(fakeFiles, reverse(func(a fileListItem, b fileListItem) int {
			if a.Type < b.Type {
				return -1
			} else if a.Type > b.Type {
				return 1
			}
			return 0
		}))
	}
}

func loadFiles(fileListWidget *tview.Table) {
	for fileListWidget.GetRowCount() > 1 {
		fileListWidget.RemoveRow(fileListWidget.GetRowCount() - 1)
	}

	for i, item := range fakeFiles {
		nameCell := &tview.TableCell{
			Text:  item.Name,
			Color: tcell.ColorWhite,
		}
		fileListWidget.SetCell(i+1, 0, nameCell)

		sizeCell := &tview.TableCell{
			Text:  fmt.Sprintf("%d", item.Size),
			Color: tcell.ColorWhite,
			Align: tview.AlignRight,
		}
		fileListWidget.SetCell(i+1, 1, sizeCell)

		typeCell := &tview.TableCell{
			Text:  item.Type,
			Color: tcell.ColorWhite,
		}
		fileListWidget.SetCell(i+1, 2, typeCell)
	}
}

func updateHeader(fileListWidget *tview.Table, sortColumn int, sortDirection int) {
	for i, item := range []string{"Name", "Size", "Type"} {
		cell := fileListWidget.GetCell(0, i)
		cell.SetText(item)
	}

	cell := fileListWidget.GetCell(0, sortColumn)

	if sortDirection == 1 {
		cell.SetText(cell.Text + " ▼")
	} else {
		cell.SetText(cell.Text + " ▲")
	}
}

func setSort(fileListWidget *tview.Table, columnIndex int, direction int) {
	sortColumn = columnIndex
	sortDirection = direction
	sortFiles(sortColumn, direction)
	updateHeader(fileListWidget, sortColumn, direction)
	loadFiles(fileListWidget)
}

func main() {
	logFile := setupLogging()
	defer logFile.Close()

	style.Init()

	app := tview.NewApplication()
	app.EnableMouse(true)
	log.Println("Starting Filelist Demo...")

	fileList := tview.NewTable()
	fileList.SetSelectable(true, false)
	fileList.SetBorder(false)
	fileList.SetFixed(1, 0)
	fileList.SetEvaluateAllRows(true)

	for i, item := range []string{"Name", "Size", "Type"} {
		cell := &tview.TableCell{
			Text:          item,
			NotSelectable: true,
			Clicked: func() bool {
				newSortDirection := sortDirection
				if i == sortColumn {
					newSortDirection *= -1 // Toggle sort direction
				}
				setSort(fileList, i, newSortDirection)
				return true
			},
		}
		fileList.SetCell(0, i, cell)
	}

	fakeFiles = []fileListItem{
		{".gitignore", 1234, "Text"},
		{"cmd/", 1, "Directory"},
		{"dinky", 8283935, "Executable"},
		{"file1.txt", 1234, "Text"},
		{"file2.txt", 1234, "Text"},
		{"file3.txt", 1234, "Text"},
		{"file4.txt", 1234, "Text"},
		{"file5.txt", 1234, "Text"},
		{"file6.txt", 1234, "Text"},
		{"file7.txt", 1234, "Text"},
		{"file8.txt", 1234, "Text"},
		{"filedialogdemo", 4769103, "Executable"},
		{".git", 4096, "Directory"},
		{".gitignore", 507, "Text"},
		{"LICENSE", 2234, "Text"},
		{"Makefile", 1234, "Text"},
		{"README.md", 36, "Markdown"},
		{"README2.md", 36, "Markdown"},
		{"README3.md", 36, "Markdown"},
		{"README4.md", 36, "Markdown"},
		{"README5.md", 36, "Markdown"},
		{"README6.md", 36, "Markdown"},
	}
	setSort(fileList, 0, 1)

	layout := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(fileList, 40, 0, true).
			AddItem(nil, 0, 1, false).
			SetDirection(tview.FlexColumn),
			20, 0, true).
		AddItem(nil, 0, 1, false).
		SetDirection(tview.FlexRow)

	app.SetRoot(layout, true)

	if err := app.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}

func setupLogging() *os.File {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic("Failed to open log file: " + err.Error())
	}
	log.SetOutput(logFile)
	return logFile
}
