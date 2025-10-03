package entities

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"time"
)

type Packet struct {
	ID        uint32    `json:"id"`
	Timestamp time.Time `json:"timeStamp"`
	Data      []byte    `json:"data"`
	Checksum  string    `json:"checksum"`
}

func (p *Packet) ComputeChecksum() string {
	var buf bytes.Buffer

	binary.Write(&buf, binary.BigEndian, p.ID)

	binary.Write(&buf, binary.BigEndian, p.Timestamp.UnixNano())

	buf.Write(p.Data)

	h := sha256.Sum256(buf.Bytes())
	return hex.EncodeToString(h[:])
}

func (p *Packet) IsValid() bool {
	return p.ComputeChecksum() == p.Checksum
}
