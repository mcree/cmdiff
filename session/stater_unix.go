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
//
// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package session

import (
	log "github.com/inconshreveable/log15"
	"golang.org/x/sys/unix"
	"time"
	"os/user"
	"strconv"
)

type Stater struct {
}

func NewStater() (*Stater) {
	stater := new(Stater)
	return stater
}

func (stater *Stater) Name() string {
	return "Stater"
}

func (stater *Stater) Concurrency() int {
	return 10
}

func (stater *Stater) Process(item interface{}) {
	i := item.(*SessionItem)

	path := i.Id.Path

	var stat unix.Stat_t

	err := unix.Lstat(path, &stat)
	if err != nil {
		i.SetAttr("stater.error", err)
		log.Info("error open", log.Ctx{"err": err, "item": i})
		return
	}

	username, _ := user.LookupId(strconv.Itoa(int(stat.Uid)))
	groupname, _ := user.LookupGroupId(strconv.Itoa(int(stat.Gid)))

	i.SetAttrMap(map[string]interface{} {
		"stat.atime": time.Unix(int64(stat.Atim.Sec), int64(stat.Atim.Nsec)),
		"stat.ctime": time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec)),
		//"stat.mtime": time.Unix(int64(stat.Mtim.Sec), int64(stat.Mtim.Nsec)), // also added by filewalker
		"stat.uid": stat.Uid,
		"stat.gid": stat.Gid,
		"stat.user": username.Username,
		"stat.group": groupname.Name,
	})

}
