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

// Package db provides disk persistence functions
package db

import (
	"github.com/peterbourgon/diskv"
	"github.com/spf13/viper"
	"github.com/mcree/cmdiff/session"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"fmt"
	"sync"
	"sort"
	"errors"
)

// database internal singleton struct
type dbStruct struct {
	diskv *diskv.Diskv
	index *Index
}

// index containing session metadata
type SessionIndex []session.SessionMeta

// database index
type Index struct {
	Sessions SessionIndex `json:"sessions"`
}

// helper for sort.Sort
func (p SessionIndex) Len() int {
	return len(p)
}

// helper for sort.Sort
func (p SessionIndex) Less(i, j int) bool {
	return p[i].Time.Before(p[j].Time)
}

// helper for sort.Sort
func (p SessionIndex) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// singleton instance
var db *dbStruct

// lock for thread safety of public functions
var lock sync.Mutex

// internal method to set up singleton instance, also initializes index if needed
func thawDB() {
	if db == nil {
		db = new(dbStruct)
		db.diskv = diskv.New(diskv.Options{
			 BasePath: viper.GetString("db.basepath"),
			 //Compression: diskv.NewGzipCompression(),
		})
		log.Debug("opened database: ", db.diskv)
		if db.diskv.Has("index") {
			data, err := db.diskv.Read("index")
			if err == nil {
				db.index = new(Index)
				err := json.Unmarshal(data, db.index)
				if err != nil {
					log.Fatal("error reading index: ", err)
				}
			} else {
				log.Fatal("error reading index: ", err)
			}
		} else {
			log.Warn("initializing empty database")
			db.index = new(Index)
			writeIndex()
		}
	}
}

// internal function for index persistence
func writeIndex() {
	data, err := json.Marshal(db.index)
	if err == nil {
		err = db.diskv.Write("index", data)
	}
	if err != nil {
		log.Fatal("error updating index: ", err)
	}
}

// get session from db by uuid
func ReadSession(uuid string) (*session.Session, error) {
	lock.Lock()
	defer lock.Unlock()
	thawDB()
	key := fmt.Sprintf("session-%s", uuid)
	data, err := db.diskv.Read(key)
	if err == nil {
		sess := new(session.Session)
		err := json.Unmarshal(data, sess)
		return sess, err
	} else {
		return nil, err
	}
}

// save session to db - also updates index with session metadata
func WriteSession(sess *session.Session) (error) {
	lock.Lock()
	defer lock.Unlock()
	thawDB()
	data, err := json.Marshal(sess)
	if err == nil {
		key := fmt.Sprintf("session-%s", sess.Meta.UUID.String())
		err := db.diskv.Write(key, data)
		if err == nil {
			db.index.Sessions = append(db.index.Sessions, sess.Meta)
			writeIndex()
		}
		return err
	} else {
		return err
	}
}

// erase session from db - also updates index by removing relevant session metadata
func EraseSession(uuid string) (error) {
	lock.Lock()
	defer lock.Unlock()
	thawDB()
	key := fmt.Sprintf("session-%s", uuid)
	err := db.diskv.Erase(key)
	if err == nil {
		//log.Debug("pre erase: ",db.index.Sessions)
		for i, s := range db.index.Sessions {
			//log.Debug("db.index.Sessions[",i,"] ", s)
			if s.UUID.String() == uuid {
				db.index.Sessions = append(db.index.Sessions[:i], db.index.Sessions[i+1:]...)
			}
		}
		//log.Debug("post erase: ",db.index.Sessions)
		writeIndex()
	}
	return err
}

// erase old database records
// works based on configuration variables:
//  db.maxSessions
//  db.maxReports
func DoHousekeeping() {
	maxSessions := viper.GetInt("db.maxSessions")
	numSessions := len(db.index.Sessions)
	if numSessions > maxSessions {
	sort.Sort(db.index.Sessions)
		for _, s := range db.index.Sessions[:numSessions-maxSessions] {
			log.Info("Housekeeping, removing session: ", s)
			EraseSession(s.UUID.String())
		}
	}
}

func NthSession(n int) (*session.Session, error) {
	if len(db.index.Sessions) < n {
		return nil, errors.New(fmt.Sprintf("%d. session requested, %d available", n, len(db.index.Sessions)))
	}
	sort.Sort(db.index.Sessions)
	uuid := db.index.Sessions[n].UUID.String()
	return ReadSession(uuid)
}

func CurrentSession () (*session.Session, error) {
	return NthSession(0)
}

func PreviousSession () (*session.Session, error) {
	return NthSession(1)
}