package application

import (
	"dinky/internal/tui/dialog"
	"dinky/internal/tui/filedialog"
	"dinky/internal/tui/style"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/renameio/v2"
	"github.com/rivo/tview"
	"github.com/sedwards2009/femto"
	"github.com/sedwards2009/femto/runtime"
)

const (
	ACTION_NEW                     = "NewFile"
	ACTION_CLOSE_FILE              = "CloseFile"
	ACTION_OPEN_FILE               = "OpenFile"
	ACTION_SAVE_FILE               = "SaveFile"
	ACTION_SAVE_FILE_AS            = "SaveFileAs"
	ACTION_OPEN_MENU               = "OpenMenu"
	ACTION_TOGGLE_SOFT_WRAP        = "ToggleSoftWrap"
	ACTION_TOGGLE_MATCH_BRACKET    = "ToggleMatchBracket"
	ACTION_SET_TAB_SIZE            = "SetTabSize"
	ACTION_SET_TAB_CHARACTER       = "SetTabCharacter"
	ACTION_SET_LINE_ENDINGS        = "SetLineEndings"
	ACTION_SET_SYNTAX_HIGHLIGHTING = "SetSyntaxHighlighting"
	ACTION_GO_TO_LINE              = "GoToLine"
	ACTION_QUIT                    = "Quit"
	ACTION_FIND                    = "Find"
	ACTION_FIND_NEXT               = "FindNext"
	ACTION_FIND_PREVIOUS           = "FindPrevious"
	ACTION_ABOUT                   = "About"
)

var dinkyActionMapping map[string]func() tview.Primitive

func init() {
	dinkyActionMapping = map[string]func() tview.Primitive{
		ACTION_NEW:                     handleNewFile,
		ACTION_CLOSE_FILE:              handleCloseFile,
		ACTION_OPEN_FILE:               handleOpenFile,
		ACTION_OPEN_MENU:               handleOpenMenu,
		ACTION_SAVE_FILE:               handleSaveFile,
		ACTION_SAVE_FILE_AS:            handleSaveFileAs,
		ACTION_TOGGLE_SOFT_WRAP:        handleSoftWrap,
		ACTION_TOGGLE_MATCH_BRACKET:    handleMatchBracket,
		ACTION_SET_TAB_SIZE:            handleSetTabSize,
		ACTION_SET_TAB_CHARACTER:       handleSetTabCharacter,
		ACTION_SET_LINE_ENDINGS:        handleSetLineEndings,
		ACTION_SET_SYNTAX_HIGHLIGHTING: handleSetSyntaxHighlighting,
		ACTION_GO_TO_LINE:              handleGoToLine,
		ACTION_QUIT:                    handleQuit,
		ACTION_ABOUT:                   handleAbout,
		ACTION_FIND:                    handleFind,
		ACTION_FIND_NEXT:               handleFindNext,
		ACTION_FIND_PREVIOUS:           handleFindPrevious,
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

func showFileDialog(title string, mode filedialog.FileDialogMode, defaultPath string, completedFunc func(accepted bool,
	filePath string)) tview.Primitive {

	if fileDialog == nil {
		fileDialog = filedialog.NewFileDialog(app)
		style.StyleFileDialog(fileDialog)
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

func bufferToBytes(buffer *femto.Buffer) []byte {
	return buffer.LineArray.Bytes(isBufferCRLF(buffer))
}

func writeFile(filename string, buffer *femto.Buffer) (ok bool, message string) {
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
	fileBuffer.buffer.IsModified = false
	if ok {
		statusBar.ShowMessage(message)
	} else {
		statusBar.ShowError(message)
	}
}

func handleSaveFileAs() tview.Primitive {
	fileBuffer := getFileBufferByID(fileBufferID)
	return showFileDialog("Save File As", filedialog.SAVE_FILE_MODE, fileBuffer.filename, func(accepted bool, filePath string) {
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
	if fileBuffer.buffer.IsModified {
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

func handleOpenMenu() tview.Primitive {
	menuBar.Open()
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

func handleFemtoAction(id string) tview.Primitive {
	if f, ok := femto.BindingActionsMapping[id]; ok {
		f(currentFileBuffer.editor)
		if id == femto.ActionToggleRuler {
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
	var buffer *femto.Buffer
	if currentFileBuffer.buffer != nil {
		buffer = currentFileBuffer.buffer
	}

	if buffer == nil {
		statusBar.ShowMessage("No file open")
		return nil
	}

	// Get all available syntax files
	syntaxFiles := runtime.Files.ListRuntimeFiles(femto.RTSyntax)

	// Create list items from syntax files
	items := []dialog.ListItem{}
	items = append(items, dialog.ListItem{Text: "Auto-detect", Value: ""})

	for _, file := range syntaxFiles {
		// Extract the name without the .yaml extension
		name := strings.TrimSuffix(file.Name(), ".yaml")
		// Capitalize the first letter for display
		displayName := strings.Title(name)
		items = append(items, dialog.ListItem{Text: displayName, Value: name})
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
				buffer.Settings["filetype"] = "Unknown"
				currentFileBuffer.editor.SetRuntimeFiles(runtime.Files)
				statusBar.ShowMessage("Syntax highlighting set to auto-detect")
			} else {
				// Set specific syntax
				buffer.Settings["filetype"] = value
				currentFileBuffer.editor.SetRuntimeFiles(runtime.Files)
				statusBar.ShowMessage("Syntax highlighting set to " + strings.Title(value))
			}
			syncMenuFromBuffer(buffer)
		},
	})
}

func handleQuit() tview.Primitive {
	var closeNextFileBuffer func() tview.Primitive
	closeNextFileBuffer = func() tview.Primitive {
		if len(fileBuffers) > 0 {
			fileBuffer := fileBuffers[0]
			if fileBuffer.buffer.IsModified {
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
