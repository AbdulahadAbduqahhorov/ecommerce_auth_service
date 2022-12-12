package repo

import (
	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/ecommerce_auth_service/genproto/auth_service"
)

type UserRepoI interface {
	CreateUser(req *auth_service.CreateUserRequest) (string, error)
	GetUserList(req *auth_service.GetUserListRequest) (*auth_service.GetUserListResponse, error)
	GetUserById(id string) (*auth_service.User, error)
	UpdateUser(req *auth_service.UpdateUserRequest) (int64,error)
	DeleteUser(id string)(int64, error)
	GetUserByUsername(username string) (*auth_service.User, error)
	Register(req *auth_service.RegisterUserRequest) (string, error)
}
