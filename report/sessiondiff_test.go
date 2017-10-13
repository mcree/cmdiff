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
	"github.com/stretchr/testify/assert"
	"net/url"
	"github.com/spf13/viper"
//	"encoding/json"
)

func newSessionItem(id string, attrs map[string]interface{}) (*session.SessionItem)  {
	item := new(session.SessionItem)
	u, _ := url.Parse(id)
	item.Id = *u
	item.Attrs = attrs
	return item
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func TestNewSessionDiff(t *testing.T) {
	viper.Reset()
	viper.Set("diff.ignore.ignoreLost", []string {"/diff"})
	viper.Set("diff.ignore.ignoreFound", []string {"/diff"})
	//t.Log(viper.AllSettings())

	old := new(session.Session)
	old.AddItem(newSessionItem("file:///oldLost", map[string]interface{} {
		"a1": "v1",
	}))
	old.AddItem(newSessionItem("file:///equal", map[string]interface{} {
		"a1": "v1",
	}))
	old.AddItem(newSessionItem("file:///diff", map[string]interface{} {
		"equal": "val",
		"lost": "val",
		"diff": "oldval",
		"ignoreLost": "val",
	}))
	n := new(session.Session)
	n.AddItem(newSessionItem("file:///newFound", map[string]interface{} {
		"a1": "v1",
	}))
	n.AddItem(newSessionItem("file:///equal", map[string]interface{} {
		"a1": "v1",
	}))
	n.AddItem(newSessionItem("file:///diff", map[string]interface{} {
		"equal": "val",
		"found": "val",
		"diff": "newval",
		"ignoreFound": "val",
	}))
	diff := NewSessionDiff(old, n)
	//t.Log(diff)
	//b,_ := json.MarshalIndent(diff.Diff[0],"", "  ")
	//t.Log(string(b))

	assert.Equal(t, 3, diff.Meta.Changes)
	assert.Equal(t, 1, diff.Meta.ItemsFound)
	assert.Equal(t, 1, diff.Meta.ItemsLost)
	assert.Equal(t, 1, diff.Meta.ItemsEqual)
	assert.Equal(t, 1, diff.Meta.ItemsDiff)

	assert.Equal(t, "file:///equal", diff.Equal[0].Id.String())
	assert.Equal(t, "file:///oldLost", diff.Lost[0].Id.String())
	assert.Equal(t, "file:///newFound", diff.Found[0].Id.String())
	assert.Equal(t, "file:///diff", diff.Diff[0].Id.String())
	assert.Equal(t, "val", diff.Diff[0].Equal["equal"])
	assert.Equal(t, 1, len(diff.Diff[0].Equal))
	assert.Equal(t, "val", diff.Diff[0].Found["found"])
	assert.Equal(t, 1, len(diff.Diff[0].Found))
	assert.Equal(t, "val", diff.Diff[0].Lost["lost"])
	assert.Equal(t, 1, len(diff.Diff[0].Lost))
	assert.Equal(t, "oldval", diff.Diff[0].Diff["diff"].Old)
	assert.Equal(t, 1, len(diff.Diff[0].Diff))
	assert.Equal(t, "newval", diff.Diff[0].Diff["diff"].New)
	assert.Equal(t, 2, len(diff.Diff[0].Ignore))
	assert.Equal(t, true, stringInSlice("ignoreLost", diff.Diff[0].Ignore))
	assert.Equal(t, true, stringInSlice("ignoreFound", diff.Diff[0].Ignore))

}