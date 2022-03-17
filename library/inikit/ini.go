package inikit

import "gopkg.in/ini.v1"

func Load(src any) (*ini.File, error) {
	return ini.LoadSources(ini.LoadOptions{
		SkipUnrecognizableLines: true,
	}, src)
}
