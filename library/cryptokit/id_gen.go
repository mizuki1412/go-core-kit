package cryptokit

import (
	"github.com/rs/xid"
)

func ID() string {
	id := xid.New()
	return id.String()
}
