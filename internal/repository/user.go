package repository

import (
	"context"

	"github.com/truongnqse05461/ewallet/internal/cx"
	"github.com/truongnqse05461/ewallet/internal/model"
	"github.com/truongnqse05461/ewallet/pkg/slice"
)

type UserRepository interface {
	Create(ctx context.Context, user model.User) error
	List(ctx context.Context, limit, offset uint64) ([]*model.User, error)
}

type userRepository struct {
}

func (w *userRepository) Create(ctx context.Context, user model.User) error {
	tx := cx.GetTx(ctx)

	_, err := tx.ExecContext(ctx, tx.Rebind(
		`INSERT INTO users (id, name, created_time, updated_time)
			VALUES (?, ?, ?, ?)`),
		user.ID, user.Name, user.CreatedTime, user.UpdatedTime,
	)
	if err != nil {
		return err
	}
	return nil
}

func (w *userRepository) List(ctx context.Context, limit uint64, offset uint64) ([]*model.User, error) {
	tx := cx.GetTx(ctx)
	var users []dbUser

	err := tx.SelectContext(ctx, &users, tx.Rebind(
		`SELECT * FROM users LIMIT ? OFFSET ?`),
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	return slice.Map(users, func(d dbUser) *model.User {
		return d.toEntity()
	}), nil
}

type dbUser struct {
	ID          string `db:"id"`
	Name        string `db:"name"`
	CreatedTime int64  `db:"created_time"`
	UpdatedTime int64  `db:"updated_time"`
}

func (d dbUser) toEntity() *model.User {
	return &model.User{
		ID:          d.ID,
		Name:        d.Name,
		CreatedTime: d.CreatedTime,
		UpdatedTime: d.UpdatedTime,
	}
}

func NewUserRepository() UserRepository {
	return &userRepository{}
}
