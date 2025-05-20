package application

import "github.com/pgavlin/femto"

const (
	ACTION_OPEN_MENU           = "OpenMenu"
	ACTION_TOGGLE_SOFT_WRAP    = "ToggleSoftWrap"
	ACTION_TOGGLE_LINE_NUMBERS = "ToggleLineNumbers"
)

func handleLineNumbers(_ string) {
	on := buffer.Settings["ruler"].(bool)
	buffer.Settings["ruler"] = !on
	syncLineNumbers(menus, !on)
}

func handleSoftWrap(_ string) {
	on := buffer.Settings["softwrap"].(bool)
	buffer.Settings["softwrap"] = !on
	syncSoftWrap(menus, !on)
}

func handleFemtoAction(id string) {
	if f, ok := femto.BindingActionsMapping[id]; ok {
		f(editor)
	}
}

func handleQuit(_ string) {
	app.Stop()
}
