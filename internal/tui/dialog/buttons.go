package dialog

import (
	nuview "github.com/rivo/tview"
)

// createButtonsRow creates and configures the buttons for the message dialog
func createButtonsRow(buttonsFlex *nuview.Flex, buttons []string, onButtonClick func(button string, index int)) []*nuview.Button {
	buttonsFlex.Clear()
	nuviewButtons := make([]*nuview.Button, len(buttons))

	if len(buttons) == 1 {
		// If it is just one button, then center it.
		btn := nuview.NewButton(buttons[0])
		nuviewButtons[0] = btn
		btn.SetSelectedFunc(func() {
			if onButtonClick != nil {
				onButtonClick(buttons[0], 0)
			}
		})
		buttonsFlex.AddItem(nil, 0, 1, false)
		buttonsFlex.AddItem(btn, 0, 2, true)
		buttonsFlex.AddItem(nil, 0, 1, false)
	} else {
		for i, button := range buttons {
			btn := nuview.NewButton(button)
			nuviewButtons[i] = btn
			btn.SetSelectedFunc(func() {
				if onButtonClick != nil {
					onButtonClick(button, i)
				}
			})

			buttonsFlex.AddItem(btn, 0, 2, i == 0)
			if i < len(buttons)-1 {
				buttonsFlex.AddItem(nil, 0, 1, false)
			}
		}
	}

	return nuviewButtons
}
