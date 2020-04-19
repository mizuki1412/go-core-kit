package middleware

import (
	"github.com/alexedwards/scs/v2"
	"net/http"
	"time"
)

var Session *scs.SessionManager

func init()  {
	Session = scs.New()
	// todo 是通过cookies expires的原理？
	Session.IdleTimeout = 3*time.Hour
	Session.Cookie.SameSite = http.SameSiteNoneMode
	// todo 无法在生成时拿到session token
}
