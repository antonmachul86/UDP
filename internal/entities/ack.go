package entities

import "time"

type Ack struct {
	ID       uint32
	OK       bool
	RecvTime time.Time
}
