package filedialog

import (
	"dinky/internal/tui/style"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rivo/tview"
)

type FileDialog struct {
	*tview.Flex
	app       *tview.Application
	pathField *tview.InputField
	fileList  *tview.List

	dirRequestsChan  chan string
	currentDirectory string
	directoryEntries []os.DirEntry
}

func NewFileDialog(app *tview.Application) *FileDialog {
	vertContentsFlex := tview.NewFlex()
	vertContentsFlex.SetTitle("Open File")
	vertContentsFlex.SetTitleAlign(tview.AlignLeft)
	vertContentsFlex.SetBorder(true)
	vertContentsFlex.SetDirection(tview.FlexRow)

	vertContentsFlex.AddItem(nil, 1, 0, false)

	pathField := tview.NewInputField().
		SetLabel("Path: ")
	// SetAcceptanceFunc(tview.InputFieldInteger).
	// SetDoneFunc(func(key tcell.Key) {
	// 	app.Stop()
	// })
	style.StyleInputField(pathField)
	vertContentsFlex.AddItem(pathField, 1, 0, false)

	vertContentsFlex.AddItem(nil, 1, 0, false)

	fileList := tview.NewList().
		SetUseStyleTags(false, false).
		ShowSecondaryText(false).
		SetHighlightFullLine(true)
	style.StyleList(fileList)

	vertContentsFlex.AddItem(fileList, 0, 1, false)
	// fileList.AddItem("..", "", 0, nil)
	// fileList.AddItem("dir/", "", 0, nil)
	// fileList.AddItem("file1.txt", "", 0, nil)

	vertContentsFlex.AddItem(nil, 1, 0, false)

	buttonFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn)
	buttonFlex.AddItem(tview.NewButton("Open"), 10, 1, false)
	buttonFlex.AddItem(nil, 0, 1, false)
	buttonFlex.AddItem(tview.NewButton("Cancel"), 10, 1, false)

	vertContentsFlex.AddItem(buttonFlex, 1, 0, false)
	vertContentsFlex.AddItem(nil, 1, 0, false)

	height := 30
	width := 60

	dirRequestsChan := make(chan string, 10)
	fileDialog := &FileDialog{
		Flex: tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(vertContentsFlex, height, 1, true).
				AddItem(nil, 0, 1, false), width, 1, true).
			AddItem(nil, 0, 1, false),
		app:              app,
		pathField:        pathField,
		fileList:         fileList,
		dirRequestsChan:  dirRequestsChan,
		currentDirectory: "/",
	}

	fileList.SetChangedFunc(fileDialog.handleListChanged)
	fileList.SetDoneFunc(func() {

	})
	fileList.SetSelectedFunc(fileDialog.handleListSelected)

	go fileDialog.runDirectoryLister(dirRequestsChan)

	return fileDialog
}

func (fileDialog *FileDialog) runDirectoryLister(dirRequests chan string) {
	for {
		dirPath, ok := <-dirRequests
		if !ok || dirPath == "" {
			return
		}
		entries, err := os.ReadDir(dirPath)
		if err != nil {

		} else {
			fileDialog.app.QueueUpdateDraw(func() {
				fileDialog.loadDirectoryEntries(entries, dirPath)
			})
		}
	}
}

func (fileDialog *FileDialog) loadDirectoryEntries(entries []os.DirEntry, dirPath string) {
	fileDialog.currentDirectory = dirPath
	fileDialog.directoryEntries = entries

	fileList := fileDialog.fileList
	fileList.Clear()

	// Sort the entries by directory first, then by name, case-insensitive
	entries = sortEntries(entries)

	// Add ".." for parent directory
	fileList.AddItem("\u2ba4 ../", "", 0, nil)

	for _, entry := range entries {
		if entry.IsDir() {
			fileList.AddItem("📁"+entry.Name()+"/", "", 0, nil)
		} else {
			fileList.AddItem("\U0001f4c4"+entry.Name(), "", 0, nil)
		}
	}
}

func sortEntries(entries []os.DirEntry) []os.DirEntry {
	// Separate directories and files
	var dirs, files []os.DirEntry
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry)
		} else {
			files = append(files, entry)
		}
	}
	// Sort directories and files by name, case-insensitive
	sort.Slice(dirs, func(i, j int) bool {
		return strings.ToLower(filepath.Clean(dirs[i].Name())) < strings.ToLower(filepath.Clean(dirs[j].Name()))
	})
	sort.Slice(files, func(i, j int) bool {
		return strings.ToLower(filepath.Clean(files[i].Name())) < strings.ToLower(filepath.Clean(files[j].Name()))
	})

	// Combine directories and files
	return append(dirs, files...)
}

func (fileDialog *FileDialog) handleListChanged(index int, mainText string, secondaryText string, shortcut rune) {
	offset := 1 // Offset for the ".." entry
	if fileDialog.currentDirectory == "/" {
		offset = 0 // No parent directory for root
	}
	entryIndex := index - offset
	if entryIndex < 0 {
		return
	}

	dirEntry := fileDialog.directoryEntries[entryIndex]
	if !dirEntry.IsDir() {
		selectedPath := filepath.Join(fileDialog.currentDirectory, dirEntry.Name())
		fileDialog.pathField.SetText(selectedPath)
	}
}

func (fileDialog *FileDialog) handleListSelected(index int, mainText string, secondaryText string, shortcut rune) {
	if fileDialog.currentDirectory != "/" && index == 0 {
		// Handle ".." entry to go up one directory
		parentDir := filepath.Dir(filepath.Clean(fileDialog.currentDirectory)) + "/"
		fileDialog.dirRequestsChan <- parentDir
		fileDialog.pathField.SetText(parentDir)
		return
	}

	offset := 1 // Offset for the ".." entry
	if fileDialog.currentDirectory == "/" {
		offset = 0 // No parent directory for root
	}
	entryIndex := index - offset
	dirEntry := fileDialog.directoryEntries[entryIndex]
	if dirEntry.IsDir() {
		// Handle entering a directory
		selectedDir := filepath.Join(fileDialog.currentDirectory, dirEntry.Name()) + "/"
		fileDialog.dirRequestsChan <- selectedDir
		fileDialog.pathField.SetText(selectedDir)
	}
}

// func (fileDialog *FileDialog) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
// 	return fileDialog.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
// 		return true, nil
// 	})
// }

func (fileDialog *FileDialog) SetPath(path string) {
	fileDialog.pathField.SetText(path)
	if fileInfo, err := os.Stat(path); err == nil {
		dir := path
		if !fileInfo.IsDir() {
			dir, _ = filepath.Split(path)
		}
		fileDialog.dirRequestsChan <- dir
	}
}
