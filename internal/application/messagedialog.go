package application

import "dinky/internal/tui/messagedialog"

var messageDialog *messagedialog.MessageDialog

const messageDialogName = "messagedialog"

func ShowMessageDialog(title string, message string, buttons []string, width int, height int, OnClose func(),
	OnButtonClick func(button string, index int)) {

	if messageDialog == nil {
		messageDialog = messagedialog.NewMessageDialog(app)
	}
	modalPages.AddPanel(messageDialogName, messageDialog, true, true)
	messageDialog.OnClose = OnClose
	messageDialog.OnButtonClick = OnButtonClick
	messageDialog.Open(title, message, buttons, width, height)
	messageDialog.FocusButton(0)
}

func CloseMessageDialog() {
	if messageDialog != nil {
		messageDialog.Close()
		modalPages.RemovePanel("messagedialog")
	}
}

func ShowConfirmDialog(message string, onConfirm func(), onCancel func()) {
	ShowMessageDialog("Confirm", message, []string{"OK", "Cancel"}, 50, 7,
		func() {
			CloseMessageDialog()
			onCancel()
		},
		func(button string, index int) {
			CloseMessageDialog()
			if index == 0 {
				onConfirm()
			} else {
				onCancel()
			}
		})
}

func ShowOkDialog(title string, message string, width int, height int, onClose func()) {
	ShowMessageDialog(title, message, []string{"OK"}, width, height,
		func() {
			CloseMessageDialog()
			onClose()
		},
		func(button string, index int) {
			CloseMessageDialog()
			onClose()
		})
}
