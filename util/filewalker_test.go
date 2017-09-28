// Copyright © 2017 Erno Rigo <erno@rigo.info>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.
package util

import (
	"testing"
)
//#   pattern:
//#       { term }
//#
//#   term:
//#       `*`         matches any sequence of non-separator characters
//#       `**`        matches any sequence of characters
//#       `?`         matches any single non-separator character
//#       `[` [ `!` ] { character-range } `]`
//#                   character class (must be non-empty)
//#       `{` pattern-list `}`
//#                   pattern alternatives
//#       c           matches character c (c != `*`, `**`, `?`, `\`, `[`, `{`, `}`)
//#       `\` c       matches character c
//#
//#   character-range:
//#       c           matches character c (c != `\\`, `-`, `]`)
//#       `\` c       matches character c
//#       lo `-` hi   matches character c for lo <= c <= hi
//#
//#   pattern-list:
//#       pattern { `,` pattern }
//#                   comma-separated (without spaces) patterns

var findRootTests = []struct {
	path     string // input
	expected string // expected result
}{
	{"", ""},
	{"\\/", "/"},
	{".", "."},
	{"/", "/"},
	{"./", "./"},
	{".a*a", "."},
	{"/a?a", "/"},
	{"./a[bc]a", "."},
	{"/var/log", "/var/log"},
	{"/var/log/", "/var/log/"},
	{"/var/log/a", "/var/log/a"},
	{"/var/log/a*", "/var/log"},
	{"/var/log/**/things", "/var/log"},
	{"var/log/**/things", "var/log"},
	{"var/log/aaa**aaa", "var/log"},
	{"var/log/aaa*aaa", "var/log"},
	{"var/log/aaa**aaa/things", "var/log"},
	{"var/log/aaa*aaa/things", "var/log"},
	{"var/log/aaa[abc]aaa/things", "var/log"},
	{"[!abc]", "."},
	{"/var/{a,b,c}", "/var"},
	{"/var/log\\*/thing/asdf**asdf", "/var/log*/thing"},
	{"/var/\\{a,b,c\\}\\[abc\\]\\*\\*\\?\\\\/thing", "/var/{a,b,c}[abc]**?\\/thing"},
}

func TestFindRoot(t *testing.T) {
	for _, tt := range findRootTests {
		actual := FindRoot(tt.path)
		if actual != tt.expected {
			t.Errorf("FindRoot(%d): expected %d, actual %d", tt.path, tt.expected, actual)
		}
	}
}
