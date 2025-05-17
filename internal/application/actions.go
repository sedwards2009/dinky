package application

func handleSelectAll(_ string) {
	editor.SelectAll()
}

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

func handleUndo(_ string) {
	editor.Undo()
}

func handleRedo(_ string) {
	editor.Redo()
}

func handleQuit(_ string) {
	app.Stop()
}
