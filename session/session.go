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
package session

import (
	"net/url"
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/satori/go.uuid"
	"time"
)

type SessionItem struct {
	Id    url.URL                `json:"id"`
	Attrs map[string]interface{} `json:"attrs"`
}

func NewSessionItem(id url.URL) (*SessionItem) {
	si := new(SessionItem)
	si.Id = id
	si.Attrs = make(map[string]interface{})
	return si
}

func NewSessionItemMap(id url.URL, attrs map[string]interface{}) (*SessionItem) {
	si := new(SessionItem)
	si.Id = id
	si.Attrs = attrs
	return si
}

func (item *SessionItem) SetAttr(key string, value interface{}) {
	item.Attrs[key] = value
}

func (item *SessionItem) String() string {
	var res bytes.Buffer
	res.WriteString(item.Id.String())
	res.WriteString(" {")
	for k, v := range item.Attrs {
		res.WriteString(fmt.Sprint(k, ":'", v, "' "))
	}
	res.WriteString("}")
	return res.String()
}

func (item *SessionItem) SetAttrMap(attrs map[string]interface{}) {
	for k, v := range attrs {
		item.Attrs[k] = v
	}
}

func (item *SessionItem) GetAttr(key string) (interface{}) {
	return item.Attrs[key]
}

func (item *SessionItem) GetAttrString(key string) (string) {
	return item.Attrs[key].(string)
}

func (item *SessionItem) GetAttrUint64(key string) (uint64) {
	return item.Attrs[key].(uint64)
}

func (item *SessionItem) HasAttr(key string) (bool) {
	_, ok := item.Attrs[key]
	return ok
}

type SessionMeta struct {
	UUID uuid.UUID `json:"uuid"`
	Time time.Time `json:"time"`
	Items int64
}

type Session struct {
	Meta SessionMeta `json:"meta"`
	Items []*SessionItem `json:"items"`
}

func NewSession() (*Session) {
	sess := new(Session)
	sess.Meta.UUID = uuid.NewV2(0x01)
	sess.Meta.Time = time.Now()
	return sess
}

func (sess *Session) AddItem(item *SessionItem) {
	sess.Items = append(sess.Items, item)
	sess.Meta.Items++
}

func (sess *Session) Name() string {
	return "Session"
}

func (sess *Session) Concurrency() int {
	return 1
}

func (sess *Session) Process(item interface{}) {
	sess.AddItem(item.(*SessionItem))
	log.Info("Session stored: ",item)
}