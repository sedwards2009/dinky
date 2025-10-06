package application

import (
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
	"strings"

	"runtime/debug"

	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"
	"github.com/sedwards2009/femto"
	"github.com/sedwards2009/femto/runtime"
)

// -----------------------------------------------------------------
var app *tview.Application
var enableLogging bool
var menus []*menu.Menu
var fileBufferID string
var tabBarLine *tabbar.TabBar
var menuBar *menu.MenuBar

var modalPages *tview.Pages
var editorPages *tview.Pages
var statusBar *statusbar.StatusBar

var colorscheme femto.Colorscheme

type FileBuffer struct {
	panelVFlex    *tview.Flex
	panelHFlex    *tview.Flex
	scrollbar     *scrollbar.Scrollbar
	findbar       *findbar.Findbar
	isFindbarOpen bool
	openFindbar   func()

	buffer   *femto.Buffer
	editor   *femto.View
	uuid     string
	filename string
}

var fileBuffers []*FileBuffer
var currentFileBuffer *FileBuffer

// -----------------------------------------------------------------
func initEditorColorScheme() {
	if monokai := runtime.Files.FindFile(femto.RTColorscheme, "monokai"); monokai != nil {
		if data, err := monokai.Data(); err == nil {
			colorscheme = femto.ParseColorscheme(string(data))
		}
	}
}

func newFile(contents string, filename string) {
	buffer := femto.NewBufferFromString(contents, "")
	editor := femto.NewView(buffer)
	buffer.Path = filename // femto uses this to determine the file type
	editor.SetRuntimeFiles(runtime.Files)
	editor.SetColorscheme(colorscheme)
	editor.SetKeybindings(femtoDefaultKeyBindings)
	editor.SetInputCapture(editorInputCapture)
	buffer.Settings["matchbrace"] = true
	buffer.Settings["tabstospaces"] = false // Default to using tab character

	panelHFlex := tview.NewFlex()
	panelHFlex.SetDirection(tview.FlexColumn)
	panelHFlex.AddItem(editor, 0, 1, true)
	vScrollbar := scrollbar.NewScrollbar()
	style.StyleScrollbar(vScrollbar)
	panelHFlex.AddItem(vScrollbar, 1, 0, false)

	panelVFlex := tview.NewFlex()
	panelVFlex.SetDirection(tview.FlexRow)
	panelVFlex.AddItem(panelHFlex, 0, 1, true)

	bufferFindbar := findbar.NewFindbar(app, editor)
	style.StyleFindbar(bufferFindbar)
	bufferFindbar.SetFemtoKeybindings(femtoSingleLineKeyBindings)

	fileBuffer := &FileBuffer{
		panelVFlex:    panelVFlex,
		panelHFlex:    panelHFlex,
		scrollbar:     vScrollbar,
		buffer:        buffer,
		findbar:       bufferFindbar,
		isFindbarOpen: false,

		editor:   editor,
		uuid:     uuid.New().String(),
		filename: filename,
	}

	fileBuffer.openFindbar = func() {
		if !fileBuffer.isFindbarOpen {
			fileBuffer.panelVFlex.AddItem(fileBuffer.findbar, 2, 0, false)
			fileBuffer.isFindbarOpen = true
		}

		selectionText := editor.Cursor.GetSelection()
		if selectionText != "" {
			// Split the text into lines and use the first line only
			// (as the findbar is a single line input)
			if idx := strings.IndexByte(selectionText, '\n'); idx > 0 {
				selectionText = selectionText[:idx]
			}
			fileBuffer.findbar.SetSearchText(selectionText)
		}
	}
	bufferFindbar.OnClose = func() {
		if fileBuffer.isFindbarOpen {
			fileBuffer.isFindbarOpen = false
			fileBuffer.panelVFlex.RemoveItem(fileBuffer.findbar)
			app.SetFocus(editor)
		}
	}

	vScrollbar.UpdateHook = func(sb *scrollbar.Scrollbar) {
		// Update the scrollbar's position and size based on the content
		_, _, _, height := editor.GetRect()
		sb.Track.SetThumbSize(height)
		sb.Track.SetMax(buffer.NumLines)
		sb.Track.SetPosition(editor.Topline)
	}
	vScrollbar.SetChangedFunc(func(position int) {
		editor.Topline = position
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

func loadFile(filename string) string {
	// Read the file contents
	contents, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Sprintf("Failed to read file '%s':\n%v", filename, err)
	}
	newFile(string(contents), filename)
	return ""
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
	statusBar.Line = fileBuffer.editor.Cursor.Y + 1
	statusBar.Col = fileBuffer.editor.Cursor.X + 1

	statusBar.IsModified = fileBuffer.buffer.IsModified

	tabSize := int(fileBuffer.buffer.Settings["tabsize"].(float64))
	statusBar.TabSize = tabSize
	statusBar.IsOverwriteMode = fileBuffer.editor.IsOverwriteMode()

	lineEndings := "LF"
	if isBufferCRLF(fileBuffer.buffer) {
		lineEndings = "CRLF"
	}
	statusBar.LineEndings = lineEndings
}

func isBufferCRLF(buffer *femto.Buffer) bool {
	return buffer.Settings["fileformat"].(string) == "dos"
}

func editorInputCapture(event *tcell.EventKey) *tcell.EventKey {
	for keyDesc, action := range dinkyKeyBindings {
		if event.Key() == keyDesc.KeyCode {
			if event.Key() == tcell.KeyRune {
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
	fmt.Printf("  file1, file2, ...  Files to open in the editor\n\n")
	fmt.Printf("If no files are specified, a new empty file will be created.\n")
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

	for _, arg := range args {
		switch arg {
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
			if len(arg) > 0 && arg[0] == '-' {
				fmt.Fprintf(os.Stderr, "Error: Unknown option '%s'\n", arg)
				fmt.Fprintf(os.Stderr, "Use 'dinky --help' for usage information.\n")
				return false
			}
			// Otherwise, it's a file to open
			fileArgs = append(fileArgs, arg)
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

	initEditorColorScheme()
	initKeyBindings()

	app = tview.NewApplication()
	tview.DoubleClickInterval = 0 // Disable tview's double-click handling
	app.EnableMouse(true)

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

	defaultStyle := colorscheme.GetColor("default")
	_, bg, _ := defaultStyle.Decompose()
	tabBarLine.SetTabBackgroundColor(bg)

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
	for _, arg := range os.Args[1:] {
		resultString := loadFile(arg)
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
	showLoadingError()

	if len(fileBuffers) == 0 {
		newFile("", "")
	}
	selectTab(fileBuffers[0].uuid)

	if err := app.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
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
