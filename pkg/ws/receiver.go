package ws

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type Receiver struct {
	conn *websocket.Conn
	quit bool
}

func NewReceiver(conn *websocket.Conn) *Receiver {
	return &Receiver{conn: conn}
}

func (r *Receiver) Run() {
	var (
		message    map[string]interface{}
		filterData FilterData
	)
	for {
		err := r.conn.ReadJSON(&message)
		if err != nil {
			log.Warnf("failed to read message: %v", err)
			continue
		}
		if r.quit {
			return
		}
		if _, ok := message["title"]; !ok {
			log.Warnf("received message %s does not contain field 'title'", message)
			continue
		}
		title := message["title"].(string)
		switch title {
		case "FILTER_DATA":
			err := reserialize(message, &filterData)
			if err != nil {
				log.Warnf("failed to get filter data from message: %v", err)
				continue
			}
			log.Debugf("received filter data: %#v", filterData)
		case "REQ_KEEP_ALIVE":
			// keep alive message. do nothing
		default:
			log.Warnf("unknown title %s in message: %#v", title, message)
			continue
		}
	}
}

func (r *Receiver) Stop() {
	r.quit = true
}
