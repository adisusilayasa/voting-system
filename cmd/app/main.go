package main

import (
	"belajar-blockchain/internal/server"
	"log"
)

const (
	ctxTimeout = 10
)

func main() {
	log.Println("Starting app server")

	s := server.NewServer()
	// s.RunServer()
	if err := s.RunServer(); err != nil {
		log.Fatal(err)
	}

}
