package processor

import (
	"context"
	"encoding/json"
	"time"

	"github.com/arjunrn/eheim-exporter/pkg/data"
	"github.com/arjunrn/eheim-exporter/pkg/wswrapper"
	log "github.com/sirupsen/logrus"
)

type Sender struct {
	interval *time.Ticker
	socket   wswrapper.Interface
	tracker  FilterIDTracker
}

func NewSender(socket wswrapper.Interface, interval time.Duration, tracker FilterIDTracker) *Sender {
	return &Sender{
		socket:   socket,
		interval: time.NewTicker(interval),
		tracker:  tracker,
	}
}

func (s *Sender) Run(ctx context.Context) {
	for {
		select {
		case <-s.interval.C:
			log.Debugf("Sending GET_FILTER_DATA message")
			for _, filterID := range s.tracker.List() {
				payload, err := json.Marshal(data.NewGetFilterDataMessage(filterID))
				if err != nil {
					log.Warnf("Failed to create GET_FILTER_DATA message for filter %s", filterID)
					continue
				}
				err = s.socket.Send(payload)
				if err != nil {
					log.Warnf("Failed to write GET_FILTER_DATA message: %v", err)
				}
			}
		case <-ctx.Done():
			log.Debugf("Stopping the Sender")
			s.interval.Stop()
			return
		}
	}
}
