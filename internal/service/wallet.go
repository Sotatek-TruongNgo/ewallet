package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/truongnqse05461/ewallet/internal/model"
	"github.com/truongnqse05461/ewallet/internal/repository"
)

type WalletService interface {
	Create(ctx context.Context, wallet model.Wallet) (*model.Wallet, error)
	List(ctx context.Context, userID string, limit, offset int) (model.WalletPage, error)
}

type walletService struct {
	walletRepo repository.WalletRepository
}

func (w *walletService) Create(ctx context.Context, wallet model.Wallet) (*model.Wallet, error) {
	wallet.Address = uuid.New().String()
	timestamp := time.Now().UnixMilli()
	wallet.CreatedTime = timestamp
	wallet.UpdatedTime = timestamp
	err := w.walletRepo.Create(ctx, wallet)
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (w *walletService) List(ctx context.Context, userID string, limit int, offset int) (model.WalletPage, error) {
	walets, err := w.walletRepo.List(ctx, userID, limit, offset)
	if err != nil {
		return model.WalletPage{}, err
	}
	total, err := w.walletRepo.Count(ctx, userID)
	if err != nil {
		return model.WalletPage{}, err
	}
	return model.WalletPage{
		Limit:   limit,
		Offset:  offset,
		Total:   total,
		Wallets: walets,
	}, nil
}

func NewWalletService() WalletService {
	return &walletService{
		walletRepo: repository.NewWalletRepository(),
	}
}
