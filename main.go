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

package main

import (
	"flag"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strings"

	"github.com/xogas/cowsay-go/assets"
	"github.com/xogas/cowsay-go/internal/cow"
)

var (
	cowFilePath string
	cowName     string
	random      bool
	wrap        int
)

func init() {
	flag.StringVar(&cowFilePath, "filepath", "cows", "Folder where cow files are stored")
	flag.StringVar(&cowName, "cow", "default", "Name of the cow")
	flag.BoolVar(&random, "random", false, "Use a random cow")
	flag.IntVar(&wrap, "wrap", 40, "Wrap text at this column")
}

func main() {
	flag.Parse()

	// message from args or default
	msg := strings.Join(flag.Args(), " ")
	if strings.TrimSpace(msg) == "" {
		msg = "Hello, World!"
	}

	// decide cow file
	var cf *cow.CowFile
	const defaultCowPath = "cows"
	if random {
		names := assets.AssetNames()
		if len(names) == 0 {
			fmt.Fprintln(os.Stderr, "no embedded cows available")
			os.Exit(1)
		}
		name := names[rand.IntN(len(names))]
		cf = &cow.CowFile{Name: name, BasePath: defaultCowPath, LocationType: cow.InBinary}
	} else if cowName != "" && cowName != "default" {
		// if user left filepath as default, use embedded named cow
		if cowFilePath == defaultCowPath {
			cf = &cow.CowFile{Name: cowName, BasePath: defaultCowPath, LocationType: cow.InBinary}
		} else {
			// user provided a custom folder -> treat as user dir
			cf = &cow.CowFile{Name: cowName, BasePath: cowFilePath, LocationType: cow.InUserDir}
		}
	} else if cowFilePath != "" && cowFilePath != defaultCowPath {
		// user provided a custom path: could be a single .cow file or a directory
		if strings.HasSuffix(strings.ToLower(cowFilePath), ".cow") {
			name := strings.TrimSuffix(filepath.Base(cowFilePath), ".cow")
			base := filepath.Dir(cowFilePath)
			cf = &cow.CowFile{Name: name, BasePath: base, LocationType: cow.InUserDir}
		} else {
			// treat as directory containing cow files (user-provided)
			// use default "default" name when no explicit cowName supplied
			cf = &cow.CowFile{Name: "default", BasePath: cowFilePath, LocationType: cow.InUserDir}
		}
	} else {
		// leave cf nil to let cow.New() use its internal default (embedded default cow)
		cf = nil
	}

	// create cow (pass cf only if non-nil)
	var cInst *cow.Cow
	var err error
	if cf != nil {
		cInst, err = cow.New(cf)
	} else {
		cInst, err = cow.New()
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to create cow:", err)
		os.Exit(1)
	}

	// apply settings
	cInst.BallonWidth = wrap

	// run
	if err := cInst.SayTo(os.Stdout, msg); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
