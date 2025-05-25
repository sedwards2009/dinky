package application

import (
	"dinky/internal/tui/filedialog"

	"github.com/google/renameio/v2"
	"github.com/pgavlin/femto"
)

const (
	ACTION_NEW                 = "NewFile"
	ACTION_OPEN_FILE           = "OpenFile"
	ACTION_SAVE_FILE           = "SaveFile"
	ACTION_SAVE_FILE_AS        = "SaveFileAs"
	ACTION_OPEN_MENU           = "OpenMenu"
	ACTION_TOGGLE_SOFT_WRAP    = "ToggleSoftWrap"
	ACTION_TOGGLE_LINE_NUMBERS = "ToggleLineNumbers"
	ACTION_QUIT                = "Quit"
)

var dinkyActionMapping map[string]func()

func init() {
	dinkyActionMapping = map[string]func(){
		ACTION_NEW:                 handleNewFile,
		ACTION_OPEN_FILE:           handleOpenFile,
		ACTION_OPEN_MENU:           handleOpenMenu,
		ACTION_SAVE_FILE:           handleSaveFile,
		ACTION_SAVE_FILE_AS:        handleSaveFileAs,
		ACTION_TOGGLE_LINE_NUMBERS: handleLineNumbers,
		ACTION_TOGGLE_SOFT_WRAP:    handleSoftWrap,
		ACTION_QUIT:                handleQuit,
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

func handleOpenFile() {
	openFileDialog := filedialog.NewFileDialog(app)
	modalPages.AddPage("openFileDialog", openFileDialog, true, true)
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
	ok, message := writeFile(fileBuffer.filename, fileBuffer.buffer)
	if ok {
		statusBar.ShowMessage(message)
	} else {
		statusBar.ShowError(message)
	}
}

func handleSaveFileAs() {

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

func handleQuit() {
	app.Stop()
}
