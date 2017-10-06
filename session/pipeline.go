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
	"github.com/jboelter/pipeline"
	log "github.com/sirupsen/logrus"
)

type Pipeline struct {

}

func NewPipeline() (*Pipeline) {
	pl := new(Pipeline)
	return pl
}

func (pl *Pipeline) Run() (*Session) {
	log.Info("Pipeline start")

	p := pipeline.New()

	fw := NewFilewalker()
	go fw.Generate()

	p.SetGenerator(fw)
	p.AddStage(NewStater())
	p.AddStage(NewFileContentLoader())
	p.AddStage(NewHasher())
	sess := NewSession()
	p.AddStage(sess)
	p.Run()
	log.Info("Pipeline done.")
	return sess
}