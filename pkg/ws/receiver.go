package ws

import (
	"github.com/arjunrn/eheim-exporter/pkg/data"
	"github.com/arjunrn/eheim-exporter/pkg/metrics"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type Receiver struct {
	conn          *websocket.Conn
	quit          bool
	filterMetrics metrics.FilterMetrics
}

func NewReceiver(conn *websocket.Conn, filterMetrics metrics.FilterMetrics) *Receiver {
	return &Receiver{
		conn:          conn,
		filterMetrics: filterMetrics,
	}
}

func (r *Receiver) Run() {
	var (
		message    map[string]interface{}
		filterData data.FilterData
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
			r.filterMetrics.RotationSpeed(filterData.From, filterData.RotationSpeed)
			r.filterMetrics.DFS(filterData.From, filterData.DFS)
			r.filterMetrics.DFSFactor(filterData.From, filterData.DFSFactor)
			r.filterMetrics.Frequency(filterData.From, filterData.Frequency)
			r.filterMetrics.PumpMode(filterData.From, filterData.PumpMode)
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
