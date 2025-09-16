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

package cowsay_test

import (
	"os"
	"reflect"
	"sort"
	"testing"

	"github.com/xogas/cowsay-go/assets"
	"github.com/xogas/cowsay-go/cowsay"
)

func TestCow(t *testing.T) {
	wantMsg, err := os.ReadFile("./testdata/default.cow")
	if err != nil {
		t.Fatalf("failed to read wantMsg: %v", err)
	}

	tests := []struct {
		name    string
		cow     *cowsay.Cow
		msg     string
		wantMsg []byte
		hasErr  bool
	}{
		{
			name:    "default cow in binary",
			cow:     cowsay.NewCow("default", "", cowsay.InBinary),
			msg:     "Hello!",
			wantMsg: wantMsg,
			hasErr:  false,
		},
		{
			name:    "default cow in directory",
			cow:     cowsay.NewCow("test", "./testdata/testdir", cowsay.InDirectory),
			msg:     "Hello!",
			wantMsg: wantMsg,
			hasErr:  false,
		},
	}

	for _, tc := range tests {
		got, err := tc.cow.Render(tc.msg)
		if tc.hasErr {
			if err == nil {
				t.Fatalf("%s: expected error but got none", tc.name)
			}
			continue
		}
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", tc.name, err)
		}
		if !reflect.DeepEqual(got, tc.wantMsg) {
			t.Fatalf("expected %v cows, got %v cows", tc.wantMsg, got)
		}
	}
}

func TestAvailableCows(t *testing.T) {
	tests := []struct {
		name     string
		basePath string
		location cowsay.LocationType
		want     []string
		hasErr   bool
	}{
		{
			name:     "builtin",
			basePath: "",
			location: cowsay.InBinary,
			want:     assets.CowInBinary(),
			hasErr:   false,
		},
		{
			name:     "form directory",
			basePath: "./testdata/testdir",
			location: cowsay.InDirectory,
			want:     []string{"test"},
			hasErr:   false,
		},
		{
			name:     "no such directory",
			basePath: "/no/such/dir",
			location: cowsay.InDirectory,
			want:     []string{},
			hasErr:   true,
		},
	}

	for _, tc := range tests {
		got, err := cowsay.AvailableCows(tc.basePath, tc.location)
		if tc.hasErr {
			if err == nil {
				t.Fatalf("%s: expected error but got none", tc.name)
			}
			continue
		}
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", tc.name, err)
		}
		sort.Strings(got)
		sort.Strings(tc.want)
		if !reflect.DeepEqual(got, tc.want) {
			t.Fatalf("expected %v cows, got %v cows", tc.want, got)
		}
	}

}
