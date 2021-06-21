package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/arjunrn/eheim-exporter/pkg/metrics"
	"github.com/arjunrn/eheim-exporter/pkg/processor"
	"github.com/arjunrn/eheim-exporter/pkg/wswrapper"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func App(websocketURL string, metricsPort int, refreshInterval time.Duration, debug bool) {
	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	promRegistry := prometheus.NewRegistry()
	filterMetrics := metrics.NewFilterMetrics(promRegistry)
	tracker := processor.NewFilterIDTracker()
	ws := wswrapper.NewReconnectingWebsocket(websocketURL, refreshInterval)
	receiver := processor.NewReceiver(ws.ReceiverChan(), filterMetrics, tracker)
	sender := processor.NewSender(ws, refreshInterval, tracker)

	terminate := make(chan os.Signal, 2)
	signal.Notify(terminate, syscall.SIGTERM, os.Interrupt)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := ws.Run(ctx)
		if err != nil {
			log.Errorf("websocket failed to connect: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		sender.Run(ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := receiver.Run(ctx)
		if err != nil {
			log.Errorf("failed to run receiver: %v", err)
		}
	}()

	promHandler := promhttp.HandlerFor(promRegistry, promhttp.HandlerOpts{})
	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", metricsPort),
		Handler: promHandler,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := server.ListenAndServe()
		if err != nil {
			log.Errorf("failed to start metrics server: %v", err)
		}
	}()

	receivedSignal := <-terminate
	log.Infof("Received Signal %s. Terminating...", receivedSignal)
	cancelFunc()
	if err := server.Shutdown(ctx); err != nil {
		log.Errorf("failed to shutdown metrics server: %v", err)
	}
	wg.Wait()
}
