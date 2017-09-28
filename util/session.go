package util

import "net/url"

type SessionItem struct {
	id url.URL
	attrs map[string]interface{}
}

func NewSessionItem() (*SessionItem) {
	si := new(SessionItem)
	si.attrs = make(map[string]interface{})
	return si
}

func (item *SessionItem) SetAttr(key string, value interface{}) {
	item.attrs[key] = value
}

func (item *SessionItem) GetAttr(key string) (interface{}) {
	return item.attrs[key]
}

func (item *SessionItem) GetAttrString(key string) (string) {
	return item.attrs[key].(string)
}


type Session struct {
	items []*SessionItem
}

func NewSession() (*Session) {
	return new(Session)
}

func (sess *Session) AddItem(item *SessionItem) {
	sess.items = append(sess.items, item)
}
