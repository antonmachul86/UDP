package usecase

import (
	"UDP/internal/entities"
	"UDP/internal/port"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type ClientUsecase struct {
	network port.Network
	clock   port.Clock
	rand    port.Rand

	totalPackets uint32
	numWorkers   int

	packetId     uint32
	ackReceived  sync.Map
	printedId    uint32
	mu           sync.Mutex
	stopPrinting chan bool
}

func NewClientUsecase(
	network port.Network,
	clock port.Clock,
	rand port.Rand,
	totalPackets uint32,
	numWorkers int,
) *ClientUsecase {
	return &ClientUsecase{
		network:      network,
		clock:        clock,
		rand:         rand,
		totalPackets: totalPackets,
		numWorkers:   numWorkers,
		stopPrinting: make(chan bool, 1),
	}
}

func (c *ClientUsecase) Run() {
	var wg sync.WaitGroup

	go c.printAcksInOrder()

	for i := 0; i < c.numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.worker()
		}()
	}

	wg.Wait()
	time.Sleep(10 * time.Second)

	c.stopPrinting <- true
	c.printRemaining()
}
func (c *ClientUsecase) worker() {
	for {
		id := atomic.AddUint32(&c.packetId, 1)
		if id > c.totalPackets {
			return
		}

		ts := c.clock.Now().UTC()
		minLen := int(id)
		maxLen := 2 * int(id)
		dataLen := minLen + c.rand.Intn(maxLen-minLen+1)
		data := make([]byte, dataLen)
		c.rand.Read(data)

		pkt := &entities.Packet{
			ID:        id,
			Timestamp: ts,
			Data:      data,
		}
		pkt.Checksum = pkt.ComputeChecksum()

		data, _ = json.Marshal(pkt)
		c.network.SendTo("", data)
		time.Sleep(100 * time.Microsecond)
	}
}
func (c *ClientUsecase) HandleAck(data []byte) {
	var ack entities.Ack
	if err := json.Unmarshal(data, &ack); err != nil {
		return
	}
	c.ackReceived.Store(ack.ID, ack)
}

func (c *ClientUsecase) printAcksInOrder() {
	for {
		select {
		case <-c.stopPrinting:
			return
		default:
			c.mu.Lock()
			next := c.printedId + 1
			if ack, ok := c.ackReceived.Load(next); ok {
				c.printedId = next
				c.mu.Unlock()

				a := ack.(entities.Ack)
				status := "OK"
				if !a.OK {
					status = "CORRUPT"
				}
				fmt.Printf("Packet %d: %s\n", next, status)
			} else {
				c.mu.Unlock()
				time.Sleep(10 * time.Microsecond)
			}

			if c.printedId >= c.totalPackets {
				return
			}
		}
	}

}

func (c *ClientUsecase) printRemaining() {
	lostPackets := 0
	for i := uint32(1); i <= c.totalPackets; i++ {
		if _, ok := c.ackReceived.Load(i); !ok {
			fmt.Printf("Packet %d: LOST\n", i)
			lostPackets++
		}
	}

	lossPercentage := float64(lostPackets) / float64(c.totalPackets) * 100
	fmt.Printf("\n--- Summary ---\n")
	fmt.Printf("Total packets: %d\n", c.totalPackets)
	fmt.Printf("Lost packets: %d\n", lostPackets)
	fmt.Printf("Packet loss: %.2f%%\n", lossPercentage)
}
