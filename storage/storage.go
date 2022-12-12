package storage

import (
	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/ecommerce_auth_service/storage/postgres"
	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/ecommerce_auth_service/storage/repo"
	"github.com/jmoiron/sqlx"
)

type StorageI interface {
	User() repo.UserRepoI
}
type storagePg struct {
	db   *sqlx.DB
	user repo.UserRepoI
}

func NewStoragePg(db *sqlx.DB) StorageI {
	return &storagePg{
		db:   db,
		user: postgres.NewUserRepo(db),
	}
}

func (p *storagePg) User() repo.UserRepoI {
	return p.user
}

