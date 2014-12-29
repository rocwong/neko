package neko

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
	"sync"
)

type (
	HandlerFunc func(*Context)
	// HtmlEngine is an interface for parsing html templates and redering HTML.
	HtmlEngine interface {
		Render(view string, context interface{}, status ...int) error
	}
	Engine struct {
		*RouterGroup
		AppName    string
		router     *httprouter.Router
		allNoRoute []HandlerFunc
		pool       sync.Pool
	}
)

func Version() string {
	return "0.0.1"
}

// Classic creates a classic Neko with some basic default middleware - neko.Logger and neko.Recovery.
func Classic(appName ...string) *Engine {
	engine := New()
	if appName != nil && appName[0] != "" {
		engine.AppName = appName[0]
	}
	engine.Use(Logger())
	engine.Use(Recovery())
	return engine
}

// New returns a new blank Engine instance without any middleware attached.
func New() *Engine {
	engine := &Engine{}
	engine.AppName = "NEKO"
	engine.RouterGroup = &RouterGroup{
		absolutePath: "/",
		engine:       engine,
	}
	engine.router = httprouter.New()
	engine.router.NotFound = engine.handle404
	engine.pool.New = func() interface{} {
		ctx := &Context{Engine: engine}
		return ctx
	}
	return engine
}

func (c *Engine) Use(middlewares ...HandlerFunc) {
	c.RouterGroup.Use(middlewares...)
	c.allNoRoute = c.combineHandlers(nil)
}

// ServeHTTP makes the router implement the http.Handler interface.
func (c *Engine) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	c.router.ServeHTTP(res, req)
}

// Run run the http server.
func (c *Engine) Run(addr string) {
	fmt.Printf("[%s] Listening and serving HTTP on %s \n", c.AppName, addr)
	if err := http.ListenAndServe(addr, c); err != nil {
		panic(err)
	}
}

// Run run the https server.
func (c *Engine) RunTLS(addr string, cert string, key string) {
	fmt.Printf("[%s] Listening and serving HTTPS on %s \n", c.AppName, addr)
	if err := http.ListenAndServeTLS(addr, cert, key, c); err != nil {
		panic(err)
	}
}

func (c *Engine) handle404(w http.ResponseWriter, req *http.Request) {
	ctx := c.createContext(w, req, nil, c.allNoRoute)
	ctx.Writer.WriteHeader(404)
	ctx.Next()
	if !ctx.Writer.Written() {
		if ctx.Writer.Status() == 404 {
			ctx.Writer.Header().Set("Content-Type", "text/html")
			ctx.Writer.Write([]byte(`<!DOCTYPE html><html><head><meta charset="UTF-8"><title>404 PAGE NOT FOUND</title></head><body style="padding:0;text-align:center;"><div style="padding-top:1em;font-size:2.5em;">404 PAGE NOT FOUND</div><div style="font-size:1em;color:#999;">Powered by Neko</div></body></html>`))
		} else {
			ctx.Writer.WriteHeader(ctx.Writer.Status())
		}
	}
	c.reuseContext(ctx)
}

func (c *Engine) reuseContext(ctx *Context) {
	c.pool.Put(ctx)
}

const (
	DEV  string = "development"
	PROD string = "production"
	TEST string = "test"
)

// NekoEnv is the environment that Neko is executing in.
// The NEKO_ENV is read on initialization to set this variable.
var NekoEnv = DEV

func init() {
	env := os.Getenv("NEKO_ENV")
	if len(env) > 0 {
		NekoEnv = env
	}
}
