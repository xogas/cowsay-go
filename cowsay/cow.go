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
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/xogas/cowsay-go/assets"
)

// LocationType represents the location of a cow asset.
type LocationType int

const (
	// InBinary indicates the CowPath binary.
	InBinary LocationType = iota
	// InDirectory indicates a directory containing cow files.
	InDirectory
)

// Cow represents a talking cow.
type Cow struct {
	Name     string
	BasePath string
	Location LocationType
	Wrap     int
}

// NewCow creates a new Cow instance.
func NewCow(name, basePath string, location LocationType) *Cow {
	if name == "" {
		name = "default"
	}
	if location == InBinary && basePath == "" {
		basePath = "cows"
	}
	return &Cow{
		Name:     name,
		BasePath: basePath,
		Location: location,
		Wrap:     40,
	}
}

// AvailableCows lists all available cow names from the specified location.
func AvailableCows(basePath string, location LocationType) ([]string, error) {
	var names []string

	if location == InBinary {
		names = assets.CowInBinary()
	} else {
		entries, err := os.ReadDir(basePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read directory %q: %w", basePath, err)
		}

		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".cow") {
				name := strings.TrimSuffix(entry.Name(), ".cow")
				names = append(names, name)
			}
		}
		sort.Strings(names)
	}

	return names, nil
}

// Render builds the speech balloon and append the cow art.
func (c *Cow) Render(msg string) ([]byte, error) {
	if strings.TrimSpace(msg) == "" {
		msg = "Hello, World!"
	}

	// balloon
	balloon := buildBalloon(msg, c.Wrap)

	// load cow data
	var data []byte
	var err error
	if c.Location == InBinary {
		assetsPath := filepath.ToSlash(filepath.Join(c.BasePath, c.Name+".cow"))
		data, err = assets.Asset(assetsPath)
		if err != nil {
			return nil, fmt.Errorf("embedded cow %q not found: %w", c.Name, err)
		}
	} else {
		// BasePath might be a file or directory
		var cowFile string
		if strings.HasSuffix(strings.ToLower(c.BasePath), ".cow") {
			cowFile = c.BasePath
		} else {
			cowFile = filepath.Join(c.BasePath, c.Name+".cow")
		}
		data, err = os.ReadFile(cowFile)
		if err != nil {
			return nil, fmt.Errorf("cow file %q not found: %w", cowFile, err)
		}
	}

	// extract art between "<<EOC" and "EOC" (fallback to whole file)
	startDelimiter := []byte("<<EOC;")
	endDelimiter := []byte("EOC")
	art := data
	if pos := bytes.Index(data, startDelimiter); pos >= 0 {
		start := pos + len(startDelimiter)
		if start < len(data) && data[start] == '\n' {
			start++
		}
		if endRelPos := bytes.Index(data[start:], endDelimiter); endRelPos >= 0 {
			art = data[start : start+endRelPos]
		} else {
			art = data[start:]
		}
	}
	art = bytes.TrimRight(art, "\n")
	if len(art) == 0 {
		return []byte{}, errors.New("invalid cow file: no art found")
	}

	var out bytes.Buffer
	out.Write(balloon)
	out.WriteByte('\n')
	out.Write(art)
	out.WriteByte('\n')
	return out.Bytes(), nil
}
