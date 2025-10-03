package main

import (
	"UDP/internal/adapter"
	"UDP/internal/usecase"
	"fmt"
)

func main() {
	network, err := adapter.NewUDPClient("127.0.0.1:9000")
	if err != nil {
		fmt.Printf("adapter.NewUDPClient: %v", err)
	}
	defer network.Close()

	clock := adapter.SystemClock{}
	rand := adapter.CryptoRand{}
	clientUsecase := usecase.NewClientUsecase(
		network,
		clock,
		rand,
		10000,
		10,
	)
	go func() {
		for {
			data, _, err := network.ReceiveFrom(0)
			if err != nil {
				continue
			}
			clientUsecase.HandleAck(data)
		}
	}()

	clientUsecase.Run()
}
