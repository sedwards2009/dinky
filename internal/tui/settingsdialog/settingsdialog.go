package settingsdialog

import (
	"dinky/internal/application/settingstype"
	"dinky/internal/tui/scrollbar"
	"dinky/internal/tui/table2"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sedwards2009/femto"
	"github.com/sedwards2009/femto/runtime"
)

type SettingsDialog struct {
	*tview.Flex
	verticalContentsFlex *tview.Flex
	closeFunc            func()
	okFunc               func(settings settingstype.Settings)
	OkButton             *tview.Button
	CancelButton         *tview.Button

	ShowLineNumbersCheckbox  *tview.Checkbox
	ShowWhitespaceCheckbox   *tview.Checkbox
	ShowMatchBracketCheckbox *tview.Checkbox
	SoftWrapCheckbox         *tview.Checkbox
	TabCharDropdown          *tview.DropDown
	TabSizeDropdown          *tview.DropDown

	// Femto color scheme list
	ColorSchemeTableField             *table2.Table
	ColorSchemeTableFlex              *tview.Flex
	ColorSchemePreviewEditor          *femto.View
	ColorSchemeTableVerticalScrollbar *scrollbar.Scrollbar
	selectedColorScheme               string

	colorFiles []femto.RuntimeFile

	colorSchemeItemTextColor               tcell.Color
	colorSchemeItemBackgroundColor         tcell.Color
	colorSchemeSelectedItemBackgroundColor tcell.Color
}

func NewSettingsDialog() *SettingsDialog {
	verticalContentsFlex := tview.NewFlex()
	verticalContentsFlex.Box = tview.NewBox() // Nasty hack to clear the `dontClear` flag inside Box.
	verticalContentsFlex.Box.Primitive = verticalContentsFlex
	verticalContentsFlex.SetTitle("Settings")
	verticalContentsFlex.SetTitleAlign(tview.AlignLeft)
	verticalContentsFlex.SetBorder(true)
	verticalContentsFlex.SetDirection(tview.FlexRow)
	verticalContentsFlex.SetBorderPadding(1, 1, 1, 1)

	colorSchemeLabel := tview.NewTextView()
	colorSchemeLabel.SetText("Color Scheme:")
	verticalContentsFlex.AddItem(colorSchemeLabel, 1, 0, false)

	colorSchemeTableFlex := tview.NewFlex()
	colorSchemeTableFlex.SetDirection(tview.FlexColumn)
	colorSchemeTableFlex.SetBorder(false)
	colorSchemeTableField := table2.NewTable()
	colorSchemeTableField.SetSelectable(true, false)
	colorSchemeTableFlex.AddItem(colorSchemeTableField, 12, 0, false)

	verticalScrollbar := scrollbar.NewScrollbar()
	colorSchemeTableFlex.AddItem(verticalScrollbar, 1, 0, false)
	colorSchemeTableFlex.AddItem(nil, 1, 0, false)

	// Color scheme preview editor
	contents := `package main
import "fmt"

func main() {
	for i := 0; i < 5; i++ {
		fmt.Printf("Hello, Go! Iteration %d\n", i)
	}
}
`
	colorSchemeBuffer := femto.NewBufferFromString(contents, "example.go")
	colorSchemePreviewEditor := femto.NewView(colorSchemeBuffer)
	colorSchemePreviewEditor.SetRuntimeFiles(runtime.Files)
	colorSchemeBuffer.Settings["matchbrace"] = false
	colorSchemeBuffer.Settings["ruler"] = false

	colorSchemeTableFlex.AddItem(colorSchemePreviewEditor, 0, 1, false)

	verticalContentsFlex.AddItem(colorSchemeTableFlex, 0, 1, false)

	verticalContentsFlex.AddItem(nil, 1, 0, false)

	defaultsLabel := tview.NewTextView()
	defaultsLabel.SetText("Defaults")
	verticalContentsFlex.AddItem(defaultsLabel, 1, 0, false)

	verticalContentsFlex.AddItem(nil, 1, 0, false)

	optionsColumnFlex := tview.NewFlex()
	optionsColumnFlex.SetDirection(tview.FlexColumn)

	// First column of checkboxes
	firstColumnFlex := tview.NewFlex()
	firstColumnFlex.SetDirection(tview.FlexRow)

	showLineNumbersCheckbox := tview.NewCheckbox()
	showLineNumbersCheckbox.SetLabel("Show Line Numbers: ")
	firstColumnFlex.AddItem(showLineNumbersCheckbox, 1, 0, false)

	showWhitespaceCheckbox := tview.NewCheckbox()
	showWhitespaceCheckbox.SetLabel("Show Whitespace:   ")
	firstColumnFlex.AddItem(showWhitespaceCheckbox, 1, 0, false)

	softWrapCheckbox := tview.NewCheckbox()
	softWrapCheckbox.SetLabel("Soft Wrap:         ")
	firstColumnFlex.AddItem(softWrapCheckbox, 1, 0, false)

	optionsColumnFlex.AddItem(firstColumnFlex, 0, 1, false)

	// Second column of checkboxes
	secondColumnFlex := tview.NewFlex()
	secondColumnFlex.SetDirection(tview.FlexRow)

	showMatchBracketCheckbox := tview.NewCheckbox()
	showMatchBracketCheckbox.SetLabel("Show Match Bracket: ")
	secondColumnFlex.AddItem(showMatchBracketCheckbox, 1, 0, false)

	tabCharDropdown := tview.NewDropDown()
	tabCharDropdown.SetLabel("Tab Character:      ")
	tabCharDropdown.SetOptions([]string{"Tab", "Spaces"}, nil)
	secondColumnFlex.AddItem(tabCharDropdown, 1, 0, false)

	tabSizeDropdown := tview.NewDropDown()
	tabSizeDropdown.SetLabel("Tab Size:           ")
	tabSizeDropdown.SetOptions([]string{"2", "4", "8", "16"}, nil)
	secondColumnFlex.AddItem(tabSizeDropdown, 1, 0, false)

	optionsColumnFlex.AddItem(secondColumnFlex, 0, 1, false)

	verticalContentsFlex.AddItem(optionsColumnFlex, 3, 0, false)

	verticalContentsFlex.AddItem(nil, 1, 0, false)

	buttonFlex := tview.NewFlex()
	buttonFlex.SetDirection(tview.FlexColumn)

	buttonFlex.AddItem(nil, 0, 1, false)
	okButton := tview.NewButton("OK")
	buttonFlex.AddItem(okButton, 10, 0, false)
	buttonFlex.AddItem(nil, 1, 0, false)
	cancelButton := tview.NewButton("Cancel")
	buttonFlex.AddItem(cancelButton, 10, 0, false)

	verticalContentsFlex.AddItem(buttonFlex, 1, 0, false)

	innerFlex := tview.NewFlex()
	innerFlex.SetDirection(tview.FlexRow)
	innerFlex.AddItem(nil, 0, 1, false)
	innerFlex.AddItem(verticalContentsFlex, 22, 0, true)
	innerFlex.AddItem(nil, 0, 1, false)

	topLayout := tview.NewFlex()
	topLayout.AddItem(nil, 0, 1, false)
	topLayout.AddItem(innerFlex, 70, 1, true)
	topLayout.AddItem(nil, 0, 1, false)

	sd := &SettingsDialog{
		Flex:                     topLayout,
		verticalContentsFlex:     verticalContentsFlex,
		OkButton:                 okButton,
		CancelButton:             cancelButton,
		ShowLineNumbersCheckbox:  showLineNumbersCheckbox,
		ShowWhitespaceCheckbox:   showWhitespaceCheckbox,
		ShowMatchBracketCheckbox: showMatchBracketCheckbox,
		SoftWrapCheckbox:         softWrapCheckbox,
		TabCharDropdown:          tabCharDropdown,
		TabSizeDropdown:          tabSizeDropdown,

		ColorSchemeTableField:             colorSchemeTableField,
		ColorSchemeTableFlex:              colorSchemeTableFlex,
		ColorSchemePreviewEditor:          colorSchemePreviewEditor,
		ColorSchemeTableVerticalScrollbar: verticalScrollbar,
	}
	okButton.SetSelectedFunc(sd.doOK)
	cancelButton.SetSelectedFunc(sd.doCancel)

	// Set up vertical scrollbar
	verticalScrollbar.Track.SetBeforeDrawFunc(func(_ tcell.Screen) {
		row, _ := colorSchemeTableField.GetOffset()
		verticalScrollbar.Track.SetMax(colorSchemeTableField.GetRowCount() - 1)
		_, _, _, height := sd.ColorSchemeTableField.GetInnerRect()
		verticalScrollbar.Track.SetThumbSize(height)
		verticalScrollbar.Track.SetPosition(row)
	})
	verticalScrollbar.SetChangedFunc(func(position int) {
		_, column := sd.ColorSchemeTableField.GetOffset()
		sd.ColorSchemeTableField.SetOffset(position, column)
	})
	sd.ColorSchemeTableField.SetSelectionChangedFunc(sd.handleColorSchemeSelected)

	sd.colorFiles = runtime.Files.ListRuntimeFiles(femto.RTColorscheme)
	sd.loadColorSchemes()
	return sd
}

func (sd *SettingsDialog) loadColorSchemes() {
	sd.ColorSchemeTableField.Clear()
	cellStyle := tcell.StyleDefault.Foreground(sd.colorSchemeItemTextColor).Background(sd.colorSchemeItemBackgroundColor)
	for rowIndex, item := range sd.colorFiles {
		cell := &table2.TableCell{
			Text:  item.Name(),
			Style: cellStyle,
		}
		sd.ColorSchemeTableField.SetCell(rowIndex, 0, cell)
	}
}

func (sd *SettingsDialog) handleColorSchemeSelected(row int, column int) {
	if row < 0 || row >= len(sd.colorFiles) {
		return
	}
	colorschemeFile := sd.colorFiles[row]
	if data, err := colorschemeFile.Data(); err == nil {
		colorscheme := femto.ParseColorscheme(string(data))
		sd.ColorSchemePreviewEditor.SetColorscheme(colorscheme)
		sd.selectedColorScheme = colorschemeFile.Name()
	}
}

func (sd *SettingsDialog) SetOkFunc(okFunc func(settings settingstype.Settings)) {
	sd.okFunc = okFunc
}

func (sd *SettingsDialog) SetCloseFunc(closeFunc func()) {
	sd.closeFunc = closeFunc
}

func (sd *SettingsDialog) SetSettings(settings settingstype.Settings) {
	for rowIndex, item := range sd.colorFiles {
		if item.Name() == settings.ColorScheme {
			sd.ColorSchemeTableField.Select(rowIndex, 0)
			sd.selectedColorScheme = item.Name()
			sd.handleColorSchemeSelected(rowIndex, 0)
			break
		}
	}

	sd.ShowLineNumbersCheckbox.SetChecked(settings.ShowLineNumbers)
	sd.ShowWhitespaceCheckbox.SetChecked(settings.ShowWhitespace)
	sd.ShowMatchBracketCheckbox.SetChecked(settings.ShowMatchBracket)
	sd.SoftWrapCheckbox.SetChecked(settings.SoftWrap)
	if settings.TabCharacter == "tab" {
		sd.TabCharDropdown.SetCurrentOption(0)
	} else {
		sd.TabCharDropdown.SetCurrentOption(1)
	}
	switch settings.TabSize {
	case 2:
		sd.TabSizeDropdown.SetCurrentOption(0)
	case 4:
		sd.TabSizeDropdown.SetCurrentOption(1)
	case 8:
		sd.TabSizeDropdown.SetCurrentOption(2)
	case 16:
		sd.TabSizeDropdown.SetCurrentOption(3)
	default:
		sd.TabSizeDropdown.SetCurrentOption(1)
	}
}

func (sd *SettingsDialog) getSettings() settingstype.Settings {
	newSettings := settingstype.DefaultSettings()
	newSettings.ColorScheme = sd.selectedColorScheme
	newSettings.ShowLineNumbers = sd.ShowLineNumbersCheckbox.IsChecked()
	newSettings.ShowWhitespace = sd.ShowWhitespaceCheckbox.IsChecked()
	newSettings.ShowMatchBracket = sd.ShowMatchBracketCheckbox.IsChecked()
	newSettings.SoftWrap = sd.SoftWrapCheckbox.IsChecked()
	_, tabChar := sd.TabCharDropdown.GetCurrentOption()
	if tabChar == "Tab" {
		newSettings.TabCharacter = "tab"
	} else {
		newSettings.TabCharacter = "space"
	}
	_, tabSizeStr := sd.TabSizeDropdown.GetCurrentOption()
	switch tabSizeStr {
	case "2":
		newSettings.TabSize = 2
	case "4":
		newSettings.TabSize = 4
	case "8":
		newSettings.TabSize = 8
	case "16":
		newSettings.TabSize = 16
	default:
		newSettings.TabSize = 4
	}
	return newSettings
}

func (sd *SettingsDialog) doOK() {
	if sd.okFunc != nil {
		sd.okFunc(sd.getSettings())
	}
	sd.doCancel()
}

func (sd *SettingsDialog) doCancel() {
	if sd.closeFunc != nil {
		sd.closeFunc()
	}
}

func (sd *SettingsDialog) SetItemTextColor(color tcell.Color) {
	sd.colorSchemeItemTextColor = color
	sd.ColorSchemeTableField.SetSelectedStyle(
		tcell.StyleDefault.Background(sd.colorSchemeSelectedItemBackgroundColor).Foreground(sd.colorSchemeItemTextColor))
	sd.loadColorSchemes()
}

func (sd *SettingsDialog) SetItemBackgroundColor(color tcell.Color) {
	sd.colorSchemeItemBackgroundColor = color
	sd.loadColorSchemes()
}

func (sd *SettingsDialog) SetSelectedItemBackgroundColor(color tcell.Color) {
	sd.colorSchemeSelectedItemBackgroundColor = color
	sd.ColorSchemeTableField.SetSelectedStyle(
		tcell.StyleDefault.Background(sd.colorSchemeSelectedItemBackgroundColor).Foreground(sd.colorSchemeItemTextColor))
}

func (sd *SettingsDialog) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return sd.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		// This exists to prevent elements behind this dialog from receiving mouse events.
		sd.verticalContentsFlex.MouseHandler()(action, event, setFocus)
		return true, nil
	})
}
