package filedialog

import (
	"dinky/internal/tui/filelist"
	"os"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	nuview "github.com/rivo/tview"
)

type FileDialogMode int

const (
	OPEN_FILE_MODE FileDialogMode = iota
	SAVE_FILE_MODE
)

type FileDialog struct {
	*nuview.Flex
	app                *nuview.Application
	DirectoryField     *nuview.InputField
	FilenameField      *nuview.InputField
	FileList           *filelist.FileList
	vertContentsFlex   *nuview.Flex
	ActionButton       *nuview.Button
	CancelButton       *nuview.Button
	ParentButton       *nuview.Button
	ShowHiddenCheckbox *nuview.Checkbox

	dirRequestsChan  chan string
	currentDirectory string
	completedFunc    func(accepted bool, path string)
	mode             FileDialogMode
}

func NewFileDialog(app *nuview.Application) *FileDialog {
	vertContentsFlex := nuview.NewFlex()

	vertContentsFlex.Box = nuview.NewBox() // Nasty hack to clear the `dontClear` flag inside Box.
	vertContentsFlex.Box.Primitive = vertContentsFlex

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

	vertContentsFlex.AddItem(fileList, 0, 1, false)
	vertContentsFlex.AddItem(nil, 1, 0, false)

	filenameField := nuview.NewInputField()
	filenameField.SetLabel("File name: ")
	vertContentsFlex.AddItem(filenameField, 1, 0, false)

	vertContentsFlex.AddItem(nil, 1, 0, false)

	buttonFlex := nuview.NewFlex()
	buttonFlex.SetDirection(nuview.FlexColumn)

	showHiddenCheckbox := nuview.NewCheckbox()
	showHiddenCheckbox.SetLabel("Show Hidden Files: ")
	buttonFlex.AddItem(showHiddenCheckbox, 0, 1, false)

	actionButton := nuview.NewButton("Open")
	buttonFlex.AddItem(actionButton, 10, 1, false)
	buttonFlex.AddItem(nil, 1, 0, false)
	cancelButton := nuview.NewButton("Cancel")
	buttonFlex.AddItem(cancelButton, 10, 1, false)

	vertContentsFlex.AddItem(buttonFlex, 1, 0, false)
	vertContentsFlex.AddItem(nil, 1, 0, false)

	dirRequestsChan := make(chan string, 10)

	padding := 4
	flex := nuview.NewFlex()
	flex.AddItem(nil, padding, 0, false)

	innerFlex := nuview.NewFlex()
	innerFlex.SetDirection(nuview.FlexRow)
	innerFlex.AddItem(nil, padding, 0, false)
	innerFlex.AddItem(vertContentsFlex, 0, 1, true)
	innerFlex.AddItem(nil, padding, 0, false)

	flex.AddItem(innerFlex, 0, 1, true)
	flex.AddItem(nil, padding, 0, false)

	fileDialog := &FileDialog{
		Flex:               flex,
		app:                app,
		vertContentsFlex:   vertContentsFlex,
		FilenameField:      filenameField,
		DirectoryField:     directoryField,
		FileList:           fileList,
		dirRequestsChan:    dirRequestsChan,
		currentDirectory:   "/",
		ActionButton:       actionButton,
		CancelButton:       cancelButton,
		ParentButton:       parentButton,
		ShowHiddenCheckbox: showHiddenCheckbox,
		mode:               OPEN_FILE_MODE,
	}

	fileList.SetChangedFunc(fileDialog.handleListChanged)
	fileList.SetSelectedFunc(fileDialog.handleListSelected)

	directoryField.SetDoneFunc(fileDialog.handleDirectoryDone)
	directoryField.SetChangedFunc(func(text string) {
		fileDialog.syncActionButton()
	})
	filenameField.SetDoneFunc(fileDialog.handleFilenameDone)
	filenameField.SetChangedFunc(func(text string) {
		fileDialog.syncActionButton()
	})

	parentButton.SetSelectedFunc(func() {
		currentPath := fileDialog.FileList.Path()
		if currentPath == "" || currentPath == "/" {
			return
		}

		newPath := filepath.Dir(filepath.Clean(currentPath))
		fileDialog.FileList.SetPath(newPath)
		fileDialog.DirectoryField.SetText(newPath)
	})

	showHiddenCheckbox.SetChangedFunc(func(checked bool) {
		fileList.SetShowHidden(checked)
	})

	actionButton.SetSelectedFunc(fileDialog.doOpen)
	cancelButton.SetSelectedFunc(fileDialog.doCancel)
	return fileDialog
}

func (fileDialog *FileDialog) SetTitle(title string) {
	fileDialog.vertContentsFlex.SetTitle(title)
}

func (fileDialog *FileDialog) SetMode(mode FileDialogMode) {
	fileDialog.mode = mode
	if mode == OPEN_FILE_MODE {
		fileDialog.ActionButton.SetLabel("Open")
		fileDialog.vertContentsFlex.SetTitle("Open File")
	} else if mode == SAVE_FILE_MODE {
		fileDialog.ActionButton.SetLabel("Save")
		fileDialog.vertContentsFlex.SetTitle("Save As")
	}
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

func (fileDialog *FileDialog) Focus(delegate func(p nuview.Primitive)) {
	delegate(fileDialog.FilenameField)
}

func (fileDialog *FileDialog) handleDirectoryDone(key tcell.Key) {
	if key != tcell.KeyEnter {
		return
	}

	newDirectory := fileDialog.DirectoryField.GetText()
	if info, err := os.Stat(newDirectory); err == nil && info.IsDir() {
		fileDialog.FileList.SetPath(newDirectory)
		fileDialog.DirectoryField.SetText(newDirectory)
		fileDialog.syncActionButton()
	}
}

func (fileDialog *FileDialog) handleFilenameDone(key tcell.Key) {
	if key != tcell.KeyEnter {
		return
	}

	filename := fileDialog.FilenameField.GetText()
	directory := fileDialog.DirectoryField.GetText()
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

	fileDialog.FileList.SetPath(directory)
	fileDialog.DirectoryField.SetText(directory)
	fileDialog.FilenameField.SetText(filename)
	fileDialog.syncActionButton()
}

func (fileDialog *FileDialog) doOpen() {
	if fileDialog.completedFunc != nil {
		completePath := filepath.Join(fileDialog.DirectoryField.GetText(), fileDialog.FilenameField.GetText())
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
		fileDialog.FileList.SetPath(path)
		fileDialog.DirectoryField.SetText(path + "/")
		fileDialog.syncActionButton()
	} else {
		fileDialog.doOpen()
	}
}

func (fileDialog *FileDialog) handleListChanged(path string, entry os.DirEntry) {
	if !entry.IsDir() {
		fileDialog.FilenameField.SetText(entry.Name())
	}
	fileDialog.syncActionButton()
}

func (fileDialog *FileDialog) SetPath(path string) {
	if info, err := os.Stat(path); err == nil {
		if info.IsDir() {
			if !strings.HasSuffix(path, "/") {
				path += "/"
			}
			fileDialog.FileList.SetPath(path)
			fileDialog.DirectoryField.SetText(path)
			fileDialog.FilenameField.SetText("")
		} else {
			dir := filepath.Dir(path)
			filename := filepath.Base(path)
			fileDialog.FileList.SetPath(dir)
			fileDialog.DirectoryField.SetText(dir)
			fileDialog.FilenameField.SetText(filename)
		}
		fileDialog.syncActionButton()
	}
}

func (fileDialog *FileDialog) syncActionButton() {
	var fullpath string
	directory := fileDialog.DirectoryField.GetText()
	filename := fileDialog.FilenameField.GetText()

	if fileDialog.mode == OPEN_FILE_MODE {
		fullpath = filepath.Join(directory, filename)

		enabled := false
		if info, err := os.Stat(fullpath); err == nil && !info.IsDir() {
			enabled = true
		}
		fileDialog.ActionButton.SetDisabled(!enabled)
		return
	}

	if fileDialog.mode == SAVE_FILE_MODE {
		if filename == "" {
			fileDialog.ActionButton.SetDisabled(true)
			return
		}

		if info, err := os.Stat(directory); err == nil {
			if !info.IsDir() {
				fileDialog.ActionButton.SetDisabled(true)
			} else {
				fileDialog.ActionButton.SetDisabled(false)
			}
		}
		return
	}

}
