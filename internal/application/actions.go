package application

import (
	"dinky/internal/application/settingstype"
	"dinky/internal/tui/dialog"
	"dinky/internal/tui/filedialog"
	"dinky/internal/tui/settingsdialog"
	"dinky/internal/tui/style"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/renameio/v2"
	"github.com/rivo/tview"
	"github.com/sedwards2009/smidgen"
	"github.com/sedwards2009/smidgen/micro/buffer"
)

const (
	ACTION_NEW                        = "NewFile"
	ACTION_CLOSE_FILE                 = "CloseFile"
	ACTION_OPEN_FILE                  = "OpenFile"
	ACTION_SAVE_FILE                  = "SaveFile"
	ACTION_SAVE_FILE_AS               = "SaveFileAs"
	ACTION_OPEN_FILE_MENU             = "OpenFileMenu"
	ACTION_OPEN_EDIT_MENU             = "OpenEditMenu"
	ACTION_OPEN_SELECTION_MENU        = "OpenSelectionMenu"
	ACTION_OPEN_TRANSFORM_MENU        = "OpenTransformMenu"
	ACTION_OPEN_VIEW_MENU             = "OpenViewMenu"
	ACTION_OPEN_HELP_MENU             = "OpenHelpMenu"
	ACTION_TOGGLE_SOFT_WRAP           = "ToggleSoftWrap"
	ACTION_TOGGLE_MATCH_BRACKET       = "ToggleMatchBracket"
	ACTION_SET_TAB_SIZE               = "SetTabSize"
	ACTION_SET_TAB_CHARACTER          = "SetTabCharacter"
	ACTION_SET_LINE_ENDINGS           = "SetLineEndings"
	ACTION_SET_SYNTAX_HIGHLIGHTING    = "SetSyntaxHighlighting"
	ACTION_SET_VERTICAL_RULER         = "SetVerticalRuler"
	ACTION_GO_TO_LINE                 = "GoToLine"
	ACTION_QUIT                       = "Quit"
	ACTION_FIND                       = "Find"
	ACTION_FIND_NEXT                  = "FindNext"
	ACTION_FIND_PREVIOUS              = "FindPrevious"
	ACTION_ABOUT                      = "About"
	ACTION_NEXT_EDITOR                = "NextEditor"
	ACTION_PREVIOUS_EDITOR            = "PreviousEditor"
	ACTION_CONVERT_TAB_SPACES         = "ConvertTabSpaces"
	ACTION_TOGGLE_WHITESPACE          = "ToggleWhitespace"
	ACTION_TOGGLE_TRAILING_WHITESPACE = "ToggleTrailingWhitespace"
	ACTION_FIND_AND_REPLACE           = "FindAndReplace"
	ACTION_SETTINGS                   = "Settings"
	ACTION_TO_UPPERCASE               = "ToUppercase"
	ACTION_TO_LOWERCASE               = "ToLowercase"
)

var dinkyActionMapping map[string]func() tview.Primitive

func init() {
	dinkyActionMapping = map[string]func() tview.Primitive{
		ACTION_NEW:                        handleNewFile,
		ACTION_CLOSE_FILE:                 handleCloseFile,
		ACTION_OPEN_FILE:                  handleOpenFile,
		ACTION_OPEN_FILE_MENU:             handleOpenFileMenu,
		ACTION_OPEN_EDIT_MENU:             handleOpenEditMenu,
		ACTION_OPEN_SELECTION_MENU:        handleOpenSelectionMenu,
		ACTION_OPEN_VIEW_MENU:             handleOpenViewMenu,
		ACTION_OPEN_HELP_MENU:             handleOpenHelpMenu,
		ACTION_SAVE_FILE:                  handleSaveFile,
		ACTION_SAVE_FILE_AS:               handleSaveFileAs,
		ACTION_TOGGLE_SOFT_WRAP:           handleSoftWrap,
		ACTION_TOGGLE_MATCH_BRACKET:       handleMatchBracket,
		ACTION_SET_TAB_SIZE:               handleSetTabSize,
		ACTION_SET_TAB_CHARACTER:          handleSetTabCharacter,
		ACTION_SET_LINE_ENDINGS:           handleSetLineEndings,
		ACTION_SET_SYNTAX_HIGHLIGHTING:    handleSetSyntaxHighlighting,
		ACTION_SET_VERTICAL_RULER:         handleSetVerticalRuler,
		ACTION_GO_TO_LINE:                 handleGoToLine,
		ACTION_QUIT:                       handleQuit,
		ACTION_ABOUT:                      handleAbout,
		ACTION_FIND:                       handleFind,
		ACTION_FIND_NEXT:                  handleFindNext,
		ACTION_FIND_PREVIOUS:              handleFindPrevious,
		ACTION_NEXT_EDITOR:                handleNextEditor,
		ACTION_PREVIOUS_EDITOR:            handlePreviousEditor,
		ACTION_CONVERT_TAB_SPACES:         handleConvertTabSpaces,
		ACTION_TOGGLE_WHITESPACE:          handleToggleWhitespace,
		ACTION_TOGGLE_TRAILING_WHITESPACE: handleToggleTrailingWhitespace,
		ACTION_FIND_AND_REPLACE:           handleFindAndReplace,
		ACTION_SETTINGS:                   handleSettings,
		ACTION_TO_UPPERCASE:               handleToUppercase,
		ACTION_TO_LOWERCASE:               handleToLowercase,
	}
}

func handleDinkyAction(id string) tview.Primitive {
	if f, ok := dinkyActionMapping[id]; ok {
		return f()
	}
	return nil
}

func handleNewFile() tview.Primitive {
	newFile("", "")
	return nil
}

var fileDialog *filedialog.FileDialog

const fileDialogName = "fileDialog"

func showFileDialog(title string, mode filedialog.FileDialogMode, defaultPath string,
	completedFunc func(accepted bool, filePath string)) tview.Primitive {

	if fileDialog == nil {
		fileDialog = filedialog.NewFileDialog(app)
		style.StyleFileDialog(fileDialog)
		fileDialog.SetSmidgenKeybindings(smidgenDefaultKeyBindings)
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
	modalPages.AddPage(fileDialogName, fileDialog, true, true)
	return fileDialog
}

func hideFileDialog() {
	if fileDialog != nil {
		modalPages.RemovePage(fileDialogName)
	}
}

var listDialog *dialog.ListDialog

const listDialogName = "listDialog"

func ShowListDialog(options dialog.ListDialogOptions) tview.Primitive {
	if listDialog == nil {
		listDialog = dialog.NewListDialog(app)
		style.StyleListDialog(listDialog)
	}
	modalPages.AddPage(listDialogName, listDialog, true, true)
	listDialog.Open(options)
	style.StyleListDialog(listDialog)
	return listDialog
}

func hideListDialog() {
	if listDialog != nil {
		listDialog.Close()
		modalPages.RemovePage(listDialogName)
	}
}

func handleOpenFile() tview.Primitive {
	return showFileDialog("Open File", filedialog.OPEN_FILE_MODE, "", func(accepted bool, filePath string) {
		hideFileDialog()
		if !accepted {
			return
		}
		loadFile(filePath)
	})
}

func bufferToBytes(buffer *buffer.Buffer) []byte {
	return buffer.LineArray.Bytes()
}

func writeFile(filename string, buffer *buffer.Buffer) (ok bool, message string) {
	contents := bufferToBytes(buffer)
	err := renameio.WriteFile(filename, contents, 0644)
	if err != nil {
		return false, "Error writing file: " + err.Error()
	}
	return true, "Wrote file " + filename
}

func handleSaveFile() tview.Primitive {
	fileBuffer := getFileBufferByID(fileBufferID)
	if fileBuffer.filename == "" {
		return handleSaveFileAs()
	} else {
		writeCurrentFileBuffer()
	}
	return nil
}

func writeCurrentFileBuffer() {
	fileBuffer := getFileBufferByID(fileBufferID)
	ok, message := writeFile(fileBuffer.filename, fileBuffer.buffer)
	fileBuffer.buffer.ClearModified()
	if ok {
		statusBar.ShowMessage(message)
	} else {
		statusBar.ShowError(message)
	}
}

func handleSaveFileAs() tview.Primitive {
	fileBuffer := getFileBufferByID(fileBufferID)
	return showFileDialog("Save File As", filedialog.SAVE_FILE_MODE, fileBuffer.filename,
		func(accepted bool, filePath string) {
			hideFileDialog()
			if !accepted {
				return
			}
			fileBuffer.filename = filePath
			tabBarLine.SetTabTitle(fileBufferID, filepath.Base(fileBuffer.filename))
			writeCurrentFileBuffer()
		})
}

func handleCloseFile() tview.Primitive {
	fileBuffer := getFileBufferByID(fileBufferID)
	if fileBuffer.buffer.Modified() {
		return ShowConfirmDialog("File has unsaved changes. Close anyway?", func() {
			closeFile(fileBufferID)
			if len(fileBuffers) == 0 {
				handleQuit()
			}
		}, func() {})
	} else {
		closeFile(fileBufferID)
		if len(fileBuffers) == 0 {
			handleQuit()
		}
	}
	return nil
}

func closeFile(fileBufferID string) {
	tabBarLine.RemoveTab(fileBufferID)
	editorPages.RemovePage(fileBufferID)

	for i, fileBuffer := range fileBuffers {
		if fileBuffer.uuid == fileBufferID {
			fileBuffers = append(fileBuffers[:i], fileBuffers[i+1:]...)
			break
		}
	}

	if len(fileBuffers) == 0 {
		fileBufferID = ""
	} else {
		selectTab(fileBuffers[0].uuid)
	}
}

func handleOpenFileMenu() tview.Primitive {
	menuBar.Open(0)
	app.SetFocus(menuBar)
	return nil
}

func handleOpenEditMenu() tview.Primitive {
	menuBar.Open(1)
	app.SetFocus(menuBar)
	return nil
}

func handleOpenSelectionMenu() tview.Primitive {
	menuBar.Open(2)
	app.SetFocus(menuBar)
	return nil
}

func handleOpenViewMenu() tview.Primitive {
	menuBar.Open(3)
	app.SetFocus(menuBar)
	return nil
}

func handleOpenHelpMenu() tview.Primitive {
	menuBar.Open(4)
	app.SetFocus(menuBar)
	return nil
}

func handleSoftWrap() tview.Primitive {
	buffer := currentFileBuffer.buffer
	on := buffer.Settings["softwrap"].(bool)
	buffer.Settings["softwrap"] = !on
	syncSoftWrap(menus, !on)
	return nil
}

func handleMatchBracket() tview.Primitive {
	buffer := currentFileBuffer.buffer
	on := buffer.Settings["matchbrace"].(bool)
	buffer.Settings["matchbrace"] = !on
	syncMatchBracket(menus, !on)
	return nil
}

func handleSmidgenAction(id string) tview.Primitive {
	action := currentFileBuffer.editor.MapActionNameToAction(id)
	if action != nil {
		action()
		if id == "ToggleRuler" {
			syncMenuFromBuffer(currentFileBuffer.buffer)
		}
	}
	return nil
}

func handleAbout() tview.Primitive {
	return ShowOkDialog("About", "Dinky - A little text editor\n"+
		"\n"+
		"Version: "+getDinkyVersion()+"\n"+
		"Version time: "+getDinkyVersionTime()+"\n"+
		"Website: https://github.com/sedwards2009/dinky\n"+
		"(c) 2025 Simon Edwards",
		nil)
}

func handleSetTabSize() tview.Primitive {
	buttons := []string{"2", "4", "8", "16", "Cancel"}
	return ShowMessageDialog("Tab Size", "Select tab size:", buttons,
		func() {
			CloseMessageDialog()
		},
		func(button string, index int) {
			CloseMessageDialog()
			if index < len(buttons)-1 { // Not cancel
				var tabSize float64
				switch index {
				case 0:
					tabSize = 2
				case 1:
					tabSize = 4
				case 2:
					tabSize = 8
				case 3:
					tabSize = 16
				}
				currentFileBuffer.buffer.Settings["tabsize"] = tabSize
				statusBar.ShowMessage("Tab size set to " + button)
			}
		})
}

func handleSetTabCharacter() tview.Primitive {
	buffer := currentFileBuffer.buffer
	buttons := []string{"Tab", "Space", "Cancel"}
	return ShowMessageDialog("Tab Character", "Select tab character:", buttons,
		func() {
			CloseMessageDialog()
		},
		func(button string, index int) {
			CloseMessageDialog()
			if index < len(buttons)-1 { // Not cancel
				switch index {
				case 0:
					buffer.Settings["tabstospaces"] = false
					statusBar.ShowMessage("Tab character set to Tab")
				case 1:
					buffer.Settings["tabstospaces"] = true
					statusBar.ShowMessage("Tab character set to Space")
				}
				syncMenuFromBuffer(buffer)
			}
		})
}

func handleSetLineEndings() tview.Primitive {
	buffer := currentFileBuffer.buffer
	buttons := []string{"LF (Unix)", "CRLF (DOS)", "Cancel"}
	return ShowMessageDialog("Line Endings", "Select line ending style:", buttons,
		func() {
			CloseMessageDialog()
		},
		func(button string, index int) {
			CloseMessageDialog()
			if index < len(buttons)-1 { // Not cancel
				switch index {
				case 0:
					buffer.Settings["fileformat"] = "unix"
					statusBar.ShowMessage("Line endings set to LF (Unix)")
				case 1:
					buffer.Settings["fileformat"] = "dos"
					statusBar.ShowMessage("Line endings set to CRLF (DOS)")
				}
				syncMenuFromBuffer(buffer)
			}
		})
}

func handleSetSyntaxHighlighting() tview.Primitive {
	var buffer *buffer.Buffer
	if currentFileBuffer.buffer != nil {
		buffer = currentFileBuffer.buffer
	}

	if buffer == nil {
		statusBar.ShowMessage("No file open")
		return nil
	}

	// Get all available syntax files
	syntaxes := smidgen.ListSyntaxes()

	// Create list items from syntax files
	items := []dialog.ListItem{}
	items = append(items, dialog.ListItem{Text: "Auto-detect", Value: ""})

	for _, syntaxName := range syntaxes {
		// Capitalize the first letter for display
		displayName := strings.Title(syntaxName)
		items = append(items, dialog.ListItem{Text: displayName, Value: syntaxName})
	}

	currentFiletype := buffer.FileType()
	if currentFiletype == "Unknown" || currentFiletype == "" {
		currentFiletype = ""
	}

	return ShowListDialog(dialog.ListDialogOptions{
		Title:           "Set Syntax Highlighting",
		Message:         "Select syntax highlighting mode:",
		Buttons:         []string{"OK", "Cancel"},
		Width:           50,
		Height:          20,
		DefaultSelected: currentFiletype,
		Items:           items,
		OnCancel: func() {
			hideListDialog()
		},
		OnAccept: func(value string, buttonIndex int) {
			hideListDialog()
			if buttonIndex == 1 { // Cancel button
				return
			}

			if value == "" {
				// Auto-detect - reset filetype and trigger detection
				buffer.Settings["filetype"] = "unknown"
				statusBar.ShowMessage("Syntax highlighting set to auto-detect")
			} else {
				// Set specific syntax
				buffer.Settings["filetype"] = value
				statusBar.ShowMessage("Syntax highlighting set to " + strings.Title(value))
			}
			buffer.UpdateRules()
			syncMenuFromBuffer(buffer)
		},
	})
}

func handleSetVerticalRuler() tview.Primitive {
	var buffer *buffer.Buffer
	if currentFileBuffer.buffer != nil {
		buffer = currentFileBuffer.buffer
	}

	currentVerticalRuler := int(buffer.Settings["colorcolumn"].(float64))
	defaultValue := ""
	if currentVerticalRuler > 0 {
		defaultValue = strconv.Itoa(currentVerticalRuler)
	}

	return ShowInputDialog("Vertical Ruler", "Enter column number (0 = off):", defaultValue,
		func() {
			// On cancel
			CloseInputDialog()
		},
		func(value string, index int) {
			CloseInputDialog()
			if index == 0 || index == -1 { // OK button or Enter key in input field
				if value == "" {
					buffer.Settings["colorcolumn"] = 0.0
					statusBar.ShowMessage("Vertical ruler disabled")
					syncMenuFromBuffer(buffer)
					return
				}

				columnNum, err := strconv.Atoi(value)
				if err != nil || columnNum < 0 {
					statusBar.ShowError("Invalid column number")
					return
				}

				buffer.Settings["colorcolumn"] = float64(columnNum)
				if columnNum == 0 {
					statusBar.ShowMessage("Vertical ruler disabled")
				} else {
					statusBar.ShowMessage("Vertical ruler set to " + value)
				}
				syncMenuFromBuffer(buffer)
			}
		}, numericInputFilter)
}

func handleQuit() tview.Primitive {
	var closeNextFileBuffer func() tview.Primitive
	closeNextFileBuffer = func() tview.Primitive {
		if len(fileBuffers) > 0 {
			fileBuffer := fileBuffers[0]
			if fileBuffer.buffer.Modified() {
				selectTab(fileBuffer.uuid)
				return ShowConfirmDialog("File has unsaved changes. Close anyway?", func() {
					closeFile(fileBuffer.uuid)
					nextFocus := closeNextFileBuffer()
					if nextFocus != nil {
						app.SetFocus(nextFocus)
					}
				}, func() {})
			} else {
				closeFile(fileBuffer.uuid)
				nextFocus := closeNextFileBuffer()
				if nextFocus != nil {
					app.SetFocus(nextFocus)
				}
			}
		} else {
			app.Stop()
		}
		return nil
	}
	return closeNextFileBuffer()
}

func handleFind() tview.Primitive {
	currentFileBuffer.openFindbar()
	return currentFileBuffer.findbar
}

// Find Next: open findbar if needed, then search forward
func handleFindNext() tview.Primitive {
	currentFileBuffer.openFindbar()
	currentFileBuffer.findbar.SearchDown()
	return nil
}

// Find Previous: open findbar if needed, then search backward
func handleFindPrevious() tview.Primitive {
	currentFileBuffer.openFindbar()
	currentFileBuffer.findbar.SearchUp()
	return nil
}

func handleNextEditor() tview.Primitive {
	if len(fileBuffers) < 2 {
		return nil
	}

	for i, fileBuffer := range fileBuffers {
		if fileBuffer.uuid == fileBufferID {
			nextIndex := (i + 1) % len(fileBuffers)
			selectTab(fileBuffers[nextIndex].uuid)
			break
		}
	}
	return nil
}

func handlePreviousEditor() tview.Primitive {
	if len(fileBuffers) < 2 {
		return nil
	}

	for i, fileBuffer := range fileBuffers {
		if fileBuffer.uuid == fileBufferID {
			prevIndex := (i - 1 + len(fileBuffers)) % len(fileBuffers)
			selectTab(fileBuffers[prevIndex].uuid)
			break
		}
	}
	return nil
}

func handleConvertTabSpaces() tview.Primitive {
	currentFileBuffer.editor.Buffer().Retab()

	expandTab, ok := currentFileBuffer.buffer.Settings["tabstospaces"]
	var msg string
	if !ok || !expandTab.(bool) {
		msg = "Converted all tabs to spaces"
	} else {
		msg = "Converted all spaces to tabs"
	}
	statusBar.ShowMessage(msg)
	return nil
}

func handleToggleWhitespace() tview.Primitive {
	buffer := currentFileBuffer.buffer
	showchars := buffer.Settings["showchars"].(string)
	on := showchars != ""
	if !on {
		buffer.Settings["showchars"] = "space=·,tab=→"
	} else {
		buffer.Settings["showchars"] = ""
	}
	syncShowWhitespace(menus, !on)
	return nil
}

func handleToggleTrailingWhitespace() tview.Primitive {
	buffer := currentFileBuffer.buffer
	on := buffer.Settings["hltrailingws"].(bool)
	buffer.Settings["hltrailingws"] = !on
	syncShowTrailingWhitespace(menus, !on)
	return nil
}

func handleFindAndReplace() tview.Primitive {
	currentFileBuffer.openFindbar()
	currentFileBuffer.findbar.Expand()
	return currentFileBuffer.findbar
}

var settingsDialog *settingsdialog.SettingsDialog

func handleSettings() tview.Primitive {
	settingsDialogName := "settings"
	if settingsDialog == nil {
		settingsDialog = settingsdialog.NewSettingsDialog(app)
		style.StyleSettingsDialog(settingsDialog)
		settingsDialog.SetCloseFunc(func() {
			modalPages.RemovePage(settingsDialogName)
		})
		settingsDialog.SetOkFunc(func(newSettings settingstype.Settings) {
			settings = newSettings
			SaveSettings(settings)
			loadEditorColorScheme(settings.ColorScheme)
		})
	}
	settingsDialog.SetSettings(settings)
	modalPages.AddPage(settingsDialogName, settingsDialog, true, true)
	return settingsDialog
}

func handleToUppercase() tview.Primitive {
	currentFileBuffer.editor.ActionController().TransformSelection(func(lines []string) []string {
		result := make([]string, len(lines))
		for i, line := range lines {
			result[i] = strings.ToUpper(line)
		}
		return result
	})
	return nil
}

func handleToLowercase() tview.Primitive {
	currentFileBuffer.editor.ActionController().TransformSelection(func(lines []string) []string {
		result := make([]string, len(lines))
		for i, line := range lines {
			result[i] = strings.ToLower(line)
		}
		return result
	})
	return nil
}
