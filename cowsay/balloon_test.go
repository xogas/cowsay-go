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
	"testing"
)

func TestBuildBalloon(t *testing.T) {
	tests := []struct {
		name    string
		msg     string
		wrap    int
		wantMsg string
	}{
		{
			name:    "short message",
			msg:     "hi",
			wrap:    40,
			wantMsg: " ____\n< hi >\n ----",
		},
		{
			name:    "long message",
			msg:     "aa bb cc dd ee",
			wrap:    4,
			wantMsg: " ____\n/ aa \\\n| bb  |\n| cc  |\n| dd  |\n\\ ee /\n ----",
		},
		{
			name:    "full-width unicode handled correctly",
			msg:     "你好",
			wrap:    40,
			wantMsg: " ______\n< 你好 >\n ------",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := buildBalloon(tc.msg, tc.wrap)
			if tc.wantMsg != string(got) {
				t.Fatalf("expected message balloon, got: %q", string(got))
			}
		})
	}
}
