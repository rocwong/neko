package neko

import (
	"github.com/gorilla/sessions"
)

// SessionStore is an interface for custom session stores.
type SessionStore interface {
	sessions.Store
	Options(SessionOptions)
}

// SessionOptions stores configuration for a session or session store.
type SessionOptions struct {
	Path     string
	Domain   string
	MaxAge   int
	Secure   bool
	HTTPOnly bool
}

// Session stores the values and optional configuration for a session.
type Session interface {
	Get(key interface{}) interface{}
	Set(key interface{}, val interface{})
	Delete(key interface{})
	Clear()
	AddFlash(value interface{}, vars ...string)
	Flashes(vars ...string) []interface{}
	Options(SessionOptions)
}
