package ws

import (
	"net/http"

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
		if err := jsoniter.Unmarshal(p, &message); err != nil {
			s.logger.Debug("failed to unmarshal message", zap.Error(err), zap.String("message", string(p)))

			if err := conn.WriteMessage(messageType, []byte("failed read message")); err != nil {
				s.logger.Error("failed to write message", zap.Error(err))
				return
			}

			continue
		}

		if err := conn.WriteMessage(messageType, p); err != nil {
			s.logger.Error("failed to write message", zap.Error(err))
			return
		}
	}
}
