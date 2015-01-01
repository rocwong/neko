#Neko
[![wercker status](https://app.wercker.com/status/2ab4b79cf2d418606e884c5d98d1ec0d/s "wercker status")](https://app.wercker.com/project/bykey/2ab4b79cf2d418606e884c5d98d1ec0d)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/rocwong/neko)
[![GoCover](http://gocover.io/_badge/github.com/rocwong/neko)](http://gocover.io/github.com/rocwong/neko)

A lightweight web application framework for Golang

**NOTE: Neko is still under development, so API might be changed in future.**

## Features

* Extremely simple to use.
* RESTful support
* Middleware support
* Unlimited nested group routers.

## Getting Started
Basic usage
~~~go
package main
import "github.com/rocwong/neko"
func main() {
  app := neko.Classic()
  app.GET("/", func(ctx *neko.Context)  {
      ctx.Text("Hello world!")
  })
  m.Run(":3000")
}
~~~
Initial Neko without middlewares
~~~go
app := neko.New()
app.Use(neko.Logger())
app.Use(neko.Recovery())
~~~

##Routing
Using GET, POST, PUT, PATCH, DELETE, HEAD and OPTIONS
~~~go
app.GET("/get", get)
app.POST("/post", post)
app.PUT("/put", put)
app.PATCH("/patch", patch)
app.DELETE("/delete", delete)
app.HEAD("/head", head)
app.OPTIONS("/options", options)
~~~
Neko uses julienschmidt's [httprouter](https://github.com/julienschmidt/httprouter) internaly.


##Group Routing
~~~go
v1 := app.Group("/v1", func(router *neko.RouterGroup) {
  //match /v1/item
  router.GET("/item", item)

  // nested group
  router.Group("/sub", func(sub *neko.RouterGroup) {
    //match /v1/sub/myitem
    sub.GET("/myitem", myitem)
  })
})
//match /v1/act
v1.GET("/act", act)
~~~

## Parameters
~~~go
// This handler will match /user/neko but will not match neither /user/ or /user
app.GET("/user/:name", func(ctx *neko.Context) {
  ctx.Text("Hello " + ctx.Params.ByName("name"))
})

// This one will match /user/neko/ and also /user/neko/3, but no match /user/neko
app.GET("/user/:name/*age", func(ctx *neko.Context) {
  name := c.Params.ByName("name")
  age := c.Params.ByName("age")
  message := name + " is " + action
  ctx.Text(message)
})
~~~

##Response

####Render
~~~go
type ExampleXml struct {
  XMLName xml.Name `xml:"example"`
  One     string   `xml:"one,attr"`
  Two     string   `xml:"two,attr"`
}

//response: <example one="hello" two="xml"/>
ctx.Xml(ExampleXml{One: "hello", Two: "xml"})
~~~

~~~go
//response: {"msg": "json render", "status": 200}
ctx.Json(neko.JSON{"msg": "json render", "status": 200})

//response: neko({"msg": "json render", "status": 200})
ctx.Jsonp("neko", neko.JSON{"msg": "json render", "status": 200})

// response: neko text
ctx.Text("neko text")
~~~

####Redirect
~~~go
//default 302
ctx.Redirect("/")

//redirect 301
ctx.Redirect("/", 301)
~~~

####Headers
~~~go
// get header
ctx.Writer.Header()

// set header
ctx.SetHeader("x-before", "before")
~~~

##Cookie
~~~ go
app.GET("/", func (ctx *neko.Context) {
  ctx.SetCookie("myvalue", "Cookies Save")
  ctx.Text("Cookies Save")
})

app.GET("/get", func (ctx *neko.Context) {
  ctx.Text(ctx.GetCookie("myvalue"))
})
~~~
####Secure cookie
~~~ go
//set cookie secret
app.SetCookieSecret("secret123")

app.GET("/set-secure", func (ctx *neko.Context) {
  ctx.SetSecureCookie("sv", "Cookies Save")
  ctx.Text("Cookies Save")
})

app.GET("/get-secure", func (ctx *neko.Context) {
  ctx.Text(ctx.GetSecureCookie("sv"))
})

~~~
Use following arguments order to set more properties: SetCookie/SetCookieSecret(name, value [, MaxAge, Path, Domain, Secure, HttpOnly]).

## Middlewares

####Using middlewares
~~~go
// Global middlewares
app.Use(neko.Logger())

// Per route middlewares, you can add as many as you desire.
app.Get("/user", mymiddleware(), mymiddleware2(), user)

//Pass middlewares to groups
v1 := app.Group("/v1", func(router *neko.RouterGroup) {
  //url: /v1/item
  router.GET("/item", item)
}, mymiddleware1(), mymiddleware2(), mymiddleware3())

v1.Use(mymiddleware4)
~~~

####Custom middlewares
~~~go
func mymiddleware() neko.HandlerFunc {
  return func (ctx *neko.Context) {
    //before request
    t := time.Now()

    ctx.Next()

    // after request
    latency := time.Since(t)
    log.Print(latency)

    // access the status we are sending
    status := c.Writer.Status()
    log.Println(status)
  }
}
~~~

#### More middleware
For more middleware and functionality, check out the repositories in the  [neko-contrib](https://github.com/neko-contrib) organization.

## Others
~~~go
// static serves
app.Static("/static", "content/static")

// get remote ip address
app.GET("/", func (ctx *neko.Context) {
  ctx.ClientIP()
}
~~~


## Credits & Thanks
I use code/got inspiration from these excellent libraries:

*  [Gin](https://github.com/gin-gonic/gin) - design based on.
*  [Httprouter](https://github.com/julienschmidt/httprouter)
*  [Martini](https://github.com/go-martini/martini)


## License
Neko is licensed under the MIT