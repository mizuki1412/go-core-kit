package cryptokit

import (
	uuid "github.com/iris-contrib/go.uuid"
	"github.com/mizuki1412/go-core-kit/class/exception"
)

func UUID() string {
	id, _ := uuid.NewV4()
	if id.String() == "" {
		panic(exception.New("uuid gen error"))
	}
	return id.String()
}
