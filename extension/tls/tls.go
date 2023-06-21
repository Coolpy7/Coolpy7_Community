package tls

import (
	"github.com/Coolpy7/Coolpy7_Community/mempool"
	"github.com/Coolpy7/Coolpy7_Community/pollio"
	"github.com/Coolpy7/Coolpy7_Community/std/crypto/tls"
)

// Conn .
type Conn = tls.Conn

// Config .
type Config = tls.Config

// Dial returns a net.Conn to be added to a Engine.
func Dial(network, addr string, config *Config) (*tls.Conn, error) {
	tlsConn, err := tls.Dial(network, addr, config, mempool.DefaultMemPool)
	if err != nil {
		return nil, err
	}

	return tlsConn, nil
}

func WrapOpen(tlsConfig *Config, isClient bool, h func(c *pollio.Conn, tlsConn *Conn)) func(c *pollio.Conn) {
	return func(c *pollio.Conn) {
		var tlsConn *tls.Conn
		sesseion := c.Session()
		if sesseion != nil {
			tlsConn = sesseion.(*tls.Conn)
		}
		if tlsConn == nil && !isClient {
			tlsConn = tls.NewConn(c, tlsConfig, isClient, true, mempool.DefaultMemPool)
			c.SetSession(tlsConn)
		}
		if h != nil {
			h(c, tlsConn)
		}
	}
}

func WrapClose(h func(c *pollio.Conn, tlsConn *Conn, err error)) func(c *pollio.Conn, err error) {
	return func(c *pollio.Conn, err error) {
		if h != nil && c != nil {
			if session := c.Session(); session != nil {
				if tlsConn, ok := session.(*Conn); ok {
					h(c, tlsConn, err)
				}
			}
		}
	}
}

func WrapData(h func(c *pollio.Conn, tlsConn *Conn, data []byte), args ...interface{}) func(c *pollio.Conn, data []byte) {
	getBuffer := func() []byte {
		return make([]byte, 2048)
	}
	if len(args) > 0 {
		if bh, ok := args[0].(func() []byte); ok {
			getBuffer = bh
		}
	}
	return func(c *pollio.Conn, data []byte) {
		if session := c.Session(); session != nil {
			if tlsConn, ok := session.(*Conn); ok {
				tlsConn.Append(data)
				buffer := getBuffer()
				for {
					n, err := tlsConn.Read(buffer)
					if err != nil {
						c.Close()
						return
					}
					if h != nil && n > 0 {
						h(c, tlsConn, buffer[:n])
					}
					if n == 0 {
						return
					}
				}
			}
		}
	}
}
