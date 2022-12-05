//go:generate protoc ./../../../api/RotatorService.proto --go_out=./pb --go-grpc_out=./pb --proto_path=./../../../
package internalgrpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/a-klimenko/go-otus-final-project/internal/server/grpc/pb"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"net"
	"net/http"
)

type Server struct {
	app    Application
	logger Logger
	server *grpc.Server
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type Application interface {
	AddBanner(ctx context.Context, bannerId uuid.UUID, slotId uuid.UUID) error
	RemoveBanner(ctx context.Context, slotId uuid.UUID, bannerId uuid.UUID) error
	ClickBanner(ctx context.Context, slotId uuid.UUID, bannerId uuid.UUID, groupId uuid.UUID) error
	ChooseBanner(ctx context.Context, slotId uuid.UUID, groupId uuid.UUID) (*uuid.UUID, error)
}

type RotatorService struct {
	App    Application
	Logger Logger
	pb.UnimplementedRotatorServer
}

func NewServer(logger Logger, app Application) *Server {
	chainInterceptor := grpc.ChainUnaryInterceptor(
		loggingMiddleware(logger),
	)
	grpcServer := grpc.NewServer(chainInterceptor)

	service := &RotatorService{
		App:    app,
		Logger: logger,
	}
	pb.RegisterRotatorServer(grpcServer, service)

	srv := &Server{
		logger: logger,
		app:    app,
		server: grpcServer,
	}

	return srv
}

func (s *Server) Start() error {
	lsn, err := net.Listen("tcp", net.JoinHostPort("", "50051"))
	if err != nil {
		return err
	}

	if err := s.server.Serve(lsn); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Stop() error {
	s.logger.Info("grpc server is shutting down")
	s.server.GracefulStop()

	return nil
}

func (s *RotatorService) AddBanner(ctx context.Context, in *pb.AddBannerRequest) (*pb.AddBannerResponse, error) {
	bannerId, err := uuid.Parse(in.BannerID)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("can not parse banner id: %s", err))
		return nil, err
	}

	slotId, err := uuid.Parse(in.SlotID)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("can not parse slot id: %s", err))
		return nil, err
	}

	err = s.App.AddBanner(ctx, bannerId, slotId)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("error trying to add banner: %s", err))
		return nil, err
	}

	return &pb.AddBannerResponse{}, nil
}

func (s *RotatorService) RemoveBanner(ctx context.Context, in *pb.RemoveBannerRequest) (*pb.RemoveBannerResponse, error) {
	slotId, err := uuid.Parse(in.SlotID)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("can not parse slot id: %s", err))
		return nil, err
	}

	bannerId, err := uuid.Parse(in.BannerID)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("can not parse banner id: %s", err))
		return nil, err
	}

	err = s.App.RemoveBanner(ctx, slotId, bannerId)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("remove banner error: %s", err))
		return nil, err
	}

	return &pb.RemoveBannerResponse{}, nil
}

func (s *RotatorService) ClickBanner(ctx context.Context, in *pb.ClickBannerRequest) (*pb.ClickBannerResponse, error) {
	slotId, err := uuid.Parse(in.SlotID)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("can not parse slot id: %s", err))
		return nil, err
	}

	bannerId, err := uuid.Parse(in.BannerID)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("can not parse banner id: %s", err))
		return nil, err
	}

	groupId, err := uuid.Parse(in.GroupID)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("can not parse group id: %s", err))
		return nil, err
	}

	err = s.App.ClickBanner(ctx, slotId, bannerId, groupId)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("click banner error: %s", err))
		return nil, err
	}

	return &pb.ClickBannerResponse{}, nil
}

func (s *RotatorService) ChooseBanner(ctx context.Context, in *pb.ChooseBannerRequest) (*pb.ChooseBannerResponse, error) {
	slotId, err := uuid.Parse(in.SlotID)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("can not parse slot id: %s", err))
		return nil, err
	}

	groupId, err := uuid.Parse(in.GroupID)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("can not parse group id: %s", err))
		return nil, err
	}

	bannerId, err := s.App.ChooseBanner(ctx, slotId, groupId)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("can not choose banner: %s", err))
		return nil, err
	}

	return &pb.ChooseBannerResponse{BannerID: bannerId.String()}, nil
}
