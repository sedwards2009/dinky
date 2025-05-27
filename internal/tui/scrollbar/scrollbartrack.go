package scrollbar

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ScrollbarTrack struct {
	*tview.Box

	// position of the scrollbar thumb.
	position int
	// thumbSize of the scrollbar thumb.
	thumbSize int
	// Maximum value of the scrollbar.
	max int

	// Minimum value of the scrollbar.
	min int

	// width of the scrollbar.
	width int

	trackColor tcell.Color
	thumbColor tcell.Color
}

func NewScrollbarTrack() *ScrollbarTrack {
	return &ScrollbarTrack{
		Box:        tview.NewBox(),
		position:   0,
		thumbSize:  10,
		max:        100,
		min:        0,
		width:      1, // Default width of the scrollbar
		trackColor: tcell.ColorDarkGray,
		thumbColor: tcell.ColorWhite,
	}
}

func (scrollbarTrack *ScrollbarTrack) Draw(screen tcell.Screen) {
	x, y, width, height := scrollbarTrack.GetInnerRect()
	if width < 1 || height < 1 {
		return
	}

	// Draw the track
	trackStyle := tcell.StyleDefault.Background(scrollbarTrack.trackColor)
	for i := 0; i < height; i++ {
		screen.SetContent(x, y+i, ' ', nil, trackStyle)
	}

	// Calculate the position and size of the scrollbar thumb.
	thumbHeight := (height * scrollbarTrack.thumbSize) / (scrollbarTrack.max - scrollbarTrack.min + 1)
	if thumbHeight < 1 {
		thumbHeight = 1
	}

	thumbY := y + height*(scrollbarTrack.position-scrollbarTrack.min)/(scrollbarTrack.max-scrollbarTrack.min)
	thumbStyle := tcell.StyleDefault.Background(scrollbarTrack.thumbColor)

	// Draw the scrollbar thumb.
	for i := 0; i < thumbHeight; i++ {
		if thumbY+i < height {
			screen.SetContent(x, thumbY+i, ' ', nil, thumbStyle)
		}
	}
}

func (scrollbarTrack *ScrollbarTrack) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return scrollbarTrack.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		rx, ry, _, height := scrollbarTrack.GetInnerRect()
		absEventX, absEventY := event.Position()
		eventX := absEventX - rx
		eventY := absEventY - ry
		if eventX < 0 || eventX >= scrollbarTrack.width || eventY < 0 || eventY >= height {
			return false, nil // Click outside the scrollbar
		}

		if action == tview.MouseLeftDown {
			// Calculate the new position based on the click
			// Assuming the scrollbar is vertical, we calculate the position based on the y coordinate
			newPosition := (absEventY-ry)*(scrollbarTrack.max-scrollbarTrack.min)/height + scrollbarTrack.min - scrollbarTrack.thumbSize/2
			if newPosition < scrollbarTrack.min {
				newPosition = scrollbarTrack.min
			} else if newPosition > scrollbarTrack.max {
				newPosition = scrollbarTrack.max
			}
			scrollbarTrack.position = newPosition
			return true, nil // Consumed the event
		}

		if action == tview.MouseScrollUp || action == tview.MouseScrollDown {
			// Handle scroll events
			if action == tview.MouseScrollUp {
				pos := scrollbarTrack.Position() - max(scrollbarTrack.ThumbSize()/2, 1)
				scrollbarTrack.SetPosition(pos)

			} else if action == tview.MouseScrollDown {
				pos := scrollbarTrack.Position() + max(scrollbarTrack.ThumbSize()/2, 1)
				scrollbarTrack.SetPosition(pos)
			}
			return true, nil // Consumed the event
		}
		return false, nil // Not consumed
	})
}

// SetPosition sets the position of the scrollbar thumb.
func (scrollbarTrack *ScrollbarTrack) SetPosition(position int) *ScrollbarTrack {
	if position < scrollbarTrack.min {
		scrollbarTrack.position = scrollbarTrack.min
	} else if position > scrollbarTrack.max-scrollbarTrack.thumbSize {
		scrollbarTrack.position = scrollbarTrack.max - scrollbarTrack.thumbSize
	} else {
		scrollbarTrack.position = position
	}
	return scrollbarTrack
}

// GetPosition returns the current position of the scrollbar thumb.
func (scrollbarTrack *ScrollbarTrack) Position() int {
	return scrollbarTrack.position
}

// SetThumbSize sets the size of the scrollbar thumb.
func (scrollbarTrack *ScrollbarTrack) SetThumbSize(size int) *ScrollbarTrack {
	if size < 1 {
		scrollbarTrack.thumbSize = 1
	} else {
		scrollbarTrack.thumbSize = size
	}
	return scrollbarTrack
}

// GetThumbSize returns the size of the scrollbar thumb.
func (scrollbarTrack *ScrollbarTrack) ThumbSize() int {
	return scrollbarTrack.thumbSize
}

// SetMax sets the maximum value of the scrollbar.
func (scrollbarTrack *ScrollbarTrack) SetMax(max int) *ScrollbarTrack {
	if max > scrollbarTrack.min {
		scrollbarTrack.max = max
	}
	return scrollbarTrack
}

// GetMax returns the maximum value of the scrollbar.
func (scrollbarTrack *ScrollbarTrack) Max() int {
	return scrollbarTrack.max
}

// SetMin sets the minimum value of the scrollbar.
func (scrollbarTrack *ScrollbarTrack) SetMin(min int) *ScrollbarTrack {
	if min < scrollbarTrack.max {
		scrollbarTrack.min = min
	}
	return scrollbarTrack
}

// GetMin returns the minimum value of the scrollbar.
func (scrollbarTrack *ScrollbarTrack) Min() int {
	return scrollbarTrack.min
}

// SetWidth sets the width of the scrollbar.
func (scrollbarTrack *ScrollbarTrack) SetWidth(width int) *ScrollbarTrack {
	if width > 0 {
		scrollbarTrack.width = width
	}
	return scrollbarTrack
}

// GetWidth returns the width of the scrollbar.
func (scrollbarTrack *ScrollbarTrack) Width() int {
	return scrollbarTrack.width
}

// SetTrackColor sets the color of the scrollbar track.
func (scrollbarTrack *ScrollbarTrack) SetTrackColor(color tcell.Color) *ScrollbarTrack {
	scrollbarTrack.trackColor = color
	return scrollbarTrack
}

// GetTrackColor returns the color of the scrollbar track.
func (scrollbarTrack *ScrollbarTrack) TrackColor() tcell.Color {
	return scrollbarTrack.trackColor
}

// SetThumbColor sets the color of the scrollbar thumb.
func (scrollbarTrack *ScrollbarTrack) SetThumbColor(color tcell.Color) *ScrollbarTrack {
	scrollbarTrack.thumbColor = color
	return scrollbarTrack
}

// GetThumbColor returns the color of the scrollbar thumb.
func (scrollbarTrack *ScrollbarTrack) ThumbColor() tcell.Color {
	return scrollbarTrack.thumbColor
}
