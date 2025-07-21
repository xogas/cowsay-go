/* MIT License
 *
 * Copyright (c) 2025 xogas <askxogas@gmail.com>
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

package cow

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/xogas/cowsay-go/assets"
)

// CowFile represents a cow file.
type CowFile struct {
	Name         string
	BasePath     string
	LocationType LocationType
}

// ReadAll reads the entire cow file.
func (cf *CowFile) ReadAll() ([]byte, error) {
	filePath := path.Join(cf.BasePath, cf.Name+".cow")

	if cf.LocationType == InBinary {
		// go embed is use "/" separator
		return assets.Asset(filePath)
	}
	return os.ReadFile(filePath)
}

// CowPath is information of cow file paths.
type CowPath struct {
	Name         string
	CowFileNames []string
	LocationType LocationType
}

// Lookup looks up a cow file by name.
func (cp *CowPath) Lookup(name string) (*CowFile, bool) {
	for _, cowFileName := range cp.CowFileNames {
		if cowFileName == name {
			return &CowFile{
				Name:         cowFileName,
				BasePath:     cp.Name,
				LocationType: cp.LocationType,
			}, true
		}
	}
	return nil, false
}

// Cow struct
type Cow struct {
	CowFile         *CowFile
	BallonWidth     int
	DisableWordWrap bool
}

// New returns a new Cow instance.
func New(cowFiles ...*CowFile) (*Cow, error) {
	var cowFile *CowFile
	if len(cowFiles) > 0 && cowFiles[0] != nil {
		cowFile = cowFiles[0]
	} else {
		cowFile = &CowFile{
			Name:         "default",
			BasePath:     "cows",
			LocationType: InBinary,
		}
	}

	cow := &Cow{
		CowFile:     cowFile,
		BallonWidth: 40,
	}

	return cow, nil
}

// Say makes the cow say a phrase.
func (c *Cow) SayTo(w io.Writer, phrase string) error {

	// say balloon
	if err := c.writeBalloon(w, phrase); err != nil {
		return fmt.Errorf("failed to say balloon: %w", err)
	}

	// say cow
	if err := c.writeCow(w); err != nil {
		return fmt.Errorf("failed to say cow: %w", err)
	}

	return nil
}

// sayBalloon
func (c *Cow) writeBalloon(w io.Writer, phrase string) error {

	// simple word wrap
	// lines := c.getLines(phrase)
	var lines []*line

	// compute max width
	maxWidth := 0
	for _, line := range lines {
		if line.runeWidth > maxWidth {
			maxWidth = line.runeWidth
		}
	}

	// top border
	if _, err := io.WriteString(w, " "+strings.Repeat("-", maxWidth+2)+"\n"); err != nil {
		return err
	}

	// content with borders
	if len(lines) == 1 {
		// single line: < ... >
		line := lines[0]
		if _, err := io.WriteString(w, "< "+line.text+" >\n"); err != nil {
			return err
		}
	} else {
		// multiple lines: / \ , | | , \ /
		first := lines[0]
		if _, err := io.WriteString(w, "/ "+first.text+" "+strings.Repeat(" ", maxWidth-first.runeWidth)+" \\\n"); err != nil {
			return err
		}
		for _, line := range lines[1 : len(lines)-1] {
			if _, err := io.WriteString(w, "| "+line.text+" "+strings.Repeat(" ", maxWidth-line.runeWidth)+" |\n"); err != nil {
				return err
			}
		}
		last := lines[len(lines)-1]
		if _, err := io.WriteString(w, "\\ "+last.text+" "+strings.Repeat(" ", maxWidth-last.runeWidth)+" /\n"); err != nil {
			return err
		}
	}

	// bottom border
	if _, err := io.WriteString(w, " "+strings.Repeat("-", maxWidth+2)+"\n"); err != nil {
		return err
	}

	return nil
}

func (c *Cow) writeCow(w io.Writer) error {
	// read cow file
	if c.CowFile == nil {
		return fmt.Errorf("no cow file selected")
	}
	data, err := c.CowFile.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read cow file: %w", err)
	}

	// extract art between "<<EOC;" and "\nEOC" when present
	startDelim := []byte("<<EOC;")
	endDelim := []byte("\nEOC")
	art := data
	if pos := bytes.Index(data, startDelim); pos >= 0 {
		start := pos + len(startDelim)
		if start < len(data) && data[start] == '\n' {
			start++
		}
		if endRel := bytes.Index(data[start:], endDelim); endRel >= 0 {
			art = data[start : start+endRel]
		} else {
			art = data[start:]
		}
	}

	// trim trailing newlines and ensure we have art
	art = bytes.TrimRight(art, "\n")
	if len(art) == 0 {
		return fmt.Errorf("invalid cow file: no art found")
	}

	out := make([]byte, 0, len(art)+2)
	out = append(out, '\n')
	out = append(out, art...)
	out = append(out, '\n')
	if _, err := w.Write(out); err != nil {
		return err
	}

	return nil
}

func (c *Cow) maxLineWidth(lines []*line) int {
	maxWidth := 0
	for _, line := range lines {
		if line.runeWidth > maxWidth {
			maxWidth = line.runeWidth
		}
		if !c.DisableWordWrap && maxWidth > c.BallonWidth {
			maxWidth = c.BallonWidth
		}
	}
	return maxWidth
}

// func (c *Cow) getLines(phrase string) []*line {
// 	text := c.canonicalizePhrase(phrase)
// 	lineTexts := strings.Split(text, "\n")
// 	lines := make([]*line, 0, len(lineTexts))
// 	for _, lineText := range lineTexts {
// 		lines = append(lines, &line{
// 			text:      lineText,
// 			runeWidth: runewidth.StringWidth(lineText),
// 		})
// 	}
// 	return lines
// }

// func (c *Cow) canonicalizePhrase(phrase string) string {
// 	// Replace tab to 8 spaces
// 	phrase = strings.Replace(phrase, "\t", "       ", -1)

// 	if c.DisableWordWrap {
// 		return phrase
// 	}
// 	width := c.BallonWidth
// 	return wordwrap.WrapString(phrase, uint(width))
// }

func simpleWordWrap(s string, width int) []string {
	if width <= 0 {
		return []string{s}
	}

	words := strings.Fields(s)
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	var currentLine strings.Builder

	for _, word := range words {
		// if current line empty, put word directly
		if currentLine.Len() == 0 {
			// if single word is longer than width, still place it (could be further split if desired)
			currentLine.WriteString(word)
			continue
		}

		// if fits with a preceding space
		if currentLine.Len()+1+len(word) <= width {
			currentLine.WriteByte(' ')
			currentLine.WriteString(word)
			continue
		}

		// doesn't fit -> push current line and start new one with word
		lines = append(lines, currentLine.String())
		currentLine.Reset()
		currentLine.WriteString(word)
	}

	// push remaining
	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return lines
}
