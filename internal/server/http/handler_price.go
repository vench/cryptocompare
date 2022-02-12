package http

import (
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/valyala/fasthttp"
)

func (s *Server) handlerPrice(rCtx *fasthttp.RequestCtx) {
	fromSymbol, toSymbol := string(rCtx.QueryArgs().Peek("fsyms")), string(rCtx.QueryArgs().Peek("tsyms"))
	if fromSymbol == "" || toSymbol == "" {
		rCtx.SetStatusCode(http.StatusBadRequest)
		return
	}

	result, err := s.storage.GetCurrencyBy(strings.Split(fromSymbol, ","), strings.Split(toSymbol, ","))
	if err != nil {
		s.logger.Error("failed to get currency", zap.Error(err))
		rCtx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	s.logger.Debug("result", zap.Reflect("result", result))

	rCtx.SetStatusCode(http.StatusOK)
	rCtx.SetBodyString(fmt.Sprintf("%s %s", fromSymbol, toSymbol))
}
