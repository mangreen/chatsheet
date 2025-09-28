package gormimpl

import (
	"chatsheet/internal/itfc"
	"chatsheet/internal/model"
	"context"
	"log/slog"

	"gorm.io/gorm"
)

type gormUnipileRepository struct {
	db *gorm.DB
}

func NewUnipileRepository(db *gorm.DB) itfc.UnipileRepository {
	return &gormUnipileRepository{db: db}
}

func (r *gormUnipileRepository) Create(ctx context.Context, acct *model.UnipileAccount) (*model.UnipileAccount, error) {
	err := r.db.WithContext(ctx).
		Create(&acct).
		Error
	if err != nil {
		slog.Error("Failed to create a UnipileAccount", "error", err)
		return nil, err
	}

	return acct, nil
}

func (r *gormUnipileRepository) ListByEmail(ctx context.Context, email string) ([]model.UnipileAccount, error) {
	var accts []model.UnipileAccount
	err := r.db.WithContext(ctx).
		Where("user_email = ?", email).
		Find(&accts).
		Error
	if err != nil {
		slog.Error("Failed to list UnipileAccount by email", "error", err)
		return nil, err
	}

	return accts, nil
}
