package authgrpc

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/mnogokotin/golang-grpc-auth/internal/grpc/auth/requests"
	"github.com/mnogokotin/golang-grpc-auth/internal/services/auth"
	"github.com/mnogokotin/golang-grpc-auth/internal/storage"
	pvalidator "github.com/mnogokotin/golang-packages/validator"
	ssov1 "github.com/mnogokotin/grpc-protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appID int64,
	) (token string, err error)
	Register(
		ctx context.Context,
		email string,
		password string,
		appID int64,
	) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

var validate *validator.Validate

func Register(grpcServer *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(grpcServer, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(
	ctx context.Context,
	req *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	loginRequest := &requests.LoginRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		AppID:    req.GetAppId(),
	}
	validate = validator.New()
	if err := validate.Struct(loginRequest); err != nil {
		errs := pvalidator.BuildValidationErrors(err)
		return nil, status.Error(codes.InvalidArgument, errs.Error())
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), req.GetAppId())
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}

		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &ssov1.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	req *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {
	registerRequest := &requests.RegisterRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		AppID:    req.GetAppId(),
	}
	validate = validator.New()
	if err := validate.Struct(registerRequest); err != nil {
		errs := pvalidator.BuildValidationErrors(err)
		return nil, status.Error(codes.InvalidArgument, errs.Error())
	}

	uid, err := s.auth.Register(ctx, req.GetEmail(), req.GetPassword(), req.GetAppId())
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "failed to register user")
	}

	return &ssov1.RegisterResponse{UserId: uid}, nil
}

func (s *serverAPI) IsAdmin(
	ctx context.Context,
	req *ssov1.IsAdminRequest,
) (*ssov1.IsAdminResponse, error) {
	registerRequest := &requests.IsAdminRequest{
		UserID: req.GetUserId(),
	}
	validate = validator.New()
	if err := validate.Struct(registerRequest); err != nil {
		errs := pvalidator.BuildValidationErrors(err)
		return nil, status.Error(codes.InvalidArgument, errs.Error())
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Error(codes.Internal, "failed to check admin status")
	}

	return &ssov1.IsAdminResponse{IsAdmin: isAdmin}, nil
}
