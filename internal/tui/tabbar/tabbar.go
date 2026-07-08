package tabbar

import (
	"dinky/internal/tui/stylecolor"
	"math"
	"slices"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/tview"
)

const tabNamePadding = 1
const tabScrollStep = 20

// While dragging a tab against either edge of an overflowing bar, scroll this
// many columns per tick so the drag can reach off-screen tabs.
const dragAutoScrollStep = 3
const dragAutoScrollInterval = 50 * time.Millisecond

// Tabs slide toward their settled position rather than jumping. Each animation
// tick moves a tab tabAnimEase of the remaining distance toward its target.
const tabAnimInterval = 16 * time.Millisecond
const tabAnimEase = 0.3

type TabBar struct {
	*tview.Box
	BackgroundStyle  tcell.Style
	ActiveTabStyle   tcell.Style
	InactiveTabStyle tcell.Style
	tabs             []Tab
	active           int
	hscroll          int
	dragging         bool
	dragIndex        int
	dragMouseX       int
	// dragGrabOffset is the distance, in layout coordinates, from the dragged
	// tab's left edge to the point the pointer grabbed it, so the tab stays
	// pinned under the cursor at that same point throughout the drag.
	dragGrabOffset  float64
	dragStop        chan struct{}
	animating       bool
	animStop        chan struct{}
	OnActive        func(id string, index int)
	OnTabCloseClick func(id string, index int)
	OnReorder       func(id string, newIndex int)
	// QueueUpdateDraw schedules f to run on the UI goroutine and redraws.
	// When set, it powers timer-driven auto-scroll during a tab drag and the
	// slide animation when tabs change position.
	QueueUpdateDraw func(f func())
}

type Tab struct {
	Title string
	ID    string
	width int
	// renderX is the tab's current left edge in layout coordinates (before
	// hscroll), as a float so it can ease toward its target. placed is false
	// until renderX has been initialised so freshly added tabs appear in place
	// rather than sliding in from the left.
	renderX float64
	placed  bool
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
			// The tabs to the right shift left into the gap; slide them.
			tabBar.startAnim()
			break
		}
	}
}

func (tabBar *TabBar) SetTabTitle(id string, title string) {
	for i, tab := range tabBar.tabs {
		if tab.ID == id {
			tabBar.tabs[i].Title = title
			newWidth := computeTabWidth(tabBar.tabs[i])
			if newWidth != tabBar.tabs[i].width {
				tabBar.tabs[i].width = newWidth
				// A width change shifts every following tab; slide them.
				tabBar.startAnim()
			}
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
	x0, y, width, _ := tabBar.GetInnerRect()

	tabBar.ensurePlaced()

	tabBarStyle := tabBar.BackgroundStyle

	// Paint the whole bar background first; tabs are then drawn on top at their
	// (possibly mid-animation) positions, which may overlap during a slide.
	for col := 0; col < width; col++ {
		screen.SetContent(x0+col, y, ' ', nil, tabBarStyle)
	}

	// Draw inactive tabs first, then the active tab, so the dragged/active tab
	// stays on top while tabs cross over one another.
	for i := range tabBar.tabs {
		if i == tabBar.active {
			continue
		}
		tabBar.drawTab(screen, i, x0, y, width)
	}
	if tabBar.active >= 0 && tabBar.active < len(tabBar.tabs) {
		tabBar.drawTab(screen, tabBar.active, x0, y, width)
	}

	if tabBar.isOverflow() {
		screen.SetContent(x0+width-2, y, '⯇', nil, tabBarStyle)
		screen.SetContent(x0+width-1, y, '⯈', nil, tabBarStyle)
	}
}

// drawTab renders a single tab's glyphs at its current animated position,
// clipped to the bar's inner rect.
func (tabBar *TabBar) drawTab(screen tcell.Screen, i, x0, y, width int) {
	tab := tabBar.tabs[i]

	_, tabBarBg, _ := tabBar.BackgroundStyle.Decompose()
	_, tabBg, _ := tabBar.ActiveTabStyle.Decompose()
	_, tabInactiveBg, _ := tabBar.InactiveTabStyle.Decompose()

	cornerStyle := tcell.Style{}.Foreground(tabInactiveBg).Background(tabBarBg)
	textStyle := tabBar.InactiveTabStyle
	if i == tabBar.active {
		cornerStyle = tcell.Style{}.Foreground(tabBg).Background(tabBarBg)
		textStyle = tabBar.ActiveTabStyle
	}

	x := x0 - tabBar.hscroll + int(math.Round(tab.renderX))
	put := func(r rune, style tcell.Style) {
		rw := runewidth.RuneWidth(r)
		for j := 0; j < rw; j++ {
			c := r
			if j > 0 {
				c = ' '
			}
			if x >= x0 && x < x0+width {
				screen.SetContent(x, y, c, nil, style)
			}
			x++
		}
	}

	put('◢', cornerStyle)
	for j := 0; j < tabNamePadding; j++ {
		put(' ', textStyle)
	}
	for _, c := range tab.Title {
		put(c, textStyle)
	}
	put(' ', textStyle)
	put('✕', textStyle)
	put('◣', cornerStyle)
}

// ensurePlaced gives any not-yet-placed tab its settled render position so new
// tabs appear in place instead of sliding in from the origin.
func (tabBar *TabBar) ensurePlaced() {
	x := 0
	for i := range tabBar.tabs {
		if !tabBar.tabs[i].placed {
			tabBar.tabs[i].renderX = float64(x)
			tabBar.tabs[i].placed = true
		}
		x += tabBar.tabs[i].width
	}
}

// startAnim eases every tab toward its settled slot. When no async redraw hook
// is available it snaps instantly so the bar stays correct in headless tests.
func (tabBar *TabBar) startAnim() {
	if tabBar.QueueUpdateDraw == nil {
		x := 0
		for i := range tabBar.tabs {
			tabBar.tabs[i].renderX = float64(x)
			tabBar.tabs[i].placed = true
			x += tabBar.tabs[i].width
		}
		return
	}
	if tabBar.animating {
		return
	}
	tabBar.animating = true
	stop := make(chan struct{})
	tabBar.animStop = stop
	go func() {
		ticker := time.NewTicker(tabAnimInterval)
		defer ticker.Stop()
		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				tabBar.QueueUpdateDraw(func() {
					if !tabBar.animateStep() {
						tabBar.stopAnim()
					}
				})
			}
		}
	}()
}

// stopAnim halts the easing ticker. Runs on the UI goroutine.
func (tabBar *TabBar) stopAnim() {
	if !tabBar.animating {
		return
	}
	tabBar.animating = false
	if tabBar.animStop != nil {
		close(tabBar.animStop)
		tabBar.animStop = nil
	}
}

// animateStep moves every tab a fraction of the way toward its target slot.
// Returns true while any tab is still in motion. Runs on the UI goroutine.
func (tabBar *TabBar) animateStep() bool {
	moving := false
	x := 0
	for i := range tabBar.tabs {
		tabBar.tabs[i].placed = true
		// The dragged tab is positioned live under the cursor, so leave it be;
		// still advance x by its width so the following tabs target the slot
		// it currently occupies.
		if tabBar.dragging && i == tabBar.dragIndex {
			x += tabBar.tabs[i].width
			continue
		}
		target := float64(x)
		delta := target - tabBar.tabs[i].renderX
		if math.Abs(delta) < 0.5 {
			tabBar.tabs[i].renderX = target
		} else {
			tabBar.tabs[i].renderX += delta * tabAnimEase
			moving = true
		}
		x += tabBar.tabs[i].width
	}
	return moving
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

		// While a drag is in progress we capture all mouse events so the tab
		// keeps following the pointer even if it strays off the tab bar row.
		if tabBar.dragging {
			switch action {
			case tview.MouseMove:
				tabBar.dragMouseX = x - rx
				tabBar.updateDragPosition()
				return true, tabBar
			case tview.MouseLeftUp:
				tabBar.stopDrag()
				return true, nil
			default:
				return true, tabBar
			}
		}

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
					} else {
						// Begin a drag-to-reorder. Capture future mouse
						// events until the button is released.
						tabBar.startDrag(index, x-rx)
						return true, tabBar
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

// startDrag begins dragging the tab at index, with relX the pointer position
// relative to the bar's left edge. If a redraw hook is available it also starts
// a ticker that auto-scrolls while the pointer is held against an edge.
func (tabBar *TabBar) startDrag(index int, relX int) {
	tabBar.dragging = true
	tabBar.dragIndex = index
	tabBar.dragMouseX = relX

	// Record where within the tab the pointer landed so the tab stays pinned
	// under the cursor at that exact spot for the whole drag.
	tabBar.ensurePlaced()
	pointerLayoutX := float64(relX + tabBar.hscroll)
	tabBar.dragGrabOffset = pointerLayoutX - tabBar.tabs[index].renderX

	if tabBar.QueueUpdateDraw == nil {
		return
	}
	stop := make(chan struct{})
	tabBar.dragStop = stop
	go func() {
		ticker := time.NewTicker(dragAutoScrollInterval)
		defer ticker.Stop()
		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				tabBar.QueueUpdateDraw(func() {
					tabBar.autoScrollTick()
				})
			}
		}
	}()
}

// stopDrag ends the current drag and tears down the auto-scroll ticker. The
// dropped tab is left wherever the cursor released it and then eased into its
// settled slot.
func (tabBar *TabBar) stopDrag() {
	tabBar.dragging = false
	tabBar.dragIndex = -1
	if tabBar.dragStop != nil {
		close(tabBar.dragStop)
		tabBar.dragStop = nil
	}
	// Snap-slide the dropped tab into its slot now that it is no longer pinned.
	tabBar.startAnim()
}

// updateDragPosition reorders the tabs if the pointer has crossed into another
// tab's slot, then pins the dragged tab live under the cursor. Runs on the UI
// goroutine (from the mouse handler or the auto-scroll tick).
func (tabBar *TabBar) updateDragPosition() {
	target := tabBar.dragTargetIndex(tabBar.dragMouseX)
	if target != tabBar.dragIndex && target != -1 {
		tabBar.moveTab(tabBar.dragIndex, target)
		tabBar.dragIndex = target
		tabBar.active = target
		// Slide the displaced tabs into their new slots. The dragged tab itself
		// is excluded from the easing (animateStep) and tracks the cursor.
		tabBar.startAnim()
		if tabBar.OnReorder != nil {
			tabBar.OnReorder(tabBar.tabs[target].ID, target)
		}
	}
	tabBar.pinDraggedTab()
}

// pinDraggedTab sets the dragged tab's render position so the grabbed point
// stays directly under the cursor, clamped to the bar's logical extent.
func (tabBar *TabBar) pinDraggedTab() {
	if !tabBar.dragging || tabBar.dragIndex < 0 || tabBar.dragIndex >= len(tabBar.tabs) {
		return
	}
	pointerLayoutX := float64(tabBar.dragMouseX + tabBar.hscroll)
	renderX := pointerLayoutX - tabBar.dragGrabOffset

	maxX := float64(tabBar.totalBarWidth() - tabBar.tabs[tabBar.dragIndex].width)
	if maxX < 0 {
		maxX = 0
	}
	if renderX < 0 {
		renderX = 0
	} else if renderX > maxX {
		renderX = maxX
	}

	tabBar.tabs[tabBar.dragIndex].renderX = renderX
	tabBar.tabs[tabBar.dragIndex].placed = true
}

// autoScrollTick is fired by the drag ticker. If the pointer is held within an
// edge zone of an overflowing bar it scrolls a step and re-evaluates the drag.
func (tabBar *TabBar) autoScrollTick() {
	if !tabBar.dragging || !tabBar.isOverflow() {
		return
	}
	_, _, width, _ := tabBar.GetInnerRect()
	before := tabBar.hscroll
	if tabBar.dragMouseX <= 1 {
		tabBar.hscroll -= dragAutoScrollStep
	} else if tabBar.dragMouseX >= width-2 {
		tabBar.hscroll += dragAutoScrollStep
	}
	tabBar.clampHScroll()
	if tabBar.hscroll != before {
		tabBar.updateDragPosition()
	}
}

// dragTargetIndex returns the index at which a dragged tab should be inserted
// for the given pointer position. Unlike tabIndexAtX it never returns -1 (the
// pointer is clamped to the first/last tab) so dragging stays responsive even
// past the ends of the bar.
func (tabBar *TabBar) dragTargetIndex(relativeX int) int {
	if len(tabBar.tabs) == 0 {
		return -1
	}
	posX := relativeX + tabBar.hscroll
	if posX < 0 {
		return 0
	}
	x := 0
	for i, tab := range tabBar.tabs {
		if posX < x+tab.width/2 {
			return i
		}
		x += tab.width
	}
	return len(tabBar.tabs) - 1
}

// moveTab moves the tab at index from to index to, shifting the others. The
// tabs keep their renderX values so the change can be animated afterwards.
func (tabBar *TabBar) moveTab(from, to int) {
	if from == to || from < 0 || to < 0 || from >= len(tabBar.tabs) || to >= len(tabBar.tabs) {
		return
	}
	tab := tabBar.tabs[from]
	tabBar.tabs = slices.Delete(tabBar.tabs, from, from+1)
	tabBar.tabs = slices.Insert(tabBar.tabs, to, tab)
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
