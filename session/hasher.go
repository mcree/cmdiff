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
	"crypto/md5"
	"crypto/sha256"
	"io"
	log "github.com/sirupsen/logrus"
	"os"
	"encoding/hex"
)

type Hasher struct {
}

func NewHasher() (*Hasher) {
	hasher := new(Hasher)
	return hasher
}

func (hasher *Hasher) Name() string {
	return "Hasher"
}

func (hasher *Hasher) Concurrency() int {
	return 10
}

func (hasher *Hasher) Process(item interface{}) {
	i := item.(*SessionItem)

	// skip non-regular files
	if ! i.GetAttr("mode").(os.FileMode).IsRegular() {
		return
	}

	path := i.Id.Path
	md5 := md5.New()
	sha256 := sha256.New()
	sha256.Reset()
	md5.Reset()

	// use already loaded file content if exists
	if i.HasAttr("content") {
		content := i.GetAttrString("content")
		if _, err := md5.Write([]byte(content)); err != nil {
			log.Fatal("md5", err)
		}
		if _, err := sha256.Write([]byte(content)); err != nil {
			log.Fatal("sha256", err)
		}

	} else {

		f, err := os.Open(path)
		if err != nil {
			i.SetAttr("hasher.error", err)
			log.WithFields(log.Fields{"err":err, "item": i}).Warn("hasher error open")
			return
		}
		defer f.Close()

		// make a buffer to keep chunks that are read
		buf := make([]byte, 4096*10)
		for {
			// read a chunk
			n, err := f.Read(buf)
			if err != nil && err != io.EOF {
				log.Fatal("file read", err)
			}
			if n == 0 {
				break
			}

			// write a chunk to md5 hasher
			if _, err := md5.Write(buf[:n]); err != nil {
				log.Fatal("md5", err)
			}
			// write a chunk to sha256 hasher
			if _, err := sha256.Write(buf[:n]); err != nil {
				log.Fatal("sha256", err)
			}
		}
	}

	i.SetAttr("md5", hex.EncodeToString(md5.Sum(nil)) )
	i.SetAttr("sha256", hex.EncodeToString(sha256.Sum(nil)))

}
