package application

import (
	"dinky/internal/tui/menu"
	"dinky/internal/tui/statusbar"
	"dinky/internal/tui/style"
	"dinky/internal/tui/tabbar"
	"log"
	"os"
	"path"

	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/pgavlin/femto"
	"github.com/pgavlin/femto/runtime"
	"github.com/sedwards2009/nuview"
)

// -----------------------------------------------------------------
var app *nuview.Application
var menus []*menu.Menu
var fileBufferID string
var buffer *femto.Buffer
var tabBarLine *tabbar.TabBar
var menuBar *menu.MenuBar
var editor *femto.View
var modalPages *nuview.Panels
var editorPages *nuview.Panels
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
	editorPages.AddPanel(fileBuffer.uuid, editor, true, false)
	tabName := "[Untitled]"
	if filename != "" {
		tabName = path.Base(filename)
	}
	tabBarLine.AddTab(tabName, fileBuffer.uuid)
	tabBarLine.SetActive(fileBuffer.uuid)
	if buffer == nil {
		buffer = fileBuffer.buffer
		editor = fileBuffer.editor
	}

	selectTab(fileBuffer.uuid)
	app.SetFocus(editor)
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
	editorPages.SetCurrentPanel(id)
	buffer = fileBuffer.buffer
	editor = fileBuffer.editor
	syncMenuFromBuffer(buffer)
}

func syncStatusBarFromFileBuffer(statusBar *statusbar.StatusBar) {
	fileBuffer := getFileBufferByID(fileBufferID)
	statusBar.Filename = fileBuffer.filename
	statusBar.Line = fileBuffer.editor.Cursor.Y + 1
	statusBar.Col = fileBuffer.editor.Cursor.X + 1

	statusBar.IsModified = fileBuffer.buffer.IsModified

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

func Main() {
	logFile := setupLogging()
	defer logFile.Close()

	style.Install()

	initEditorColorScheme()
	initKeyBindings()

	app = nuview.NewApplication()
	app.EnableMouse(true)
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC {
			return tcell.NewEventKey(tcell.KeyCtrlC, 0, tcell.ModCtrl)
		}
		return event
	})

	modalPages = nuview.NewPanels()

	mainUiFlex := nuview.NewFlex()
	mainUiFlex.SetDirection(nuview.FlexRow)

	menuBar = menu.NewMenuBar()
	style.StyleMenuBar(menuBar)
	menus = createMenus()
	syncMenuKeyBindings(menus, actionToKeyMapping)
	menuBar.SetMenus(menus)

	mainUiFlex.AddItem(menuBar, 1, 0, false)

	tabBarLine = tabbar.NewTabBar()
	style.StyleTabBar(tabBarLine)
	tabBarLine.OnActive = func(id string, index int) {
		selectTab(id)
	}
	tabBarLine.OnTabCloseClick = func(id string, index int) {
		fileBufferID = id
		handleCloseFile()
	}

	defaultStyle := colorscheme.GetColor("default")
	_, bg, _ := defaultStyle.Decompose()
	tabBarLine.SetTabBackgroundColor(bg)
	// tabBarLine.SetTabInactiveBackgroundColor(bg)

	mainUiFlex.AddItem(tabBarLine, 1, 0, false)

	editorPages = nuview.NewPanels()
	mainUiFlex.AddItem(editorPages, 0, 1, true)

	statusBar = statusbar.NewStatusBar(app)
	statusBar.UpdateHook = syncStatusBarFromFileBuffer
	mainUiFlex.AddItem(statusBar, 1, 0, false)

	modalPages.AddPanel("workspace", mainUiFlex, true, true)

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
