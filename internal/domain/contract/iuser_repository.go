package contract

import (
	"blog/internal/domain/entity"
	"context"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, user entity.User) (*entity.User, error)
	GetUserByID(ctx context.Context, id string) (*entity.User, error)
	GetByUserName(ctx context.Context, username string) (*entity.User, error)
	UpdateUser(ctx context.Context, user entity.User) (*entity.User, error)
	DeleteUser(ctx context.Context, id string) error
}