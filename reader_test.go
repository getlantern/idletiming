package idletiming

import (
	"io"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReader(t *testing.T) {
	testBytes := []byte("aaaa")

	l, err := net.Listen("tcp", ":0")
	if !assert.NoError(t, err) {
		return
	}
	defer l.Close()

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			go func() {
				for i := 0; i < len(testBytes); i++ {
					time.Sleep(100 * time.Millisecond)
					conn.Write(testBytes[i : i+1])
				}
				conn.Close()
			}()
		}
	}()

	conn, err := net.Dial("tcp", l.Addr().String())
	if !assert.NoError(t, err) {
		return
	}
	defer conn.Close()

	r := NewReader(conn, 50*time.Millisecond)
	b := make([]byte, len(testBytes))
	n, err := io.ReadFull(r, b)
	if !assert.NoError(t, err) {
		return
	}
	if assert.Equal(t, len(testBytes), n) {
		assert.EqualValues(t, testBytes, b)
	}
}
