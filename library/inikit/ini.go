package inikit

import "gopkg.in/ini.v1"

func Load(src interface{}) (*ini.File, error) {
	return ini.LoadSources(ini.LoadOptions{
		SkipUnrecognizableLines: true,
	}, src)
}
