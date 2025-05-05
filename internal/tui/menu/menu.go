package menu

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/tview"
)

type MenuBar struct {
	*tview.Box
	MenuBarStyle tcell.Style
	menus        []*Menu
	selectedPath []int
}

type Menu struct {
	Title string
	Items []*MenuItem
}

type MenuItem struct {
	Title    string
	Shortcut string
	Callback func()
}

const (
	MENU_BAR_SPACING = 2
	MENU_BAR_PADDING = 1
)

func NewMenuBar() *MenuBar {
	fg := tcell.NewHexColor(0xf3f3f3)
	bg := tcell.NewHexColor(0x007ace)
	return &MenuBar{
		Box: tview.NewBox(),

		MenuBarStyle: tcell.StyleDefault.Foreground(fg).Background(bg).Bold(true),
		selectedPath: []int{-1},
	}
}

func (menuBar *MenuBar) SetMenus(menus []*Menu) {
	menuBar.menus = menus
}

func (menuBar *MenuBar) Draw(screen tcell.Screen) {
	menuBar.Box.DrawForSubclass(screen, menuBar)
	x, y, width, _ := menuBar.GetInnerRect()

	for i := 0; i < width; i += 1 {
		screen.SetContent(x+i, y, ' ', nil, menuBar.MenuBarStyle)
	}

	padding := ""
	for range MENU_BAR_PADDING {
		padding += " "
	}

	reverse := menuBar.MenuBarStyle.Reverse(true)
	dx := MENU_BAR_SPACING
	for i, menu := range menuBar.menus {
		title := menu.Title
		style := menuBar.MenuBarStyle
		if i == menuBar.selectedPath[0] {
			style = reverse
		}
		drawText(screen, dx, y, padding, style)
		dx += MENU_BAR_PADDING
		drawText(screen, dx, y, title, style)
		dx += len(title)
		drawText(screen, dx, y, padding, style)
		dx += MENU_BAR_PADDING
		dx += MENU_BAR_SPACING
	}
}

func drawText(screen tcell.Screen, x int, y int, text string, style tcell.Style) {
	for _, r := range text {
		screen.SetContent(x, y, r, nil, style)
		x++
	}
}

func drawHorizontalLine(screen tcell.Screen, x int, y int, width int, borderStyle tcell.Style, middleStyle tcell.Style, left rune,
	middle rune, right rune) {

	screen.SetContent(x, y, left, nil, borderStyle)
	for i := 1; i < width-1; i++ {
		screen.SetContent(x+i, y, middle, nil, middleStyle)
	}
	screen.SetContent(x+width-1, y, right, nil, borderStyle)
}

func (menuBar *MenuBar) AfterDraw() func(tcell.Screen) {
	return func(screen tcell.Screen) {
		if menuBar.selectedPath[0] == -1 {
			return
		}
		selectedIndex := menuBar.selectedPath[0]
		menu := menuBar.menus[selectedIndex]

		mx := menuBar.menuIndexLeft(selectedIndex)

		rx, ry, _, _ := menuBar.GetRect()
		menuBar.drawMenuItems(screen, rx+mx, ry, menu.Items, menuBar.selectedPath[1])
	}
}

func (menuBar *MenuBar) drawMenuItems(screen tcell.Screen, menuX int, menuY int, items []*MenuItem, selectedIndex int) {
	titleWidth, _ := measureWidths(items)
	y := menuY + 1
	menuWidth := menuWidthInCells(items)

	borderStyle := menuBar.MenuBarStyle

	drawHorizontalLine(screen, menuX, y, menuWidth, borderStyle, borderStyle, '┌', '─', '┐')
	y++

	for i, item := range items {
		if item.Title == "" {
			drawHorizontalLine(screen, menuX, y, menuWidth, borderStyle, borderStyle, '├', '─', '┤')
		} else {
			textStyle := borderStyle
			if i == selectedIndex {
				textStyle = textStyle.Reverse(true)
			}
			drawHorizontalLine(screen, menuX, y, menuWidth, borderStyle, textStyle, '│', ' ', '│')
			drawText(screen, menuX+2, y, item.Title, textStyle)
			drawText(screen, menuX+2+titleWidth+2, y, item.Shortcut, textStyle)
		}
		y++
	}

	drawHorizontalLine(screen, menuX, y, menuWidth, borderStyle, borderStyle, '└', '─', '┘')

	// Draw the drop shadow
	// drawDimVerticalLine(d.X+menuWidth, d.Y+1, len(*d.MenuDefinition)+1)
	// drawDimVerticalLine(d.X+menuWidth+1, d.Y+1, len(*d.MenuDefinition)+1)
	// drawDimHorizontalLine(d.X+2, y+1, menuWidth)
}

func menuWidthInCells(items []*MenuItem) int {
	titleWidth, shortcutWidth := measureWidths(items)
	return 1 + 1 + titleWidth + 2 + shortcutWidth + 1 + 1
}

func measureWidths(items []*MenuItem) (int, int) {
	maxTitleWidth := 0
	maxShortcutWidth := 0
	for _, item := range items {
		width := runewidth.StringWidth(item.Title)
		if width > maxTitleWidth {
			maxTitleWidth = width
		}
		shortCutWidth := runewidth.StringWidth(item.Shortcut)
		if shortCutWidth > maxShortcutWidth {
			maxShortcutWidth = shortCutWidth
		}
	}
	return maxTitleWidth, maxShortcutWidth
}

func (menuBar *MenuBar) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return menuBar.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Printf("InputHandler %v", event)
		if menuBar.selectedPath[0] == -1 {
			return
		}

		selectedMenuIndex := menuBar.selectedPath[0]
		menu := menuBar.menus[selectedMenuIndex]
		selectedItemIndex := menuBar.selectedPath[1]
		item := menu.Items[selectedItemIndex]

		switch event.Key() {
		case tcell.KeyEscape:
			menuBar.selectMenuBarItem(-1)

		case tcell.KeyLeft:
			selectedMenuIndex--
			if selectedMenuIndex < 0 {
				selectedMenuIndex = 0
			}
			menuBar.selectMenuBarItem(selectedMenuIndex)

		case tcell.KeyRight:
			selectedMenuIndex++
			if selectedMenuIndex >= len(menuBar.menus) {
				selectedMenuIndex = len(menuBar.menus) - 1
			}
			menuBar.selectMenuBarItem(selectedMenuIndex)

		case tcell.KeyUp:
			menuBar.selectedPath[1] = nextMenuItem(menu.Items, selectedItemIndex, -1)

		case tcell.KeyDown:
			menuBar.selectedPath[1] = nextMenuItem(menu.Items, selectedItemIndex, 1)

		case tcell.KeyEnter:
			if item.Title != "" {
				if item.Callback != nil {
					item.Callback()
				}
				menuBar.selectMenuBarItem(-1)
			}
		}
	})
}

func nextMenuItem(items []*MenuItem, selectedIndex int, direction int) int {
	next := func() {
		selectedIndex += direction

		if selectedIndex < 0 {
			selectedIndex = len(items) - 1
		}
		if selectedIndex >= len(items) {
			selectedIndex = 0
		}
	}
	next()

	for items[selectedIndex].Title == "" {
		next()
	}
	return selectedIndex
}

func (menuBar *MenuBar) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return menuBar.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		if !menuBar.InRect(event.Position()) {
			return false, nil
		}

		rx, ry, _, _ := menuBar.GetRect()
		x, y := event.Position()

		if action == tview.MouseLeftDown {
			if y == ry {
				// Clicked on the menu bar itself
				index, _ := menuBar.menuItemIndexAtX(x - rx)
				if index != -1 {
					menuBar.selectMenuBarItem(index)
					return true, nil
				}
			}
		}

		if action == tview.MouseLeftUp || action == tview.MouseLeftDown {
			if menuBar.selectedPath[0] != -1 {
				selectedIndex := menuBar.selectedPath[0]
				// left := menuBar.menuIndexLeft(selectedIndex)
				items := menuBar.menus[selectedIndex].Items
				index := y - ry - 2
				if index >= 0 && index < len(items) {

					if action == tview.MouseLeftUp {
						if items[index].Title != "" {
							if items[index].Callback != nil {
								items[index].Callback()
							}
							menuBar.selectMenuBarItem(-1)
						}
					}

					return true, nil

				} else {
					return false, nil
				}
			}
		}

		if action == tview.MouseLeftDown {
			menuBar.selectMenuBarItem(-1)
		}

		return true, nil
	})
}

func (menuBar *MenuBar) selectMenuBarItem(index int) {
	if index == -1 {
		menuBar.selectedPath = []int{index}
	} else {
		menuBar.selectedPath = []int{index, 0}
	}
}

func (m *MenuBar) menuItemIndexAtX(posX int) (index int, leftX int) {
	x := MENU_BAR_SPACING
	for i, menu := range m.menus {
		if posX < x {
			return -1, -1
		}

		left := x
		x += MENU_BAR_PADDING
		x += runewidth.StringWidth(menu.Title)
		x += MENU_BAR_PADDING
		if posX < x {
			return i, left
		}

		x += MENU_BAR_SPACING
	}
	return -1, -1
}

func (m *MenuBar) menuIndexLeft(index int) int {
	x := MENU_BAR_SPACING
	for i, menu := range m.menus {
		if i == index {
			return x
		}

		x += MENU_BAR_PADDING
		x += runewidth.StringWidth(menu.Title)
		x += MENU_BAR_PADDING
		x += MENU_BAR_SPACING
	}
	return -1
}
