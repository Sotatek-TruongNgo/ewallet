package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/truongnqse05461/ewallet/internal/model"
	"github.com/truongnqse05461/ewallet/internal/repository"
)

type UserService interface {
	Create(ctx context.Context, user model.User) (*model.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func (u *userService) Create(ctx context.Context, user model.User) (*model.User, error) {
	user.ID = uuid.New().String()
	timestamp := time.Now().UnixMilli()
	user.CreatedTime = timestamp
	user.UpdatedTime = timestamp
	err := u.userRepo.Create(ctx, user)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &user, nil
}

func NewUserService() UserService {
	return &userService{
		userRepo: repository.NewUserRepository(),
	}
}
