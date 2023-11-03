package rpc

import "google.golang.org/grpc"

type ServiceMetadata struct {
	ServiceName string
}

type Service interface {
	Initialization()
	BindEngine(conn *grpc.Server)
}
