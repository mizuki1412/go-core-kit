package cryptokit

import (
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/mizuki1412/go-core-kit/class/exception"
)

func NanoID() string {
	id, err := gonanoid.New()
	if err != nil {
		panic(exception.New(err.Error()))
	}
	return id
}
