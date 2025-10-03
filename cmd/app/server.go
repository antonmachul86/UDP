package main

import (
	"UDP/internal/adapter"
	"UDP/internal/usecase"
	"log"
	"runtime"
)

func main() {
	network, err := adapter.NewUDPServer(":9000")
	if err != nil {
		log.Fatal(err)
	}
	defer network.Close()

	if err := network.SetReadBufferSize(25 * 1024 * 1024); err != nil {
		log.Printf("Warning: Failed to set read buffer size: %v", err)
	}

	clock := adapter.SystemClock{}

	usecase := usecase.NewServerUsecase(network, clock, runtime.NumCPU())
	usecase.Run()
}
