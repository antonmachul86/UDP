package entities

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"
)

type Packet struct {
	ID        uint32    `json:"id"`
	Timestamp time.Time `json:"timeStamp"`
	Data      []byte    `json:"data"`
	Checksum  string    `json:"checksum"`
}

func (p *Packet) ComputeChecksum() string {
	tmp := p.Checksum
	p.Checksum = ""
	data, _ := json.Marshal(p)
	p.Checksum = tmp

	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

func (p *Packet) IsValid() bool {
	return p.ComputeChecksum() == p.Checksum
}
