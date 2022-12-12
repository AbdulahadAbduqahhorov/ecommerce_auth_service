package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/ecommerce_auth_service/genproto/auth_service"
	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/ecommerce_auth_service/pkg/helper"
	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/ecommerce_auth_service/storage/repo"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) repo.UserRepoI {
	return userRepo{
		db: db,
	}
}

func (a userRepo) CreateUser(req *auth_service.CreateUserRequest) (string, error) {
	id := uuid.New().String()
	query := `INSERT INTO 
	"user" (
		"id",
		"full_name",
		"login",
		"phone",
		"email",
		"password",
		"user_type"
		) 
	VALUES (
		$1, 
		$2,
		$3,
		$4,
		$5,
		$6,
		$7
		)`
	_, err := a.db.Exec(query, id, req.FullName, req.Login, req.Phone, req.Email, req.Password, req.UserType)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (a userRepo) GetUserList(req *auth_service.GetUserListRequest) (*auth_service.GetUserListResponse, error) {
	res := &auth_service.GetUserListResponse{}
	params := make(map[string]interface{})
	var arr []interface{}
	query := `SELECT
		id,
		full_name,
		login,
		phone,
		email,
		password,
		user_type,
		created_at,
		updated_at
	FROM
		"user"`
	filter := " WHERE 1=1"
	order := " ORDER BY created_at"
	offset := " OFFSET 0"
	limit := " LIMIT 10"

	if len(req.Search) > 0 {
		params["search"] = req.Search
		filter += " AND ((full_name || phone || login) ILIKE ('%' || :search || '%'))"
	}
	if len(req.Type) > 0 {
		params["user_type"] = req.Type
		filter += " AND user_type=:user_type"
	}
	if req.Offset > 0 {
		params["offset"] = req.Offset
		offset = " OFFSET :offset"
	}

	if req.Limit > 0 {
		params["limit"] = req.Limit
		limit = " LIMIT :limit"
	}

	cQ := `SELECT count(1) FROM "user"` + filter
	cQ, arr = helper.ReplaceQueryParams(cQ, params)
	err := a.db.QueryRow(cQ, arr...).Scan(
		&res.Count,
	)
	if err != nil {
		return res, err
	}

	q := query + filter + order + offset + limit

	q, arr = helper.ReplaceQueryParams(q, params)
	rows, err := a.db.Query(q, arr...)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		obj := &auth_service.User{}
		var updatedAt sql.NullString

		err = rows.Scan(
			&obj.Id,
			&obj.FullName,
			&obj.Login,
			&obj.Phone,
			&obj.Email,
			&obj.Password,
			&obj.UserType,
			&obj.CreatedAt,
			&updatedAt,
		)

		if err != nil {
			return res, err
		}
		if updatedAt.Valid {
			obj.UpdatedAt = updatedAt.String
		}
		res.Users = append(res.Users, obj)
	}

	return res, nil
}

func (a userRepo) GetUserById(id string) (*auth_service.User, error) {
	res := &auth_service.User{}
	var updatedAt sql.NullString
	query := `SELECT
		"id",
		"full_name",
		"login",
		"phone",
		"email",
		"password",
		"user_type",
		"created_at",
		"updated_at"
	FROM
		"user"
	WHERE
		id = $1`

	err := a.db.QueryRow(query, id).Scan(
		&res.Id,
		&res.FullName,
		&res.Login,
		&res.Phone,
		&res.Email,
		&res.Password,
		&res.UserType,
		&res.CreatedAt,
		&updatedAt,
	)
	if err != nil {
		return res, err
	}
	if updatedAt.Valid {
		res.UpdatedAt = updatedAt.String
	}

	return res, nil
}

func (a userRepo) UpdateUser(req *auth_service.UpdateUserRequest) (int64, error) {

	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1
	if len(strings.Trim(req.Login, " ")) > 0 {
		setValues = append(setValues, fmt.Sprintf("login=$%d ", argId))
		args = append(args, req.Login)
		argId++
	}
	if len(strings.Trim(req.FullName, " ")) > 0 {
		setValues = append(setValues, fmt.Sprintf("full_name=$%d ", argId))
		args = append(args, req.FullName)
		argId++
	}
	if len(strings.Trim(req.Email, " ")) > 0 {
		setValues = append(setValues, fmt.Sprintf("email=$%d ", argId))
		args = append(args, req.Email)
		argId++
	}
	if len(strings.Trim(req.Phone, " ")) > 0 {
		setValues = append(setValues, fmt.Sprintf("phone=$%d ", argId))
		args = append(args, req.Phone)
		argId++
	}

	s := strings.Join(setValues, ",")
	query := fmt.Sprintf(`
			UPDATE "user"
			SET %s ,updated_at = now()
			WHERE id = $%d`,
		s, argId)

	args = append(args, req.Id)

	result, err := a.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil

}

func (a userRepo) DeleteUser(id string) (int64, error) {

	query := `DELETE FROM "user" WHERE id = $1`

	result, err := a.db.Exec(query, id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()

	return rowsAffected, err
}
func (u userRepo) Register(req *auth_service.RegisterUserRequest) (string, error) {
	id := uuid.New().String()
	_, err := u.db.Exec(`INSERT INTO 
		"user" (
			id,
			full_name,
			login,
			phone,
			email,
			password,
			user_type
			) 
		VALUES (
			$1, 
			$2,
			$3,
			$4,
			$5,
			$6,
			$7
			)`,
		id,
		req.FullName,
		req.Login,
		req.Phone,
		req.Email,
		req.Password,
		"USER",
	)
	if err != nil {
		return "", err
	}
	return id, nil

}
func (u userRepo) GetUserByUsername(username string) (*auth_service.User, error) {
	res := &auth_service.User{}
	var (
		updatedAt sql.NullString
		userType  string
	)
	err := u.db.QueryRow(`
	SELECT 
		id,
		login,
		password,
		user_type,
		created_at,
		updated_at 
	FROM "user"
	WHERE login=$1`, username).Scan(
		&res.Id,
		&res.Login,
		&res.Password,
		&userType,
		&res.CreatedAt,
		&updatedAt,
	)
	if updatedAt.Valid {
		res.UpdatedAt = updatedAt.String
	}
	if err != nil {
		return nil, err
	}

	return res, nil
}
