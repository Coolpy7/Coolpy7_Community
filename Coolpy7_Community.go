// Copyright 2022 LiDonghai Email:5241871@qq.com Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"flag"
	"fmt"
	"github.com/Coolpy7/Coolpy7_Community/broker"
	ntls "github.com/Coolpy7/Coolpy7_Community/extension/tls"
	"github.com/Coolpy7/Coolpy7_Community/packet"
	"github.com/Coolpy7/Coolpy7_Community/pollio"
	"github.com/Coolpy7/Coolpy7_Community/std/crypto/tls"
	"github.com/Coolpy7/Coolpy7_Community/wsProxy"
	"golang.org/x/sys/unix"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var version = "1.0.1"
var goversion = "1.19.4"

var eng *broker.Engine

func main() {
	setLimit()
	var (
		addr         = flag.String("l", ":1883", "host port (default 1883)")
		pem          = flag.String("p", "", "tls pem file path")
		key          = flag.String("k", "", "tls key file path")
		wsAddr       = flag.String("w", ":8083", "host ws port (default 8083)")
		wsPem        = flag.String("wp", "", "wss pem file path")
		wsKey        = flag.String("wk", "", "wss key file path")
		jwtSecretKey = flag.String("j", "", "jwt secret key(multiple split by ,)")
		SelfDdos     = flag.Int("sd", 59, "self ddos deny seconds")
		LazyMsg      = flag.Int("lm", 30, "lazy qos message resend seconds")
		NilConn      = flag.Int("nc", 2, "nil connect deny seconds")
	)
	flag.Parse()

	eng = broker.NewEngine()
	if *jwtSecretKey != "" {
		jsks := strings.Split(*jwtSecretKey, ",")
		if len(jsks) > 0 {
			for _, jsk := range jsks {
				eng.JwtSecretKey = append(eng.JwtSecretKey, []byte(jsk))
			}
		}
	}
	eng.SelfDdosDeny = *SelfDdos
	eng.LazyMsgSend = *LazyMsg
	eng.NilConnDeny = *NilConn

	if err := os.RemoveAll(wsProxy.SockAddr); err != nil {
		log.Fatal(err)
	}
	pollUnix := pollio.NewEngine(pollio.Config{
		Network:      "unix", //"udp", "unix"
		Addrs:        []string{wsProxy.SockAddr},
		EPOLLONESHOT: unix.EPOLLONESHOT,
		EpollMod:     unix.EPOLLET,
	})
	// hanlde new connection
	pollUnix.OnOpen(OnConnect)
	// hanlde connection closed
	pollUnix.OnClose(OnDisconnect)
	// handle data
	pollUnix.OnData(OnMessage)

	poll := pollio.NewEngine(pollio.Config{
		Network:      "tcp", //"udp", "unix"
		Addrs:        []string{*addr},
		EPOLLONESHOT: unix.EPOLLONESHOT,
		EpollMod:     unix.EPOLLET,
	})

	if *pem != "" && *key != "" {
		btsKey, err := os.ReadFile(*key)
		if err != nil {
			log.Fatal(err)
		}
		btsPem, err := os.ReadFile(*pem)
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
		tlsConfig.BuildNameToCertificate()
		poll.OnOpen(ntls.WrapOpen(tlsConfig, false, func(conn *pollio.Conn, tlsConn *tls.Conn) {
			OnTlsConnect(conn, tlsConn)
		}))
		poll.OnClose(ntls.WrapClose(func(conn *pollio.Conn, tlsConn *tls.Conn, err error) {
			OnDisconnect(conn, err)
		}))
		poll.OnData(ntls.WrapData(func(conn *pollio.Conn, tlsConn *tls.Conn, data []byte) {
			OnMessage(conn, data)
		}))
	} else {
		// hanlde new connection
		poll.OnOpen(OnConnect)
		// hanlde connection closed
		poll.OnClose(OnDisconnect)
		// handle data
		poll.OnData(OnMessage)
	}

	go func() {
		if err := poll.Start(); err != nil {
			log.Println(err)
		}
	}()
	go func() {
		if err := pollUnix.Start(); err != nil {
			log.Println(err)
		}
	}()
	log.Println("Coolpy7 Community On Port", *addr)

	wsp := wsProxy.NewWsProxy()
	err := wsp.Start(*wsAddr, *wsPem, *wsKey)
	if err != nil {
		log.Printf("Coolpy7 ws host error %s", err)
	}
	log.Printf("Coolpy7 v%s build golang v%s", version, goversion)

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for range signalChan {
			log.Println("Waiting For Coolpy7 Community Close...")
			wsp.Srv.Stop()
			poll.Stop()
			pollUnix.Stop()
			eng.Close()
			fmt.Println("safe exit")
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}

func OnConnect(conn *pollio.Conn) {
	eng.AddClient(conn, nil)
}

func OnTlsConnect(conn *pollio.Conn, tlsConn *tls.Conn) {
	eng.AddClient(conn, tlsConn)
}

func OnDisconnect(conn *pollio.Conn, err error) {
	c, ok := eng.GetClient(conn)
	if ok {
		go c.Wills()
		c.Clear()
		eng.DelClient(conn)
	}
}

func OnMessage(conn *pollio.Conn, data []byte) {
	c, ok := eng.GetClient(conn)
	if !ok {
		_ = conn.Close()
		return
	}
	c.Cache.Write(data)
	packetLength, packetType := packet.DetectPacket(c.Cache.Bytes())
	if packetLength <= 0 {
		return
	}
	if uint64(packetLength) > packet.MaxVarint {
		c.Close()
		return
	}
	if c.Cache.Len() < packetLength {
		return
	}
	defer c.Cache.Truncate(0)
	pkt, err := packetType.New()
	if err != nil {
		c.Close()
		return
	}
	_, err = pkt.Decode(c.Cache.Bytes())
	if err != nil {
		c.Close()
		return
	}
	err = c.Receive(pkt)
	if err != nil {
		c.Close()
	}
}

func setLimit() {
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		log.Fatal(err)
	}
	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		log.Fatal(err)
	}
}
