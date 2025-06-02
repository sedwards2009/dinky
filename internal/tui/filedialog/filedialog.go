package filedialog

import (
	"dinky/internal/tui/filelist"
	"dinky/internal/tui/style"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/rivo/tview"
)

type FileDialog struct {
	*tview.Flex
	app       *tview.Application
	pathField *tview.InputField
	fileList  *filelist.FileList

	dirRequestsChan  chan string
	currentDirectory string
	completedFunc    func(accepted bool, path string)
}

func NewFileDialog(app *tview.Application) *FileDialog {
	vertContentsFlex := tview.NewFlex()
	vertContentsFlex.SetTitle("Open File")
	vertContentsFlex.SetTitleAlign(tview.AlignLeft)
	vertContentsFlex.SetBorder(true)
	vertContentsFlex.SetDirection(tview.FlexRow)

	vertContentsFlex.AddItem(nil, 1, 0, false)

	pathFlex := tview.NewFlex()
	pathFlex.SetDirection(tview.FlexColumn)
	pathFlex.SetBorder(false)

	parentButton := tview.NewButton("\u2ba4")
	pathFlex.AddItem(parentButton, 3, 0, false)
	pathFlex.AddItem(nil, 1, 0, false)

	pathField := tview.NewInputField()
	style.StyleInputField(pathField)
	pathFlex.AddItem(pathField, 0, 1, false)

	// SetAcceptanceFunc(tview.InputFieldInteger).
	// SetDoneFunc(func(key tcell.Key) {
	// 	app.Stop()
	// })
	vertContentsFlex.AddItem(pathFlex, 1, 0, false)
	vertContentsFlex.AddItem(nil, 1, 0, false)

	fileList := filelist.NewFileList(app)

	vertContentsFlex.AddItem(fileList, 0, 1, false)
	vertContentsFlex.AddItem(nil, 1, 0, false)

	buttonFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn)
	buttonFlex.AddItem(nil, 0, 1, false)
	openButton := tview.NewButton("Open")
	buttonFlex.AddItem(openButton, 10, 1, false)
	buttonFlex.AddItem(nil, 1, 0, false)
	cancelButton := tview.NewButton("Cancel")
	buttonFlex.AddItem(cancelButton, 10, 1, false)

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
	fileList.SetSelectedFunc(fileDialog.handleListSelected)

	parentButton.SetSelectedFunc(func() {
		currentPath := fileDialog.fileList.Path()
		if currentPath == "" || currentPath == "/" {
			return
		}

		newPath := filepath.Dir(filepath.Clean(currentPath))
		fileDialog.fileList.SetPath(newPath)
		fileDialog.pathField.SetText(newPath)
	})

	openButton.SetSelectedFunc(func() {
		if fileDialog.completedFunc != nil {
			fileDialog.completedFunc(true, fileDialog.pathField.GetText())
		}
	})
	cancelButton.SetSelectedFunc(func() {
		if fileDialog.completedFunc != nil {
			fileDialog.completedFunc(false, "")
		}
	})
	return fileDialog
}

func (fileDialog *FileDialog) SetCompletedFunc(completedFunc func(accepted bool, path string)) {
	fileDialog.completedFunc = completedFunc
}

func (fileDialog *FileDialog) handleListChanged(path string, entry os.DirEntry) {
	if entry.IsDir() {
		fileDialog.pathField.SetText(path + "/")
	} else {
		fileDialog.pathField.SetText(path)
	}
}

func (fileDialog *FileDialog) handleListSelected(path string, entry os.DirEntry) {
	log.Printf("FileDialog handleListSelected: %s, isDir: %v", path, entry.IsDir())
	if entry.IsDir() {
		fileDialog.fileList.SetPath(path)
	} else {
		// // Handle file selection
		// fileDialog.pathField.SetText(path)
	}
}

func (fileDialog *FileDialog) SetPath(path string) {
	if info, err := os.Stat(path); err == nil {
		if info.IsDir() && !strings.HasSuffix(path, "/") {
			path += "/"
		}
	} else {
		return
	}
	fileDialog.fileList.SetPath(path)
	fileDialog.pathField.SetText(path)
}
