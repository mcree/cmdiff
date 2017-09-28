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

package set

import "encoding/json"

/**
 * Minor rewrite of https://raw.githubusercontent.com/agonopol/goset/master/set.go
 */

type Set map[interface{}]bool

func New() *Set {
	set := make(Set)
	return &set
}

func (this *Set) Add(x interface{}) {
	(*this)[x] = true
}

func (this *Set) Remove(x interface{}) {
	delete((*this), x)
}

func (this *Set) Reset() {
	(*this) = make(map[interface{}]bool)
}

func (this *Set) Has(x interface{}) bool {
	_, found := (*this)[x]
	return found
}

func (this *Set) Do(f func(interface{})) {
	for k, _ := range *this {
		f(k)
	}
}

func (this *Set) Len() int {
	return len(*this)
}

func (this *Set) Map(f func(interface{}) interface{}) *Set {
	set := New()
	for k, _ := range *this {
		set.Add(f(k))
	}
	return set
}

func (this *Set) Items() ([]interface{}) {
	set := make([]interface{}, 0)
	for k, _ := range *this {
		set = append(set, k)
	}
	return set
}

func (this *Set) StringItems() ([]string) {
	set := make([]string, 0)
	for k, _ := range *this {
		set = append(set, k.(string))
	}
	return set
}

func (this *Set) MarshalJSON() ([]byte, error) {
	return json.Marshal(this.Items())
}

func (this *Set) UnmarshalJSON(body []byte) error {
	var set []interface{}
	err := json.Unmarshal(body, &set)
	if err != nil {
		return err
	} else {
		for _, v := range set {
			this.Add(v)
		}
	}
	return nil
}

