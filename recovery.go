package neko

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
)

const panicHtml = `
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>PANIC</title>
<style>
body{margin:0;background:#333;color:#fff;font-size:14px;}
h1{margin:0;padding:20px 30px 15px;border-bottom:2px solid #000;background:#222;color:#d04526;font-size:30px;}
</style>
</head>
<body>
<h1>PANIC</h1>
<pre style="margin:20px 30px;font-size:20px;">%s</pre>
<pre style="margin:10px 30px 10px;padding:25px 30px;border:2px solid #222;background:#444;color:#ccc">%s</pre>
</body>
</html>
`

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// stack returns a nicely formated stack frame, skipping skip frames
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

// Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.
func Recovery() HandlerFunc {
	return func(ctx *Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := stack(3)
				log.Printf("[%s] PANIC: %s\n%s", ctx.Engine.AppName, err, stack)
				ctx.Writer.WriteHeader(http.StatusInternalServerError)
				if NekoEnv != PROD {
					ctx.SetHeader("Content-Type", "text/html")
					ctx.Writer.Write([]byte([]byte(fmt.Sprintf(panicHtml, err, stack))))
				}
				ctx.Abort()
			}
		}()

		ctx.Next()
	}
}
