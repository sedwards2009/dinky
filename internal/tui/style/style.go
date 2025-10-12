package style

import (
	"dinky/internal/tui/dialog"
	"dinky/internal/tui/femtoinputfield"
	"dinky/internal/tui/filedialog"
	"dinky/internal/tui/filelist"
	"dinky/internal/tui/findbar"
	"dinky/internal/tui/menu"
	"dinky/internal/tui/scrollbar"
	"dinky/internal/tui/stylecolor"
	"dinky/internal/tui/tabbar"
	"dinky/internal/tui/table2"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func StyleButton(button *tview.Button) {
	button.SetActivatedStyle(tcell.StyleDefault.Background(stylecolor.ButtonBackgroundFocusedColor).Foreground(stylecolor.ButtonLabelFocusedColor))
	button.SetStyle(tcell.StyleDefault.Background(stylecolor.ButtonBackgroundColor).Foreground(stylecolor.ButtonLabelColor))
	button.SetDisabledStyle(tcell.StyleDefault.Background(stylecolor.ButtonBackgroundDisabledColor).Foreground(stylecolor.ButtonLabelDisabledColor))
}

func StyleCheckbox(checkbox *tview.Checkbox) {
	checkbox.SetLabelStyle(stylecolor.CheckboxLabelStyle)
	checkbox.SetCheckedString(stylecolor.CheckboxCheckedString)
	checkbox.SetUncheckedString(stylecolor.CheckboxUncheckedString)
	checkbox.SetCheckedStyle(stylecolor.CheckboxCheckedStyle)
	checkbox.SetUncheckedStyle(stylecolor.CheckboxUncheckedStyle)
	checkbox.SetActivatedStyle(stylecolor.CheckboxFocusStyle)
}

func StyleScrollbarTrack(scrollbarTrack *scrollbar.ScrollbarTrack) {
	scrollbarTrack.SetTrackColor(stylecolor.DarkGray)
	scrollbarTrack.SetThumbColor(stylecolor.White)
	scrollbarTrack.SetWidth(1)
}

func StyleScrollbar(scrollbar *scrollbar.Scrollbar) {
	StyleScrollbarTrack(scrollbar.Track)
	StyleButton(scrollbar.UpButton)
	StyleButton(scrollbar.DownButton)
}

func StyleFileList(fileList *filelist.FileList) {
	fileList.SetTextColor(stylecolor.White)
	fileList.SetBackgroundColor(stylecolor.Black)
	fileList.SetSelectedBackgroundColor(stylecolor.Blue)
	fileList.SetHeaderLabelColor(stylecolor.Black)
	fileList.SetHeaderBackgroundColor(stylecolor.White)

	StyleScrollbar(fileList.VerticalScrollbar)
	StyleScrollbar(fileList.HorizontalScrollbar)
}

func StyleInputField(inputField *tview.InputField) {
	inputField.SetFieldStyle(tcell.StyleDefault.Background(stylecolor.InputFieldFieldBackgroundColor).Foreground(stylecolor.InputFieldFieldTextColor))
	inputField.SetLabelStyle(tcell.StyleDefault.Foreground(stylecolor.InputFieldLabelColor))
}

func StyleFemtoInputField(femtoInputField *femtoinputfield.FemtoInputField) {
	femtoInputField.SetTextColor(stylecolor.InputFieldFieldTextColor, stylecolor.InputFieldFieldBackgroundColor)
}

func StyleTabBar(tabBar *tabbar.TabBar) {
	tabBar.ActiveTabStyle = tcell.StyleDefault.Foreground(stylecolor.White).Background(stylecolor.Black).Bold(true)
	tabBar.InactiveTabStyle = tcell.StyleDefault.Foreground(stylecolor.LightGray).Background(stylecolor.DarkGray).Bold(false)
	tabBar.BackgroundStyle = tcell.StyleDefault.Foreground(stylecolor.White).Background(stylecolor.Blue)
}

func StyleMenuBar(menuBar *menu.MenuBar) {
	menuBar.MenuBarStyle = tcell.StyleDefault.Foreground(stylecolor.White).Background(stylecolor.Blue)
	menuBar.MenuStyle = tcell.StyleDefault.Foreground(stylecolor.Black).Background(stylecolor.LightGray)
	menuBar.MenuSelectedStyle = tcell.StyleDefault.Foreground(stylecolor.White).Background(stylecolor.Blue)
}

func StyleTable(table *table2.Table) {
	table.SetBackgroundColor(stylecolor.Black)
	table.SetSelectedStyle(tcell.StyleDefault.Background(stylecolor.Blue).Foreground(stylecolor.White))
}

func StyleTableCell(cell *table2.TableCell) {
	cell.SetStyle(tcell.StyleDefault.Foreground(stylecolor.White).Background(stylecolor.Black))
}

func StyleMessageDialog(messageDialog *dialog.MessageDialog) {
	messageDialog.SetBackgroundColor(stylecolor.LightGray)
	for _, button := range messageDialog.Buttons {
		StyleButton(button)
	}
}

func StyleInputDialog(inputDialog *dialog.InputDialog) {
	StyleFemtoInputField(inputDialog.InputField)
	inputDialog.SetBackgroundColor(stylecolor.LightGray)
	for _, button := range inputDialog.Buttons {
		StyleButton(button)
	}
}

func StyleListDialog(d *dialog.ListDialog) {
	d.SetBackgroundColor(stylecolor.LightGray)

	d.SetItemTextColor(stylecolor.White)
	d.SetItemBackgroundColor(stylecolor.Black)
	d.SetSelectedItemBackgroundColor(stylecolor.Blue)

	for _, button := range d.Buttons {
		StyleButton(button)
	}
	StyleTable(d.TableField)
	StyleScrollbar(d.VerticalScrollbar)
}

func StyleFileDialog(fileDialog *filedialog.FileDialog) {
	fileDialog.SetBackgroundColor(stylecolor.LightGray)
	StyleFemtoInputField(fileDialog.DirectoryField)
	StyleFemtoInputField(fileDialog.FilenameField)
	StyleFileList(fileDialog.FileList)
	StyleButton(fileDialog.ActionButton)
	StyleButton(fileDialog.CancelButton)
	StyleButton(fileDialog.ParentButton)
	StyleCheckbox(fileDialog.ShowHiddenCheckbox)
}

func StyleFindbar(findBar *findbar.Findbar) {
	findBar.SetBackgroundColor(stylecolor.LightGray)
	StyleFemtoInputField(findBar.SearchStringField)
	StyleButton(findBar.SearchUpButton)
	StyleButton(findBar.SearchDownButton)
	StyleButton(findBar.CloseButton)

	StyleCheckbox(findBar.RegexCheckbox)
	findBar.RegexCheckbox.SetCheckedString("[✓Regex ]")
	findBar.RegexCheckbox.SetUncheckedString("[ Regex ]")

	StyleCheckbox(findBar.CaseSensitiveCheckbox)
	findBar.CaseSensitiveCheckbox.SetCheckedString("[✓Aa ]")
	findBar.CaseSensitiveCheckbox.SetUncheckedString("[ Aa ]")
}
