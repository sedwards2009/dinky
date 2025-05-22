package application

import (
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

}

func bufferToBytes(buffer *femto.Buffer) []byte {

}

func writeFile(filename string, buffer *femto.Buffer) string {
	contents := bufferToBytes(buffer)
	err := renameio.WriteFile(filename, contents, 0644)
	if err != nil {

	}
	return "Wrote file " + filename
}

func handleSaveFile() {

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
