package tabbar

import (
	"dinky/internal/tui/stylecolor"
	"math"
	"slices"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/tview"
)

const tabNamePadding = 1
const tabScrollStep = 20

type TabBar struct {
	*tview.Box
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
	width int
}

func NewTabBar() *TabBar {
	fg := stylecolor.White
	bg := stylecolor.Blue
	tabBg := stylecolor.Black
	inactiveTabBg := stylecolor.InactiveGrey
	return &TabBar{
		Box:              tview.NewBox(),
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
	pos := 0
	for i, tab := range tabBar.tabs {
		if tab.ID == id {
			tabBar.active = i

			_, _, width, _ := tabBar.GetInnerRect()
			if tabBar.isOverflow() {
				width -= 2
			}
			overhang := pos - tabBar.hscroll + tab.width - width
			if overhang > 0 {
				tabBar.hscroll += overhang
			}

			if tabBar.hscroll > pos {
				tabBar.hscroll = pos
			}
			break
		}
		pos += tab.width
	}
}

func (tabBar *TabBar) AddTab(title string, id string) {
	newTab := Tab{Title: title, ID: id}
	newTab.width = computeTabWidth(newTab)
	tabBar.tabs = append(tabBar.tabs, newTab)
}

func computeTabWidth(tab Tab) int {
	return 1 + tabNamePadding + runewidth.StringWidth(tab.Title) + 1 + 1 + 1 + 2
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

func (tabBar *TabBar) totalBarWidth() int {
	totalWidth := 0
	for _, tab := range tabBar.tabs {
		totalWidth += tab.width
	}
	return totalWidth
}

func (tabBar *TabBar) isOverflow() bool {
	_, _, width, _ := tabBar.GetInnerRect()
	isOverflow := tabBar.totalBarWidth() > width
	return isOverflow
}

func (tabBar *TabBar) Draw(screen tcell.Screen) {
	x, y, width, _ := tabBar.GetInnerRect()

	x = x - tabBar.hscroll

	tabBarStyle := tabBar.BackgroundStyle
	_, tabBarBg, _ := tabBarStyle.Decompose()

	tabStyle := tabBar.ActiveTabStyle
	_, tabBg, _ := tabStyle.Decompose()

	tabInactiveStyle := tabBar.InactiveTabStyle
	_, tabInactiveBg, _ := tabInactiveStyle.Decompose()

	tabCornerStyle := tcell.Style{}.Foreground(tabBg).Background(tabBarBg)
	tabCornerInactiveStyle := tcell.Style{}.Foreground(tabInactiveBg).Background(tabBarBg)

	isOverflow := tabBar.isOverflow()

	clipWidth := width
	draw := func(r rune, n int, style tcell.Style) {
		for range n {
			rw := runewidth.RuneWidth(r)
			for j := range rw {
				c := r
				if j > 0 {
					c = ' '
				}
				if x >= 0 && x < clipWidth {
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
		draw('◣', 1, currentTabCornerStyle)
		draw(' ', 2, tabBarStyle)

		if x >= width {
			break
		}
	}

	if x < width {
		draw(' ', width-x, tabBarStyle)
	}

	if isOverflow {
		screen.SetContent(width-2, y, '\u2BC7', nil, tabBarStyle)
		screen.SetContent(width-1, y, '\u2BC8', nil, tabBarStyle)
	}
}

func (tabBar *TabBar) clampHScroll() {
	_, _, width, _ := tabBar.GetInnerRect()
	if tabBar.hscroll < 0 {
		tabBar.hscroll = 0
	}
	overflowWidth := width - tabBar.totalBarWidth()
	if overflowWidth < 0 && tabBar.hscroll > -overflowWidth {
		tabBar.hscroll = -overflowWidth
	}
}

func (tabBar *TabBar) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return tabBar.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		rx, ry, _, _ := tabBar.GetRect()
		x, y := event.Position()

		if y == ry {
			if action == tview.MouseLeftDown {
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

				if tabBar.isOverflow() {
					relX := x - rx
					_, _, width, _ := tabBar.GetInnerRect()
					if relX == width-2 {
						// Clicked on left overflow indicator
						tabBar.hscroll -= tabScrollStep
						tabBar.clampHScroll()
						return true, nil
					} else if relX == width-1 {
						// Clicked on right overflow indicator
						tabBar.hscroll += tabScrollStep
						tabBar.clampHScroll()
						return true, nil
					}

				}
			} else if action == tview.MouseScrollUp {
				tabBar.hscroll -= tabScrollStep
				tabBar.clampHScroll()
				return true, nil
			} else if action == tview.MouseScrollDown {
				tabBar.hscroll += tabScrollStep
				tabBar.clampHScroll()
			}
			return true, nil
		}

		return false, nil
	})
}

func (tabBar *TabBar) tabIndexAtX(relativeX int) (index int, leftX int, closeClick bool) {
	posX := relativeX + tabBar.hscroll

	overflowLeftPos := math.MaxInt
	if tabBar.isOverflow() {
		_, _, width, _ := tabBar.GetInnerRect()
		overflowLeftPos = width - 2 + tabBar.hscroll
	}

	x := 0
	for i, tab := range tabBar.tabs {
		if posX < x {
			return -1, -1, false
		}

		if posX >= overflowLeftPos {
			return -1, -1, false
		}

		if posX < x+tab.width {
			closeClick := posX == x+tab.width-4
			return i, x, closeClick
		}
		x += tab.width
	}

	return -1, -1, false
}
