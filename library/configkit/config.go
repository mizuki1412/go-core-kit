package configkit

import (
	"log"
	"mizuki/project/core-kit/library/filekit"
	"mizuki/project/core-kit/library/jsonkit"
)

// Deprecated
func ConfigInit(filepath string, config interface{})  {
	f,err := filekit.ReadString(filepath)
	if err!=nil{
		log.Fatalln(err)
	}
	jsonkit.ParseObj(f,config)
}
