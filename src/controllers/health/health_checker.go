package health

import "github.com/valyala/fasthttp"

func CheckHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody([]byte(`{"status":"ok"}`))
}
