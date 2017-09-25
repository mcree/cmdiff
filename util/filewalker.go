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
	"github.com/spf13/viper"
	log "github.com/sirupsen/logrus"
	glob "github.com/gobwas/glob"
)

// collect (local) file list into current session based on configuration
func Filewalker() {
	log.Info("Filewalker start...")
	include := viper.GetStringSlice("filewalker.include")
	log.Debug("Include: ", include)
	for _, i := range include {
		log.Debug(glob.Compile(i))
	}
	log.Info("Filewalker done.")
}
