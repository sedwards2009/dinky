package menu

import (
	"dinky/internal/tui/utils"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/tview"
)

type MenuBar struct {
	*tview.Box
	MenuBarStyle tcell.Style
	menus        []*Menu
	selectedPath []int
	onClose      func()
}

type Menu struct {
	ID    string
	Title string
	Items []*MenuItem
}

type MenuItem struct {
	ID       string
	Title    string
	Shortcut string
	Callback func(ID string)
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

func (menuBar *MenuBar) Open() {
	menuBar.selectMenuBarItem(0)
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
	dx := 0
	for i, menu := range menuBar.menus {
		title := menu.Title
		style := menuBar.MenuBarStyle
		if i == menuBar.selectedPath[0] {
			style = reverse
		}
		utils.DrawText(screen, dx, y, padding, style)
		dx += MENU_BAR_PADDING
		utils.DrawText(screen, dx, y, title, style)
		dx += len(title)
		utils.DrawText(screen, dx, y, padding, style)
		dx += MENU_BAR_PADDING
		dx += MENU_BAR_SPACING
	}
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
	topY := menuY + 1
	y := topY
	menuWidth := menuWidthInCells(items)

	borderStyle := menuBar.MenuBarStyle

	utils.DrawCappedHorizontalLine(screen, menuX, y, menuWidth, borderStyle, borderStyle, '┌', '─', '┐')
	y++

	for i, item := range items {
		if item.Title == "" {
			utils.DrawCappedHorizontalLine(screen, menuX, y, menuWidth, borderStyle, borderStyle, '├', '─', '┤')
		} else {
			textStyle := borderStyle
			if i == selectedIndex {
				textStyle = textStyle.Reverse(true)
			}
			utils.DrawCappedHorizontalLine(screen, menuX, y, menuWidth, borderStyle, textStyle, '│', ' ', '│')
			utils.DrawText(screen, menuX+2, y, item.Title, textStyle)
			utils.DrawText(screen, menuX+2+titleWidth+2, y, item.Shortcut, textStyle)
		}
		y++
	}

	utils.DrawCappedHorizontalLine(screen, menuX, y, menuWidth, borderStyle, borderStyle, '└', '─', '┘')

	// Draw the drop shadow
	utils.DrawDimVerticalLine(screen, menuX+menuWidth, topY+1, len(items)+1)
	utils.DrawDimVerticalLine(screen, menuX+menuWidth+1, topY+1, len(items)+1)
	utils.DrawDimHorizontalLine(screen, menuX+2, y+1, menuWidth)
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
		if menuBar.selectedPath[0] == -1 {
			return
		}

		selectedMenuIndex := menuBar.selectedPath[0]
		menu := menuBar.menus[selectedMenuIndex]
		selectedItemIndex := menuBar.selectedPath[1]
		item := menu.Items[selectedItemIndex]

		switch event.Key() {
		case tcell.KeyEscape:
			menuBar.Close()

		case tcell.KeyLeft:
			selectedMenuIndex--
			if selectedMenuIndex < 0 {
				selectedMenuIndex = 0
			}
			menuBar.selectMenuBarItem(selectedMenuIndex)

		case tcell.KeyRight:
			selectedMenuIndex++
			if selectedMenuIndex == len(menuBar.menus) {
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
					item.Callback(item.ID)
				}
				menuBar.Close()
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
		rx, ry, _, _ := menuBar.GetRect()
		x, y := event.Position()

		if y == ry {
			if action == tview.MouseLeftDown {
				// Clicked on the menu bar itself
				index, _ := menuBar.menuItemIndexAtX(x - rx)
				if index != -1 {
					selectedIndex := menuBar.selectedPath[0]
					if selectedIndex != index {
						menuBar.selectMenuBarItem(index)
						setFocus(menuBar)
						return true, nil
					} else {
						menuBar.Close()
						return true, nil
					}
				}
			}
			return true, nil
		}

		if menuBar.selectedPath[0] == -1 {
			return false, nil
		}

		if action == tview.MouseLeftUp || action == tview.MouseLeftDown {
			if menuBar.selectedPath[0] != -1 {
				selectedIndex := menuBar.selectedPath[0]
				items := menuBar.menus[selectedIndex].Items
				index := y - ry - 2

				menuLeft := menuBar.menuIndexLeft(selectedIndex)
				width := menuWidthInCells(items)
				if x >= menuLeft && x < (menuLeft+width) && index >= 0 && index < len(items) {
					if action == tview.MouseLeftUp {
						if items[index].Title != "" {
							if items[index].Callback != nil {
								items[index].Callback(items[index].ID)
							}
							menuBar.Close()
						}
					}

					return true, nil
				}

				menuBar.Close()
			}
		}

		if action == tview.MouseLeftDown {
			menuBar.Close()
		}

		return true, nil
	})
}

func (menuBar *MenuBar) selectMenuBarItem(index int) {
	if index == -1 {
		menuBar.selectedPath = []int{-1}
	} else {
		menuBar.selectedPath = []int{index, 0}
	}
}

func (m *MenuBar) menuItemIndexAtX(posX int) (index int, leftX int) {
	x := 0
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
	x := 0
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

func (menuBar *MenuBar) Close() {
	menuBar.selectMenuBarItem(-1)
	if menuBar.onClose != nil {
		menuBar.onClose()
	}
}

func (menuBar *MenuBar) SetOnClose(callback func()) {
	menuBar.onClose = callback
}
