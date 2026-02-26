package filtercommandaction

import (
	"bytes"
	"dinky/internal/tui/filterdialog"
	"dinky/internal/tui/statusbar"
	"dinky/internal/tui/style"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/rivo/tview"
	"github.com/sedwards2009/smidgen"
)

var filterDialog *filterdialog.FilterDialog

const filterDialogName = "filterDialog"

var recentFilterCommands []string
var recentFilterDirectories []string

func HandleFilterExternalCommand(app *tview.Application, modalPages *tview.Pages, editor *smidgen.View,
	smidgenSingleLineKeyBindings smidgen.Keybindings, statusBar *statusbar.StatusBar) tview.Primitive {

	if !editor.Cursor().HasSelection() {
		statusBar.ShowWarning("No text selected")
		return nil
	}

	return showFilterDialog(app, modalPages, smidgenSingleLineKeyBindings, recentFilterCommands,
		func() {
			// On cancel
			closeFilterDialog(modalPages)
		},
		func(command string, directory string, index int) {
			closeFilterDialog(modalPages)
			if index == 1 {
				return
			}

			if !slices.Contains(recentFilterCommands, command) {
				recentFilterCommands = append(recentFilterCommands, command)
				if len(recentFilterCommands) > 10 {
					recentFilterCommands = recentFilterCommands[1:]
				}
			}

			if !slices.Contains(recentFilterDirectories, directory) {
				recentFilterDirectories = append(recentFilterDirectories, directory)
				if len(recentFilterDirectories) > 10 {
					recentFilterDirectories = recentFilterDirectories[1:]
				}
			}

			selectionBytes := editor.Cursor().GetSelection()

			// Check if the directory is ok before trying to run the command.
			if directory != "" {
				info, err := os.Stat(directory)
				if err != nil {
					statusBar.ShowError("Directory '" + directory + "'isn't valid: " + err.Error())
					return
				} else {
					if !info.IsDir() {
						statusBar.ShowError("Directory '" + directory + "' isn't valid: not a directory")
						return
					}
				}
			}

			// Run external command with selection as stdin
			output, err := runExternalShellCommandWithInput(command, directory, selectionBytes)
			if err != nil {
				statusBar.ShowError("Error running shell command: " + err.Error())
				return
			}

			editor.ActionController().TransformSelection(func(lines []string) []string {
				stringOutput := string(output)
				return strings.Split(strings.TrimRight(stringOutput, "\n"), "\n")
			})
		})
}

func showFilterDialog(app *tview.Application, modalPages *tview.Pages,
	smidgenSingleLineKeyBindings smidgen.Keybindings, recentFilterCommands []string,
	onCancel func(), onAccept func(command string, directory string, index int)) tview.Primitive {

	if filterDialog == nil {
		filterDialog = filterdialog.NewFilterDialog(app)
		filterDialog.SetSmidgenKeybindings(smidgenSingleLineKeyBindings)
	}
	filterDialog.SetRecentCommands(recentFilterCommands)
	filterDialog.SetRecentDirectories(recentFilterDirectories)
	modalPages.AddPage(filterDialogName, filterDialog, true, true)

	options := filterdialog.FilterDialogOptions{
		OnCancel: onCancel,
		OnAccept: onAccept,
	}
	filterDialog.Open(options)
	style.StyleFilterDialog(filterDialog)
	return filterDialog
}

func closeFilterDialog(modalPages *tview.Pages) {
	if filterDialog != nil {
		filterDialog.Close()
		modalPages.RemovePage(filterDialogName)
	}
}

func runExternalShellCommandWithInput(command string, directory string, input []byte) ([]byte, error) {
	// Run the command via `sh` -c to allow for shell features like pipes and redirection
	cmd := exec.Command("sh", "-c", command)
	if directory != "" {
		cmd.Dir = directory
	}

	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	return out.Bytes(), err
}
