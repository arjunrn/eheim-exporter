package ws

import (
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type Sender struct {
	conn     *websocket.Conn
	interval *time.Ticker
	quit     chan struct{}
	messages []GetFilterData
}

func NewWSSender(conn *websocket.Conn, interval time.Duration, filters ...string) *Sender {
	quit := make(chan struct{})
	sender := &Sender{
		conn:     conn,
		interval: time.NewTicker(interval),
		quit:     quit,
	}
	sender.messages = make([]GetFilterData, len(filters))
	for i, f := range filters {
		sender.messages[i] = NewGetFilterDataMessage(f)
	}
	return sender
}

func (s *Sender) Stop() {
	close(s.quit)
}

func (s *Sender) Run() {
	for {
		select {
		case <-s.interval.C:
			log.Debugf("Sending GET_FILTER_DATA message")
			for _, m := range s.messages {
				err := s.conn.WriteJSON(m)
				if err != nil {
					log.Warnf("failed to send GET_FILTER_DATA message: %v", err)
				}
			}
		case <-s.quit:
			s.interval.Stop()
			return
		}
	}
}
