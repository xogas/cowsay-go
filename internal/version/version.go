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

// Package version provides information about the version of the software.
package version

import (
	"bytes"
	"fmt"
	"runtime"
	"text/tabwriter"
)

var (
	AppVersion = "--"
	GitCommit  = "--"
	BuildTime  = "--"
	GoVersion  = runtime.Version()
)

// Version returns the version information of the software.
func Version() string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 1, ' ', 0)

	fmt.Fprintf(w, "AppVersion:\t%s\n", AppVersion)
	fmt.Fprintf(w, "GitCommit:\t%s\n", GitCommit)
	fmt.Fprintf(w, "BuildTime:\t%s\n", BuildTime)
	fmt.Fprintf(w, "GoVersion:\t%s\n", GoVersion)
	if err := w.Flush(); err != nil {
		return fmt.Sprintf("failed to format version: %v", err)
	}

	return buf.String()
}
