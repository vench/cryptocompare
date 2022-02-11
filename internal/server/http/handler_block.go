package http

import (
	"net/http"

	"github.com/valyala/fasthttp"
)

func (s *Server) blockId(rCtx *fasthttp.RequestCtx) {
	ID, ok := userValueUint64(rCtx, "id")
	if !ok {
		rCtx.SetStatusCode(http.StatusBadRequest)
	}
	_ = ID

	rCtx.SetStatusCode(http.StatusOK)
	rCtx.SetBodyString("foo")
}
