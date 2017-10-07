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
	"github.com/mcree/cmdiff/session"
	"net/url"
	"github.com/satori/go.uuid"
	"time"
	"gopkg.in/fatih/set.v0"
	log "github.com/sirupsen/logrus"
	"encoding/json"
)

// record representing attribute and value differences between two SessionItems
type SessionItemDiff struct {
	Id    url.URL                `json:"id"`
	Lost  map[string]interface{} `json:"lost"`
	Found map[string]interface{} `json:"found"`
	Equal map[string]interface{} `json:"equal"`
	Diff map[string]struct {
		Old interface{} `json:"old"`
		New interface{} `json:"new"`
	} `json:"diff"`
	Ignore []string `json:"ignore"`
}

// metadata of SessionDiff records - mainly for indexing purposes
type SessionDiffMeta struct {
	UUID       uuid.UUID `json:"uuid"`
	OldSession uuid.UUID `json:"oldSessionUuid"`
	NewSession uuid.UUID `json:"newSessionUuid"`
	Time       time.Time `json:"time"`
	Changes    int       `json:"changes"`
	ItemsLost  int       `json:"itemsLost"`
	ItemsFound int       `json:"itemsFound"`
	ItemsEqual int       `json:"itemsEqual"`
	ItemsDiff  int       `json:"itemsDiff"`
}

// record representing difference between two Sessions
type SessionDiff struct {
	Meta  SessionDiffMeta        `json:"meta"`
	Lost  []*session.SessionItem `json:"lost"`
	Found []*session.SessionItem `json:"found"`
	Equal []*session.SessionItem `json:"equal"`
	Diff  []SessionItemDiff      `json:"diff"`
}

func (diff *SessionDiff) String() (string) {
	json, _ := json.MarshalIndent(diff, "", "  ")
	return string(json)
}

// extract set of urls from a session
func urlSet(sess *session.Session) (*set.Set) {
	urls := set.New()
	for _, i := range sess.Items {
		urls.Add(i.Id)
	}
	return urls
}

// get items from a session for item urls that are found in the set
func getItems(sess *session.Session, set *set.Set) ([]*session.SessionItem) {
	var res []*session.SessionItem
	for _, i := range sess.Items {
		if set.Has(i.Id) {
			res = append(res, i)
		}
	}
	return res
}

// get item from a session with a given item URL
func getItem(sess *session.Session, url url.URL) (*session.SessionItem) {
	for _, i := range sess.Items {
		if i.Id == url {
			return i
		}
	}
	return nil
}

func diffItems(oldItem *session.SessionItem, newItem *session.SessionItem) (SessionItemDiff) {
	var diff SessionItemDiff
	diff.Lost = make(map[string]interface{})
	diff.Found = make(map[string]interface{})
	diff.Equal = make(map[string]interface{})
	diff.Diff = make(map[string]struct {
		Old interface{} `json:"old"`
		New interface{} `json:"new"`
	})

	if oldItem.Id != newItem.Id {
		log.Fatal("diffItems ", oldItem.Id.String(), "!=", newItem.Id.String())
	}
	diff.Id = newItem.Id

	oldAttrs := set.New()
	for k, _ := range oldItem.Attrs {
		oldAttrs.Add(k)
	}

	newAttrs := set.New()
	for k, _ := range newItem.Attrs {
		newAttrs.Add(k)
	}

	set.Difference(oldAttrs, newAttrs).Each(func(lost interface{}) (bool) {
		key := lost.(string)
		diff.Lost[key] = oldItem.Attrs[key]
		return true
	})

	set.Difference(newAttrs, oldAttrs).Each(func(found interface{}) (bool) {
		key := found.(string)
		diff.Found[key] = newItem.Attrs[key]
		return true
	})

	set.Intersection(newAttrs, oldAttrs).Each(func(found interface{}) (bool) {
		key := found.(string)
		oldVal := oldItem.Attrs[key]
		newVal := newItem.Attrs[key]
		if oldVal == newVal {
			diff.Equal[key] = newVal
		} else {
			diff.Diff[key] = struct {
				Old interface{} `json:"old"`
				New interface{} `json:"new"`
			}{
				oldVal, newVal,
			}
		}
		return true
	})

	return diff
}

// calculate difference between two Sessions
func NewSessionDiff(oldSession *session.Session, newSession *session.Session) (*SessionDiff) {
	diff := new(SessionDiff)
	diff.Meta.UUID = uuid.NewV1()
	diff.Meta.Time = time.Now()
	oldUrls := urlSet(oldSession)
	newUrls := urlSet(newSession)

	diff.Lost = getItems(oldSession, set.Difference(oldUrls, newUrls).(*set.Set))
	diff.Meta.ItemsLost = len(diff.Lost)
	diff.Found = getItems(newSession, set.Difference(newUrls, oldUrls).(*set.Set))
	diff.Meta.ItemsFound = len(diff.Found)

	commonUrls := set.Intersection(oldUrls, newUrls).(*set.Set)
	commonUrls.Each(func(i interface{}) (bool) {
		oldItem := getItem(oldSession, i.(url.URL))
		newItem := getItem(newSession, i.(url.URL))
		idiff := diffItems(oldItem, newItem)
		if len(idiff.Diff) == 0 {
			diff.Equal = append(diff.Equal, newItem)
		} else {
			diff.Diff = append(diff.Diff, idiff)
		}
		return true
	})
	diff.Meta.ItemsEqual = len(diff.Equal)
	diff.Meta.ItemsDiff = len(diff.Diff)
	diff.Meta.Changes = diff.Meta.ItemsLost + diff.Meta.ItemsFound + diff.Meta.ItemsDiff

	return diff
}
