package application

import (
	"dinky/internal/tui/filedialog"
	"os"
	"path/filepath"

	"github.com/google/renameio/v2"
	"github.com/pgavlin/femto"
)

const (
	ACTION_NEW                 = "NewFile"
	ACTION_CLOSE_FILE          = "CloseFile"
	ACTION_OPEN_FILE           = "OpenFile"
	ACTION_SAVE_FILE           = "SaveFile"
	ACTION_SAVE_FILE_AS        = "SaveFileAs"
	ACTION_OPEN_MENU           = "OpenMenu"
	ACTION_TOGGLE_SOFT_WRAP    = "ToggleSoftWrap"
	ACTION_TOGGLE_LINE_NUMBERS = "ToggleLineNumbers"
	ACTION_QUIT                = "Quit"
	ACTION_ABOUT               = "About"
)

var dinkyActionMapping map[string]func()

func init() {
	dinkyActionMapping = map[string]func(){
		ACTION_NEW:                 handleNewFile,
		ACTION_CLOSE_FILE:          handleCloseFile,
		ACTION_OPEN_FILE:           handleOpenFile,
		ACTION_OPEN_MENU:           handleOpenMenu,
		ACTION_SAVE_FILE:           handleSaveFile,
		ACTION_SAVE_FILE_AS:        handleSaveFileAs,
		ACTION_TOGGLE_LINE_NUMBERS: handleLineNumbers,
		ACTION_TOGGLE_SOFT_WRAP:    handleSoftWrap,
		ACTION_QUIT:                handleQuit,
		ACTION_ABOUT:               handleAbout,
	}
}

func handleDinkyAction(id string) {
	if f, ok := dinkyActionMapping[id]; ok {
		f()
	}
}

func handleNewFile() {
	newFile("", "")
}

var fileDialog *filedialog.FileDialog

const fileDialogName = "fileDialog"

func showFileDialog(title string, mode filedialog.FileDialogMode, defaultPath string, completedFunc func(accepted bool, filePath string)) {
	if fileDialog == nil {
		fileDialog = filedialog.NewFileDialog(app)
	}
	fileDialog.SetTitle(title)
	if defaultPath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			cwd = "/"
		}
		fileDialog.SetPath(cwd)
	} else {
		fileDialog.SetPath(defaultPath)
	}
	fileDialog.SetMode(mode)
	fileDialog.SetCompletedFunc(completedFunc)
	modalPages.AddPanel(fileDialogName, fileDialog, true, true)
}

func hideFileDialog() {
	if fileDialog != nil {
		modalPages.RemovePanel(fileDialogName)
	}
}

func handleOpenFile() {
	showFileDialog("Open File", filedialog.OPEN_FILE_MODE, "", func(accepted bool, filePath string) {
		hideFileDialog()
		if !accepted {
			return
		}
		loadFile(filePath)
	})
}

func bufferToBytes(buffer *femto.Buffer) []byte {
	str := buffer.LineArray.SaveString(isBufferCRLF(buffer))
	return []byte(str)
}

func writeFile(filename string, buffer *femto.Buffer) (ok bool, message string) {
	contents := bufferToBytes(buffer)
	err := renameio.WriteFile(filename, contents, 0644)
	if err != nil {
		return false, "Error writing file: " + err.Error()
	}
	return true, "Wrote file " + filename
}

func handleSaveFile() {
	fileBuffer := getFileBufferByID(fileBufferID)
	if fileBuffer.filename == "" {
		handleSaveFileAs()
	} else {
		writeCurrentFileBuffer()
	}
}

func writeCurrentFileBuffer() {
	fileBuffer := getFileBufferByID(fileBufferID)
	ok, message := writeFile(fileBuffer.filename, fileBuffer.buffer)
	if ok {
		statusBar.ShowMessage(message)
	} else {
		statusBar.ShowError(message)
	}
}

func handleSaveFileAs() {
	fileBuffer := getFileBufferByID(fileBufferID)
	showFileDialog("Save File As", filedialog.SAVE_FILE_MODE, fileBuffer.filename, func(accepted bool, filePath string) {
		hideFileDialog()
		if !accepted {
			return
		}
		fileBuffer.filename = filePath
		tabBarLine.SetTabTitle(fileBufferID, filepath.Base(fileBuffer.filename))
		writeCurrentFileBuffer()
	})
}

func handleCloseFile() {
	fileBuffer := getFileBufferByID(fileBufferID)
	if fileBuffer.buffer.IsModified {
		ShowConfirmDialog("File has unsaved changes. Close anyway?", func() {
			closeFile(fileBufferID)
		}, func() {})
	} else {
		closeFile(fileBufferID)
	}
}

func closeFile(fileBufferID string) {
	tabBarLine.RemoveTab(fileBufferID)
	editorPages.RemovePanel(fileBufferID)

	for i, fileBuffer := range fileBuffers {
		if fileBuffer.uuid == fileBufferID {
			fileBuffers = append(fileBuffers[:i], fileBuffers[i+1:]...)
			break
		}
	}

	if len(fileBuffers) == 0 {
		fileBufferID = ""
		handleQuit()
	} else {
		selectTab(fileBuffers[0].uuid)
	}
}

func handleOpenMenu() {
	menuBar.Open()
	app.SetFocus(menuBar)
}

func handleLineNumbers() {
	on := buffer.Settings["ruler"].(bool)
	buffer.Settings["ruler"] = !on
	syncLineNumbers(menus, !on)
}

func handleSoftWrap() {
	on := buffer.Settings["softwrap"].(bool)
	buffer.Settings["softwrap"] = !on
	syncSoftWrap(menus, !on)
}

func handleFemtoAction(id string) {
	if f, ok := femto.BindingActionsMapping[id]; ok {
		f(editor)
	}
}

func handleAbout() {
	ShowOkDialog("About", "Dinky - A little text editor\nVersion 0.1.0\n"+
		"\n"+
		"Website: https://github.com/sedwards2009/dinky\n"+
		"(c) 2025 Simon Edwards",
		40, 11, func() {})
}

func handleQuit() {
	app.Stop()
}
