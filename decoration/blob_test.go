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

package decoration_test

import (
	"testing"

	"github.com/xogas/cowsay-go/decoration"
)

func TestBlob(t *testing.T) {
	tests := []struct {
		name    string
		msg     string
		wantMsg string
	}{
		{
			name:    "simple letters ans space",
			msg:     "x y\n",
			wantMsg: "\x1b[1mx\x1b[0m \x1b[1my\x1b[0m\n",
		},
		{
			name:    "only whitespace preserved",
			msg:     " \t\n",
			wantMsg: " \t\n",
		},
		{
			name:    "unicode full-width runes",
			msg:     "你好\n",
			wantMsg: "\x1b[1m你\x1b[0m\x1b[1m好\x1b[0m\n",
		},
		{
			name:    "empty input",
			msg:     "",
			wantMsg: "",
		},
	}

	for _, tc := range tests {
		got := decoration.Blob([]byte(tc.msg))
		if tc.wantMsg != string(got) {
			t.Fatalf("Test %q failed: want %q, got %q", tc.name, tc.wantMsg, string(got))
		}
	}
}
