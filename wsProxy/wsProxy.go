// Copyright 2022 LiDonghai Email:5241871@qq.com Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package wsProxy

import (
	"fmt"
	"github.com/Coolpy7/Coolpy7_Community/iohttp"
	"github.com/Coolpy7/Coolpy7_Community/iohttp/websocket"
	"github.com/Coolpy7/Coolpy7_Community/std/crypto/tls"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

const SockAddr = "/tmp/cp9.sock"

type WsProxy struct {
	Clients  *sync.Map
	Srv      *iohttp.Server
	upgrader *websocket.Upgrader
}

func NewWsProxy() *WsProxy {
	ws := WsProxy{
		Clients: &sync.Map{},
	}
	upgrader := websocket.NewUpgrader()
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	upgrader.Subprotocols = []string{"mqtt", "mqttv3.1"}
	upgrader.OnOpen(ws.OnConnect)
	upgrader.OnMessage(ws.OnMessage)
	upgrader.OnClose(ws.OnDisconnect)
	ws.upgrader = upgrader
	return &ws
}

func (w *WsProxy) AddClient(conn *websocket.Conn) (*WsBridge, error) {
	wp, err := NewWsbridge(conn)
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	w.Clients.Store(conn, wp)
	return wp, nil
}

func (w *WsProxy) GetClient(conn *websocket.Conn) (*WsBridge, bool) {
	c, ok := w.Clients.Load(conn)
	if ok {
		return c.(*WsBridge), ok
	}
	return nil, ok
}

func (w *WsProxy) DelClient(conn *websocket.Conn) {
	w.Clients.Delete(conn)
}

func (w *WsProxy) OnConnect(conn *websocket.Conn) {
	wsb, err := w.AddClient(conn)
	if err != nil {
		wsb.Close()
	}
}

func (w *WsProxy) OnDisconnect(conn *websocket.Conn, err error) {
	wsb, ok := w.GetClient(conn)
	if ok {
		w.DelClient(conn)
		wsb.Close()
	}
}

func (w *WsProxy) OnMessage(conn *websocket.Conn, messageType websocket.MessageType, data []byte) {
	_ = conn.SetDeadline(time.Time{})
	wsb, ok := w.GetClient(conn)
	if ok {
		err := wsb.Receive(data)
		if err != nil {
			wsb.Close()
		}
	}
}

func (ws *WsProxy) onWebsocket(w http.ResponseWriter, r *http.Request) {
	_, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (ws *WsProxy) Start(addr, pem, key string) error {
	mux := &http.ServeMux{}
	mux.HandleFunc("/", ws.onWebsocket)

	fcg := iohttp.Config{
		Network:                 "tcp",
		Addrs:                   []string{addr},
		MaxLoad:                 1000000,
		ReleaseWebsocketPayload: true,
		Handler:                 mux,
		ReadBufferSize:          1024 * 4,
		IOMod:                   iohttp.IOModMixed,
		MaxBlockingOnline:       100000,
	}

	if pem != "" && key != "" {
		btsKey, err := os.ReadFile(key)
		if err != nil {
			log.Fatal(err)
		}
		btsPem, err := os.ReadFile(pem)
		if err != nil {
			log.Fatal(err)
		}
		cert, err := tls.X509KeyPair(btsPem, btsKey)
		if err != nil {
			log.Fatalf("tls.X509KeyPair failed: %v", err)
		}
		tlsConfig := &tls.Config{
			Certificates:       []tls.Certificate{cert},
			InsecureSkipVerify: true,
		}
		fcg.Addrs = nil
		fcg.AddrsTLS = []string{addr}
		fcg.TLSConfig = tlsConfig
	}

	ws.Srv = iohttp.NewServer(fcg)

	err := ws.Srv.Start()
	if err != nil {
		fmt.Printf("Start failed: %v\n", err)
		return err
	}
	log.Println("Coolpy7 Community Websocket On Port", addr)
	return nil
}
