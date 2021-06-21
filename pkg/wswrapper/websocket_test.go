package wswrapper

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type mockWSHandler struct {
	t            *testing.T
	requestCount int
}

func (t *mockWSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.requestCount++
	conn, err := upgrader.Upgrade(w, r, nil)
	require.NoError(t.t, err)
	err = conn.WriteMessage(websocket.TextMessage, []byte("1"))
	require.NoError(t.t, err)
	err = conn.WriteMessage(websocket.CloseMessage, nil)
	require.NoError(t.t, err)
	err = conn.Close()
	require.NoError(t.t, err)
}

func makeTestWSServer(t *testing.T) {
	t.Helper()
	mockHandler := &mockWSHandler{t: t}
	server := httptest.NewServer(mockHandler)
	t.Log(server.URL)
	defer server.Close()
	wsURL := fmt.Sprintf("ws%s", strings.TrimPrefix(server.URL, "http"))
	client := NewReconnectingWebsocket(wsURL, time.Second*10)
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	go func() {
		t.Log("before run")
		err := client.Run(ctx)
		t.Log(err)
		require.NoError(t, err)
	}()
	for i := 0; i < 3; i++ {
		val := <-client.ReceiverChan()
		require.Equal(t, "1", string(val))
	}
	require.Equal(t, 3, mockHandler.requestCount)
}

func TestBasic(t *testing.T) {
	makeTestWSServer(t)
}
