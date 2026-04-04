package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"storage/internal/config"
)

func main() {
	cfg := config.Load()

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Storage running on :" + cfg.Port)
	s.Serve(lis)
}
