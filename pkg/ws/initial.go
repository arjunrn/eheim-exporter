package ws

import (
	"encoding/json"
	"fmt"

	"github.com/arjunrn/eheim-exporter/pkg/data"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type initialMessageParser struct {
	conn *websocket.Conn
}

type InitialMessageParser interface {
	Parse() (*data.UserData, *data.NetworkDevice, *data.AccessPoint, *data.FilterData, error)
}

func NewInitialMessageParser(conn *websocket.Conn) InitialMessageParser {
	return &initialMessageParser{conn: conn}
}

func (p *initialMessageParser) Parse() (*data.UserData, *data.NetworkDevice, *data.AccessPoint, *data.FilterData, error) {
	var messages []map[string]interface{}
	err := p.conn.ReadJSON(&messages)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to read first messages on websocket: %w", err)
	}
	var (
		userData   data.UserData
		netSt      data.NetworkDevice
		netAp      data.AccessPoint
		filterData data.FilterData
	)
	for _, m := range messages {
		if title, ok := m["title"]; !ok {
			log.Warnf("received message with no title field: %#v", m)
			continue
		} else {
			switch title {
			case "USRDTA":
				err := reserialize(m, &userData)
				if err != nil {
					log.Warnf("failed to parse userdata: %v", err)
				}
			case "NET_ST":
				err := reserialize(m, &netSt)
				if err != nil {
					log.Warnf("failed to parse network status: %v", err)
				}
			case "NET_AP":
				err := reserialize(m, &netAp)
				if err != nil {
					log.Warnf("failed to parse access point status: %v", err)
				}
			case "CLOCK":
				// ignore
			default:
				log.Warnf("unknown message type: %s", title)
			}
		}
	}

	err = p.conn.ReadJSON(&messages)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to read second messages on websocket: %w", err)
	}
	for _, m := range messages {
		if title, ok := m["title"]; !ok {
			log.Warnf("received message with no title field: %#v", m)
			continue
		} else {
			switch title {
			case "FILTER_DATA":
				err := reserialize(m, &filterData)
				if err != nil {
					log.Warnf("failed to parse filter data: %v", err)
				}
			case "MESH_NETWORK":
				// TODO: Add parsing for this
			default:
				log.Warnf("unknown message type: %s", title)
			}
		}
	}
	return &userData, &netSt, &netAp, &filterData, nil
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
