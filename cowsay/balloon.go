/* MIT License
 *
 * Copyright (c) 2025 xogas <57179186+xogas@users.noreply.github.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package cowsay

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

// stringWidth calculates the display width of a string, considering full-width characters.
func stringWidth(s string) int {
	width := 0
	for _, r := range s {
		if unicode.Is(unicode.Han, r) || r > 0xFF {
			width += 2 // full-width
		} else {
			width += 1 // half-width
		}
	}
	return width
}

// buildBalloon wraps the message in a speech balloon.
func buildBalloon(msg string, wrap int) []byte {
	if wrap <= 0 {
		wrap = 40
	}
	words := strings.Fields(msg)
	if len(words) == 0 {
		return []byte("< >")
	}

	var lines []string
	var cur strings.Builder
	cur.WriteString(words[0])
	for _, word := range words[1:] {
		if stringWidth(cur.String())+1+stringWidth(word) <= wrap {
			cur.WriteByte(' ')
			cur.WriteString(word)
			continue
		}
		lines = append(lines, cur.String())
		cur.Reset()
		cur.WriteString(word)
	}
	lines = append(lines, cur.String())

	// compute max width
	max := 0
	for _, line := range lines {
		width := stringWidth(line)
		if width > max {
			max = width
		}
	}

	var out bytes.Buffer
	out.WriteByte(' ')
	out.Write(bytes.Repeat([]byte("_"), max+2))
	out.WriteByte('\n')

	if len(lines) == 1 {
		line := lines[0]
		padding := max - stringWidth(line)
		out.WriteString(fmt.Sprintf("< %s %s>\n", line, strings.Repeat(" ", padding)))
	} else {
		// first
		first := lines[0]
		padding := max - stringWidth(first)
		out.WriteString(fmt.Sprintf("/ %s %s\\\n", first, strings.Repeat(" ", padding)))
		// middle
		for _, line := range lines[1 : len(lines)-1] {
			padding = max - stringWidth(line)
			out.WriteString(fmt.Sprintf("| %s  %s|\n", line, strings.Repeat(" ", padding)))
		}
		// last
		last := lines[len(lines)-1]
		padding = max - stringWidth(last)
		out.WriteString(fmt.Sprintf("\\ %s %s/\n", last, strings.Repeat(" ", padding)))
	}

	out.WriteByte(' ')
	out.WriteString(strings.Repeat("-", max+2))
	return out.Bytes()
}
