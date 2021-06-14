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
	message  *GetFilterData
}

func NewWSSender(conn *websocket.Conn, interval time.Duration) *Sender {
	quit := make(chan struct{})
	return &Sender{
		conn:     conn,
		interval: time.NewTicker(interval),
		quit:     quit,
		message: NewGetFilterDataMessage("FC:F5:C4:93:C5:0A"),
	}
}

func (s *Sender) Stop() {
	close(s.quit)
}

func (s *Sender) Run() {

	for {
		select {
		case <-s.interval.C:
			log.Debugf("Sending GET_FILTER_DATA message")
			err := s.conn.WriteJSON(s.message)
			if err != nil {
				log.Warnf("failed to send GET_FILTER_DATA message: %v", err)
			}
		case <-s.quit:
			s.interval.Stop()
			return
		}
	}
}
