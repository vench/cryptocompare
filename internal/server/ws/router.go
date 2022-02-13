package ws

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) router() http.Handler {
	mux := chi.NewMux()

	mux.Get("/service/price", s.handlerPrice)

	return mux
}
