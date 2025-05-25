package application

import (
	"dinky/internal/tui/menu"
	"dinky/internal/tui/statusbar"
	"dinky/internal/tui/tabbar"
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/pgavlin/femto"
	"github.com/pgavlin/femto/runtime"
	"github.com/rivo/tview"
)

// -----------------------------------------------------------------
var app *tview.Application
var menus []*menu.Menu
var fileBufferID string
var buffer *femto.Buffer
var tabBarLine *tabbar.TabBar
var menuBar *menu.MenuBar
var editor *femto.View
var modalPages *tview.Pages
var editorPages *tview.Pages
var statusBar *statusbar.StatusBar

var colorscheme femto.Colorscheme

type FileBuffer struct {
	buffer   *femto.Buffer
	editor   *femto.View
	uuid     string
	filename string
}

var fileBuffers []*FileBuffer

// -----------------------------------------------------------------
func initEditorColorScheme() {
	if monokai := runtime.Files.FindFile(femto.RTColorscheme, "monokai"); monokai != nil {
		if data, err := monokai.Data(); err == nil {
			colorscheme = femto.ParseColorscheme(string(data))
		}
	}
}

func newFile(contents string, filename string) {
	buffer = femto.NewBufferFromString(contents, "")
	editor = femto.NewView(buffer)
	editor.SetRuntimeFiles(runtime.Files)
	editor.SetColorscheme(colorscheme)
	editor.SetKeybindings(femtoDefaultKeyBindings)
	editor.SetInputCapture(editorInputCapture)

	fileBuffer := &FileBuffer{
		buffer:   buffer,
		editor:   editor,
		uuid:     uuid.New().String(),
		filename: filename,
	}

	fileBuffers = append(fileBuffers, fileBuffer)
	editorPages.AddPage(fileBuffer.uuid, editor, true, false)
	tabBarLine.AddTab(fileBuffer.filename, fileBuffer.uuid)
	if buffer == nil {
		buffer = fileBuffer.buffer
		editor = fileBuffer.editor
	}
}

func loadFile(filename string) {
	// Read the file contents
	contents, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	newFile(string(contents), filename)
}

func getFileBufferByID(id string) *FileBuffer {
	for _, fileBuffer := range fileBuffers {
		if fileBuffer.uuid == id {
			return fileBuffer
		}
	}
	return nil
}

func selectTab(id string) {
	fileBuffer := getFileBufferByID(id)
	fileBufferID = id
	editorPages.SwitchToPage(id)
	buffer = fileBuffer.buffer
	editor = fileBuffer.editor
	syncMenuFromBuffer(buffer)
	syncStatusBarFromFileBuffer(fileBuffer)
}

func syncStatusBarFromFileBuffer(fileBuffer *FileBuffer) {
	statusBar.Filename = fileBuffer.filename
	statusBar.Line = fileBuffer.editor.Cursor.Y + 1
	statusBar.Col = fileBuffer.editor.Cursor.X + 1

	tabSize := int(fileBuffer.buffer.Settings["tabsize"].(float64))
	statusBar.TabSize = tabSize

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
				dinkyActionMapping[action]()
				return nil
			}
		}
	}
	return event
}

func updateStatusBar(screen tcell.Screen) bool {
	fileBuffer := getFileBufferByID(fileBufferID)
	syncStatusBarFromFileBuffer(fileBuffer)
	return false
}

func Main() {
	logFile := setupLogging()
	defer logFile.Close()

	initEditorColorScheme()
	initKeyBindings()

	app = tview.NewApplication()
	app.EnableMouse(true)
	app.SetBeforeDrawFunc(updateStatusBar)
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC {
			return tcell.NewEventKey(tcell.KeyCtrlC, 0, tcell.ModCtrl)
		}
		return event
	})

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)

	menuBar = menu.NewMenuBar()
	menus = createMenus()
	syncMenuKeyBindings(menus, actionToKeyMapping)
	menuBar.SetMenus(menus)

	flex.AddItem(menuBar, 1, 0, false)

	// buffer = femto.NewBufferFromString("Hello Smoe\nSome words to click on\n", "/home/sbe/smoe.txt")
	// editor = femto.NewView(buffer)

	// editor.SetRuntimeFiles(runtime.Files)
	// editor.SetColorscheme(colorscheme)
	// editor.SetInputCapture(editorInputCapture)

	tabBarLine = tabbar.NewTabBar()
	tabBarLine.OnActive = func(id string, index int) {
		selectTab(id)
	}

	defaultStyle := colorscheme.GetColor("default")
	_, bg, _ := defaultStyle.Decompose()
	tabBarLine.SetTabBackgroundColor(bg)
	// tabBarLine.SetTabInactiveBackgroundColor(bg)

	flex.AddItem(tabBarLine, 1, 0, false)

	editorPages = tview.NewPages()
	flex.AddItem(editorPages, 0, 1, true)

	statusBar = statusbar.NewStatusBar(app)
	flex.AddItem(statusBar, 1, 0, false)

	modalPages = tview.NewPages()
	modalPages.AddPage("workspace", flex, true, true)

	app.SetRoot(modalPages, true)
	app.SetAfterDrawFunc(menuBar.AfterDraw())

	menuBar.SetOnClose(func() {
		app.SetFocus(editor)
	})

	for _, arg := range os.Args[1:] {
		loadFile(arg)
	}
	if len(fileBuffers) == 0 {
		newFile("Hello Dinky\nSome words to click on\n", "")
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
