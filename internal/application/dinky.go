package application

import (
	"dinky/internal/tui/menu"
	"dinky/internal/tui/statusbar"
	"dinky/internal/tui/tabbar"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/pgavlin/femto"
	"github.com/pgavlin/femto/runtime"
	"github.com/rivo/tview"
)

// -----------------------------------------------------------------
var app *tview.Application
var menus []*menu.Menu
var buffer *femto.Buffer
var tabBarLine *tabbar.TabBar
var editor *femto.View
var pages *tview.Pages
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

	fileBuffer := &FileBuffer{
		buffer:   buffer,
		editor:   editor,
		uuid:     uuid.New().String(),
		filename: filename,
	}

	fileBuffers = append(fileBuffers, fileBuffer)
	pages.AddPage(fileBuffer.uuid, editor, true, true)
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
	pages.SwitchToPage(id)
	buffer = fileBuffer.buffer
	editor = fileBuffer.editor
	syncMenuFromBuffer(buffer)
	syncStatusBarFromFileBuffer(fileBuffer)
}

func syncStatusBarFromFileBuffer(fileBuffer *FileBuffer) {
	statusBar.Filename = fileBuffer.filename
	// statusbar.SetLineCount(fileBuffer.buffer.LineCount())
	// statusbar.SetColumnCount(fileBuffer.buffer.ColumnCount())
	// statusbar.SetCursorPosition(fileBuffer.buffer.CursorX(), fileBuffer.buffer.CursorY())
	// statusbar.SetReadOnly(fileBuffer.buffer.ReadOnly())
	// statusbar.SetModified(fileBuffer.buffer.Modified())
}

func Main() {
	logFile := setupLogging()
	defer logFile.Close()

	initEditorColorScheme()
	initKeyBindings()

	app = tview.NewApplication()
	app.EnableMouse(true)

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)

	menuBar := menu.NewMenuBar()

	menus = createMenus()
	syncMenuKeyBindings(menus, femtoDefaultKeyBindings)
	menuBar.SetMenus(menus)

	flex.AddItem(menuBar, 1, 0, false)

	// buffer = femto.NewBufferFromString("Hello Smoe\nSome words to click on\n", "/home/sbe/smoe.txt")
	// editor = femto.NewView(buffer)

	// editor.SetRuntimeFiles(runtime.Files)
	// editor.SetColorscheme(colorscheme)
	// editor.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
	// 	switch event.Key() {
	// 	case tcell.KeyCtrlS:
	// 		return nil
	// 	case tcell.KeyCtrlQ:
	// 		app.Stop()
	// 		return nil
	// 	}
	// 	return event
	// })

	tabBarLine = tabbar.NewTabBar()
	tabBarLine.OnActive = func(id string, index int) {
		selectTab(id)
	}

	defaultStyle := colorscheme.GetColor("default")
	_, bg, _ := defaultStyle.Decompose()
	tabBarLine.SetTabBackgroundColor(bg)
	// tabBarLine.SetTabInactiveBackgroundColor(bg)

	flex.AddItem(tabBarLine, 1, 0, false)

	pages = tview.NewPages()
	flex.AddItem(pages, 0, 1, true)

	statusBar = statusbar.NewStatusBar()
	flex.AddItem(statusBar, 1, 0, false)

	app.SetRoot(flex, true)
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
