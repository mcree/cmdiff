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

package db

import (
	"github.com/peterbourgon/diskv"
	"github.com/spf13/viper"
	"github.com/mcree/cmdiff/session"
	"encoding/json"
)

type dbStruct struct {
	diskv *diskv.Diskv
}

// singleton instance
var db *dbStruct

// set up singleton instance if needed
func thawDB() {
	if db == nil {
		db = new(dbStruct)
		db.diskv = diskv.New(diskv.Options{
			 BasePath: viper.GetString("db.basepath"),
		})
	}
}

func ReadSession(uuid string) (*session.Session) {
	thawDB()
	data, err := db.diskv.Read(uuid)
	if err != nil {
		sess := new(session.Session)
		json.Unmarshal(data, sess)
		return sess
	} else {
		return nil
	}
}