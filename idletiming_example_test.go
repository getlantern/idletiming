package idletiming

import (
	"log"
	"net"
	"time"
)

func ExampleDial() {
	c, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatalf("Unable to dial %s", err)
	}

	ic := WithIdleTimeout(c, 5*time.Second, func() {
		log.Printf("Connection was idled")
	})

	ic.Write([]byte("My data"))
}

func ExampleListen() {
	l, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatalf("Unable to listen %s", err)
	}

	il := &IdleTimingListener{
		Orig:        l,
		IdleTimeout: 5 * time.Second,
		OnIdle: func() {
			log.Printf("Connection was idled")
		},
	}

	il.Accept()
}
