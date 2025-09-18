// MIT License
//
// Copyright (c) 2025 xogas <57179186+xogas@users.noreply.github.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package assets

import (
	"embed"
	"sort"
	"strings"
)

//go:embed cows/*
var cowsFS embed.FS

// Asset loads and returns the asset for the given name.
func Asset(path string) ([]byte, error) {
	return cowsFS.ReadFile(path)
}

// AssetNames returns the names of all assets.
func AssetNames() []string {
	entries, err := cowsFS.ReadDir("cows")
	if err != nil {
		panic(err)
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			name := strings.TrimSuffix(entry.Name(), ".cow")
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names
}

var cowsInBinary = AssetNames()

// CowInBinary returns the names of all cow assets in binary format.
func CowInBinary() []string {
	return cowsInBinary
}
