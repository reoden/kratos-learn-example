package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

// Greeter is a Greeter model.
type User struct {
	ID int64
	NickName string
	Password string
	Mobile string
	Birthday int64
	Gender string
	Role int
}

// GreeterRepo is a Greater repo.
type UserRepo interface {
	CreateUser(context.Context, *User) (*User, error)
}

// GreeterUsecase is a Greeter usecase.
type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateGreeter creates a Greeter, and returns the new Greeter.
func (uc *UserUsecase) Create(ctx context.Context, u *User) (*User, error) {
	return uc.repo.CreateUser(ctx, u)
}
