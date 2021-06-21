package wswrapper

import (
	"context"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

const (
	pingInterval = time.Second * 5
)

type Interface interface {
	Run(context.Context) error
	ReceiverChan() chan []byte
	Send([]byte) error
}

func NewReconnectingWebsocket(address string, resyncInterval time.Duration) Interface {
	return &reconnectingWS{
		address:        address,
		receiver:       make(chan []byte, 100),
		resyncInterval: resyncInterval,
	}
}

type reconnectingWS struct {
	address        string
	conn           *websocket.Conn
	stop           bool
	receiver       chan []byte
	resyncInterval time.Duration
}

func (r *reconnectingWS) reconnect(ctx context.Context) error {
	var err error
	r.conn, _, err = websocket.DefaultDialer.DialContext(ctx, r.address, nil)
	if err != nil {
		return fmt.Errorf("failed to open connection to %s: %w", r.address, err)
	}
	return nil
}

func (r *reconnectingWS) Run(ctx context.Context) error {
	if err := r.reconnect(ctx); err != nil {
		return err
	}

	go r.read(ctx)
	ticker := time.NewTicker(pingInterval)
	for {
		select {
		case <-ticker.C:
			func() {
				log.Debug("Sending ping message")
				err := r.conn.WriteMessage(websocket.PingMessage, nil)
				if err != nil && (websocket.IsUnexpectedCloseError(err) || websocket.IsCloseError(err)) {
					_ = r.conn.Close()
					err = r.reconnect(ctx)
					if err != nil {
						log.Errorf("failed to reconnect: %v", err)
					}
				}
			}()
		case <-ctx.Done():
			log.Debugf("Shutting down websocket connection")
			ticker.Stop()
			r.stop = true
			return nil
		}
	}
}

func (r *reconnectingWS) read(ctx context.Context) {
	for {
		if r.stop {
			return
		}
		var (
			messageType int
			message     []byte
			err         error
		)
		func() {
			defer func() {
				if err := recover(); err != nil {
					log.Warnf("Recovered from panic while reading message: %v", err)
					time.Sleep(10 * time.Second)
				}
			}()
			log.Debugf("Waiting for read message")
			messageType, message, err = r.conn.ReadMessage()
			log.Debugf("Received message: %#v", string(message))
			if err != nil {
				if websocket.IsCloseError(err) || websocket.IsUnexpectedCloseError(err) {
					err = r.reconnect(ctx)
					if err != nil {
						log.Warnf("failed to reconnect: %v", err)
						return
					}
				}
				log.Warnf("Failed to read message: %v", err)
			}
		}()
		if message == nil {
			log.Warnf("received message is empty")
			continue
		}
		switch messageType {
		case websocket.TextMessage:
			r.receiver <- message
		case -1:
			log.Debugf("Receoved nothing")
			// do nothing
		default:
			log.Warnf("Received unknown message type %d", messageType)
		}
	}
}

func (r *reconnectingWS) Send(message []byte) error {
	if r.conn == nil {
		panic("tried to send message on non-running websocket client")
	}
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("recovered from panic after writing to websocket: %v", r)
		}
	}()
	log.Debugf("Acquiring lock for writing")
	err := r.conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}
	log.Debugf("Finished writing message")
	return nil
}

func (r *reconnectingWS) ReceiverChan() chan []byte {
	return r.receiver
}
