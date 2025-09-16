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

package main

import (
	"bytes"
	"flag"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/xogas/cowsay-go/appversion"
	"github.com/xogas/cowsay-go/assets"
	"github.com/xogas/cowsay-go/cowsay"
	"github.com/xogas/cowsay-go/decoration"
)

// Options struct for parse command line arguments
type Options struct {
	CowFilePath string
	CowName     string
	Random      bool
	Rainbow     bool
	Blob        bool
	Wrap        int
	ListCows    bool
	Version     bool
	Help        bool
}

func (opts *Options) Usage() []byte {
	buf := new(bytes.Buffer)

	fmt.Fprintf(buf, "Usage: cowsay [options] [message]\n\n")
	fmt.Fprintf(buf, "Options:\n")

	w := tabwriter.NewWriter(buf, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintf(w, "  --filepath\tstring\tFolder where cow files are stored\n")
	_, _ = fmt.Fprintf(w, "  --cow\tstring\tName of the cow\n")
	_, _ = fmt.Fprintf(w, "  --random\t \tUse a random cow\n")
	_, _ = fmt.Fprintf(w, "  --rainbow\t \tRainbow output\n")
	_, _ = fmt.Fprintf(w, "  --blob\t \tBlob output\n")
	_, _ = fmt.Fprintf(w, "  --wrap\tint\tWrap text at this column\n")
	_, _ = fmt.Fprintf(w, "  --list\t \tList all available cows\n")
	_, _ = fmt.Fprintf(w, "  --version\t \tShow version information\n")
	_, _ = fmt.Fprintf(w, "  --help\t \tShow help message\n")

	_ = w.Flush()

	return buf.Bytes()
}

var opts Options

func init() {
	flag.StringVar(&opts.CowFilePath, "filepath", "", "Folder where cow files are stored")
	flag.StringVar(&opts.CowName, "cow", "default", "Name of the cow")
	flag.BoolVar(&opts.Random, "random", false, "Use a random cow")
	flag.BoolVar(&opts.Rainbow, "rainbow", false, "Rainbow output")
	flag.BoolVar(&opts.Blob, "blob", false, "Blob output")
	flag.IntVar(&opts.Wrap, "wrap", 40, "Wrap text at this column")
	flag.BoolVar(&opts.ListCows, "list", false, "List all available cows")
	flag.BoolVar(&opts.Version, "version", false, "Show version information")
	flag.BoolVar(&opts.Help, "help", false, "Show help message")

	flag.Parse()
}

func main() {
	if opts.Help {
		_, _ = os.Stdout.Write(opts.Usage())
		os.Exit(0)
	}

	if opts.Version {
		_, _ = os.Stdout.Write(appversion.Info())
		os.Exit(0)
	}

	if opts.ListCows {
		out, err := listCows()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		_, _ = os.Stdout.Write(out)
		os.Exit(0)
	}

	msg := strings.Join(flag.Args(), " ")
	if strings.TrimSpace(msg) == "" {
		msg = "Hello, World!"
	}

	os.Exit(run(msg))
}

func listCows() ([]byte, error) {
	var buf bytes.Buffer
	if opts.CowFilePath == "" {
		// list embedded
		names := assets.CowInBinary()
		for _, n := range names {
			fmt.Fprintln(&buf, n)
		}
		return buf.Bytes(), nil
	}
	// user path: could be file or dir
	info, err := os.Stat(opts.CowFilePath)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		entries, err := os.ReadDir(opts.CowFilePath)
		if err != nil {
			return nil, err
		}
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			if strings.HasSuffix(strings.ToLower(e.Name()), ".cow") {
				fmt.Fprintln(&buf, strings.TrimSuffix(e.Name(), ".cow"))
			}
		}
		return buf.Bytes(), nil
	}
	// file
	fmt.Fprintln(&buf, strings.TrimSuffix(filepath.Base(opts.CowFilePath), ".cow"))
	return buf.Bytes(), nil
}

func run(msg string) int {
	location, basePath := determineLocationAndBase()

	if opts.Random {
		cowNames, err := cowsay.AvailableCows(basePath, location)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return 1
		}
		cowName := cowNames[rand.Int32N(int32(len(cowNames)))]
		opts.CowName = cowName
	}

	out, err := renderCow(msg, location, basePath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	if opts.Rainbow {
		out = decoration.Rainbow(out)
	}

	if opts.Blob {
		out = decoration.Blob(out)
	}

	_, _ = os.Stdout.Write(out)
	return 0
}

func determineLocationAndBase() (cowsay.LocationType, string) {
	var location cowsay.LocationType
	basePath := opts.CowFilePath
	if basePath == "" {
		location = cowsay.InBinary
		basePath = "cows"
	} else {
		// if it's a single .cow file -> directory mode but point to file
		lower := strings.ToLower(basePath)
		if strings.HasSuffix(lower, ".cow") {
			location = cowsay.InDirectory
		} else {
			// check FS
			if info, err := os.Stat(basePath); err == nil && info.IsDir() {
				location = cowsay.InDirectory
			} else {
				// treat as directory even if not exists (error will surface on read)
				location = cowsay.InDirectory
			}
		}
	}
	return location, basePath
}

func renderCow(msg string, location cowsay.LocationType, basePath string) ([]byte, error) {
	cowName := opts.CowName
	// if user provided a specific .cow file and left default name, use file basename
	if strings.HasSuffix(strings.ToLower(basePath), ".cow") && cowName == "default" {
		cowName = strings.TrimSuffix(filepath.Base(basePath), ".cow")
	}

	c := cowsay.NewCow(cowName, basePath, location)
	c.Wrap = opts.Wrap

	out, err := c.Render(msg)
	if err != nil {
		return nil, err
	}
	return out, nil
}
