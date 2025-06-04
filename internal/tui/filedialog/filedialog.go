package filedialog

import (
	"dinky/internal/tui/filelist"
	"dinky/internal/tui/style"
	"os"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type FileDialog struct {
	*tview.Flex
	app            *tview.Application
	directoryField *tview.InputField
	filenameField  *tview.InputField
	fileList       *filelist.FileList

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

	directoryFlex := tview.NewFlex()
	directoryFlex.SetDirection(tview.FlexColumn)
	directoryFlex.SetBorder(false)

	directoryField := tview.NewInputField()
	directoryField.SetLabel("Directory: ")
	style.StyleInputField(directoryField)
	directoryFlex.AddItem(directoryField, 0, 1, false)

	directoryFlex.AddItem(nil, 1, 0, false)
	parentButton := tview.NewButton("\u2ba4")
	directoryFlex.AddItem(parentButton, 3, 0, false)

	vertContentsFlex.AddItem(directoryFlex, 1, 0, false)
	vertContentsFlex.AddItem(nil, 1, 0, false)

	fileList := filelist.NewFileList(app)

	vertContentsFlex.AddItem(fileList, 0, 1, false)
	vertContentsFlex.AddItem(nil, 1, 0, false)

	filenameField := tview.NewInputField()
	filenameField.SetLabel("File name: ")
	style.StyleInputField(filenameField)
	vertContentsFlex.AddItem(filenameField, 1, 0, false)

	vertContentsFlex.AddItem(nil, 1, 0, false)

	buttonFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	showHiddenCheckbox := tview.NewCheckbox()
	showHiddenCheckbox.SetLabel("Show Hidden Files")
	buttonFlex.AddItem(showHiddenCheckbox, 20, 0, false)

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
		filenameField:    filenameField,
		directoryField:   directoryField,
		fileList:         fileList,
		dirRequestsChan:  dirRequestsChan,
		currentDirectory: "/",
	}

	fileList.SetChangedFunc(fileDialog.handleListChanged)
	fileList.SetSelectedFunc(fileDialog.handleListSelected)

	directoryField.SetDoneFunc(fileDialog.handleDirectoryDone)
	filenameField.SetDoneFunc(fileDialog.handleFilenameDone)

	parentButton.SetSelectedFunc(func() {
		currentPath := fileDialog.fileList.Path()
		if currentPath == "" || currentPath == "/" {
			return
		}

		newPath := filepath.Dir(filepath.Clean(currentPath))
		fileDialog.fileList.SetPath(newPath)
		fileDialog.directoryField.SetText(newPath)
	})

	showHiddenCheckbox.SetChangedFunc(func(checked bool) {
		fileList.SetShowHidden(checked)
	})

	openButton.SetSelectedFunc(fileDialog.doOpen)
	cancelButton.SetSelectedFunc(fileDialog.doCancel)
	return fileDialog
}

func (fileDialog *FileDialog) handleDirectoryDone(key tcell.Key) {
	if key != tcell.KeyEnter {
		return
	}

	newDirectory := fileDialog.directoryField.GetText()
	if info, err := os.Stat(newDirectory); err == nil && info.IsDir() {
		fileDialog.fileList.SetPath(newDirectory)
		fileDialog.directoryField.SetText(newDirectory)
	}
}

func (fileDialog *FileDialog) handleFilenameDone(key tcell.Key) {
	if key != tcell.KeyEnter {
		return
	}

	filename := fileDialog.filenameField.GetText()
	directory := fileDialog.directoryField.GetText()
	if strings.HasPrefix(filename, "/") {
		directory = "/"
		filename = strings.TrimPrefix(filename, "/")
	} else {
		completePath := filepath.Join(directory, filename)
		if info, err := os.Stat(completePath); err == nil {
			if info.IsDir() {
				directory = completePath
				filename = ""
			} else {
				// Got a valid file name and directory combination
				fileDialog.doOpen()
				return
			}
		}
	}

	for strings.Contains(filename, "/") {
		parts := strings.SplitN(filename, "/", 2)
		newDirectory := filepath.Join(directory, parts[0])
		if info, err := os.Stat(newDirectory); err == nil && info.IsDir() {
			directory = newDirectory
			filename = parts[1]
		} else {
			break
		}
	}

	fileDialog.fileList.SetPath(directory)
	fileDialog.directoryField.SetText(directory)
	fileDialog.filenameField.SetText(filename)
}

func (fileDialog *FileDialog) doOpen() {
	if fileDialog.completedFunc != nil {
		completePath := filepath.Join(fileDialog.directoryField.GetText(), fileDialog.filenameField.GetText())
		fileDialog.completedFunc(true, completePath)
	}
}

func (fileDialog *FileDialog) doCancel() {
	if fileDialog.completedFunc != nil {
		fileDialog.completedFunc(false, "")
	}
}

func (fileDialog *FileDialog) SetCompletedFunc(completedFunc func(accepted bool, path string)) {
	fileDialog.completedFunc = completedFunc
}

func (fileDialog *FileDialog) handleListSelected(path string, entry os.DirEntry) {
	if entry.IsDir() {
		fileDialog.fileList.SetPath(path)
		fileDialog.directoryField.SetText(path + "/")
	} else {
		fileDialog.doOpen()
	}
}

func (fileDialog *FileDialog) handleListChanged(path string, entry os.DirEntry) {
	if !entry.IsDir() {
		fileDialog.filenameField.SetText(entry.Name())
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
	fileDialog.directoryField.SetText(path)
}
