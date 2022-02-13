package ws

import (
	"net/http"
	"strings"

	"github.com/vench/cryptocompare/internal/server"

	jsoniter "github.com/json-iterator/go"

	"go.uber.org/zap"

	"github.com/fasthttp/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (s *Server) handlerPrice(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Error("failed to upgrade upgrader", zap.Error(err))
		return
	}

	s.logger.Debug("client connected")

	err = ws.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
		s.logger.Error("failed to write message", zap.Error(err))
		return
	}

	s.reader(ws)
}

type messageRequest struct {
	FromSymbols string `json:"fsyms"`
	ToSymbols   string `json:"tsyms"`
}

var (
	messageFailedMarshalResponse  = []byte("failed to marshal response")
	messageFailedGetCurrency      = []byte("failed to get currency")
	messageFailedUnmarshalMessage = []byte("failed to unmarshal message")
)

func (s *Server) reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				s.logger.Debug("connection has been closed")
				return
			}

			s.logger.Error("failed to read message", zap.Error(err))
			return
		}

		s.logger.Debug("print out that message", zap.String("message", string(p)))

		var message messageRequest
		if err = jsoniter.Unmarshal(p, &message); err != nil {
			s.logger.Debug("failed to unmarshal message", zap.Error(err), zap.String("message", string(p)))

			s.writeMessage(conn, messageType, messageFailedUnmarshalMessage)
			continue
		}

		result, err := s.storage.GetCurrencyBy(
			strings.Split(message.FromSymbols, ","),
			strings.Split(message.ToSymbols, ","))
		if err != nil {
			s.logger.Error("failed to get currency", zap.Error(err))

			s.writeMessage(conn, messageType, messageFailedGetCurrency)
			return
		}

		response := server.MakeCurrencyResponse(result)
		data, err := jsoniter.Marshal(response)
		if err != nil {
			s.logger.Error("failed to marshal response", zap.Error(err))

			s.writeMessage(conn, messageType, messageFailedMarshalResponse)
			return
		}

		s.writeMessage(conn, messageType, data)
	}
}

func (s *Server) writeMessage(conn *websocket.Conn, messageType int, data []byte) {
	if err := conn.WriteMessage(messageType, data); err != nil {
		s.logger.Error("failed to write message", zap.Error(err))
	}
}
