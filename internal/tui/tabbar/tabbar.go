package tabbar

import (
	"slices"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/sedwards2009/nuview"
)

const tabNamePadding = 1

type TabBar struct {
	*nuview.Box
	BackgroundStyle  tcell.Style
	ActiveTabStyle   tcell.Style
	InactiveTabStyle tcell.Style
	tabs             []Tab
	active           int
	hscroll          int
	OnActive         func(id string, index int)
	OnTabCloseClick  func(id string, index int)
}

type Tab struct {
	Title string
	ID    string
}

func NewTabBar() *TabBar {
	fg := tcell.NewHexColor(0xf3f3f3) // White
	bg := tcell.NewHexColor(0x007ace) // Blue
	tabBg := tcell.NewHexColor(0x000000)
	inactiveTabBg := tcell.NewHexColor(0x404040)
	return &TabBar{
		Box:              nuview.NewBox(),
		BackgroundStyle:  tcell.StyleDefault.Foreground(fg).Background(bg).Bold(true),
		ActiveTabStyle:   tcell.StyleDefault.Foreground(fg).Background(tabBg).Bold(true),
		InactiveTabStyle: tcell.StyleDefault.Foreground(fg).Background(inactiveTabBg).Bold(false),
	}
}

func (tabBar *TabBar) SetTabBackgroundColor(color tcell.Color) {
	tabBar.ActiveTabStyle = tabBar.ActiveTabStyle.Background(color)
}

func (tabBar *TabBar) SetTabInactiveBackgroundColor(color tcell.Color) {
	tabBar.InactiveTabStyle = tabBar.InactiveTabStyle.Background(color)
}

func (tabBar *TabBar) Active() (int, string) {
	if len(tabBar.tabs) == 0 {
		return -1, ""
	}
	return tabBar.active, tabBar.tabs[tabBar.active].ID
}

func (tabBar *TabBar) SetActive(id string) {
	for i, tab := range tabBar.tabs {
		if tab.ID == id {
			tabBar.active = i
		}
	}
}

func (tabBar *TabBar) AddTab(title string, id string) {
	tabBar.tabs = append(tabBar.tabs, Tab{Title: title, ID: id})
}

func (tabBar *TabBar) RemoveTab(id string) {
	for i, tab := range tabBar.tabs {
		if tab.ID == id {
			tabBar.tabs = slices.Delete(tabBar.tabs, i, i+1)
			break
		}
	}
}

func (tabBar *TabBar) SetTabTitle(id string, title string) {
	for i, tab := range tabBar.tabs {
		if tab.ID == id {
			tabBar.tabs[i].Title = title
			break
		}
	}
}

func (tabBar *TabBar) Draw(screen tcell.Screen) {
	x, y, width, _ := tabBar.GetInnerRect()

	x = x - tabBar.hscroll
	done := false

	tabBarStyle := tabBar.BackgroundStyle
	_, tabBarBg, _ := tabBarStyle.Decompose()

	tabStyle := tabBar.ActiveTabStyle
	_, tabBg, _ := tabStyle.Decompose()

	tabInactiveStyle := tabBar.InactiveTabStyle
	_, tabInactiveBg, _ := tabInactiveStyle.Decompose()

	tabCornerStyle := tcell.Style{}.Foreground(tabBg).Background(tabBarBg)
	tabCornerInactiveStyle := tcell.Style{}.Foreground(tabInactiveBg).Background(tabBarBg)

	draw := func(r rune, n int, style tcell.Style) {
		for range n {
			rw := runewidth.RuneWidth(r)
			for j := range rw {
				c := r
				if j > 0 {
					c = ' '
				}
				if x == width-1 && !done {
					screen.SetContent(width-1, y, '>', nil, tabBarStyle)
					x++
					break
				} else if x == 0 && tabBar.hscroll > 0 {
					screen.SetContent(0, y, '<', nil, tabBarStyle)
				} else if x >= 0 && x < width {
					screen.SetContent(x, y, c, nil, style)
				}
				x++
			}
		}
	}

	for i, tab := range tabBar.tabs {
		currentTabCornerStyle := tabCornerInactiveStyle
		currentTabTextStyle := tabInactiveStyle
		if i == tabBar.active {
			currentTabCornerStyle = tabCornerStyle
			currentTabTextStyle = tabStyle
		}

		draw('◢', 1, currentTabCornerStyle)

		for j := 0; j < tabNamePadding; j++ {
			draw(' ', 1, currentTabTextStyle)
		}

		for _, c := range tab.Title {
			draw(c, 1, currentTabTextStyle)
		}
		draw(' ', 1, currentTabTextStyle)
		draw('\u2715', 1, currentTabTextStyle)

		if i == len(tabBar.tabs)-1 {
			done = true
		}

		draw('◣', 1, currentTabCornerStyle)
		draw(' ', 2, tabBarStyle)

		if x >= width {
			break
		}
	}

	if x < width {
		draw(' ', width-x, tabBarStyle)
	}
}

func (tabBar *TabBar) MouseHandler() func(action nuview.MouseAction, event *tcell.EventMouse, setFocus func(p nuview.Primitive)) (consumed bool, capture nuview.Primitive) {
	return tabBar.WrapMouseHandler(func(action nuview.MouseAction, event *tcell.EventMouse, setFocus func(p nuview.Primitive)) (consumed bool, capture nuview.Primitive) {
		rx, ry, _, _ := tabBar.GetRect()
		x, y := event.Position()

		if y == ry {
			if action == nuview.MouseLeftDown {
				index, _, closeClick := tabBar.tabIndexAtX(x - rx)
				if index != -1 {
					tabBar.active = index
					if tabBar.OnActive != nil {
						tabBar.OnActive(tabBar.tabs[index].ID, index)
					}
					if closeClick && tabBar.OnTabCloseClick != nil {
						tabBar.OnTabCloseClick(tabBar.tabs[index].ID, index)
					}
					return true, nil
				}
			}
			return true, nil
		}

		return false, nil
	})
}

func (tabBar *TabBar) tabIndexAtX(posX int) (index int, leftX int, closeClick bool) {
	x := 0
	for i, tab := range tabBar.tabs {
		if posX < x {
			return -1, -1, false
		}

		left := x
		x += 1 // '◢'
		x += tabNamePadding
		x += runewidth.StringWidth(tab.Title)
		x += 1 // ' '
		x += 1 // '✕'
		x += 1 // '◣'
		if posX < x {
			closeClick = posX == x-2
			return i, left, closeClick
		}
		x += 2 // gap
	}

	return -1, -1, false
}
