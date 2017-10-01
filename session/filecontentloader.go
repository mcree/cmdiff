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
	"os"
	"github.com/spf13/viper"
	"io/ioutil"
	"github.com/sirupsen/logrus"
	"fmt"
    "code.cloudfoundry.org/bytefmt"
)

type FileContentLoader struct {
	sizeLimit uint64
	incex *IncEx
}

func NewFileContentLoader() (*FileContentLoader) {
	loader := new(FileContentLoader)
	var err error
	loader.sizeLimit, err = bytefmt.ToBytes(viper.GetString("filecontentloader.sizeLimit"))
	if err != nil {
		logrus.Fatal("cannot parse sizelimit", err)
	}
	include := viper.GetStringSlice("filecontentloader.include")
	exclude := viper.GetStringSlice("filecontentloader.exclude")
	loader.incex = NewIncEx(include,exclude)
	return loader
}

func (loader *FileContentLoader) Name() string {
	return "FileContentLoader"
}

func (loader *FileContentLoader) Concurrency() int {
	return 10
}

func (loader *FileContentLoader) Process(item interface{}) {
	i := item.(*SessionItem)

	// skip non-regular files
	if ! i.GetAttr("mode").(os.FileMode).IsRegular() {
		return
	}

	path := i.Id.Path
	if loader.incex.Includes(path) {
		size := i.GetAttrUint64("size")
		if size <= loader.sizeLimit {
			content, err := ioutil.ReadFile(path)
			if err == nil {
				i.SetAttr("content", string(content))
			} else {
				logrus.Info("error loading contents: ", err)
				i.SetAttr("content.error", err)
			}
		} else {
			i.SetAttr("content.error", fmt.Sprintf("size limit exceeded %d >= %d", size, loader.sizeLimit))
		}
	}


}
