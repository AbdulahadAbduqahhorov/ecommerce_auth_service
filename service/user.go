package service

import (
	"context"
	"fmt"
	"regexp"

	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/ecommerce_auth_service/config"
	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/ecommerce_auth_service/genproto/auth_service"
	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/ecommerce_auth_service/pkg/logger"
	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/ecommerce_auth_service/pkg/util"
	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/ecommerce_auth_service/storage"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type authService struct {
	log logger.LoggerI
	cfg config.Config
	auth_service.UnimplementedAuthServiceServer
	stg storage.StorageI
}

func NewAuthService(log logger.LoggerI,cfg config.Config, db *sqlx.DB) *authService {
	return &authService{
		log:log,
		cfg: cfg,
		stg: storage.NewStoragePg(db),
	}
}
func (s *authService) CreateUser(ctx context.Context, req *auth_service.CreateUserRequest) (*auth_service.User, error) {
	s.log.Info("---CreateUser--->", logger.Any("req", req))
	if len(req.Password) < 6 {
		err := fmt.Errorf("password must not be less than 6 characters")
		s.log.Error("!!!CreateUser--->", logger.Error(err))
		return nil, err
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		s.log.Error("!!!CreateUser--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	req.Password = hashedPassword

	phoneRegex := regexp.MustCompile(`^[+]?(\d{1,2})?[\s.-]?\(?\d{3}\)?[\s.-]?\d{3}[\s.-]?\d{4}$`)
	phone := phoneRegex.MatchString(req.Phone)
	if !phone {

		err = fmt.Errorf("phone number is not valid")
		s.log.Error("!!!CreateUser--->", logger.Error(err))
		return nil, err
	}
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	email := emailRegex.MatchString(req.Email)
	if !email {
		err = fmt.Errorf("email is not valid")
		s.log.Error("!!!CreateUser--->", logger.Error(err))
		return nil, err
	}

	pKey, err := s.stg.User().CreateUser(req)

	if err != nil {
		s.log.Error("!!!CreateUser--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return s.stg.User().GetUserById(pKey)

}

func (s *authService) GetUserList(ctx context.Context, req *auth_service.GetUserListRequest) (*auth_service.GetUserListResponse, error) {
	s.log.Info("---GetUserList--->", logger.Any("req", req))
	res, err := s.stg.User().GetUserList(req)

	if err != nil {
		s.log.Error("!!!GetUserList--->", logger.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return res, err
}

func (s *authService) GetUserById(ctx context.Context, req *auth_service.GetUserByIdRequest) (*auth_service.User, error) {
	s.log.Info("---GetUserById--->", logger.Any("req", req))
	res, err := s.stg.User().GetUserById(req.Id)

	if err != nil {
		s.log.Error("!!!GetUserById--->", logger.Error(err))
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return res, nil
}

func (s *authService) UpdateUser(ctx context.Context, req *auth_service.UpdateUserRequest) (*auth_service.User, error) {
	s.log.Info("---UpdateUser--->", logger.Any("req", req))

	if len(req.Email) > 0 {
		emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
		email := emailRegex.MatchString(req.Email)
		if !email {
			err := fmt.Errorf("email is not valid")
			s.log.Error("!!!UpdateUser--->", logger.Error(err))
			return nil, err
		}
	}

	if len(req.Phone) > 0 {
		phoneRegex := regexp.MustCompile(`^[+]?(\d{1,2})?[\s.-]?\(?\d{3}\)?[\s.-]?\d{3}[\s.-]?\d{4}$`)
		phone := phoneRegex.MatchString(req.Phone)
		if !phone {
			err := fmt.Errorf("phone number is not valid")
			s.log.Error("!!!UpdateUser--->", logger.Error(err))
			return nil, err
		}
	}

	rowsAffected, err := s.stg.User().UpdateUser(req)

	if err != nil {
		s.log.Error("!!!UpdateUser--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if rowsAffected <= 0 {
		return nil, status.Error(codes.NotFound, "no rows were affected")
	}

	res, err := s.stg.User().GetUserById(req.Id)
	if err != nil {
		s.log.Error("!!!UpdateUser--->", logger.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return res, err
}

func (s *authService) DeleteUser(ctx context.Context, req *auth_service.DeleteUserRequest) (*emptypb.Empty, error) {
	s.log.Info("---DeleteUser--->", logger.Any("req", req))

	rowsAffected, err := s.stg.User().DeleteUser(req.Id)

	if err != nil {
		s.log.Error("!!!DeleteUser--->", logger.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	if rowsAffected <= 0 {
		return nil, status.Error(codes.NotFound, "no rows were affected")
	}

	return &emptypb.Empty{}, nil
}
