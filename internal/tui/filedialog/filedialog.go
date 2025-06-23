package filedialog

import (
	"dinky/internal/tui/filelist"
	"dinky/internal/tui/style"
	"os"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/sedwards2009/nuview"
)

type FileDialog struct {
	*nuview.Flex
	app              *nuview.Application
	directoryField   *nuview.InputField
	filenameField    *nuview.InputField
	fileList         *filelist.FileList
	vertContentsFlex *nuview.Flex

	dirRequestsChan  chan string
	currentDirectory string
	completedFunc    func(accepted bool, path string)
}

func NewFileDialog(app *nuview.Application) *FileDialog {
	vertContentsFlex := nuview.NewFlex()
	vertContentsFlex.SetBackgroundTransparent(false)
	vertContentsFlex.SetTitle("Open File")
	vertContentsFlex.SetTitleAlign(nuview.AlignLeft)
	vertContentsFlex.SetBorder(true)
	vertContentsFlex.SetDirection(nuview.FlexRow)

	vertContentsFlex.AddItem(nil, 1, 0, false)

	directoryFlex := nuview.NewFlex()
	directoryFlex.SetDirection(nuview.FlexColumn)
	directoryFlex.SetBorder(false)

	directoryField := nuview.NewInputField()
	directoryField.SetLabel("Directory: ")
	directoryFlex.AddItem(directoryField, 0, 1, false)

	directoryFlex.AddItem(nil, 1, 0, false)
	parentButton := nuview.NewButton("\u2ba4")
	directoryFlex.AddItem(parentButton, 3, 0, false)

	vertContentsFlex.AddItem(directoryFlex, 1, 0, false)
	vertContentsFlex.AddItem(nil, 1, 0, false)

	fileList := filelist.NewFileList(app)
	style.StyleFileList(fileList)

	vertContentsFlex.AddItem(fileList, 0, 1, false)
	vertContentsFlex.AddItem(nil, 1, 0, false)

	filenameField := nuview.NewInputField()
	filenameField.SetLabel("File name: ")
	vertContentsFlex.AddItem(filenameField, 1, 0, false)

	vertContentsFlex.AddItem(nil, 1, 0, false)

	buttonFlex := nuview.NewFlex()
	buttonFlex.SetDirection(nuview.FlexColumn)

	showHiddenCheckbox := nuview.NewCheckbox()
	showHiddenCheckbox.SetLabelRight(" Show Hidden Files")
	buttonFlex.AddItem(showHiddenCheckbox, 0, 1, false)

	openButton := nuview.NewButton("Open")
	buttonFlex.AddItem(openButton, 10, 1, false)
	buttonFlex.AddItem(nil, 1, 0, false)
	cancelButton := nuview.NewButton("Cancel")
	buttonFlex.AddItem(cancelButton, 10, 1, false)

	vertContentsFlex.AddItem(buttonFlex, 1, 0, false)
	vertContentsFlex.AddItem(nil, 1, 0, false)

	height := 30
	width := 60

	dirRequestsChan := make(chan string, 10)

	flex := nuview.NewFlex()
	flex.AddItem(nil, 0, 1, false)

	innerFlex := nuview.NewFlex()
	innerFlex.SetDirection(nuview.FlexRow)
	innerFlex.AddItem(nil, 0, 1, false)
	innerFlex.AddItem(vertContentsFlex, height, 1, true)
	innerFlex.AddItem(nil, 0, 1, false)

	flex.AddItem(innerFlex, width, 1, true)
	flex.AddItem(nil, 0, 1, false)

	fileDialog := &FileDialog{
		Flex:             flex,
		app:              app,
		vertContentsFlex: vertContentsFlex,
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

func (fileDialog *FileDialog) MouseHandler() func(action nuview.MouseAction, event *tcell.EventMouse, setFocus func(p nuview.Primitive)) (consumed bool, capture nuview.Primitive) {
	return fileDialog.WrapMouseHandler(func(action nuview.MouseAction, event *tcell.EventMouse, setFocus func(p nuview.Primitive)) (consumed bool, capture nuview.Primitive) {
		fileDialog.vertContentsFlex.MouseHandler()(action, event, setFocus)
		return true, nil
	})
}

func (fileDialog *FileDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p nuview.Primitive)) {
	return fileDialog.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p nuview.Primitive)) {
		fileDialog.vertContentsFlex.InputHandler()(event, setFocus)
	})
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
