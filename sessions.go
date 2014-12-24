package neko

import (
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
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

// Sessions is a Middleware that maps a session.Session service into the neko handler chain.
func Sessions(name string, store SessionStore) HandlerFunc {
	return func(ctx *Context) {
		sess := &session{name: name, request: ctx.Req, store: store, written: false, writer: ctx.Writer}
		ctx.Session = sess
		ctx.Writer.Before(func(ResponseWriter) {
			if sess.Written() {
				err := sess.Session().Save(ctx.Req, ctx.Writer)
				if err != nil {
					log.Printf("[NEKO] SESSION ERROR! %s\n", err)
				}
			}
		})
		defer context.Clear(ctx.Req)
		ctx.Next()
	}
}

type session struct {
	name    string
	request *http.Request
	store   SessionStore
	session *sessions.Session
	written bool
	writer  http.ResponseWriter
}

func (s *session) Get(key interface{}) interface{} {
	return s.Session().Values[key]
}

func (s *session) Set(key interface{}, val interface{}) {
	s.Session().Values[key] = val
	s.written = true
}

func (s *session) Delete(key interface{}) {
	delete(s.Session().Values, key)
	s.written = true
}

func (s *session) Clear() {
	for key := range s.Session().Values {
		s.Delete(key)
	}
}

func (s *session) AddFlash(value interface{}, vars ...string) {
	s.Session().AddFlash(value, vars...)
	s.written = true
}

func (s *session) Flashes(vars ...string) []interface{} {
	s.written = true
	return s.Session().Flashes(vars...)
}

func (s *session) Options(options SessionOptions) {
	s.Session().Options = &sessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HTTPOnly,
	}
}

func (s *session) Session() *sessions.Session {
	if s.session == nil {
		var err error
		if s.session, err = s.store.Get(s.request, s.name); err != nil {
			log.Printf("[NEKO] SESSION ERROR! %s\n", err)
		}
	}

	return s.session
}

func (s *session) Written() bool {
	return s.written
}
