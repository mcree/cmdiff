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
	"gopkg.in/flosch/pongo2.v3"
)

func (diff *SessionDiff) StringTemplate(template string) (string, error) {
	tpl, err := pongo2.FromString(template)
	res, err := tpl.Execute(pongo2.Context {
		"diff": diff,
	})
	return res, err
}