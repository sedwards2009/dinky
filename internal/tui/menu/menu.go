package menu

import (
	"dinky/internal/tui/utils"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/sedwards2009/nuview"
)

type MenuBar struct {
	*nuview.Box
	MenuBarStyle      tcell.Style
	MenuStyle         tcell.Style
	MenuSelectedStyle tcell.Style
	menus             []*Menu
	selectedPath      []int
	onClose           func()
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
	defaultStyle := tcell.StyleDefault.Foreground(fg).Background(bg).Bold(true)
	return &MenuBar{
		Box: nuview.NewBox(),

		MenuBarStyle:      defaultStyle,
		MenuStyle:         defaultStyle,
		MenuSelectedStyle: defaultStyle.Reverse(true),
		selectedPath:      []int{-1},
	}
}

func (menuBar *MenuBar) SetMenus(menus []*Menu) {
	menuBar.menus = menus
}

func (menuBar *MenuBar) Open() {
	menuBar.selectMenuBarItem(0)
}

func (menuBar *MenuBar) Draw(screen tcell.Screen) {
	// menuBar.Box.DrawForSubclass(screen, menuBar)
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

	borderStyle := menuBar.MenuStyle

	utils.DrawCappedHorizontalLine(screen, menuX, y, menuWidth, borderStyle, borderStyle, '┌', '─', '┐')
	y++

	for i, item := range items {
		if item.Title == "" {
			utils.DrawCappedHorizontalLine(screen, menuX, y, menuWidth, borderStyle, borderStyle, '├', '─', '┤')
		} else {
			textStyle := borderStyle
			if i == selectedIndex {
				textStyle = menuBar.MenuSelectedStyle
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
	screen.HideCursor()
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

func (menuBar *MenuBar) InputHandler() func(event *tcell.EventKey, setFocus func(p nuview.Primitive)) {
	return menuBar.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p nuview.Primitive)) {
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
				menuBar.Close()
				if item.Callback != nil {
					item.Callback(item.ID)
				}
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

func (menuBar *MenuBar) MouseHandler() func(action nuview.MouseAction, event *tcell.EventMouse, setFocus func(p nuview.Primitive)) (consumed bool, capture nuview.Primitive) {
	return menuBar.WrapMouseHandler(func(action nuview.MouseAction, event *tcell.EventMouse, setFocus func(p nuview.Primitive)) (consumed bool, capture nuview.Primitive) {
		rx, ry, _, _ := menuBar.GetRect()
		x, y := event.Position()

		if y == ry {
			if action == nuview.MouseLeftDown {
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

		selectedIndex := menuBar.selectedPath[0]
		if selectedIndex == -1 { // Is a menu open?
			return false, nil
		}

		var selectedMenuItem *MenuItem = nil
		items := menuBar.menus[selectedIndex].Items
		selectedMenuItemIndex := y - ry - 2
		menuLeft := menuBar.menuIndexLeft(selectedIndex)
		width := menuWidthInCells(items)
		if x >= menuLeft && x < (menuLeft+width) && selectedMenuItemIndex >= 0 && selectedMenuItemIndex < len(items) {
			selectedMenuItem = items[selectedMenuItemIndex]
		}

		if action == nuview.MouseLeftUp {
			if selectedMenuItem != nil && selectedMenuItem.Title != "" {
				items[selectedMenuItemIndex].Callback(items[selectedMenuItemIndex].ID)
			}
			menuBar.Close()
			return true, nil
		}

		if action == nuview.MouseLeftDown && selectedMenuItem == nil {
			menuBar.Close()
			return true, nil
		}

		if selectedMenuItem != nil {
			menuBar.selectedPath[1] = selectedMenuItemIndex
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
