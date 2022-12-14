package microcore

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
)

var (
	output = log.New(os.Stdout, "", 0)
)

var (
	green  = string([]byte{27, 91, 48, 48, 58, 51, 50, 109})
	yellow = string([]byte{27, 91, 48, 48, 59, 51, 51, 109})
	red    = string([]byte{27, 91, 48, 48, 59, 51, 49, 109})
	blue   = string([]byte{27, 91, 48, 48, 59, 51, 52, 109})
	white  = string([]byte{27, 91, 48, 109})
)

type Logger struct{}

func (l Logger) getColorByStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return green
	case code >= 300 && code < 400:
		return blue
	case code >= 400 && code < 500:
		return yellow
	default:
		return red
	}
}

func (l Logger) colorStatus(code int) string {
	return l.getColorByStatus(code) + strconv.Itoa(code) + white
}

func (l Logger) colorMethod(method []byte, code int) string {
	return l.getColorByStatus(code) + string(method) + white
}

func (l Logger) getHttp(ctx *fasthttp.RequestCtx) string {
	if ctx.Response.Header.IsHTTP11() {
		return "HTTP/1.1"
	}
	return "HTTP/1.0"
}

/* ========================== Predefined Formats =========================== */

// Tiny format:
// <method> <url> - <status> - <response-time us>
// GET / - 200 - 11.925 us
func (l Logger) Tiny(req fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		begin := time.Now()
		req(ctx)
		end := time.Now()
		output.Printf("%s %s - %v - %v",
			ctx.Method(),
			ctx.RequestURI(),
			ctx.Response.Header.StatusCode(),
			end.Sub(begin),
		)
	})
}

// TinyColored is same as Tiny but colored
func (l Logger) TinyColored(req fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		begin := time.Now()
		req(ctx)
		end := time.Now()
		output.Printf("%s %s - %v - %v",
			l.colorMethod(ctx.Method(), ctx.Response.Header.StatusCode()),
			ctx.RequestURI(),
			l.colorStatus(ctx.Response.Header.StatusCode()),
			end.Sub(begin),
		)
	})
}

// Short format:
// <remote-addr> | <HTTP/:http-version> | <method> <url> - <status> - <response-time us>
// 127.0.0.1:53324 | HTTP/1.1 | GET /hello - 200 - 44.8??s
func (l Logger) Short(req fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		begin := time.Now()
		req(ctx)
		end := time.Now()
		output.Printf("%v | %s | %s %s - %v - %v",
			ctx.RemoteAddr(),
			l.getHttp(ctx),
			ctx.Method(),
			ctx.RequestURI(),
			ctx.Response.Header.StatusCode(),
			end.Sub(begin),
		)
	})
}

// ShortColored is same as Short but colored
func (l Logger) ShortColored(req fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		begin := time.Now()
		req(ctx)
		end := time.Now()
		output.Printf("%v | %s | %s %s - %v - %v",
			ctx.RemoteAddr(),
			l.getHttp(ctx),
			l.colorMethod(ctx.Method(), ctx.Response.Header.StatusCode()),
			ctx.RequestURI(),
			l.colorStatus(ctx.Response.Header.StatusCode()),
			end.Sub(begin),
		)
	})
}

// Combined format:
// [<time>] <remote-addr> | <HTTP/http-version> | <method> <url> - <status> - <response-time us> | <user-agent>
// [2017/05/31 - 13:27:28] 127.0.0.1:54082 | HTTP/1.1 | GET /hello - 200 - 48.279??s | Paw/3.1.1 (Macintosh; OS X/10.12.5) GCDHTTPRequest
func (l Logger) Combined(req fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		begin := time.Now()
		req(ctx)
		end := time.Now()
		output.Printf("[%v] %v | %s | %s %s - %v - %v | %s",
			end.Format("2006/01/02 - 15:04:05"),
			ctx.RemoteAddr(),
			l.getHttp(ctx),
			ctx.Method(),
			ctx.RequestURI(),
			ctx.Response.Header.StatusCode(),
			end.Sub(begin),
			ctx.UserAgent(),
		)
	})
}

// CombinedColored is same as Combined but colored
func (l Logger) CombinedColored(req fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		begin := time.Now()
		req(ctx)
		end := time.Now()
		output.Printf("[%v] %v | %s | %s %s - %v - %v | %s",
			end.Format("2006/01/02 - 15:04:05"),
			ctx.RemoteAddr(),
			l.getHttp(ctx),
			l.colorMethod(ctx.Method(), ctx.Response.Header.StatusCode()),
			ctx.RequestURI(),
			l.colorStatus(ctx.Response.Header.StatusCode()),
			end.Sub(begin),
			ctx.UserAgent(),
		)
	})
}
