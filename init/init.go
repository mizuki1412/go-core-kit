package corekit

import (
	_ "go.uber.org/automaxprocs"
)

func init() {

}

func Waiting() {
	select {}
}
