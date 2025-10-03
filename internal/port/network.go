package port

import (
	"net"
	"time"
)

type Network interface {
	SendTo(addr string, data []byte) error
	ReceiveFrom(timeout time.Duration) ([]byte, net.Addr, error)
	Close() error
	LocalAddr() string
	SetReadBufferSize(bytes int) error
}
