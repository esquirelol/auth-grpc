package auth

import (
	"context"

	ssov1 "github.com/esquirelol/protos-grpc/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appID int64) (string, error)
	Register(ctx context.Context, email, password string) (int64, error)
}
type ServerAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func RegisterServerAPI(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &ServerAPI{auth: auth})
}

func (s *ServerAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required ")
	}
	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required ")
	}

	if req.GetAppId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "app_id is required ")
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), req.GetAppId())

	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func (s *ServerAPI) Register(ctx context.Context, request *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	if request.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if request.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	id, err := s.auth.Register(ctx, request.Email, request.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.RegisterResponse{UserId: id}, nil
}
