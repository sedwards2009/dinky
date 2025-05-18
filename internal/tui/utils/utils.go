package utils

import "github.com/gdamore/tcell/v2"

func DrawText(screen tcell.Screen, x int, y int, text string, style tcell.Style) {
	for _, r := range text {
		screen.SetContent(x, y, r, nil, style)
		x++
	}
}

func DrawCappedHorizontalLine(screen tcell.Screen, x int, y int, width int, borderStyle tcell.Style, middleStyle tcell.Style, left rune,
	middle rune, right rune) {

	screen.SetContent(x, y, left, nil, borderStyle)
	for i := 1; i < width-1; i++ {
		screen.SetContent(x+i, y, middle, nil, middleStyle)
	}
	screen.SetContent(x+width-1, y, right, nil, borderStyle)
}

func DrawHorizontalLine(screen tcell.Screen, x int, y int, width int, style tcell.Style, char rune) {
	for i := 0; i < width; i++ {
		screen.SetContent(x+i, y, char, nil, style)
	}
}
