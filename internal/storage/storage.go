package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/esquirelol/auth-grpc/internal/domain/models"
	"github.com/jackc/pgx/v5"
)

var (
	UserIsNotExists = errors.New("user not is exists")
	AppIsNotExists  = errors.New("app not is exists")
)

type Storage struct {
	conn *pgx.Conn
}

func New(ctx context.Context, storagePath string) (Storage, error) {
	conn, err := pgx.Connect(ctx, storagePath)
	if err != nil {
		return Storage{}, errors.New("failed to connection storage")
	}

	return Storage{conn: conn}, nil
}

func (s *Storage) UserSaver(ctx context.Context, email string, passHash []byte) (int64, error) {
	sqlQuery := `
	INSERT INTO users(email,pass_hash)
	VALUES($1,$2)
	RETURNING id
`
	var userId int64
	if err := s.conn.QueryRow(ctx, sqlQuery, email, string(passHash)).Scan(&userId); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, UserIsNotExists
		}
		return 0, errors.New("failed to scan")
	}

	return userId, nil

}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	var outUser models.User
	sqlQuery := `
	SELECT id,email,pass_hash FROM users
	WHERE email = $1
`

	err := s.conn.QueryRow(ctx, sqlQuery, email).Scan(&outUser.ID, &outUser.Email, &outUser.PassHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, UserIsNotExists
		}
		return models.User{}, fmt.Errorf("failed to select user: %w", err)
	}

	return outUser, nil
}

func (s *Storage) App(ctx context.Context, appId int64) (models.App, error) {
	var outApp models.App
	sqlQuery := `
	SELECT id,name_app FROM apps
	WHERE id = $1
`

	if err := s.conn.QueryRow(ctx, sqlQuery, appId).Scan(&outApp.ID, &outApp.Name); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.App{}, AppIsNotExists
		}
		return models.App{}, fmt.Errorf("failed to select apps: %w", err)
	}

	return outApp, nil
}
