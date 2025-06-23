package scrollbar

import (
	"github.com/sedwards2009/nuview"
)

type Scrollbar struct {
	*nuview.Flex
	Track       *ScrollbarTrack
	upButton    *nuview.Button
	downButton  *nuview.Button
	changedFunc func(position int)
}

func NewScrollbar() *Scrollbar {
	scrollbarTrack := NewScrollbarTrack()
	upButton := nuview.NewButton("▲")
	upButton.SetRect(0, 0, 1, 1)
	downButton := nuview.NewButton("▼")
	downButton.SetRect(0, 0, 1, 1)

	flex := nuview.NewFlex()
	flex.SetDirection(nuview.FlexRow)
	flex.AddItem(scrollbarTrack, 0, 1, false)
	flex.AddItem(upButton, 1, 0, false)
	flex.AddItem(downButton, 1, 0, false)

	scrollbar := &Scrollbar{
		Flex:       flex,
		Track:      scrollbarTrack,
		upButton:   upButton,
		downButton: downButton,
	}

	upButton.SetSelectedFunc(func() {
		pos := scrollbar.Track.Position() - max(scrollbar.Track.ThumbSize()/2, 1)
		scrollbar.Track.SetPosition(pos)
		if scrollbar.changedFunc != nil {
			scrollbar.changedFunc(pos)
		}
	})

	downButton.SetSelectedFunc(func() {
		pos := scrollbar.Track.Position() + max(scrollbar.Track.ThumbSize()/2, 1)
		scrollbar.Track.SetPosition(pos)
		if scrollbar.changedFunc != nil {
			scrollbar.changedFunc(pos)
		}
	})

	return scrollbar
}

func (scrollbar *Scrollbar) SetChangedFunc(changedFunc func(position int)) {
	scrollbar.changedFunc = changedFunc
	scrollbar.Track.SetChangedFunc(changedFunc)
}
