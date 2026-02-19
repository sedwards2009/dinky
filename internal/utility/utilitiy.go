package utility

import (
	"strings"

	"github.com/rivo/uniseg"
)

// WrapParagraph takes a slice of lines that form a logical paragraph
// and wraps them to the specified width
func WrapParagraph(lines []string, width int) []string {
	// Join all lines into a single paragraph with spaces
	paragraph := strings.Join(lines, " ")

	// Normalize whitespace: replace multiple spaces with single space
	paragraph = strings.Join(strings.Fields(paragraph), " ")

	// Use the WordWrap function from table2/strings
	wrapped := WordWrap(paragraph, width)

	return wrapped
}

// WordWrap splits a text such that each resulting line does not exceed the
// given screen width. Split points are determined using the algorithm described
// in [Unicode Standard Annex #14].
//
// [Unicode Standard Annex #14]: https://www.unicode.org/reports/tr14/
func WordWrap(text string, maxWidth int) (lines []string) {
	if maxWidth <= 0 {
		return
	}

	var currentLine strings.Builder
	var currentWidth int
	var lastBreakPos int   // String position of last break opportunity
	var lastBreakWidth int // Width at last break opportunity

	originalText := text
	gr := uniseg.NewGraphemes(text)

	for gr.Next() {
		cluster := gr.Str()
		clusterWidth := gr.Width()

		// Check if this is a mandatory line break
		if cluster == "\n" || cluster == "\r\n" || cluster == "\r" {
			lines = append(lines, strings.TrimSpace(currentLine.String()))
			currentLine.Reset()
			currentWidth = 0
			lastBreakPos = 0
			lastBreakWidth = 0
			continue
		}

		// Determine if we can break after this cluster
		canBreak := gr.LineBreak() != uniseg.LineDontBreak

		// Check if adding this cluster would exceed the width
		if currentWidth+clusterWidth > maxWidth && currentWidth > 0 {
			// Need to wrap
			if lastBreakPos > 0 {
				// We have a break opportunity, use it
				line := currentLine.String()[:lastBreakPos]
				lines = append(lines, strings.TrimSpace(line))

				// Start new line with text after the break point
				remainder := currentLine.String()[lastBreakPos:]
				currentLine.Reset()
				currentLine.WriteString(remainder)
				currentLine.WriteString(cluster)
				currentWidth = currentWidth - lastBreakWidth + clusterWidth
				lastBreakPos = 0
				lastBreakWidth = 0
			} else {
				// No break opportunity, force break before this cluster
				lines = append(lines, strings.TrimSpace(currentLine.String()))
				currentLine.Reset()
				currentLine.WriteString(cluster)
				currentWidth = clusterWidth
				lastBreakPos = 0
				lastBreakWidth = 0
			}
		} else {
			// Add cluster to current line
			currentLine.WriteString(cluster)
			currentWidth += clusterWidth

			// Remember break opportunity (typically after spaces)
			if canBreak {
				lastBreakPos = currentLine.Len()
				lastBreakWidth = currentWidth
			}
		}
	}

	// Add final line if not empty
	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	// Ensure at least one line is returned
	if len(lines) == 0 && len(originalText) == 0 {
		lines = []string{""}
	}

	return
}
