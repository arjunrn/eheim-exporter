package processor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/arjunrn/eheim-exporter/pkg/data"
	"github.com/arjunrn/eheim-exporter/pkg/metrics"
	log "github.com/sirupsen/logrus"
)

type Receiver struct {
	filterMetrics metrics.FilterMetrics
	messageChan   chan []byte
	tracker       FilterIDTracker
}

func NewReceiver(messageChan chan []byte, filterMetrics metrics.FilterMetrics, tracker FilterIDTracker) *Receiver {
	return &Receiver{
		messageChan:   messageChan,
		filterMetrics: filterMetrics,
		tracker:       tracker,
	}
}

func (r *Receiver) Run(ctx context.Context) error {
	for {
		select {
		case message := <-r.messageChan:
			r.processMessage(message)
		case <-ctx.Done():
			log.Debugf("Shutting down Receiver")
			return nil
		}
	}
}

func (r *Receiver) processMessage(message []byte) {
	var (
		multiMessage  []map[string]interface{}
		singleMessage map[string]interface{}
	)

	multiMessageErr := json.Unmarshal(message, &multiMessage)
	if multiMessageErr != nil {
		singleMessageErr := json.Unmarshal(message, &singleMessage)
		if singleMessageErr != nil {
			log.Warnf("received message %#v could not be unmarshalled either as multi or single message: %v %v", message, multiMessageErr, singleMessageErr)
		}
	}

	for _, m := range multiMessage {
		r.parseMessage(m)
	}
	if singleMessage != nil {
		r.parseMessage(singleMessage)
	}
}

func (r *Receiver) parseMessage(message map[string]interface{}) {
	var (
		filterData data.FilterData
		userData   data.UserData
		netSt      data.NetworkDevice
		netAp      data.AccessPoint
	)

	if _, ok := message["title"]; !ok {
		log.Warnf("received message %s does not contain field 'title'", message)
		return
	}
	title := message["title"].(string)
	switch title {
	case "FILTER_DATA":
		err := reserialize(message, &filterData)
		if err != nil {
			log.Warnf("failed to get filter data from message: %v", err)
			return
		}
		r.tracker.Add(filterData.From)
		r.filterMetrics.FilterData(filterData)
		log.Debugf("Received filter data %#v", filterData)
	case "REQ_KEEP_ALIVE":
		// keep alive message. do nothing
	case "USRDTA":
		err := reserialize(message, &userData)
		if err != nil {
			log.Warnf("failed to parse userdata: %v", err)
		}
		log.Debugf("Received user data %#v", userData)
		r.filterMetrics.UserData(userData)
	case "NET_ST":
		err := reserialize(message, &netSt)
		if err != nil {
			log.Warnf("failed to parse network status: %v", err)
		}
		log.Debugf("Received network client data %#v", netSt)
		r.filterMetrics.NetworkClient(netSt)
	case "NET_AP":
		err := reserialize(message, &netAp)
		if err != nil {
			log.Warnf("failed to parse access point status: %v", err)
		}
		log.Debugf("Received network access point data %#v", netAp)
		r.filterMetrics.NetworkAccessPoint(netAp)
	case "CLOCK":
		// ignore
	case "MESH_NETWORK":
		// ignore
	case "GET_FILTER_DATA":
		// ignore
	case "GET_MESH_NETWORK":
		// ignore
	default:
		log.Warnf("unknown title %s in message: %#v", title, message)
	}
}

func reserialize(message map[string]interface{}, target interface{}) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshall message %#v: %w", message, err)
	}
	err = json.Unmarshal(payload, target)
	if err != nil {
		return fmt.Errorf("failed to unmarshall marshalled data %s: %w", payload, err)
	}
	return nil
}
