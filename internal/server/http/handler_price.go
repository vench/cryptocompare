package http

import (
	"net/http"
	"strings"

	jsoniter "github.com/json-iterator/go"

	"go.uber.org/zap"

	"github.com/valyala/fasthttp"
)

func (s *Server) handlerPrice(rCtx *fasthttp.RequestCtx) {
	fromSymbol, toSymbol := string(rCtx.QueryArgs().Peek("fsyms")), string(rCtx.QueryArgs().Peek("tsyms"))
	if fromSymbol == "" || toSymbol == "" {
		rCtx.Error("empty query", fasthttp.StatusBadRequest)
		return
	}

	result, err := s.storage.GetCurrencyBy(strings.Split(fromSymbol, ","), strings.Split(toSymbol, ","))
	if err != nil {
		s.logger.Error("failed to get currency", zap.Error(err))
		rCtx.Error("failed to get currency", fasthttp.StatusInternalServerError)
		return
	}

	s.logger.Debug("result", zap.Reflect("result", result))

	rCtx.SetStatusCode(http.StatusOK)
	rCtx.Response.Header.SetCanonical(strContentType, strApplicationJSON)

	if err := jsoniter.NewEncoder(rCtx).Encode(result); err != nil {
		s.logger.Error("failed to encode response", zap.Error(err))
		rCtx.Error("failed to encode response", fasthttp.StatusInternalServerError)
	}
}
