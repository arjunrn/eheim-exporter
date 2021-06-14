package ws

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Receiver struct {
	conn *websocket.Conn
	quit bool
}

func NewReceiver(conn *websocket.Conn) *Receiver {
	return &Receiver{conn: conn}
}

func (r *Receiver) Run() {
	var message interface{}
	for {
		err := r.conn.ReadJSON(&message)
		if err != nil {
			logrus.Warnf("failed to read message: %v", err)
			continue
		}
		if r.quit {
			return
		}
		logrus.Infof("Received Message: %#v", message)
	}
}

func (r *Receiver) Stop() {
	r.quit = true
}
