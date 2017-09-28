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

package util

import (
	log "github.com/sirupsen/logrus"
	glob "github.com/gobwas/glob"
	"github.com/gobwas/glob/syntax/lexer"
	"os"
	"bytes"
	"path/filepath"
	"github.com/spf13/viper"
	"github.com/mcree/cmdiff/set"
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
	includeGlobs []glob.Glob
	excludeGlobs []glob.Glob
	includeRoots []string
}

// determines if the given path would be allowed by include and exclude rules
func (fw *Filewalker) Includes(path string) (bool) {
	included := false
	// check include rules
	for _, i := range fw.includeGlobs {
		if i.Match(path) {
			included = true
			break
		}
	}
	if !included {
		return false
	}

	// check exclude rules
	excluded := false
	for _, e := range fw.excludeGlobs {
		if e.Match(path) {
			excluded = true
			break
		}
	}
	return !excluded
}

// collect file items into session
func (fw *Filewalker) Collect(sess *Session) {
	for _, root := range fw.includeRoots {
		err := filepath.Walk(root,  func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			if fw.Includes(path) {

				log.Debug("Filewalker collected: ", path, " ", info, " ", err)
			}
			return nil
		})
		if err != nil {
			log.Fatal("walk error: ", err)
		}
	}
}

// collect (local) file list into current session based on configuration
func NewFilewalker() (*Filewalker) {
	log.Info("Filewalker init.")

	fw := new(Filewalker)

	// process include directives
	include := viper.GetStringSlice("filewalker.include")
	log.Debug("Include: ", include)
	includeRoot := set.New()

	for _, i := range include {
		includeRoot.Add(FindRoot(i))
		gl, err := glob.Compile(i,os.PathSeparator)
		if err != nil {
			log.Fatal("error compiling glob path: ",i," err:",err)
		}
		fw.includeGlobs = append(fw.includeGlobs, gl)
	}
	fw.includeRoots = includeRoot.StringItems()
	log.Debug("roots: ", fw.includeRoots)

	// process exclude directives
	exclude := viper.GetStringSlice("filewalker.exclude")
	log.Debug("Exclude: ", exclude)
	for _, e := range exclude {
		gl, err := glob.Compile(e,os.PathSeparator)
		if err != nil {
			log.Fatal("error compiling glob path: ",e," err:",err)
		}
		fw.excludeGlobs = append(fw.excludeGlobs, gl)
	}

	log.Debug(fw.excludeGlobs)

	log.Info("Filewalker init done.")

	return fw
}
