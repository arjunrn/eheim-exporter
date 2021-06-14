package app

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/arjunrn/eheim-exporter/pkg/ws"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

func App(ctx context.Context, websocketUrl string) {
	log.SetLevel(log.DebugLevel)
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, websocketUrl, nil)
	if err != nil {
		panic(err)
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Errorf("failed to close connection correctly: %v", err)
		}
	}(conn)

	parser := ws.NewInitialMessageParser(conn)
	userData, networkStatus, accessPoint, filterData, err := parser.Parse()
	if err != nil {
		return
	}
	log.Infof("%#v", *userData)
	log.Infof("%#v", *networkStatus)
	log.Infof("%#v", *accessPoint)
	log.Infof("%#v", *filterData)

	terminate := make(chan os.Signal)
	signal.Notify(terminate, os.Kill, os.Interrupt)

	sender := ws.NewWSSender(conn, time.Second*10)
	receiver := ws.NewReceiver(conn)

	go sender.Run()
	go receiver.Run()

	receivedSignal := <-terminate
	log.Infof("Received Signal %s. Terminating...", receivedSignal)
	sender.Stop()
	receiver.Stop()
}
