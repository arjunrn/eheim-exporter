package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	"github.com/arjunrn/eheim-exporter/pkg/ws"
)

func App(ctx context.Context, websocketURL string, metricsPort int) {
	log.SetLevel(log.DebugLevel)
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, websocketURL, nil)
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

	terminate := make(chan os.Signal, 2)
	signal.Notify(terminate, syscall.SIGTERM, os.Interrupt)

	sender := ws.NewWSSender(conn, time.Second*10, userData.From)
	receiver := ws.NewReceiver(conn)

	go sender.Run()
	go receiver.Run()

	promRegistry := prometheus.NewRegistry()
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", metricsPort),
		Handler: promhttp.HandlerFor(promRegistry, promhttp.HandlerOpts{}),
	}
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Errorf("failed to start metrics server: %v", err)
		}
	}()
	receivedSignal := <-terminate
	log.Infof("Received Signal %s. Terminating...", receivedSignal)
	sender.Stop()
	receiver.Stop()
	err = server.Shutdown(context.TODO())
	if err != nil {
		log.Warnf("failed to shutdown metrics server: %v", err)
	}
}
