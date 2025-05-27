package scrollbar

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Scrollbar struct {
	*tview.Flex
	Track      *ScrollbarTrack
	upButton   *tview.Button
	downButton *tview.Button
}

func NewScrollbar() *Scrollbar {
	scrollbarTrack := NewScrollbarTrack()
	upButton := tview.NewButton("▲")
	upButton.SetRect(0, 0, 1, 1)
	downButton := tview.NewButton("▼")
	downButton.SetRect(0, 0, 1, 1)

	scrollbar := &Scrollbar{
		Flex: tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(scrollbarTrack, 0, 1, false).
			AddItem(upButton, 1, 0, false).
			AddItem(downButton, 1, 0, false),
		Track:      scrollbarTrack,
		upButton:   upButton,
		downButton: downButton,
	}

	upButton.SetSelectedFunc(func() {
		pos := scrollbar.Track.Position() - scrollbar.Track.ThumbSize()
		scrollbar.Track.SetPosition(pos)
	})

	downButton.SetSelectedFunc(func() {
		pos := scrollbar.Track.Position() + scrollbar.Track.ThumbSize()
		scrollbar.Track.SetPosition(pos)
	})

	return scrollbar
}

func (scrollbar *Scrollbar) SetButtonStyle(style tcell.Style) {
	scrollbar.upButton.SetStyle(style)
	scrollbar.downButton.SetStyle(style)
}
