package cryptokit

import (
	uuid "github.com/google/uuid"
	"github.com/mizuki1412/go-core-kit/class/exception"
)

func UUID() string {
	id := uuid.New()
	if id.String() == "" {
		panic(exception.New("uuid gen error"))
	}
	return id.String()
}
