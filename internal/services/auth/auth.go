package auth

import (
	"context"
	"errors"
	"time"

	"github.com/esquirelol/auth-grpc/internal/domain/models"
	jwt_lib "github.com/esquirelol/auth-grpc/internal/lib/jwt"
	"github.com/esquirelol/auth-grpc/internal/storage"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log      *zap.Logger
	storage  Storage
	tokenTTL time.Duration
}

type Storage interface {
	UserSaver(ctx context.Context, email string, passHash []byte) (int64, error)
	User(ctx context.Context, email string) (models.User, error)
	App(ctx context.Context, appId int64) (models.App, error)
}

func New(log *zap.Logger, storage Storage, tokenTTL time.Duration) *Auth {
	return &Auth{
		log:      log,
		storage:  storage,
		tokenTTL: tokenTTL,
	}
}

func (a *Auth) Login(ctx context.Context, email, password string, appID int64) (string, error) {
	const op = "auth.Login"

	log := a.log.With(zap.String("op", op))

	user, err := a.storage.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.UserIsNotExists) {
			log.Warn("user is not exists ")
			return "", err
		}
		log.Error("failed to login")
		return "", err

	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Info("invalid password")
		return "", err
	}

	app, err := a.storage.App(ctx, appID)

	if err != nil {
		if errors.Is(err, storage.AppIsNotExists) {
			log.Warn("app is not exists ")
			return "", err
		}

		log.Error("failed to login app")
		return "", err
	}
	token, err := jwt_lib.NewToken(user, app, a.tokenTTL)
	if err != nil {
		log.Error("failed to create token")
		return "", err
	}
	return token, nil
}

func (a *Auth) Register(ctx context.Context, email, password string) (int64, error) {
	const op = "auth.RegisterNewUser"
	log := a.log.With(zap.String("op", op))
	hashPass, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		log.Error("failed to generate hash pass")
		return 0, err
	}

	idUser, err := a.storage.UserSaver(ctx, email, hashPass)
	if err != nil {
		log.Error("failed to save user")
		return 0, err
	}
	return idUser, nil
}
