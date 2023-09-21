package api

import (
	"fmt"
	"log"
	"net"
	"stream-service/pkg/config"
	"stream-service/pkg/pb"

	"google.golang.org/grpc"
)

type Server struct {
	lis  net.Listener
	gsr  *grpc.Server
	port string
}

func NewServerGRPC(cfg config.Config, srv pb.StreamServiceServer) (*Server, error) {

	addr := fmt.Sprintf("%s:%s", cfg.StreamServiceHost, cfg.StreamServicePort)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %v", err)
	}

	gsr := grpc.NewServer()

	pb.RegisterStreamServiceServer(gsr, srv)

	return &Server{
		lis:  lis,
		gsr:  gsr,
		port: cfg.StreamServicePort,
	}, err
}

func (c *Server) Start() error {
	log.Println("Stream service listening on port: ", c.port)
	return c.gsr.Serve(c.lis)
}
