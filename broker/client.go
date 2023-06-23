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
	"bytes"
	"errors"
	"github.com/Coolpy7/Coolpy7_Community/packet"
	"github.com/Coolpy7/Coolpy7_Community/pollio"
	"github.com/Coolpy7/Coolpy7_Community/std/crypto/tls"
	"github.com/Coolpy7/Coolpy7_Community/topic"
	"github.com/antlabs/timer"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"time"
)

type Client struct {
	conn    *pollio.Conn
	tlsConn *tls.Conn
	isTls   bool
	resp    *packet.Pingresp

	ClientID     string
	cleanSession bool

	will *packet.Message

	Subscribed *topic.Tree

	eng   *Engine
	tm    timer.TimeNoder
	Cache bytes.Buffer

	isLogin  bool
	pingTime int64
}

func NewClient(conn *pollio.Conn, eng *Engine) *Client {
	c := &Client{
		pingTime:   0,
		isLogin:    false,
		isTls:      false,
		conn:       conn,
		eng:        eng,
		resp:       packet.NewPingresp(),
		Subscribed: topic.NewStandardTree(),
		Cache:      bytes.Buffer{},
	}
	c.tm = eng.tm.ScheduleFunc(30*time.Second, c.TickFunc)
	//客户端连接后2秒种内接收不到登陆包直接关闭连接
	time.AfterFunc(2*time.Second, func() {
		if !c.isLogin {
			_ = c.conn.Close()
		}
	})
	return c
}

func (c *Client) ClearWill() {
	if c.will != nil {
		c.will = nil
	}
}

func (c *Client) Wills() {
	if c.will == nil || c.will.QOS > 0 { //只允许qos为0的will
		c.will = nil
		return
	}
	publish := packet.NewPublish()
	publish.Message = *c.will
	c.will = nil
	cs := c.eng.Subscribed.Match(publish.Message.Topic)
	for _, sender := range cs {
		subClient := sender.(*Client)
		ss := c.Subscribed.Match(publish.Message.Topic)
		for _, s := range ss {
			subction := s.(packet.Subscription)
			if publish.Message.QOS > subction.QOS {
				publish.Message.QOS = subction.QOS
			}
		}
		err := subClient.send(publish)
		if err != nil {
			subClient.Clear()
		}
	}
}

func (c *Client) TickFunc() {
	if data, ok := c.eng.Qos1Store.GetBody(c.ClientID); ok {
		for item := data.Front(); item != nil; item = item.Next() {
			pub := item.Value.(*packet.Publish)
			pub.Dup = true
			err := c.send(pub)
			if err != nil {
				c.Clear()
			}
		}
		data = nil
	}
	if data, ok := c.eng.Qos2PubStore.GetBody(c.ClientID); ok {
		for item := data.Front(); item != nil; item = item.Next() {
			pub := item.Value.(*packet.Publish)
			pub.Dup = true
			err := c.send(pub)
			if err != nil {
				c.Clear()
			}
		}
		data = nil
	}
	if data, ok := c.eng.Qos2RelStore.GetBody(c.ClientID); ok {
		for item := data.Front(); item != nil; item = item.Next() {
			pub := item.Value.(*packet.Pubrel)
			err := c.send(pub)
			if err != nil {
				c.Clear()
			}
		}
		data = nil
	}
}

func (c *Client) Receive(pkt packet.Generic) (err error) {
	connect, ok := pkt.(*packet.Connect)
	if ok {
		err = c.processConnect(connect)
		if err != nil {
			return err
		}
		c.isLogin = true
		return nil
	}

	switch pkt.Type() {
	case packet.SUBSCRIBE:
		err = c.processSubscribe(pkt.(*packet.Subscribe))
	case packet.UNSUBSCRIBE:
		err = c.processUnsubscribe(pkt.(*packet.Unsubscribe))
	case packet.PUBLISH:
		err = c.processPublish(pkt.(*packet.Publish))
	case packet.PUBACK:
		err = c.processPuback(pkt.(*packet.Puback).ID)
	case packet.PUBCOMP:
		err = c.processPubcomp(pkt.(*packet.Pubcomp).ID)
	case packet.PUBREC:
		err = c.processPubrec(pkt.(*packet.Pubrec).ID)
	case packet.PUBREL:
		err = c.processPubrel(pkt.(*packet.Pubrel).ID)
	case packet.PINGREQ:
		err = c.processPingreq()
	case packet.DISCONNECT:
		err = c.processDisconnect()
	}

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) processConnect(pkt *packet.Connect) error {
	c.cleanSession = pkt.CleanSession
	c.ClientID = pkt.ClientID

	connack := packet.NewConnack()
	connack.ReturnCode = packet.ConnectionAccepted
	connack.SessionPresent = true

	if len(c.eng.JwtSecretKey) > 0 {
		valid := false
		if len(pkt.Password) > 0 {
			for _, jsk := range c.eng.JwtSecretKey {
				tk, err := jwt.Parse(pkt.Password, func(token *jwt.Token) (interface{}, error) {
					return jsk, nil
				})
				if err == nil && tk.Valid {
					valid = true
					break
				}
			}
		}
		if !valid {
			connack.ReturnCode = packet.NotAuthorized
			_ = c.send(connack)
			return errors.New("ErrNotAuthorized")
		}
	}

	if pkt.CleanSession {
		c.eng.Qos1Store.Clear(c.ClientID)
		c.eng.Qos2PubStore.Clear(c.ClientID)
		c.eng.Qos2RelStore.Clear(c.ClientID)
		connack.SessionPresent = false
	}

	if pkt.Will != nil {
		_, err := topic.Parse(pkt.Will.Topic, false)
		if err != nil {
			return err
		}
		if len(pkt.Will.Payload) == 0 {
			return errors.New("will payload")
		}
		c.will = pkt.Will
	}

	err := c.send(connack)
	if err != nil {
		return err
	}

	if !pkt.CleanSession {
		c.TickFunc()
	}

	return nil
}

func (c *Client) processPingreq() error {
	err := c.send(c.resp)
	if err != nil {
		return err
	}
	if c.pingTime == 0 {
		c.pingTime = time.Now().Unix()
	} else {
		if int(time.Now().Unix()-c.pingTime) < c.eng.SelfDdosDeny {
			return errors.New("self ddos")
		} else {
			c.pingTime = time.Now().Unix()
		}
	}
	return nil
}

func (c *Client) processSubscribe(pkt *packet.Subscribe) error {
	suback := packet.NewSuback()
	suback.ReturnCodes = make([]packet.QOS, len(pkt.Subscriptions))
	suback.ID = pkt.ID

	for i, subscription := range pkt.Subscriptions {
		stp, err := topic.Parse(subscription.Topic, true)
		if err != nil {
			return err
		}
		c.eng.Subscribed.Add(stp, c)        //全局
		c.Subscribed.Add(stp, subscription) //私有
		suback.ReturnCodes[i] = subscription.QOS
	}

	err := c.send(suback)
	if err != nil {
		return err
	}

	for _, subscription := range pkt.Subscriptions {
		//retain
		rt, ok := c.eng.Retained.Load(subscription.Topic)
		if ok {
			publish := rt.(*packet.Publish)
			if publish.Message.QOS > subscription.QOS {
				publish.Message.QOS = subscription.QOS
			}
			err = c.send(publish)
			if err != nil {
				c.Clear()
			}
		}
	}

	return nil
}

func (c *Client) processUnsubscribe(pkt *packet.Unsubscribe) error {
	unsuback := packet.NewUnsuback()
	unsuback.ID = pkt.ID
	for _, tp := range pkt.Topics {
		stp, err := topic.Parse(tp, true)
		if err != nil {
			return err
		}
		c.eng.Subscribed.Remove(stp, c)
		ss := c.Subscribed.Match(stp)
		for _, s := range ss {
			c.Subscribed.Clear(s)
		}
	}

	err := c.send(unsuback)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) processPublish(publish *packet.Publish) error {
	_, err := topic.Parse(publish.Message.Topic, false)
	if err != nil {
		return err
	}
	publish.Dup = false

	if publish.Message.Retain {
		if publish.Message.QOS == 0 {
			if len(publish.Message.Payload) == 0 {
				c.eng.Retained.Delete(publish.Message.Topic)
			} else {
				c.eng.Retained.Store(publish.Message.Topic, publish)
			}
		} else {
			publish.Message.Retain = false
		}
	}

	if publish.Message.QOS == 2 {
		c.eng.Qos2PubStore.Push(c.ClientID, publish)
		pubrec := packet.NewPubrec()
		pubrec.ID = publish.ID
		err = c.send(pubrec)
		if err != nil {
			return err
		}
		return nil
	}

	cs := c.eng.Subscribed.Match(publish.Message.Topic)
	for _, sender := range cs {
		subClient := sender.(*Client)
		ss := c.Subscribed.Match(publish.Message.Topic)
		for _, s := range ss {
			subction := s.(packet.Subscription)
			if publish.Message.QOS > subction.QOS {
				publish.Message.QOS = subction.QOS
			}
		}
		if publish.Message.QOS == 1 {
			c.eng.Qos1Store.Push(subClient.ClientID, publish)
		}
		err = subClient.send(publish)
		if err != nil {
			subClient.Clear()
		}
	}

	if publish.Message.QOS == 1 {
		puback := packet.NewPuback()
		puback.ID = publish.ID
		err = c.send(puback)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) processPuback(id packet.ID) error {
	if item, ok := c.eng.Qos1Store.Pop(c.ClientID); ok {
		if pkt, ok := item.Value.(*packet.Publish); ok {
			if pkt.ID == id {
				c.eng.Qos1Store.Del(c.ClientID, item)
				return nil
			}
		}
	}
	return errors.New("msgId")
}

func (c *Client) processPubrel(id packet.ID) error {
	item, ok := c.eng.Qos2PubStore.Pop(c.ClientID)
	if !ok {
		return errors.New("pubrec")
	}
	c.eng.Qos2PubStore.Del(c.ClientID, item)
	publish, _ := item.Value.(*packet.Publish)

	cs := c.eng.Subscribed.Match(publish.Message.Topic)
	for _, sender := range cs {
		subClient := sender.(*Client)
		ss := c.Subscribed.Match(publish.Message.Topic)
		for _, s := range ss {
			subction := s.(packet.Subscription)
			if publish.Message.QOS > subction.QOS {
				publish.Message.QOS = subction.QOS
			}
			if publish.Message.QOS == 1 {
				c.eng.Qos1Store.Push(subClient.ClientID, publish)
			}
		}
		err := subClient.send(publish)
		if err != nil {
			return err
		}
	}

	pubComp := packet.NewPubcomp()
	pubComp.ID = id
	err := c.send(pubComp)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) processPubrec(id packet.ID) error {
	pubRel := packet.NewPubrel()
	pubRel.ID = id
	c.eng.Qos2RelStore.Push(c.ClientID, pubRel)
	err := c.send(pubRel)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) processPubcomp(id packet.ID) error {
	if item, ok := c.eng.Qos2RelStore.Pop(c.ClientID); ok {
		if pkt, ok := item.Value.(*packet.Pubrel); ok {
			if pkt.ID == id {
				c.eng.Qos2RelStore.Del(c.ClientID, item)
				return nil
			}
		}
	}
	return errors.New("msgId")
}

func (c *Client) processDisconnect() error {
	c.ClearWill()
	//ws conn abort
	if strings.Contains(c.conn.RemoteAddr().String(), "@") {
		_ = c.conn.Close()
	}
	_ = c.conn.Close()
	return nil
}

func (c *Client) Clear() {
	c.tm.Stop()
	c.eng.Subscribed.Clear(c)
}

func (c *Client) Close() {
	_ = c.conn.Close()
}
