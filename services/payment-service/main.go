package main

import (
"context"
"fmt"
"log"
"net"

"google.golang.org/grpc"
)

func main() {
fmt.Println(" - Starting payment-service service...")

// TODO: Initialize config, DI, repositories

// Start gRPC server (example)
go func() {
lis, err := net.Listen("tcp", ":9090")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
s := grpc.NewServer()
// TODO: register gRPC services here, e.g. pb.RegisterYourServiceServer(s, &handler.YourService{})
		if err := s.Serve(lis); err != nil {
			log.Fatalf("gRPC server failed: %v", err)
		}
}()

// Placeholder: keep process alive
select {}
}
