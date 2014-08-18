package idletiming

import (
	"net"
	"testing"
	"time"
)

func Test(t *testing.T) {
	addr := "localhost:14000"
	timeout := 1 * time.Second
	listenerIdled := false
	connIdled := false

	l, err := net.Listen("tcp", addr)
	if err != nil {
		t.Fatalf("Unable to listen at %s: %s", addr, err)
	}

	il := Listener(l, timeout, func() {
		listenerIdled = true
	})

	go func() {
		il.Accept()
	}()

	c, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("Unable to dial %s: %s", addr, err)
	}

	Conn(c, timeout, func() {
		connIdled = true
	})

	time.Sleep(timeout * 2)
	if listenerIdled == false {
		t.Errorf("Listener failed to idle!")
	}
	if connIdled == false {
		t.Errorf("Conn failed to idle!")
	}
}
