package application

import "github.com/pgavlin/femto"

const (
	ACTION_OPEN_MENU           = "OpenMenu"
	ACTION_TOGGLE_SOFT_WRAP    = "ToggleSoftWrap"
	ACTION_TOGGLE_LINE_NUMBERS = "ToggleLineNumbers"
	ACTION_QUIT                = "Quit"
)

var dinkyActionMapping = map[string]func(){
	ACTION_OPEN_MENU:           handleOpenMenu,
	ACTION_TOGGLE_SOFT_WRAP:    handleSoftWrap,
	ACTION_TOGGLE_LINE_NUMBERS: handleLineNumbers,
	ACTION_QUIT:                handleQuit,
}

func handleDinkyAction(id string) {
	if f, ok := dinkyActionMapping[id]; ok {
		f()
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

func handleQuit() {
	app.Stop()
}
