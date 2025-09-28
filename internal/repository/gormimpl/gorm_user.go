package gormimpl

import (
	"chatsheet/internal/itfc"
	"chatsheet/internal/model"
	"context"
	"log/slog"

	"gorm.io/gorm"
)

type gormUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) itfc.UserRepository {
	return &gormUserRepository{db: db}
}

func (r *gormUserRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	err := r.db.WithContext(ctx).
		Create(&user).
		Error
	if err != nil {
		slog.Error("Failed to create user", "error", err)
		return nil, err
	}

	return user, nil
}

func (r *gormUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user *model.User
	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&user).
		Error
	if err != nil {
		slog.Error("Failed to get user by email", "error", err)
		return nil, err
	}

	return user, nil
}
