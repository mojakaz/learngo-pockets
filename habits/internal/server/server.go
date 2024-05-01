package server

import (
	"fmt"
	"google.golang.org/grpc"
	"learngo-pockets/habits/api"
	"net"
	"strconv"
)

// Server is the implementation of the gRPC server.
type Server struct {
	api.UnimplementedHabitsServer
	lgr Logger
}

// New returns a Server that can ListenAndServe.
func New(lgr Logger) *Server {
	return &Server{
		lgr: lgr,
	}
}

type Logger interface {
	Logf(format string, args ...any)
}

// ListenAndServe starts listening to the port and serving requests.
func (s *Server) ListenAndServe(port int) error {
	const addr = "127.0.0.1"

	listener, err := net.Listen("tcp", net.JoinHostPort(addr, strconv.Itoa(port)))
	if err != nil {
		return fmt.Errorf("unable to listen to tcp port %d: %w", port, err)
	}

	grpcServer := grpc.NewServer()
	api.RegisterHabitsServer(grpcServer, s)

	s.lgr.Logf("starting server on port %d\n", port)

	err = grpcServer.Serve(listener)
	if err != nil {
		return fmt.Errorf("error while listening: %w", err)
	}

	// Stop or GracefulStop was called, no reason to be alarmed.
	return nil
}