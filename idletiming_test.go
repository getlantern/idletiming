package idletiming

import (
	"io"
	"io/ioutil"
	"net"
	"testing"
	"time"
)

var (
	msg = []byte("HelloThere")

	dataLoops = 10

	clientTimeout                 = 25 * time.Millisecond
	serverTimeout                 = 2 * clientTimeout
	slightlyLessThanClientTimeout = time.Duration(int64(float64(clientTimeout.Nanoseconds()) * 0.9))
	slightlyMoreThanClientTimeout = time.Duration(int64(float64(clientTimeout.Nanoseconds()) * 1.1))
)

func TestWrite(t *testing.T) {
	listenerIdled := false
	connIdled := false

	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Unable to listen: %s", err)
	}

	addr := l.Addr().String()
	il := Listener(l, serverTimeout, func() {
		listenerIdled = true
	})

	go func() {
		conn, err := il.Accept()
		if err != nil {
			t.Fatalf("Unable to accept: %s", err)
		}
		go func() {
			// Discard data
			io.Copy(ioutil.Discard, conn)
		}()
	}()

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("Unable to dial %s: %s", addr, err)
	}

	c := Conn(conn, clientTimeout, func() {
		connIdled = true
	})

	// Write messages
	for i := 0; i < dataLoops; i++ {
		n, err := c.Write(msg)
		if err != nil || n != len(msg) {
			t.Fatalf("Problem writing.  n: %d  err: %s", n, err)
		}
		time.Sleep(slightlyLessThanClientTimeout)
	}

	// Now write msg with a really short deadline
	c.SetWriteDeadline(time.Now().Add(1 * time.Nanosecond))
	_, err = c.Write(msg)
	if netErr, ok := err.(net.Error); ok {
		if !netErr.Timeout() {
			t.Fatalf("Short deadline should have resulted in Timeout, but didn't: %s", err)
		}
	} else {
		t.Fatalf("Short deadline should have resulted in Timeout, but didn't: %s", err)
	}

	time.Sleep(slightlyMoreThanClientTimeout)
	if connIdled == false {
		t.Errorf("Conn failed to idle!")
	}

	connTimesOutIn := c.TimesOutIn()
	if connTimesOutIn > 0 {
		t.Errorf("TimesOutIn returned bad value, should have been negative, but was: %s", connTimesOutIn)
	}

	time.Sleep(slightlyMoreThanClientTimeout)
	if listenerIdled == false {
		t.Errorf("Listener failed to idle!")
	}
}

func TestRead(t *testing.T) {
	listenerIdled := false
	connIdled := false

	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Unable to listen: %s", err)
	}

	il := Listener(l, serverTimeout, func() {
		listenerIdled = true
	})

	addr := l.Addr().String()

	go func() {
		conn, err := il.Accept()
		if err != nil {
			t.Fatalf("Unable to accept: %s", err)
		}
		go func() {
			// Feed data
			for i := 0; i < dataLoops; i++ {
				_, err := conn.Write(msg)
				if err != nil {
					return
				}
				time.Sleep(slightlyLessThanClientTimeout)
			}
		}()
	}()

	c, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("Unable to dial %s: %s", addr, err)
	}

	c = Conn(c, clientTimeout, func() {
		connIdled = true
	})

	// Read messages
	b := make([]byte, 1024)
	totalN := 0
	for i := 0; i < dataLoops; i++ {
		n, err := c.Read(b)
		if err != nil {
			t.Fatalf("Problem reading.  err: %s", n, err)
		}
		totalN += n
		time.Sleep(slightlyLessThanClientTimeout)
	}

	if totalN == 0 {
		t.Fatal("Didn't read any data!")
	}

	// Now read with a really short deadline
	c.SetReadDeadline(time.Now().Add(1 * time.Nanosecond))
	_, err = c.Read(msg)
	if netErr, ok := err.(net.Error); ok {
		if !netErr.Timeout() {
			t.Fatalf("Short deadline should have resulted in Timeout, but didn't: %s", err)
		}
	} else {
		t.Fatalf("Short deadline should have resulted in Timeout, but didn't: %s", err)
	}

	time.Sleep(slightlyMoreThanClientTimeout)
	if connIdled == false {
		t.Errorf("Conn failed to idle!")
	}

	time.Sleep(slightlyMoreThanClientTimeout)
	if listenerIdled == false {
		t.Errorf("Listener failed to idle!")
	}
}
