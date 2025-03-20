package exchangerates

import (
	"net"

	"github.com/Krchnk/currency-wallet-proto/exchangerates"
	"github.com/Krchnk/exchange-rates-service/pkg/config"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
	logger     *logrus.Logger
}

func NewServer(cfg *config.Config, service *Service) (*Server, error) {
	logger := logrus.New()
	logger.Info("Initializing gRPC server")

	grpcServer := grpc.NewServer()
	exchangerates.RegisterExchangeRatesServiceServer(grpcServer, service)

	return &Server{
		grpcServer: grpcServer,
		logger:     logger,
	}, nil
}

func (s *Server) Start(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		s.logger.WithError(err).Error("failed to listen")
		return err
	}

	s.logger.WithField("port", port).Info("gRPC server started")
	if err := s.grpcServer.Serve(lis); err != nil {
		s.logger.WithError(err).Error("failed to serve")
		return err
	}
	return nil
}

func (s *Server) Stop() {
	s.logger.Info("Stopping gRPC server")
	s.grpcServer.GracefulStop()
}
