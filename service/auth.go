package service

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/ecommerce_auth_service/genproto/auth_service"
	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/ecommerce_auth_service/pkg/logger"
	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/ecommerce_auth_service/pkg/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *authService) Register(ctx context.Context, req *auth_service.RegisterUserRequest) (*auth_service.User, error) {
	s.log.Info("---Register--->", logger.Any("req", req))
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		s.log.Error("!!!Register--->", logger.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "util.HashPassword() : %v", err)
	}
	if len(req.Password) < 6 {
		err := fmt.Errorf("password must not be less than 6 characters")
		s.log.Error("!!!Register--->", logger.Error(err))
		return nil, err
	}
	req.Password = hashedPassword
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	email := emailRegex.MatchString(req.Email)
	if !email {
		err = fmt.Errorf("email is not valid")
		s.log.Error("!!!Register--->", logger.Error(err))
		return nil, err
	}

	phoneRegex := regexp.MustCompile(`^[+]?(\d{1,2})?[\s.-]?\(?\d{3}\)?[\s.-]?\d{3}[\s.-]?\d{4}$`)
	phone := phoneRegex.MatchString(req.Phone)
	if !phone {
		err = fmt.Errorf("phone number is not valid")
		s.log.Error("!!!Register--->", logger.Error(err))
		return nil, err
	}
	id, err := s.stg.User().Register(req)
	if err != nil {
		s.log.Error("!!!Register--->", logger.Error(err))
		return nil, status.Errorf(codes.Internal, "method Register: %v", err)

	}

	return s.stg.User().GetUserById(id)
}

func (s *authService) Login(ctx context.Context, req *auth_service.LoginRequest) (*auth_service.TokenResponse, error) {
	s.log.Info("---Login--->", logger.Any("req", req))
	user, err := s.stg.User().GetUserByUsername(req.Login)
	if err != nil {
		s.log.Error("!!!Login--->", logger.Error(err))
		return nil, status.Errorf(codes.Unauthenticated, "login or password is wrong")
	}
	match, err := util.ComparePassword(user.Password, req.Password)
	if err != nil {
		s.log.Error("!!!Login--->", logger.Error(err))
		return nil, status.Errorf(codes.Internal, "method ComparePassword: %v", err)

	}
	if !match {
		s.log.Error("!!!Login--->", logger.Error(err))
		return nil, status.Errorf(codes.Unauthenticated, "login or password is wrong")
	}
	m := map[string]interface{}{
		"user_id":   user.Id,
		"user_type": user.UserType,
	}
	tokenStr, err := util.GenerateJWT(m, time.Minute*10, s.cfg.SecretKey)
	if err != nil {
		s.log.Error("!!!Login--->", logger.Error(err))
		return nil, status.Errorf(codes.Internal, "method GenerateJWT: %v", err)

	}
	return &auth_service.TokenResponse{Token: tokenStr}, nil
}

func (s *authService) HasAccess(ctx context.Context, req *auth_service.TokenRequest) (*auth_service.HasAccessResponse, error) {
	s.log.Info("---HasAccess--->", logger.Any("req", req))
	res, err := util.ParseClaims(req.Token, s.cfg.SecretKey)
	if err != nil {
		s.log.Error("!!!HasAccess--->", logger.Error(err))
		return &auth_service.HasAccessResponse{
			UserId:    "",
			UserType:  "",
			HasAccess: false,
		}, nil
	}
	user, err := s.stg.User().GetUserById(res.UserId)
	if err != nil {
		s.log.Error("!!!HasAccess--->", logger.Error(err))
		return &auth_service.HasAccessResponse{
			UserId:    "",
			UserType:  "",
			HasAccess: false,
		}, nil
	}
	return &auth_service.HasAccessResponse{
		UserId:    user.Id,
		UserType:  user.UserType,
		HasAccess: true,
	}, nil
}
