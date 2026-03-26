package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"storage/internal/config"
	"storage/internal/grpcserver"
)

func main() {
	cfg := config.Load()

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	grpcserver.Register(s)

	log.Println("Storage running on :" + cfg.Port)
	s.Serve(lis)
}
