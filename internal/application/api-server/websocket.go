package apiserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (a *Application) websocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		ErrorWithCode(c, err, ErrInternalServer)
		return
	}
	conn.Close()
	// a.handleWebsocket(conn)
}

// func (a *Application) handleWebsocket(conn *websocket.Conn) {
// 	ticker := time.NewTicker(1 * time.Second)
// 	defer ticker.Stop()
// 	defer conn.Close()

// 	for {
// 		select {
// 		case <-a.ctx.Done():
// 			return
// 		case <-ticker.C:
// 			orderbook, ok := a.binance.GetOrderbook("BTCUSDT")
// 			if !ok {
// 				continue
// 			}
// 			if err := conn.WriteJSON(map[string]any{
// 				"s": "b",
// 				"d": orderbook,
// 			}); err != nil {
// 				return
// 			}
// 			if err := conn.WriteJSON(map[string]any{
// 				"s": "p",
// 				"d": struct {
// 					Symbol string `json:"s"`
// 					Price float64 `json:"p"`
// 				}{
// 					Symbol: orderbook.Symbol,
// 					Price: orderbook.IndexPrice,
// 				},
// 			}); err != nil {
// 				return
// 			}
// 		}
// 	}
// }
