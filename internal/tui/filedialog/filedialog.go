package filedialog

import (
	"dinky/internal/tui/femtoinputfield"
	"dinky/internal/tui/filelist"
	"os"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sedwards2009/femto"
)

type FileDialogMode int

const (
	OPEN_FILE_MODE FileDialogMode = iota
	SAVE_FILE_MODE
)

type FileDialog struct {
	*tview.Flex
	app                *tview.Application
	DirectoryField     *femtoinputfield.FemtoInputField
	FilenameField      *femtoinputfield.FemtoInputField
	FileList           *filelist.FileList
	vertContentsFlex   *tview.Flex
	ActionButton       *tview.Button
	CancelButton       *tview.Button
	ParentButton       *tview.Button
	ShowHiddenCheckbox *tview.Checkbox

	dirRequestsChan  chan string
	currentDirectory string
	completedFunc    func(accepted bool, path string)
	mode             FileDialogMode
}

func NewFileDialog(app *tview.Application) *FileDialog {
	vertContentsFlex := tview.NewFlex()

	vertContentsFlex.Box = tview.NewBox() // Nasty hack to clear the `dontClear` flag inside Box.
	vertContentsFlex.Box.Primitive = vertContentsFlex

	vertContentsFlex.SetTitle("Open File")
	vertContentsFlex.SetTitleAlign(tview.AlignLeft)
	vertContentsFlex.SetBorder(true)
	vertContentsFlex.SetDirection(tview.FlexRow)

	vertContentsFlex.SetBorderPadding(1, 1, 1, 1)

	directoryFlex := tview.NewFlex()
	directoryFlex.SetDirection(tview.FlexColumn)
	directoryFlex.SetBorder(false)

	directoryLabel := tview.NewTextView()
	directoryLabel.SetDynamicColors(true)
	directoryLabel.SetText("[::u]D[::U]irectory: ")
	directoryFlex.AddItem(directoryLabel, 11, 0, false)

	directoryField := femtoinputfield.NewFemtoInputField()
	directoryFlex.AddItem(directoryField, 0, 1, false)

	directoryFlex.AddItem(nil, 1, 0, false)
	parentButton := tview.NewButton("\u2ba4[::u].[::U].")
	directoryFlex.AddItem(parentButton, 5, 0, false)

	vertContentsFlex.AddItem(directoryFlex, 1, 0, false)
	vertContentsFlex.AddItem(nil, 1, 0, false)

	fileList := filelist.NewFileList(app)

	vertContentsFlex.AddItem(fileList, 0, 1, false)
	vertContentsFlex.AddItem(nil, 1, 0, false)

	filenameFlex := tview.NewFlex()
	filenameLabel := tview.NewTextView()
	filenameLabel.SetLabel("[::u]F[::U]ile name: ")
	filenameFlex.AddItem(filenameLabel, 11, 0, false)
	filenameField := femtoinputfield.NewFemtoInputField()
	filenameFlex.AddItem(filenameField, 0, 1, false)

	vertContentsFlex.AddItem(filenameFlex, 1, 0, false)

	vertContentsFlex.AddItem(nil, 1, 0, false)

	buttonFlex := tview.NewFlex()
	buttonFlex.SetDirection(tview.FlexColumn)

	showHiddenCheckbox := tview.NewCheckbox()
	showHiddenCheckbox.SetLabel("Show [::u]H[::U]idden Files: ")
	buttonFlex.AddItem(showHiddenCheckbox, 0, 1, false)

	actionButton := tview.NewButton("[::u]O[::U]pen")
	buttonFlex.AddItem(actionButton, 10, 1, false)
	buttonFlex.AddItem(nil, 1, 0, false)
	cancelButton := tview.NewButton("[::u]C[::U]ancel")
	buttonFlex.AddItem(cancelButton, 10, 1, false)

	vertContentsFlex.AddItem(buttonFlex, 1, 0, false)

	dirRequestsChan := make(chan string, 10)

	padding := 4
	flex := tview.NewFlex()
	flex.AddItem(nil, padding, 0, false)

	innerFlex := tview.NewFlex()
	innerFlex.SetDirection(tview.FlexRow)
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
	fileList.SetInputCapture(fileDialog.handleShortcuts)

	directoryField.SetDoneFunc(fileDialog.handleDirectoryDone)
	directoryField.SetChangedFunc(func(text string) {
		fileDialog.syncActionButton()
	})
	directoryField.SetInputCapture(fileDialog.handleShortcuts)

	filenameField.SetDoneFunc(fileDialog.handleFilenameDone)
	filenameField.SetChangedFunc(func(text string) {
		fileDialog.syncActionButton()
	})
	filenameField.SetInputCapture(fileDialog.handleShortcuts)

	parentButton.SetSelectedFunc(fileDialog.handleParentButtonClick)

	showHiddenCheckbox.SetChangedFunc(func(checked bool) {
		fileList.SetShowHidden(checked)
	})

	actionButton.SetSelectedFunc(fileDialog.doOpen)
	cancelButton.SetSelectedFunc(fileDialog.doCancel)
	return fileDialog
}

func (fileDialog *FileDialog) handleParentButtonClick() {
	currentPath := fileDialog.FileList.Path()
	if currentPath == "" || currentPath == "/" {
		return
	}

	newPath := filepath.Dir(filepath.Clean(currentPath))
	fileDialog.FileList.SetPath(newPath)
	fileDialog.DirectoryField.SetText(newPath)
}

func (fileDialog *FileDialog) handleShortcuts(event *tcell.EventKey) *tcell.EventKey {
	isAlt := event.Modifiers()&tcell.ModAlt != 0
	if isAlt {
		switch event.Rune() {
		case 'f', 'F':
			fileDialog.app.SetFocus(fileDialog.FilenameField)
			return nil
		case 'd', 'D':
			fileDialog.app.SetFocus(fileDialog.DirectoryField)
			return nil
		case 'h', 'H':
			checked := fileDialog.ShowHiddenCheckbox.IsChecked()
			fileDialog.ShowHiddenCheckbox.SetChecked(!checked)
			return nil
		case '.':
			fileDialog.handleParentButtonClick()
			return nil
		case 'c', 'C':
			fileDialog.doCancel()
			return nil
		case 'o', 'O':
			if fileDialog.mode == OPEN_FILE_MODE && !fileDialog.ActionButton.IsDisabled() {
				fileDialog.doOpen()
				return nil
			}
		case 's', 'S':
			if fileDialog.mode == SAVE_FILE_MODE && !fileDialog.ActionButton.IsDisabled() {
				fileDialog.doOpen()
				return nil
			}
		}
	}
	if event.Key() == tcell.KeyTab {
		fileDialog.app.SetFocus(fileDialog.FileList)
		return nil
	}

	return event
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

func (fileDialog *FileDialog) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return fileDialog.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		fileDialog.vertContentsFlex.MouseHandler()(action, event, setFocus)
		return true, nil
	})
}

func (fileDialog *FileDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return fileDialog.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		fileDialog.vertContentsFlex.InputHandler()(event, setFocus)
	})
}

func (fileDialog *FileDialog) Focus(delegate func(p tview.Primitive)) {
	delegate(fileDialog.FilenameField)
}

func (fileDialog *FileDialog) handleDirectoryDone(key tcell.Key) {
	if key != tcell.KeyEnter {
		return
	}

	newDirectory := fileDialog.DirectoryField.GetText()
	newDirectory = filepath.Clean(newDirectory)
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

func (fileDialog *FileDialog) SetFemtoKeybindings(keybindings femto.KeyBindings) {
	fileDialog.DirectoryField.SetKeybindings(keybindings)
	fileDialog.FilenameField.SetKeybindings(keybindings)
}
