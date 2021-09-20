package simplegin

import (
	"fmt"
	"os"
	"time"
)

func Logger() HandlerFunc {
	return func(ctx *Context) {
		start := time.Now()

		ctx.Next()

		path := ctx.Path
		query := ctx.Req.URL.RawQuery
		if query != "" {
			path += "?" + query
		}
		method := ctx.Method
		statusCode := ctx.StatusCode
		end := time.Now()
		latency := end.Sub(start)
		clientIP := ctx.Req.RemoteAddr

		fmt.Fprintln(
			os.Stdout,
			fmt.Sprintf(
				"[SIMPLEGIN] %v | %3d | %13v | %15s | %-7s %#v",
				end.Format("2006/01/02 - 15:04:05"),
				statusCode,
				latency,
				clientIP,
				method, // len("OPTIONS") = 7
				path,
			),
		)
	}
}
