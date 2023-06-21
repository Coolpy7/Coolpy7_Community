//go:build linux
// +build linux

package pollio

import (
	"github.com/Coolpy7/Coolpy7_Community/timer"
	"net"
	"runtime"
)

// Start init and start pollers.
func (g *Engine) Start() error {
	udpListeners := make([]*net.UDPConn, len(g.addrs))[0:0]
	switch g.network {
	case "unix", "tcp", "tcp4", "tcp6":
		for i := range g.addrs {
			ln, err := newPoller(g, true, i)
			if err != nil {
				for j := 0; j < i; j++ {
					g.listeners[j].stop()
				}
				return err
			}
			g.addrs[i] = ln.listener.Addr().String()
			g.listeners = append(g.listeners, ln)
		}
	case "udp", "udp4", "udp6":
		for i, addrStr := range g.addrs {
			addr, err := net.ResolveUDPAddr(g.network, addrStr)
			if err != nil {
				for j := 0; j < i; j++ {
					udpListeners[j].Close()
				}
				return err
			}
			ln, err := g.listenUDP("udp", addr)
			if err != nil {
				for j := 0; j < i; j++ {
					udpListeners[j].Close()
				}
				return err
			}
			g.addrs[i] = ln.LocalAddr().String()
			udpListeners = append(udpListeners, ln)
		}
	}

	for i := 0; i < g.pollerNum; i++ {
		p, err := newPoller(g, false, i)
		if err != nil {
			for j := 0; j < len(g.listeners); j++ {
				g.listeners[j].stop()
			}

			for j := 0; j < i; j++ {
				g.pollers[j].stop()
			}
			return err
		}
		g.pollers[i] = p
	}

	for i := 0; i < g.pollerNum; i++ {
		g.pollers[i].ReadBuffer = make([]byte, g.readBufferSize)
		g.Add(1)
		go g.pollers[i].start()
	}

	for _, l := range g.listeners {
		g.Add(1)
		go l.start()
	}

	for _, ln := range udpListeners {
		_, err := g.AddConn(ln)
		if err != nil {
			for j := 0; j < len(g.listeners); j++ {
				g.listeners[j].stop()
			}

			for j := 0; j < len(g.pollers); j++ {
				g.pollers[j].stop()
			}

			for j := 0; j < len(udpListeners); j++ {
				udpListeners[j].Close()
			}

			return err
		}
	}

	g.Timer.Start()
	return nil
}

// NewEngine is a factory impl.
func NewEngine(conf Config) *Engine {
	cpuNum := runtime.NumCPU()
	if conf.Name == "" {
		conf.Name = "EP"
	}
	if conf.NPoller <= 0 {
		conf.NPoller = cpuNum
	}
	if conf.ReadBufferSize <= 0 {
		conf.ReadBufferSize = DefaultReadBufferSize
	}
	if conf.MaxConnReadTimesPerEventLoop <= 0 {
		conf.MaxConnReadTimesPerEventLoop = DefaultMaxConnReadTimesPerEventLoop
	}
	if conf.Listen == nil {
		conf.Listen = net.Listen
	}
	if conf.ListenUDP == nil {
		conf.ListenUDP = net.ListenUDP
	}

	g := &Engine{
		Timer:                        timer.New(conf.Name, conf.TimerExecute),
		Name:                         conf.Name,
		network:                      conf.Network,
		addrs:                        conf.Addrs,
		listen:                       conf.Listen,
		listenUDP:                    conf.ListenUDP,
		pollerNum:                    conf.NPoller,
		readBufferSize:               conf.ReadBufferSize,
		maxWriteBufferSize:           conf.MaxWriteBufferSize,
		maxConnReadTimesPerEventLoop: conf.MaxConnReadTimesPerEventLoop,
		udpReadTimeout:               conf.UDPReadTimeout,
		epollMod:                     conf.EpollMod,
		epollOneshot:                 conf.EPOLLONESHOT,
		lockListener:                 conf.LockListener,
		lockPoller:                   conf.LockPoller,
		listeners:                    make([]*poller, len(conf.Addrs))[0:0],
		pollers:                      make([]*poller, conf.NPoller),
		connsUnix:                    make([]*Conn, MaxOpenFiles),
	}

	g.initHandlers()

	return g
}
