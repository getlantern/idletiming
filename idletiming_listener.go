package idletiming

import (
	"net"
	"time"
)

// Listener creates a net.Listener that wraps the connections obtained from an
// original net.Listener with idle timing connections that time out after the
// specified duration.
//
// idleTimeout specifies how long to wait for inactivity before considering
// connection idle.  Note - the actual timeout may be up to twice idleTimeout,
// depending on timing.
//
// onIdle is an optional function to call after the connection has been closed
func Listener(listener net.Listener, idleTimeout time.Duration, onIdle func()) net.Listener {
	return &idleTimingListener{listener, idleTimeout, onIdle}
}

type idleTimingListener struct {
	orig        net.Listener
	idleTimeout time.Duration
	onIdle      func()
}

func (l *idleTimingListener) Accept() (c net.Conn, err error) {
	c, err = l.orig.Accept()
	if err == nil {
		c = Conn(c, l.idleTimeout, l.onIdle)
	}
	return
}

func (l *idleTimingListener) Close() error {
	return l.orig.Close()
}

func (l *idleTimingListener) Addr() net.Addr {
	return l.orig.Addr()
}
