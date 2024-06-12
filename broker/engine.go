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

package broker

import (
	"github.com/Coolpy7/Coolpy7_Community/ipblocker"
	"github.com/Coolpy7/Coolpy7_Community/multimap"
	"github.com/Coolpy7/Coolpy7_Community/std/crypto/tls"
	"github.com/Coolpy7/Coolpy7_Community/topic"
	"github.com/antlabs/timer"
	"sync"
)

type Engine struct {
	Subscribed   *topic.Tree
	Retained     *sync.Map
	Clients      *sync.Map
	Qos1Store    *multimap.MultiMap
	Qos2PubStore *multimap.MultiMap
	Qos2RelStore *multimap.MultiMap
	tm           timer.Timer
	JwtSecretKey [][]byte
	SelfDdosDeny int
	LazyMsgSend  int
	NilConnDeny  int
	BlockTime    int
	MaxAttempts  int
	IpBlocker    *ipblocker.Blocker // 增加Blocker字段
}

func NewEngine() *Engine {
	eng := Engine{
		SelfDdosDeny: 60,
		LazyMsgSend:  30,
		NilConnDeny:  2,
		BlockTime:    1800,
		MaxAttempts:  5,
		Subscribed:   topic.NewStandardTree(),
		Retained:     &sync.Map{},
		Clients:      &sync.Map{},
		JwtSecretKey: make([][]byte, 0),
		Qos1Store:    multimap.NewMultiMap(),
		Qos2PubStore: multimap.NewMultiMap(),
		Qos2RelStore: multimap.NewMultiMap(),
		tm:           timer.NewTimer(),
	}
	go eng.tm.Run()
	return &eng
}

func (e *Engine) InitIpBlocker() {
	e.IpBlocker = ipblocker.NewBlocker(e.BlockTime, e.MaxAttempts, "cp7")
}

func (e *Engine) AddClient(conn interface{}, tlsConn *tls.Conn) *Client {
	c := NewClient(conn, e)
	if tlsConn != nil {
		c.tlsConn = tlsConn
		c.isTls = true
	}
	e.Clients.Store(conn, c)
	return c
}

func (e *Engine) GetClient(conn interface{}) (*Client, bool) {
	c, ok := e.Clients.Load(conn)
	if ok {
		return c.(*Client), ok
	}
	return nil, ok
}

func (e *Engine) DelClient(conn interface{}) {
	e.Clients.Delete(conn)
}

func (e *Engine) Close() {
	e.tm.Stop()
	e.IpBlocker.Close()
}
