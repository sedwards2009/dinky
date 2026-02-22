package filtercommandaction

import (
	"bytes"
	"dinky/internal/tui/filterdialog"
	"dinky/internal/tui/statusbar"
	"dinky/internal/tui/style"
	"os/exec"
	"strings"

	"github.com/rivo/tview"
	"github.com/sedwards2009/smidgen"
)

var filterDialog *filterdialog.FilterDialog

const filterDialogName = "filterDialog"

func HandleFilterExternalCommand(app *tview.Application, modalPages *tview.Pages, editor *smidgen.View,
	smidgenSingleLineKeyBindings smidgen.Keybindings, statusBar *statusbar.StatusBar) tview.Primitive {

	if !editor.Cursor().HasSelection() {
		statusBar.ShowWarning("No text selected")
		return nil
	}

	return showFilterDialog(app, modalPages, smidgenSingleLineKeyBindings,
		func() {
			// On cancel
			closeFilterDialog(modalPages)
		},
		func(value string, index int) {
			closeFilterDialog(modalPages)
			if index == 1 {
				return
			}

			selectionBytes := editor.Cursor().GetSelection()
			// Run external command with selection as stdin
			output, err := runExternalShellCommandWithInput(value, selectionBytes)
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
	smidgenSingleLineKeyBindings smidgen.Keybindings,
	onCancel func(), onAccept func(value string, index int)) tview.Primitive {

	if filterDialog == nil {
		filterDialog = filterdialog.NewFilterDialog(app)
		filterDialog.SetSmidgenKeybindings(smidgenSingleLineKeyBindings)
	}
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

func runExternalShellCommandWithInput(command string, input []byte) ([]byte, error) {
	// Run the command via `sh` -c to allow for shell features like pipes and redirection
	cmd := exec.Command("sh", "-c", command)

	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	return out.Bytes(), err
}
