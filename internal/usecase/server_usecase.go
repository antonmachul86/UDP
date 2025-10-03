package usecase

import (
	"UDP/internal/entities"
	"UDP/internal/port"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type serverJob struct {
	data []byte
	addr net.Addr
}

type ServerUsecase struct {
	network    port.Network
	clock      port.Clock
	numWorkers int
	jobs       chan serverJob
	logChan    chan string
}

func NewServerUsecase(network port.Network, clock port.Clock, numWorkers int) *ServerUsecase {
	return &ServerUsecase{
		network:    network,
		clock:      clock,
		numWorkers: numWorkers,
		jobs:       make(chan serverJob, 10000),
		logChan:    make(chan string, 1024),
	}
}

func (s *ServerUsecase) Run() {
	go s.logger()

	for i := 0; i < s.numWorkers; i++ {
		go s.worker()
	}

	for {
		data, clientAddr, err := s.network.ReceiveFrom(0)
		if err != nil {
			continue
		}
		s.jobs <- serverJob{data: data, addr: clientAddr}
	}
}

func (s *ServerUsecase) logger() {
	for msg := range s.logChan {
		fmt.Println(msg)
	}
}

func (s *ServerUsecase) worker() {
	for job := range s.jobs {
		var pkt entities.Packet
		if err := json.Unmarshal(job.data, &pkt); err != nil {
			continue
		}

		recvTime := s.clock.Now().UTC()
		ok := pkt.IsValid()

		logMsg := fmt.Sprintf("Packet %d | Sent: %s | Received: %s | Valid: %t",
			pkt.ID, pkt.Timestamp.Format(time.RFC3339), recvTime.Format(time.RFC3339), ok)
		select {
		case s.logChan <- logMsg:
			// Sent
		default:
			// Drop log message if channel is full
		}

		ack := entities.Ack{
			ID:       pkt.ID,
			OK:       ok,
			RecvTime: recvTime,
		}
		ackBytes, _ := json.Marshal(ack)
		s.network.SendTo(job.addr.String(), ackBytes)
	}
}
