/*
   Copyright 2022 dexenrage

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package config

import (
	"dexbot/catcherr"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

var cfg *koanf.Koanf

const DefaultSeparator = ` `

func init() {
	cfg = koanf.New(`.`)
	file := file.Provider(`config.yml`)

	err := cfg.Load(file, yaml.Parser())
	catcherr.HandleError(err)
}

func String(path string) string {
	return cfg.String(path)
}

func StringSlice(path string, sep string) []string {
	s := String(path)
	return strings.Split(s, sep)
}
