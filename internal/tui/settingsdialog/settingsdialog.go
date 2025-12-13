package settingsdialog

import (
	"dinky/internal/application/settingstype"
	"dinky/internal/tui/scrollbar"
	"dinky/internal/tui/smidgeninputfield"
	"dinky/internal/tui/table2"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sedwards2009/smidgen"
)

type SettingsDialog struct {
	*tview.Flex
	verticalContentsFlex *tview.Flex
	closeFunc            func()
	okFunc               func(settings settingstype.Settings)
	OkButton             *tview.Button
	CancelButton         *tview.Button

	ShowLineNumbersCheckbox        *tview.Checkbox
	ShowWhitespaceCheckbox         *tview.Checkbox
	ShowMatchBracketCheckbox       *tview.Checkbox
	ShowTrailingWhitespaceCheckbox *tview.Checkbox
	SoftWrapCheckbox               *tview.Checkbox
	TabCharList                    *tview.List
	TabSizeList                    *tview.List
	VerticalRulerInputField        *smidgeninputfield.SmidgenInputField

	// Smidgen color scheme list
	ColorSchemeTableField             *table2.Table
	ColorSchemeTableFlex              *tview.Flex
	ColorSchemePreviewEditor          *smidgen.View
	ColorSchemeTableVerticalScrollbar *scrollbar.Scrollbar
	selectedColorScheme               string

	colorFiles []string

	colorSchemeItemTextColor               tcell.Color
	colorSchemeItemBackgroundColor         tcell.Color
	colorSchemeSelectedItemBackgroundColor tcell.Color
}

func NewSettingsDialog(app *tview.Application) *SettingsDialog {
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
	colorSchemeTableFlex.AddItem(colorSchemeTableField, 13, 0, false)

	colorSchemeVerticalScrollbar := scrollbar.NewScrollbar()
	colorSchemeTableFlex.AddItem(colorSchemeVerticalScrollbar, 1, 0, false)
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
	colorSchemeBuffer := smidgen.NewBufferFromString(contents, "example.go")
	colorSchemePreviewEditor := smidgen.NewView(app, colorSchemeBuffer)
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
	showLineNumbersCheckbox.SetLabel("Show Line Numbers:        ")
	firstColumnFlex.AddItem(showLineNumbersCheckbox, 1, 0, false)

	showWhitespaceCheckbox := tview.NewCheckbox()
	showWhitespaceCheckbox.SetLabel("Show Whitespace:          ")
	firstColumnFlex.AddItem(showWhitespaceCheckbox, 1, 0, false)

	showTrailingWhitespaceCheckbox := tview.NewCheckbox()
	showTrailingWhitespaceCheckbox.SetLabel("Show Trailing Whitespace: ")
	firstColumnFlex.AddItem(showTrailingWhitespaceCheckbox, 1, 0, false)

	softWrapCheckbox := tview.NewCheckbox()
	softWrapCheckbox.SetLabel("Soft Wrap:                ")
	firstColumnFlex.AddItem(softWrapCheckbox, 1, 0, false)

	showMatchBracketCheckbox := tview.NewCheckbox()
	showMatchBracketCheckbox.SetLabel("Show Match Bracket:       ")
	firstColumnFlex.AddItem(showMatchBracketCheckbox, 1, 0, false)

	// Second column of checkboxes
	secondColumnFlex := tview.NewFlex()
	secondColumnFlex.SetDirection(tview.FlexRow)

	tabCharLabel := tview.NewTextView()
	tabCharLabel.SetText("Tab Character: ")
	secondColumnFlex.AddItem(tabCharLabel, 2, 0, false)

	secondColumnFlex.AddItem(nil, 1, 0, false)

	verticalRulerInputFieldFlex := tview.NewFlex()
	verticalRulerInputFieldFlex.SetDirection(tview.FlexColumn)
	verticalRulerInputFieldFlex.SetBorder(false)

	verticalRulerInputField := smidgeninputfield.NewSmidgenInputField(app)

	verticalRulerLabel := tview.NewTextView()
	verticalRulerLabel.SetText("Vertical Ruler: ")
	secondColumnFlex.AddItem(verticalRulerLabel, 1, 0, false)

	verticalRulerLabel2 := tview.NewTextView()
	verticalRulerLabel2.SetText("(column, 0=off) ")
	secondColumnFlex.AddItem(verticalRulerLabel2, 1, 0, false)

	// Third column of options
	thirdColumnFlex := tview.NewFlex()
	thirdColumnFlex.SetDirection(tview.FlexRow)

	tabCharList := tview.NewList()
	tabCharList.ShowSecondaryText(false)
	tabCharList.AddItem(" Tab    ", "", 0, nil)
	tabCharList.AddItem(" Spaces ", "", 0, nil)
	thirdColumnFlex.AddItem(tabCharList, 2, 0, false)

	thirdColumnFlex.AddItem(nil, 1, 0, false)
	thirdColumnFlex.AddItem(verticalRulerInputField, 1, 0, false)
	verticalRulerInputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if numericInputFilter(event) {
			return event
		}
		return nil
	})

	// Fourth column of options
	forthColumnFlex := tview.NewFlex()
	forthColumnFlex.SetDirection(tview.FlexRow)

	tabSizeLabel := tview.NewTextView()
	tabSizeLabel.SetText("  Tab Size:")
	forthColumnFlex.AddItem(tabSizeLabel, 1, 0, false)

	// Fifth column of options
	fifthColumnFlex := tview.NewFlex()
	fifthColumnFlex.SetDirection(tview.FlexRow)

	tabSizeList := tview.NewList()
	tabSizeList.ShowSecondaryText(false)
	tabSizeList.AddItem(" 2  ", "", 0, nil)
	tabSizeList.AddItem(" 4  ", "", 0, nil)
	tabSizeList.AddItem(" 8  ", "", 0, nil)
	tabSizeList.AddItem(" 16 ", "", 0, nil)
	fifthColumnFlex.AddItem(tabSizeList, 4, 0, false)

	optionsColumnFlex.AddItem(firstColumnFlex, 31, 0, false)
	optionsColumnFlex.AddItem(secondColumnFlex, 16, 0, false)
	optionsColumnFlex.AddItem(thirdColumnFlex, 8, 0, false)
	optionsColumnFlex.AddItem(forthColumnFlex, 12, 0, false)
	optionsColumnFlex.AddItem(fifthColumnFlex, 0, 1, false)

	verticalContentsFlex.AddItem(optionsColumnFlex, 5, 0, false)

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
	innerFlex.AddItem(verticalContentsFlex, 23, 0, true)
	innerFlex.AddItem(nil, 0, 1, false)

	topLayout := tview.NewFlex()
	topLayout.AddItem(nil, 0, 1, false)
	topLayout.AddItem(innerFlex, 75, 1, true)
	topLayout.AddItem(nil, 0, 1, false)

	sd := &SettingsDialog{
		Flex:                           topLayout,
		verticalContentsFlex:           verticalContentsFlex,
		OkButton:                       okButton,
		CancelButton:                   cancelButton,
		ShowLineNumbersCheckbox:        showLineNumbersCheckbox,
		ShowWhitespaceCheckbox:         showWhitespaceCheckbox,
		ShowMatchBracketCheckbox:       showMatchBracketCheckbox,
		ShowTrailingWhitespaceCheckbox: showTrailingWhitespaceCheckbox,
		SoftWrapCheckbox:               softWrapCheckbox,
		TabCharList:                    tabCharList,
		TabSizeList:                    tabSizeList,
		VerticalRulerInputField:        verticalRulerInputField,

		ColorSchemeTableField:             colorSchemeTableField,
		ColorSchemeTableFlex:              colorSchemeTableFlex,
		ColorSchemePreviewEditor:          colorSchemePreviewEditor,
		ColorSchemeTableVerticalScrollbar: colorSchemeVerticalScrollbar,
	}
	okButton.SetSelectedFunc(sd.doOK)
	cancelButton.SetSelectedFunc(sd.doCancel)

	// Set up vertical scrollbar
	colorSchemeVerticalScrollbar.Track.SetBeforeDrawFunc(func(_ tcell.Screen) {
		row, _ := colorSchemeTableField.GetOffset()
		colorSchemeVerticalScrollbar.Track.SetMax(colorSchemeTableField.GetRowCount() - 1)
		_, _, _, height := sd.ColorSchemeTableField.GetInnerRect()
		colorSchemeVerticalScrollbar.Track.SetThumbSize(height)
		colorSchemeVerticalScrollbar.Track.SetPosition(row)
	})
	colorSchemeVerticalScrollbar.SetChangedFunc(func(position int) {
		_, column := sd.ColorSchemeTableField.GetOffset()
		sd.ColorSchemeTableField.SetOffset(position, column)
	})
	sd.ColorSchemeTableField.SetSelectionChangedFunc(sd.handleColorSchemeSelected)

	sd.colorFiles = smidgen.ListColorschemes()
	sd.loadColorSchemes()
	return sd
}

func (sd *SettingsDialog) loadColorSchemes() {
	sd.ColorSchemeTableField.Clear()
	cellStyle := tcell.StyleDefault.Foreground(sd.colorSchemeItemTextColor).Background(sd.colorSchemeItemBackgroundColor)
	for rowIndex, item := range sd.colorFiles {
		cell := &table2.TableCell{
			Text:  item,
			Style: cellStyle,
		}
		sd.ColorSchemeTableField.SetCell(rowIndex, 0, cell)
	}
}

func (sd *SettingsDialog) handleColorSchemeSelected(row int, column int) {
	if row < 0 || row >= len(sd.colorFiles) {
		return
	}
	colorschemeName := sd.colorFiles[row]
	if colorscheme, ok := smidgen.LoadInternalColorscheme(colorschemeName); ok {
		sd.ColorSchemePreviewEditor.SetColorscheme(colorscheme)
		sd.selectedColorScheme = colorschemeName
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
		if item == settings.ColorScheme {
			sd.ColorSchemeTableField.Select(rowIndex, 0)
			sd.selectedColorScheme = item
			sd.handleColorSchemeSelected(rowIndex, 0)
			break
		}
	}

	sd.ShowLineNumbersCheckbox.SetChecked(settings.ShowLineNumbers)
	sd.ShowWhitespaceCheckbox.SetChecked(settings.ShowWhitespace)
	sd.ShowTrailingWhitespaceCheckbox.SetChecked(settings.ShowTrailingWhitespace)
	sd.ShowMatchBracketCheckbox.SetChecked(settings.ShowMatchBracket)
	sd.SoftWrapCheckbox.SetChecked(settings.SoftWrap)
	if settings.TabCharacter == "tab" {
		sd.TabCharList.SetCurrentItem(0)
	} else {
		sd.TabCharList.SetCurrentItem(1)
	}
	switch settings.TabSize {
	case 2:
		sd.TabSizeList.SetCurrentItem(0)
	case 4:
		sd.TabSizeList.SetCurrentItem(1)
	case 8:
		sd.TabSizeList.SetCurrentItem(2)
	case 16:
		sd.TabSizeList.SetCurrentItem(3)
	default:
		sd.TabSizeList.SetCurrentItem(1)
	}

	sd.VerticalRulerInputField.SetText(strconv.Itoa(int(settings.VerticalRuler)))
}

func (sd *SettingsDialog) getSettings() settingstype.Settings {
	newSettings := settingstype.DefaultSettings()
	newSettings.ColorScheme = sd.selectedColorScheme
	newSettings.ShowLineNumbers = sd.ShowLineNumbersCheckbox.IsChecked()
	newSettings.ShowWhitespace = sd.ShowWhitespaceCheckbox.IsChecked()
	newSettings.ShowTrailingWhitespace = sd.ShowTrailingWhitespaceCheckbox.IsChecked()
	newSettings.ShowMatchBracket = sd.ShowMatchBracketCheckbox.IsChecked()
	newSettings.SoftWrap = sd.SoftWrapCheckbox.IsChecked()
	tabCharIndex := sd.TabCharList.GetCurrentItem()
	if tabCharIndex == 0 {
		newSettings.TabCharacter = "tab"
	} else {
		newSettings.TabCharacter = "space"
	}

	tabSizeIndex := sd.TabSizeList.GetCurrentItem()
	switch tabSizeIndex {
	case 0:
		newSettings.TabSize = 2
	case 1:
		newSettings.TabSize = 4
	case 2:
		newSettings.TabSize = 8
	case 3:
		newSettings.TabSize = 16
	default:
		newSettings.TabSize = 4
	}

	value, _ := strconv.Atoi(sd.VerticalRulerInputField.GetText())
	newSettings.VerticalRuler = float64(value)
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

func numericInputFilter(event *tcell.EventKey) bool {
	key := event.Key()
	// Allow digits and basic editing keys
	if key == tcell.KeyBackspace || key == tcell.KeyDelete ||
		key == tcell.KeyLeft || key == tcell.KeyRight ||
		key == tcell.KeyHome || key == tcell.KeyEnd ||
		key == tcell.KeyDEL {
		return true
	}
	if event.Rune() >= '0' && event.Rune() <= '9' {
		return true
	}
	return false
}
