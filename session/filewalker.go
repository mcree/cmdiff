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
	log "github.com/sirupsen/logrus"
	"github.com/gobwas/glob/syntax/lexer"
	"os"
	"bytes"
	"path/filepath"
	"github.com/spf13/viper"
	"github.com/mcree/cmdiff/set"
	"net/url"
)

// calculates longest non-wildcard prefix for a path
// based on github.com/gobwas/glob/syntax/lexer
func FindRoot(path string) (root string) {
	var res *bytes.Buffer = new(bytes.Buffer)

	l := lexer.NewLexer(path)

	cont:
	for {
		t := l.Next()
		//log.Info(t.Type,t.Raw)
		switch t.Type {
		case lexer.Text:
			res.WriteString(t.Raw)
		case lexer.EOF:
			break cont
		default:

			res = bytes.NewBufferString(filepath.Dir(res.String()))
			break cont
		}
	}

	root = res.String()

	return
}

type Filewalker struct {
	includeRoots []string
	items chan *SessionItem
	incex *IncEx
}

// collect file Items into session
func (fw *Filewalker) Generate() {
	log.Info("Filewalker Generate start")
	defer close(fw.items)
	for _, root := range fw.includeRoots {
		err := filepath.Walk(root,  func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			if fw.incex.Includes(path) {
				var id url.URL
				id.Scheme = "file"
				id.Path = path
				item := NewSessionItemMap(id, map[string]interface{} {
					"size": uint64(info.Size()),
					"modTime": info.ModTime(),
					"mode": info.Mode(),
				})
				if err != nil {
					item.SetAttr("filewalker.error", err)
					log.Warn("Filewalker: ", err)
				}
				log.Debug("Filewalker collected: ", item)
				fw.items <- item
			}
			return nil
		})
		if err != nil {
			log.Fatal("Filewalker walk: ", err)
		}
	}
	log.Info("Filewalker Generate done")
}

// collect (local) file list into current session based on configuration
func NewFilewalker() (*Filewalker) {
	log.Info("Filewalker init.")

	fw := new(Filewalker)

	fw.items = make(chan *SessionItem, 10)

	include := viper.GetStringSlice("filewalker.include")
	exclude := viper.GetStringSlice("filewalker.exclude")

	fw.incex = NewIncEx(include, exclude)

	// process include directives
	log.Debug("Include: ", include)
	includeRoot := set.New()

	for _, i := range include {
		includeRoot.Add(FindRoot(i))
	}
	fw.includeRoots = includeRoot.StringItems()
	log.Debug("roots: ", fw.includeRoots)


	log.Info("Filewalker init done.")

	return fw
}

func (fw *Filewalker) Name() string {
	return "Filewalker"
}

func  (fw *Filewalker) Next() interface{} {
	item, ok := <- fw.items
	if ok {
		return item
	}
	return nil
}

func  (fw *Filewalker) Abort() {
	log.Info("Filewalker abort in progress.")
}
