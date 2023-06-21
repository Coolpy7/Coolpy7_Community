package iohttp

import (
	"github.com/Coolpy7/Coolpy7_Community/std/crypto/tls"
	"net/http"
)

// Server .
type Server struct {
	*Engine
}

// NewServer .
func NewServer(conf Config, v ...interface{}) *Server {
	if len(v) > 0 {
		if handler, ok := v[0].(http.Handler); ok {
			conf.Handler = handler
		}
	}
	if len(v) > 1 {
		if messageHandlerExecutor, ok := v[1].(func(f func())); ok {
			conf.ServerExecutor = messageHandlerExecutor
		}
	}
	return &Server{Engine: NewEngine(conf)}
}

// NewServerTLS .
func NewServerTLS(conf Config, v ...interface{}) *Server {
	if len(v) > 0 {
		if handler, ok := v[0].(http.Handler); ok {
			conf.Handler = handler
		}
	}
	if len(v) > 1 {
		if messageHandlerExecutor, ok := v[1].(func(f func())); ok {
			conf.ServerExecutor = messageHandlerExecutor
		}
	}
	if len(v) > 2 {
		if tlsConfig, ok := v[2].(*tls.Config); ok {
			conf.TLSConfig = tlsConfig
		}
	}
	conf.AddrsTLS = append(conf.AddrsTLS, conf.Addrs...)
	conf.Addrs = nil
	return &Server{Engine: NewEngine(conf)}
}
