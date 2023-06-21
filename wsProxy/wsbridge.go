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
	"github.com/Coolpy7/Coolpy7_Community/iohttp/websocket"
	"github.com/Coolpy7/Coolpy7_Community/pollio"
	"net"
	"time"
)

type WsBridge struct {
	wsConn  *websocket.Conn
	cpConn  *pollio.Conn
	cpEngin *pollio.Engine
}

func NewWsbridge(conn *websocket.Conn) (*WsBridge, error) {
	wp := &WsBridge{
		wsConn: conn,
	}

	wp.cpEngin = pollio.NewGopher(pollio.Config{})
	wp.cpEngin.OnData(wp.OnRequest)
	wp.cpEngin.OnClose(wp.OnDisconnect)
	err := wp.cpEngin.Start()
	if err != nil {
		return nil, err
	}
	c, err := net.DialTimeout("unix", SockAddr, 2*time.Second)
	if err != nil {
		return nil, err
	}
	nbc, err := wp.cpEngin.AddConn(c)
	if err != nil {
		return nil, err
	}
	wp.cpConn = nbc
	return wp, nil
}

func (u *WsBridge) OnDisconnect(conn *pollio.Conn, err error) {
	_ = u.wsConn.Close()
}

func (u *WsBridge) OnRequest(conn *pollio.Conn, data []byte) {
	err := u.wsConn.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		u.Close()
	}
}

func (u *WsBridge) Receive(data []byte) error {
	_, err := u.cpConn.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (u *WsBridge) Close() {
	if u.cpConn != nil {
		u.cpEngin.Stop()
		_ = u.cpConn.Close()
	}
	_ = u.wsConn.Close()
}
