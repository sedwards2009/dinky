package application

import (
	"bytes"
	"dinky/internal/application/settingstype"
	"dinky/internal/gpm"
	"dinky/internal/tui/findbar"
	"dinky/internal/tui/menu"
	"dinky/internal/tui/scrollbar"
	"dinky/internal/tui/statusbar"
	"dinky/internal/tui/style"
	"dinky/internal/tui/tabbar"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"slices"

	"runtime/debug"

	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"
	"github.com/sedwards2009/smidgen"
	"github.com/sedwards2009/smidgen/micro/buffer"
	"github.com/sedwards2009/smidgen/micro/display"
	"github.com/sedwards2009/smidgen/micro/util"
)

// -----------------------------------------------------------------
var app *tview.Application
var gpmClient *gpm.Client
var enableLogging bool
var mouseX, mouseY = -1, -1
var menus []*menu.Menu
var fileBufferID string
var tabBarLine *tabbar.TabBar
var menuBar *menu.MenuBar

var modalPages *tview.Pages
var editorPages *tview.Pages
var statusBar *statusbar.StatusBar

var settings settingstype.Settings
var colorscheme smidgen.Colorscheme

type FileBuffer struct {
	panelVFlex    *tview.Flex
	panelHFlex    *tview.Flex
	editorVFlex   *tview.Flex
	scrollbar     *scrollbar.Scrollbar
	hScrollbar    *scrollbar.Scrollbar
	findbar       *findbar.Findbar
	isFindbarOpen bool
	openFindbar   func()

	buffer   *buffer.Buffer
	editor   *smidgen.View
	uuid     string
	filename string

	// Cached visual width of the longest line, used to size the horizontal
	// scrollbar. Recomputed when the buffer is modified or the tab size changes.
	maxLineWidth        int
	maxLineWidthValid   bool
	maxLineWidthTabSize int
}

var fileBuffers []*FileBuffer
var currentFileBuffer *FileBuffer

var findbarRecentSearchTextHistory []string
var findbarRecentReplaceTextHistory []string

// -----------------------------------------------------------------
func loadEditorColorScheme(colorSchemeName string) {
	var ok bool
	colorscheme, ok = smidgen.LoadInternalColorscheme(colorSchemeName)
	if !ok {
		colorscheme, _ = smidgen.LoadInternalColorscheme("monokai")
	}

	defaultStyle := colorscheme.GetColor("default")
	_, bg, _ := defaultStyle.Decompose()

	for _, fileBuffer := range fileBuffers {
		fileBuffer.editor.SetColorscheme(colorscheme)
		// Keep the thin horizontal scrollbar's empty half blending with the
		// editor background.
		fileBuffer.hScrollbar.Track.SetBackgroundColor(bg)
	}

	tabBarLine.SetTabBackgroundColor(bg)
}

func newFile(contents string, filename string) {
	buffer := smidgen.NewBufferFromString(contents, filename)
	editor := smidgen.NewView(app, buffer)
	buffer.Path = filename // femto uses this to determine the file type
	editor.SetColorscheme(colorscheme)
	editor.SetKeybindings(smidgenDefaultKeyBindings)
	editor.SetInputCapture(editorInputCapture)
	buffer.Settings["matchbrace"] = settings.ShowMatchBracket
	buffer.Settings["ruler"] = settings.ShowLineNumbers
	buffer.Settings["showwhitespace"] = settings.ShowWhitespace
	buffer.Settings["softwrap"] = settings.SoftWrap
	buffer.Settings["tabsize"] = float64(settings.TabSize)
	buffer.Settings["tabstospaces"] = settings.TabCharacter == "space"
	buffer.Settings["hltrailingws"] = settings.ShowTrailingWhitespace
	buffer.Settings["colorcolumn"] = settings.VerticalRuler

	// The editor and the horizontal scrollbar are stacked vertically so that
	// the horizontal scrollbar spans only the editor's width (not the column
	// occupied by the vertical scrollbar).
	editorVFlex := tview.NewFlex()
	editorVFlex.SetDirection(tview.FlexRow)
	editorVFlex.AddItem(editor, 0, 1, true)
	hScrollbar := scrollbar.NewScrollbar()
	hScrollbar.SetHorizontal(true)
	style.StyleScrollbar(hScrollbar)
	hScrollbar.SetThin(true)
	_, editorBg, _ := colorscheme.GetColor("default").Decompose()
	hScrollbar.Track.SetBackgroundColor(editorBg)
	editorVFlex.AddItem(hScrollbar, 1, 0, false)

	panelHFlex := tview.NewFlex()
	panelHFlex.SetDirection(tview.FlexColumn)
	panelHFlex.AddItem(editorVFlex, 0, 1, true)
	vScrollbar := scrollbar.NewScrollbar()
	style.StyleScrollbar(vScrollbar)
	panelHFlex.AddItem(vScrollbar, 1, 0, false)

	panelVFlex := tview.NewFlex()
	panelVFlex.SetDirection(tview.FlexRow)
	panelVFlex.AddItem(panelHFlex, 0, 1, true)

	bufferFindbar := findbar.NewFindbar(app, editor)
	style.StyleFindbar(bufferFindbar)
	bufferFindbar.SetSmidgenKeybindings(smidgenSingleLineKeyBindings)
	bufferFindbar.SetSearchTextHistory(findbarRecentSearchTextHistory)
	bufferFindbar.SetReplaceTextHistory(findbarRecentReplaceTextHistory)
	bufferFindbar.OnSearchTextHistoryChange = func(history []string) {
		syncFindbarSearchTextHistory(history, bufferFindbar)
	}
	bufferFindbar.OnReplaceTextHistoryChange = func(history []string) {
		syncFindbarReplaceTextHistory(history, bufferFindbar)
	}

	bufferFindbar.SetOnError(func(err error) {
		statusBar.ShowMessage(err.Error())
	})
	bufferFindbar.SetOnMessage(func(message string) {
		statusBar.ShowMessage(message)
	})
	bufferFindbar.SetOnWarning(func(message string) {
		statusBar.ShowWarning(message)
	})

	fileBuffer := &FileBuffer{
		panelVFlex:    panelVFlex,
		panelHFlex:    panelHFlex,
		editorVFlex:   editorVFlex,
		scrollbar:     vScrollbar,
		hScrollbar:    hScrollbar,
		buffer:        buffer,
		findbar:       bufferFindbar,
		isFindbarOpen: false,

		editor:   editor,
		uuid:     uuid.New().String(),
		filename: filename,
	}

	fileBuffer.openFindbar = func() {
		if !fileBuffer.isFindbarOpen {
			fileBuffer.panelVFlex.AddItem(fileBuffer.findbar, 1, 0, false)
			fileBuffer.isFindbarOpen = true
		}

		selectionText := editor.Cursor().GetSelection()
		if len(selectionText) != 0 {
			// Split the text into lines and use the first line only
			// (as the findbar is a single line input)
			if idx := bytes.IndexByte(selectionText, '\n'); idx > 0 {
				selectionText = selectionText[:idx]
			}
			fileBuffer.findbar.SetSearchText(string(selectionText))
		}
	}
	bufferFindbar.OnClose = func() {
		if fileBuffer.isFindbarOpen {
			fileBuffer.isFindbarOpen = false
			fileBuffer.panelVFlex.RemoveItem(fileBuffer.findbar)
			app.SetFocus(editor)
		}
	}
	bufferFindbar.SetOnExpand(func(expanded bool) {
		if fileBuffer.isFindbarOpen {
			newSize := 1
			if expanded {
				newSize += 1
			}
			fileBuffer.panelVFlex.ResizeItem(fileBuffer.findbar, newSize, 0)
		}
	})

	vScrollbar.UpdateHook = func(sb *scrollbar.Scrollbar) {
		// Update the scrollbar's position and size based on the content
		_, _, _, height := editor.GetRect()
		sb.Track.SetThumbSize(height)
		sb.Track.SetMax(buffer.LinesNum())
		sloc := editor.ActionController().GetView().StartLine
		sb.Track.SetPosition(sloc.Line)
	}
	vScrollbar.SetChangedFunc(func(position int) {
		editor.ActionController().SetStartLine(display.SLoc{Line: position, Row: 0})
	})

	hScrollbar.UpdateHook = func(sb *scrollbar.Scrollbar) {
		// Size the thumb to the visible text width and position it at the
		// current horizontal scroll offset.
		bufWidth := editor.ActionController().BufView().Width
		sb.Track.SetMax(max(fileBuffer.contentWidth(), 1))
		sb.Track.SetThumbSize(bufWidth)
		sb.Track.SetPosition(editor.ActionController().GetView().StartCol)
	}
	hScrollbar.SetChangedFunc(func(position int) {
		editor.ActionController().GetView().StartCol = position
	})

	// Decide which scrollbars are visible just before the panel lays out its
	// children, so they only appear when the content overflows the viewport.
	panelHFlex.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		fileBuffer.updateScrollbarVisibility()
		return x, y, width, height
	})

	fileBuffers = append(fileBuffers, fileBuffer)

	editorPages.AddPage(fileBuffer.uuid, panelVFlex, true, false)
	tabName := "[Untitled]"
	if filename != "" {
		tabName = path.Base(filename)
	}
	tabBarLine.AddTab(tabName, fileBuffer.uuid)
	tabBarLine.SetActive(fileBuffer.uuid)

	selectTab(fileBuffer.uuid)
	app.SetFocus(editor)
}

func syncFindbarSearchTextHistory(history []string, originalFindbar *findbar.Findbar) {
	findbarRecentSearchTextHistory = history
	for _, fileBuffer := range fileBuffers {
		if fileBuffer.findbar == originalFindbar {
			continue
		}
		fileBuffer.findbar.SetSearchTextHistory(history)
	}
}

func syncFindbarReplaceTextHistory(history []string, originalFindbar *findbar.Findbar) {
	findbarRecentReplaceTextHistory = history
	for _, fileBuffer := range fileBuffers {
		if fileBuffer.findbar == originalFindbar {
			continue
		}
		fileBuffer.findbar.SetReplaceTextHistory(history)
	}
}

func loadFile(filename string) string {
	contents, err := os.ReadFile(filename)
	if err != nil {
		newFile("", filename)
		return fmt.Sprintf("Failed to read file '%s':\n%v", filename, err)
	}
	newFile(string(contents), filename)
	return ""
}

func loadStdin() string {
	contents, err := io.ReadAll(os.Stdin)
	if err != nil {
		newFile("", "")
		return fmt.Sprintf("Failed to read from stdin:\n%v", err)
	}
	newFile(string(contents), "")
	return ""
}

// contentWidth returns the visual width of the longest line in the buffer.
// The result is cached and only recomputed when the buffer has been modified
// or the tab size changed, so it is cheap to call every frame.
func (fb *FileBuffer) contentWidth() int {
	tabSize := int(fb.buffer.Settings["tabsize"].(float64))
	if fb.maxLineWidthValid && !fb.buffer.ModifiedThisFrame && fb.maxLineWidthTabSize == tabSize {
		return fb.maxLineWidth
	}

	maxWidth := 0
	for i := 0; i < fb.buffer.LinesNum(); i++ {
		line := fb.buffer.LineBytes(i)
		w := util.StringWidth(line, util.CharacterCount(line), tabSize)
		if w > maxWidth {
			maxWidth = w
		}
	}

	fb.maxLineWidth = maxWidth
	fb.maxLineWidthValid = true
	fb.maxLineWidthTabSize = tabSize
	return maxWidth
}

// updateScrollbarVisibility shows each scrollbar only when its content
// overflows the viewport, and hides it otherwise.
func (fb *FileBuffer) updateScrollbarVisibility() {
	_, _, _, editorHeight := fb.editor.GetRect()

	// The vertical scrollbar is only needed when there are more lines than fit
	// on screen.
	vNeeded := fb.buffer.LinesNum() > editorHeight
	fb.panelHFlex.ResizeItem(fb.scrollbar, boolToInt(vNeeded), 0)

	// The horizontal scrollbar is only needed when soft wrap is off and the
	// longest line is wider than the visible text area (excluding the gutter).
	softWrap := fb.buffer.Settings["softwrap"].(bool)
	bufWidth := fb.editor.ActionController().BufView().Width
	hNeeded := !softWrap && fb.contentWidth() > bufWidth
	fb.editorVFlex.ResizeItem(fb.hScrollbar, boolToInt(hNeeded), 0)
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// handleHorizontalWheelScroll scrolls the editor horizontally in response to a
// horizontal mouse wheel event. It returns true if the event was handled.
func handleHorizontalWheelScroll(action tview.MouseAction, event *tcell.EventMouse) bool {
	fb := currentFileBuffer
	if fb == nil {
		return false
	}
	if fb.buffer.Settings["softwrap"].(bool) {
		return false
	}

	x, y := event.Position()
	if !fb.editor.InRect(x, y) {
		return false
	}

	view := fb.editor.ActionController().GetView()
	bufWidth := fb.editor.ActionController().BufView().Width

	scrollSpeed := util.IntOpt(fb.buffer.Settings["scrollspeed"])
	if scrollSpeed < 1 {
		scrollSpeed = 1
	}

	newStartCol := view.StartCol
	if action == tview.MouseScrollLeft {
		newStartCol -= scrollSpeed
	} else {
		newStartCol += scrollSpeed
	}

	if maxStartCol := fb.contentWidth() - bufWidth; newStartCol > maxStartCol {
		newStartCol = maxStartCol
	}
	if newStartCol < 0 {
		newStartCol = 0
	}
	view.StartCol = newStartCol
	return true
}

func getFileBufferByID(id string) *FileBuffer {
	for _, fileBuffer := range fileBuffers {
		if fileBuffer.uuid == id {
			return fileBuffer
		}
	}
	return nil
}

func showTabPage(id string) {
	fileBuffer := getFileBufferByID(id)
	fileBufferID = id
	editorPages.SwitchToPage(id)
	currentFileBuffer = fileBuffer
	syncMenuFromBuffer(currentFileBuffer.buffer)
}

func selectTab(id string) {
	showTabPage(id)
	tabBarLine.SetActive(id)
}

func syncStatusBarFromFileBuffer(statusBar *statusbar.StatusBar) {
	fileBuffer := getFileBufferByID(fileBufferID)
	if fileBuffer == nil {
		return
	}
	statusBar.Filename = fileBuffer.filename
	statusBar.Line = fileBuffer.editor.Cursor().Y + 1
	statusBar.Col = fileBuffer.editor.Cursor().X + 1

	statusBar.IsModified = fileBuffer.buffer.Modified()

	tabSize := int(fileBuffer.buffer.Settings["tabsize"].(float64))
	statusBar.TabSize = tabSize
	statusBar.IsOverwriteMode = fileBuffer.buffer.OverwriteMode

	lineEndings := "LF"
	if isBufferCRLF(fileBuffer.buffer) {
		lineEndings = "CRLF"
	}
	statusBar.LineEndings = lineEndings
}

func isBufferCRLF(buffer *buffer.Buffer) bool {
	return buffer.Settings["fileformat"].(string) == "dos"
}

func editorInputCapture(event *tcell.EventKey) *tcell.EventKey {
	for keyDesc, action := range dinkyKeyBindings {
		if event.Key() == keyDesc.KeyCode {
			if event.Key() == tcell.KeyRune && keyDesc.R != event.Rune() {
				continue
			}

			if keyDesc.Modifiers == event.Modifiers() {
				p := dinkyActionMapping[action]()
				if p != nil {
					app.SetFocus(p)
				}
				return nil
			}
		}
	}
	return event
}

func showHelp() {
	fmt.Printf("Dinky - A little text editor\n\n")
	fmt.Printf("Usage: dinky [options] [file1] [file2] ...\n\n")
	fmt.Printf("Options:\n")
	fmt.Printf("  -h, --help     Show this help message and exit\n")
	fmt.Printf("  -v, --version  Show version information and exit\n")
	fmt.Printf("  --log          Enable logging to app.log file\n\n")
	fmt.Printf("Arguments:\n")
	fmt.Printf("  file1, file2, ...  Files to open in the editor\n")
	fmt.Printf("  -                  Use dash to read from standard input\n\n")
	fmt.Printf("Use the option `--` to terminate option parsing.\n\n")
}

func showVersion() {
	fmt.Printf("Version: %s\n", getDinkyVersion())
	fmt.Printf("Version time: %s\n", getDinkyVersionTime())
}

func getDinkyVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok || info == nil {
		return "unknown"
	}
	var tag, commit string
	for _, s := range info.Settings {
		if s.Key == "vcs.revision" {
			commit = s.Value
		} else if s.Key == "vcs.tag" {
			tag = s.Value
		}
	}
	if tag == "" {
		tag = "untagged"
	}
	if commit == "" {
		commit = "unknown"
	}
	if len(commit) > 7 {
		commit = commit[:7]
	}
	return tag + " (" + commit + ")"
}

func getDinkyVersionTime() string {
	info, ok := debug.ReadBuildInfo()
	if !ok || info == nil {
		return "unknown"
	}
	var buildTime string
	for _, s := range info.Settings {
		if s.Key == "vcs.time" {
			buildTime = s.Value
		}
	}
	if buildTime == "" {
		buildTime = "unknown"
	}
	return buildTime
}

func parseCommandLine() bool {
	args := os.Args[1:]
	fileArgs := []string{}
	endOfOptions := false

	for _, arg := range args {
		if endOfOptions {
			fileArgs = append(fileArgs, arg)
			continue
		}
		switch arg {
		case "--":
			endOfOptions = true
		case "-h", "--help":
			showHelp()
			return false
		case "-v", "--version":
			showVersion()
			return false
		case "--log":
			enableLogging = true
		default:
			// If it starts with a dash, it's an unknown option
			if len(arg) > 1 && arg[0] == '-' {
				fmt.Fprintf(os.Stderr, "Error: Unknown option '%s'\n", arg)
				fmt.Fprintf(os.Stderr, "Use 'dinky --help' for usage information.\n")
				return false
			}
			if arg == "-" {
				fileArgs = append(fileArgs, "")
				// Use empty string to indicate stdin. Empty string isn't a valid
				// filename, so it won't conflict with any real file.
			} else {
				fileArgs = append(fileArgs, arg)
			}
		}
	}

	// Update os.Args to contain only the program name and file arguments
	os.Args = append([]string{os.Args[0]}, fileArgs...)

	return true
}

func Main() {
	// Parse command line arguments first
	if !parseCommandLine() {
		return
	}

	var logFile *os.File
	if enableLogging {
		logFile = setupLogging()
		defer logFile.Close()
		log.Println("Dinky starting with logging enabled")
	} else {
		// Disable logging by setting output to discard
		log.SetOutput(io.Discard)
	}

	settings = LoadUserSettings()

	initKeyBindings()

	app = tview.NewApplication()
	tview.DoubleClickInterval = 0 // Disable tview's double-click handling
	app.EnableMouse(true)
	app.EnablePaste(true)
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Disable Ctrl-C quitting the app
		if event.Key() == tcell.KeyCtrlC {
			return tcell.NewEventKey(tcell.KeyCtrlC, event.Rune(), event.Modifiers())
		}
		return event
	})

	modalPages = tview.NewPages()

	mainUiFlex := tview.NewFlex()
	mainUiFlex.SetDirection(tview.FlexRow)

	menuBar = menu.NewMenuBar()
	style.StyleMenuBar(menuBar)
	menus = createMenus()
	syncMenuKeyBindings(menus, actionToKeyMapping)
	menuBar.SetMenus(menus)

	mainUiFlex.AddItem(menuBar, 1, 0, false)

	tabBarLine = tabbar.NewTabBar()
	style.StyleTabBar(tabBarLine)
	tabBarLine.OnActive = func(id string, index int) {
		showTabPage(id)
	}
	tabBarLine.OnTabCloseClick = func(id string, index int) {
		fileBufferID = id
		handleCloseFile()
	}
	tabBarLine.QueueUpdateDraw = func(f func()) {
		app.QueueUpdateDraw(f)
	}
	tabBarLine.OnReorder = func(id string, newIndex int) {
		// Keep the fileBuffers order in sync with the tab order so that
		// next/previous tab navigation follows the on-screen layout.
		for i, fileBuffer := range fileBuffers {
			if fileBuffer.uuid == id {
				fileBuffers = slices.Delete(fileBuffers, i, i+1)
				fileBuffers = slices.Insert(fileBuffers, newIndex, fileBuffer)
				break
			}
		}
	}

	loadEditorColorScheme(settings.ColorScheme)
	mainUiFlex.AddItem(tabBarLine, 1, 0, false)

	editorPages = tview.NewPages()
	mainUiFlex.AddItem(editorPages, 0, 1, true)

	statusBar = statusbar.NewStatusBar(app)
	statusBar.UpdateHook = syncStatusBarFromFileBuffer
	mainUiFlex.AddItem(statusBar, 1, 0, false)

	modalPages.AddPage("workspace", mainUiFlex, true, true)

	app.SetRoot(modalPages, true)
	app.SetAfterDrawFunc(menuBar.AfterDraw())

	menuBar.SetOnClose(func(nextFocus tview.Primitive) {
		if nextFocus != nil {
			app.SetFocus(nextFocus)
		} else {
			app.SetFocus(currentFileBuffer.editor)
		}
	})

	errorMessages := []string{}
	readFromStdin := false
	for _, arg := range os.Args[1:] {
		var resultString string
		if arg == "" { // Empty string indicates stdin (see parseCommandLine).
			if !readFromStdin {
				fmt.Println("Reading from stdin... (press Ctrl-D to finish input)")
				resultString = loadStdin()
				readFromStdin = true
			}
		} else {
			resultString = loadFile(arg)
		}
		if resultString != "" {
			errorMessages = append(errorMessages, resultString)
		}
	}

	var showLoadingError func()
	showLoadingError = func() {
		CloseMessageDialog()
		if len(errorMessages) > 0 {
			errorMessage := errorMessages[0]
			errorMessages = errorMessages[1:]
			ShowOkDialog("Error loading file", errorMessage, showLoadingError)
		}
	}
	if len(fileBuffers) == 0 {
		newFile("", "")
	}
	selectTab(fileBuffers[0].uuid)

	// Shown after the tabs are set up so the dialog keeps focus.
	showLoadingError()

	// Track the mouse position (used by the GPM cursor renderer) and turn
	// horizontal wheel events into horizontal scrolling of the editor. The
	// capture is always called from the main loop, so mouseX/mouseY are safe to
	// read in afterDraw.
	app.SetMouseCapture(func(event *tcell.EventMouse, action tview.MouseAction) (*tcell.EventMouse, tview.MouseAction) {
		if event == nil {
			return event, action
		}
		mouseX, mouseY = event.Position()

		switch {
		case action == tview.MouseScrollLeft || action == tview.MouseScrollRight:
			// Dedicated horizontal wheel events: scroll as a side effect and let
			// the event continue to the editor, which ignores them. (Returning a
			// nil event here would leave the shared event nil for any follow-up
			// action fired from the same mouse event.)
			handleHorizontalWheelScroll(action, event)
		case (action == tview.MouseScrollUp || action == tview.MouseScrollDown) && event.Modifiers()&tcell.ModShift != 0:
			// Many terminals report horizontal scrolling as Shift + vertical
			// wheel. Translate it and swallow the event so the editor does not
			// also scroll vertically.
			horiz := tview.MouseScrollLeft
			if action == tview.MouseScrollDown {
				horiz = tview.MouseScrollRight
			}
			if handleHorizontalWheelScroll(horiz, event) {
				return nil, action
			}
		}
		return event, action
	})

	// Connect to GPM before Run() (handshake happens now; events injected via QueueEvent).
	if client, err := gpm.Connect(func(ev tcell.Event) {
		app.QueueEvent(ev)
		// Pure-move events (ButtonNone) are not consumed by any tview component, so
		// they don't trigger a redraw on their own. Force one so the cursor updates.
		if me, ok := ev.(*tcell.EventMouse); ok && me.Buttons() == tcell.ButtonNone {
			app.QueueUpdateDraw(func() {})
		}
	}); err != nil {
		log.Printf("GPM connect failed (non-fatal): %v", err)
	} else if client != nil {
		gpmClient = client
		gpmClient.Start()

		// Chain afterDraw to render a cursor by inverting the cell under the pointer.
		existingAfterDraw := app.GetAfterDrawFunc()
		app.SetAfterDrawFunc(func(screen tcell.Screen) {
			if existingAfterDraw != nil {
				existingAfterDraw(screen)
			}
			if mouseX < 0 || mouseY < 0 {
				return
			}
			w, h := screen.Size()
			if mouseX >= w || mouseY >= h {
				return
			}
			mainc, combc, _, _ := screen.GetContent(mouseX, mouseY)
			screen.SetContent(mouseX, mouseY, mainc, combc, tcell.StyleDefault.Background(tcell.ColorYellow).Foreground(tcell.ColorBlack))
		})
	}

	if err := app.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}

	if gpmClient != nil {
		gpmClient.Stop()
	}
}

func setupLogging() *os.File {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic("Failed to open log file: " + err.Error())
	}
	log.SetOutput(logFile)
	return logFile
}
