package rpc

import "google.golang.org/grpc"

type Service interface {
	Initialization()
	BindEngine(conn *grpc.Server)
}
