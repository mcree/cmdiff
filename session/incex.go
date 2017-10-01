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
	"github.com/gobwas/glob"
	"log"
	"os"
)

type IncEx struct {
	includeGlobs []glob.Glob
	excludeGlobs []glob.Glob
}

// Include and Exclude files by glob pattern lists
func NewIncEx(include []string, exclude []string) (*IncEx) {
	incex := new(IncEx)

	// compile incliude directives
	for _, i := range include {
		gl, err := glob.Compile(i,os.PathSeparator)
		if err != nil {
			log.Fatal("error compiling glob path: ",i," err:",err)
		}
		incex.includeGlobs = append(incex.includeGlobs, gl)
	}

	// compile exclude directives
	for _, e := range exclude {
		gl, err := glob.Compile(e,os.PathSeparator)
		if err != nil {
			log.Fatal("error compiling glob path: ",e," err:",err)
		}
		incex.excludeGlobs = append(incex.excludeGlobs, gl)
	}

	return incex
}

// determines if the given path would be allowed by include and exclude rules
// returns true for included paths, false otherwise
func (incex *IncEx) Includes(path string) (bool) {
	included := false
	// check include rules
	for _, i := range incex.includeGlobs {
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
	for _, e := range incex.excludeGlobs {
		if e.Match(path) {
			excluded = true
			break
		}
	}
	return !excluded
}