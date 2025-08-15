package scrollbar

import (
	"github.com/gdamore/tcell/v2"
	"github.com/sedwards2009/nuview"
)

type ScrollbarTrack struct {
	*nuview.Box

	// position of the scrollbar thumb.
	position int
	// thumbSize of the scrollbar thumb.
	thumbSize int
	// Maximum value of the scrollbar.
	max int

	// width of the scrollbar.
	width int

	trackColor tcell.Color
	thumbColor tcell.Color

	beforeDrawFunc func(screen tcell.Screen)
	changedFunc    func(position int)

	isHorizontal bool // Indicates if the scrollbar is horizontal instead of vertical
}

func NewScrollbarTrack() *ScrollbarTrack {
	return &ScrollbarTrack{
		Box:        nuview.NewBox(),
		position:   0,
		thumbSize:  10,
		max:        100,
		width:      1, // Default width of the scrollbar
		trackColor: tcell.ColorDarkGray,
		thumbColor: tcell.ColorWhite,
	}
}

func (scrollbarTrack *ScrollbarTrack) SetBeforeDrawFunc(beforeDrawFunc func(screen tcell.Screen)) {
	scrollbarTrack.beforeDrawFunc = beforeDrawFunc
}

func (scrollbarTrack *ScrollbarTrack) SetHorizontal(isHorizontal bool) {
	scrollbarTrack.isHorizontal = isHorizontal
}

func (scrollbarTrack *ScrollbarTrack) Draw(screen tcell.Screen) {
	if scrollbarTrack.beforeDrawFunc != nil {
		scrollbarTrack.beforeDrawFunc(screen)
	}

	innerX, innerY, width, height := scrollbarTrack.GetInnerRect()
	if width < 1 || height < 1 {
		return
	}

	firstHalfCellRune := '\u2584'
	secondHalfCellRune := '\u2580'

	majorLength := height
	majorPos := innerY
	if scrollbarTrack.isHorizontal {
		majorLength = width
		majorPos = innerX
		firstHalfCellRune = '\u2590'
		secondHalfCellRune = '\u258c'
	}

	setContent := func(n int, ch rune, style tcell.Style) {
		if scrollbarTrack.isHorizontal {
			screen.SetContent(n, innerY, ch, nil, style)
		} else {
			screen.SetContent(innerX, n, ch, nil, style)
		}
	}

	// Draw the track
	trackStyle := tcell.StyleDefault.Background(scrollbarTrack.trackColor)
	for i := 0; i < majorLength; i++ {
		setContent(majorPos+i, ' ', trackStyle)
	}

	position := scrollbarTrack.position
	thumbSize := scrollbarTrack.thumbSize
	if thumbSize > scrollbarTrack.max {
		thumbSize = scrollbarTrack.max
		position = 0
	}

	doubleMajorLength := majorLength * 2
	// Calculate the position and size of the scrollbar thumb.
	doubleThumbSizeFloat := float64(doubleMajorLength) * float64(thumbSize) / float64(scrollbarTrack.max)
	doubleThumbSize := int(doubleThumbSizeFloat + 0.5) // Round to nearest integer

	doubleThumbMajor := doubleMajorLength * position / scrollbarTrack.max
	thumbStyle := tcell.StyleDefault.Foreground(scrollbarTrack.trackColor).Background(scrollbarTrack.thumbColor)
	thumbReverseStyle := thumbStyle.Reverse(true)

	if position == scrollbarTrack.max-thumbSize {
		// Special case for when the thumb is at the bottom and we don't want to show a gap due to rounding
		doubleThumbSize += 2
	}

	if doubleThumbMajor&1 == 1 {
		// Draw the top of the thumb in the bottom half of the cell
		setContent(majorPos+doubleThumbMajor>>1, firstHalfCellRune, thumbReverseStyle)
		doubleThumbMajor++
		doubleThumbSize--
	}

	if doubleThumbSize&1 == 1 {
		// Draw the bottom part of the thumb in the top half of the cell
		setContent(majorPos+doubleThumbMajor>>1+doubleThumbSize>>1, secondHalfCellRune, thumbReverseStyle)
		doubleThumbSize--
	}

	// Draw the scrollbar thumb.
	thumbPos := doubleThumbMajor >> 1 // Convert back to single height
	for i := 0; i < doubleThumbSize>>1; i++ {
		if thumbPos+i < majorLength {
			setContent(majorPos+thumbPos+i, ' ', thumbStyle)
		}
	}
}

func (scrollbarTrack *ScrollbarTrack) MouseHandler() func(action nuview.MouseAction, event *tcell.EventMouse, setFocus func(p nuview.Primitive)) (consumed bool, capture nuview.Primitive) {
	return scrollbarTrack.WrapMouseHandler(func(action nuview.MouseAction, event *tcell.EventMouse, setFocus func(p nuview.Primitive)) (consumed bool, capture nuview.Primitive) {
		rx, ry, width, height := scrollbarTrack.GetInnerRect()
		absEventX, absEventY := event.Position()
		eventX := absEventX - rx
		eventY := absEventY - ry

		majorLength := height
		eventMinorAxis := eventX
		eventMajorAxis := eventY
		if scrollbarTrack.isHorizontal {
			majorLength = width
			eventMinorAxis = eventY
			eventMajorAxis = eventX
		}

		if eventMinorAxis < 0 || eventMinorAxis >= scrollbarTrack.width || eventMajorAxis < 0 || eventMajorAxis >= majorLength {
			return false, nil // Click outside the scrollbar
		}

		if scrollbarTrack.thumbSize >= scrollbarTrack.max {
			return false, nil
		}

		if action == nuview.MouseLeftDown || (action == nuview.MouseMove && event.Buttons() == tcell.Button1) {
			// Calculate the new position based on the click
			// Assuming the scrollbar is vertical, we calculate the position based on the y coordinate
			newPosition := eventMajorAxis*scrollbarTrack.max/majorLength - scrollbarTrack.thumbSize/2
			if newPosition < 0 {
				newPosition = 0
			} else if newPosition > scrollbarTrack.max-scrollbarTrack.thumbSize {
				newPosition = scrollbarTrack.max - scrollbarTrack.thumbSize
			}
			scrollbarTrack.position = newPosition
			if scrollbarTrack.changedFunc != nil {
				scrollbarTrack.changedFunc(newPosition)
			}
			return true, nil // Consumed the event
		}

		// Handle scroll events
		if action == nuview.MouseScrollUp {
			pos := scrollbarTrack.Position() - max(scrollbarTrack.ThumbSize()/2, 1)
			scrollbarTrack.SetPosition(pos)
			if scrollbarTrack.changedFunc != nil {
				scrollbarTrack.changedFunc(pos)
			}
			return true, nil // Consumed the event
		}
		if action == nuview.MouseScrollDown {
			pos := scrollbarTrack.Position() + max(scrollbarTrack.ThumbSize()/2, 1)
			scrollbarTrack.SetPosition(pos)
			if scrollbarTrack.changedFunc != nil {
				scrollbarTrack.changedFunc(pos)
			}
			return true, nil // Consumed the event
		}

		return false, nil // Not consumed
	})
}

// SetPosition sets the position of the scrollbar thumb.
func (scrollbarTrack *ScrollbarTrack) SetPosition(position int) {
	if position < 0 {
		scrollbarTrack.position = 0
	} else if position > scrollbarTrack.max-scrollbarTrack.thumbSize {
		scrollbarTrack.position = scrollbarTrack.max - scrollbarTrack.thumbSize
	} else {
		scrollbarTrack.position = position
	}
}

// GetPosition returns the current position of the scrollbar thumb.
func (scrollbarTrack *ScrollbarTrack) Position() int {
	return scrollbarTrack.position
}

// SetThumbSize sets the size of the scrollbar thumb.
func (scrollbarTrack *ScrollbarTrack) SetThumbSize(size int) {
	if size < 1 {
		scrollbarTrack.thumbSize = 1
	} else {
		scrollbarTrack.thumbSize = size
	}
}

// GetThumbSize returns the size of the scrollbar thumb.
func (scrollbarTrack *ScrollbarTrack) ThumbSize() int {
	return scrollbarTrack.thumbSize
}

// SetMax sets the maximum value of the scrollbar.
func (scrollbarTrack *ScrollbarTrack) SetMax(max int) {
	if max > 0 {
		scrollbarTrack.max = max
	}
}

// GetMax returns the maximum value of the scrollbar.
func (scrollbarTrack *ScrollbarTrack) Max() int {
	return scrollbarTrack.max
}

// SetWidth sets the width of the scrollbar.
func (scrollbarTrack *ScrollbarTrack) SetWidth(width int) {
	if width > 0 {
		scrollbarTrack.width = width
	}
}

// GetWidth returns the width of the scrollbar.
func (scrollbarTrack *ScrollbarTrack) Width() int {
	return scrollbarTrack.width
}

// SetTrackColor sets the color of the scrollbar track.
func (scrollbarTrack *ScrollbarTrack) SetTrackColor(color tcell.Color) {
	scrollbarTrack.trackColor = color
}

// GetTrackColor returns the color of the scrollbar track.
func (scrollbarTrack *ScrollbarTrack) TrackColor() tcell.Color {
	return scrollbarTrack.trackColor
}

// SetThumbColor sets the color of the scrollbar thumb.
func (scrollbarTrack *ScrollbarTrack) SetThumbColor(color tcell.Color) {
	scrollbarTrack.thumbColor = color
}

// GetThumbColor returns the color of the scrollbar thumb.
func (scrollbarTrack *ScrollbarTrack) ThumbColor() tcell.Color {
	return scrollbarTrack.thumbColor
}

func (scrollbarTrack *ScrollbarTrack) SetChangedFunc(changedFunc func(position int)) {
	scrollbarTrack.changedFunc = changedFunc
}
