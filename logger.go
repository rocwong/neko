package neko

import (
	"log"
	"os"
	"time"
)

var (
	green   = "\033[32m"
	white   = "\033[37m"
	yellow  = "\033[33m"
	red     = "\033[31m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
	reset   = "\033[0m"
)

func Logger() HandlerFunc {
	stdlogger := log.New(os.Stdout, "", 0)

	return func(ctx *Context) {
		// Start timer
		start := time.Now()

		// Process request
		ctx.Next()
		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		clientIP := ctx.ClientIP()
		method := ctx.Req.Method
		statusCode := ctx.Writer.Status()
		statusColor := colorForStatus(statusCode)
		methodColor := colorForMethod(method)

		stdlogger.Printf("%s[%s]%s %v |%s %3d %s| %12v | %s |%s %-5s %s %s",
			blue, ctx.Engine.AppName, reset,
			end.Format("2006/01/02 - 15:04:05"),
			statusColor, statusCode, reset,
			latency,
			clientIP,
			methodColor, method, reset,
			ctx.Req.URL.Path,
		)
	}
}

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code <= 299:
		return green
	case code >= 300 && code <= 399:
		return white
	case code >= 400 && code <= 499:
		return yellow
	default:
		return red
	}
}

func colorForMethod(method string) string {
	switch {
	case method == "GET":
		return blue
	case method == "POST":
		return cyan
	case method == "PUT":
		return yellow
	case method == "DELETE":
		return red
	case method == "PATCH":
		return green
	case method == "HEAD":
		return magenta
	default:
		return white
	}
}
