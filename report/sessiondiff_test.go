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

package report

import (
	"testing"
	"github.com/mcree/cmdiff/session"
	"net/url"
)

func newSessionItem(id string, attrs map[string]interface{}) (*session.SessionItem)  {
	item := new(session.SessionItem)
	url, _ := url.Parse(id)
	item.Id = *url
	item.Attrs = attrs
	return item
}

func TestNewSessionDiff(t *testing.T) {
	old := new(session.Session)
	old.AddItem(newSessionItem("file://oldLost", map[string]interface{} {
		"a1": "v1",
	}))
	old.AddItem(newSessionItem("file://equal", map[string]interface{} {
		"a1": "v1",
	}))
	old.AddItem(newSessionItem("file://diff", map[string]interface{} {
		"equal": "val",
		"lost": "val",
		"diff": "oldval",
	}))
	new := new(session.Session)
	new.AddItem(newSessionItem("file://newFound", map[string]interface{} {
		"a1": "v1",
	}))
	new.AddItem(newSessionItem("file://equal", map[string]interface{} {
		"a1": "v1",
	}))
	new.AddItem(newSessionItem("file://diff", map[string]interface{} {
		"equal": "val",
		"found": "val",
		"diff": "newval",
	}))
	diff := NewSessionDiff(old, new)
	t.Log(diff)
}