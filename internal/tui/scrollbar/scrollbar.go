package scrollbar

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Scrollbar struct {
	*tview.Flex
	Track       *ScrollbarTrack
	UpButton    *tview.Button
	DownButton  *tview.Button
	changedFunc func(position int)

	isHorizontal bool                // Indicates if the scrollbar is horizontal instead of vertical
	UpdateHook   func(sb *Scrollbar) // Hook for updating the scrollbar just before it is drawn
}

func NewScrollbar() *Scrollbar {
	scrollbarTrack := NewScrollbarTrack()
	upButton := tview.NewButton("▲")
	upButton.SetRect(0, 0, 1, 1)
	downButton := tview.NewButton("▼")
	downButton.SetRect(0, 0, 1, 1)

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(scrollbarTrack, 0, 1, false)
	flex.AddItem(upButton, 1, 0, false)
	flex.AddItem(downButton, 1, 0, false)

	scrollbar := &Scrollbar{
		Flex:       flex,
		Track:      scrollbarTrack,
		UpButton:   upButton,
		DownButton: downButton,
	}

	upButton.SetSelectedFunc(func() {
		pos := scrollbar.Track.Position() - max(scrollbar.Track.ThumbSize()/2, 1)
		newPos := scrollbar.Track.SetPosition(pos)
		if scrollbar.changedFunc != nil {
			scrollbar.changedFunc(newPos)
		}
	})

	downButton.SetSelectedFunc(func() {
		pos := scrollbar.Track.Position() + max(scrollbar.Track.ThumbSize()/2, 1)
		newPos := scrollbar.Track.SetPosition(pos)
		if scrollbar.changedFunc != nil {
			scrollbar.changedFunc(newPos)
		}
	})

	return scrollbar
}

func (scrollbar *Scrollbar) SetHorizontal(isHorizontal bool) {
	scrollbar.isHorizontal = isHorizontal
	scrollbar.Track.SetHorizontal(isHorizontal)

	if isHorizontal {
		scrollbar.Flex.SetDirection(tview.FlexColumn)
		scrollbar.UpButton.SetLabel("\u25c4")
		scrollbar.DownButton.SetLabel("\u25ba")
	} else {
		scrollbar.Flex.SetDirection(tview.FlexRow)
		scrollbar.UpButton.SetLabel("\u25b2")
		scrollbar.DownButton.SetLabel("\u25bc")
	}
}

func (scrollbar *Scrollbar) SetChangedFunc(changedFunc func(position int)) {
	scrollbar.changedFunc = changedFunc
	scrollbar.Track.SetChangedFunc(changedFunc)
}

func (scrollbar *Scrollbar) Draw(screen tcell.Screen) {
	if scrollbar.UpdateHook != nil {
		scrollbar.UpdateHook(scrollbar)
	}
	scrollbar.Flex.Draw(screen)
}
