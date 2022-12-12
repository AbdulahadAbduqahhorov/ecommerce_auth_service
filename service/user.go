package service

import (
	"context"
	"fmt"
	"regexp"

	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/ecommerce_auth_service/config"
	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/ecommerce_auth_service/genproto/auth_service"
	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/ecommerce_auth_service/storage"
	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/ecommerce_auth_service/util"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authService struct {
	cfg config.Config
	auth_service.UnimplementedAuthServiceServer
	stg storage.StorageI
}

func NewAuthService(cfg config.Config, db *sqlx.DB) *authService {
	return &authService{
		cfg: cfg,
		stg: storage.NewStoragePg(db),
	}
}
func (s *authService) CreateUser(ctx context.Context, req *auth_service.CreateUserRequest) (*auth_service.User, error) {
	if len(req.Password) < 6 {
		err := fmt.Errorf("password must not be less than 6 characters")
		return nil, err
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	req.Password = hashedPassword

	phoneRegex := regexp.MustCompile(`^[+]?(\d{1,2})?[\s.-]?\(?\d{3}\)?[\s.-]?\d{3}[\s.-]?\d{4}$`)
	phone := phoneRegex.MatchString(req.Phone)
	if !phone {
		err = fmt.Errorf("phone number is not valid")
		return nil, err
	}

	pKey, err := s.stg.User().CreateUser(req)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return s.stg.User().GetUserById(pKey)

}

func (s *authService) GetUserList(ctx context.Context, req *auth_service.GetUserListRequest) (*auth_service.GetUserListResponse, error) {
	res, err := s.stg.User().GetUserList(req)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return res, err
}

func (s *authService) GetUserById(ctx context.Context, req *auth_service.GetUserByIdRequest) (*auth_service.User, error) {

	res, err := s.stg.User().GetUserById(req.Id)

	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return res, nil
}

func (s *authService) UpdateUserById(ctx context.Context, req *auth_service.UpdateUserRequest) (*auth_service.User, error) {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	email := emailRegex.MatchString(req.Email)
	if !email {
		err := fmt.Errorf("email is not valid")
		return nil, err
	}

	phoneRegex := regexp.MustCompile(`^[+]?(\d{1,2})?[\s.-]?\(?\d{3}\)?[\s.-]?\d{3}[\s.-]?\d{4}$`)
	phone := phoneRegex.MatchString(req.Phone)
	if !phone {
		err := fmt.Errorf("phone number is not valid")
		return nil, err
	}


	rowsAffected, err := s.stg.User().UpdateUser(req)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if rowsAffected <= 0 {
		return nil, status.Error(codes.InvalidArgument, "no rows were affected")
	}
	
	res, err := s.stg.User().GetUserById(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return res, err
}

func (s *authService) DeleteUser(ctx context.Context, req *auth_service.DeleteUserRequest) (*auth_service.Empty, error) {
	rowsAffected, err := s.stg.User().DeleteUser(req.Id)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if rowsAffected <= 0 {
		return nil, status.Error(codes.InvalidArgument, "no rows were affected")
	}

	return &auth_service.Empty{}, nil
}
