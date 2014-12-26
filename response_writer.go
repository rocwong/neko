package neko

import (
	"bufio"
	"errors"
	"log"
	"net"
	"net/http"
)

const noWritten = -1

type (
	ResponseWriter interface {
		http.ResponseWriter
		http.Flusher
		Status() int
		// Size returns the size of the response body.
		Size() int
		Written() bool
		WriteHeaderNow()
		// Before allows for a function to be called before the ResponseWriter has been written to. This is
		// useful for setting headers or any other operations that must happen before a response has been written.
		Before(func(ResponseWriter))
	}
	writer struct {
		http.ResponseWriter
		status      int
		size        int
		beforeFuncs []beforeFunc
	}
	beforeFunc func(ResponseWriter)
)

func (c *writer) Status() int {
	return c.status
}

func (c *writer) Size() int {
	return c.size
}

func (c *writer) Written() bool {
	return c.size != noWritten
}

func (c *writer) WriteHeaderNow() {
	if !c.Written() {
		c.size = 0
		c.callBefore()
		c.ResponseWriter.WriteHeader(c.status)
	}
}
func (c *writer) Before(before func(ResponseWriter)) {
	c.beforeFuncs = append(c.beforeFuncs, before)
}

func (c *writer) Write(data []byte) (size int, err error) {
	c.WriteHeaderNow()
	size, err = c.ResponseWriter.Write(data)
	c.size += size
	return
}

func (c *writer) WriteHeader(code int) {
	if code > 0 {
		c.status = code
		if c.Written() {
			log.Println("[NEKO] WARNING. Headers were already written!")
		}
	}
}

func (c *writer) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := c.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("the ResponseWriter doesn't support the Hijacker interface")
	}
	return hijacker.Hijack()
}

func (c *writer) CloseNotify() <-chan bool {
	return c.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

func (c *writer) Flush() {
	flusher, ok := c.ResponseWriter.(http.Flusher)
	if ok {
		flusher.Flush()
	}
}

func (c *writer) callBefore() {
	for i := len(c.beforeFuncs) - 1; i >= 0; i-- {
		c.beforeFuncs[i](c)
	}
}

func (c *writer) reset(writer http.ResponseWriter) {
	c.ResponseWriter = writer
	c.status = http.StatusOK
	c.beforeFuncs = nil
	c.size = noWritten
}
