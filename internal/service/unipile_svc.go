package service

import (
	"chatsheet/internal/itfc"
	"chatsheet/internal/model"
	"context"
)

// UnipileService 包含業務邏輯
type UnipileService struct {
	unipileRepo itfc.UnipileRepository // 依賴介面，而非實作
}

func NewUnipileService(repo itfc.UnipileRepository) *UnipileService {
	return &UnipileService{
		unipileRepo: repo,
	}
}

func (s *UnipileService) Create(ctx context.Context, email, provider, accountID string) (*model.UnipileAccount, error) {
	acct := &model.UnipileAccount{
		UserEmail: email,
		Provider:  provider,
		AccountID: accountID,
	}

	newAcct, err := s.unipileRepo.Create(ctx, acct)
	if err != nil {
		return nil, err
	}

	return newAcct, nil
}

func (s *UnipileService) ListByEmail(ctx context.Context, email string) ([]model.UnipileAccount, error) {
	accts, err := s.unipileRepo.ListByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return accts, nil
}
