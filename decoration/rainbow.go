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

package decoration

import (
	"bytes"
	"fmt"
	"math"
	"unicode"
	"unicode/utf8"
)

const (
	freq       = 0.35
	redPhase   = 0
	greenPhase = math.Pi * 2 / 3
	bluePhase  = math.Pi * 4 / 3
	step       = 0.9
	offset     = 3.5
)

// rgb returns 1...255 r,g,b for position p.
func rgb(i float64) (red, green, blue int64) {
	red = int64(math.Sin(freq*i+redPhase)*127 + 128)
	green = int64(math.Sin(freq*i+greenPhase)*127 + 128)
	blue = int64(math.Sin(freq*i+bluePhase)*127 + 128)
	return
}

// Rainbow applies rainbow colors to the input text.
func Rainbow(input []byte) []byte {
	var buf bytes.Buffer
	lineIndex := 0
	pos := float64(lineIndex)*offset + 1

	for len(input) > 0 {
		rn, size := utf8.DecodeRune(input)
		if rn == '\n' {
			lineIndex++
			pos = float64(lineIndex)*offset + 1
			buf.WriteRune(rn)
			input = input[size:]
			continue
		}
		if unicode.IsSpace(rn) {
			buf.WriteRune(rn)
			input = input[size:]
			pos += step
			continue
		}

		r, g, b := rgb(pos)

		fmt.Fprintf(&buf, "\x1b[38;2;%d;%d;%dm%c\x1b[0m", r, g, b, rn)

		pos += step
		input = input[size:]
	}

	return buf.Bytes()
}
