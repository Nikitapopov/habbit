package user

import "context"

type Repository interface {
	Create(ctx context.Context, user *User) (string, error)
	FindAll(ctx context.Context) ([]User, error)
	FindOne(ctx context.Context, id string) (User, error)
	Update(ctx context.Context, userId string, user User) error
	Delete(ctx context.Context, id string) error
}
