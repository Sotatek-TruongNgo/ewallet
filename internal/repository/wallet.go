package repository

import (
	"context"

	"github.com/truongnqse05461/ewallet/internal/cx"
	"github.com/truongnqse05461/ewallet/internal/model"
	"github.com/truongnqse05461/ewallet/pkg/slice"
)

type WalletRepository interface {
	Create(ctx context.Context, wallet model.Wallet) error
	List(ctx context.Context, userID string, limit, offset int) ([]*model.Wallet, error)
	Count(ctx context.Context, userID string) (int, error)
	Get(ctx context.Context, address string) (*model.Wallet, error)
	Update(ctx context.Context, address string, patch model.WalletPatch) error
}

type walletRepository struct {
}

func (w *walletRepository) Create(ctx context.Context, wallet model.Wallet) error {
	tx := cx.GetTx(ctx)

	_, err := tx.ExecContext(ctx, tx.Rebind(
		`INSERT INTO wallets (address, user_id, balance, created_time, updated_time)
			VALUES (?, ?, ?, ?, ?)`),
		wallet.Address, wallet.UserID, wallet.Balance, wallet.CreatedTime, wallet.UpdatedTime,
	)
	if err != nil {
		return err
	}
	return nil
}

func (w *walletRepository) List(ctx context.Context, userID string, limit int, offset int) ([]*model.Wallet, error) {
	tx := cx.GetTx(ctx)
	var wallets []dbWallet

	err := tx.SelectContext(ctx, &wallets, tx.Rebind(
		`SELECT * FROM wallets WHERE user_id = ? LIMIT ? OFFSET ?`),
		userID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	return slice.Map(wallets, func(d dbWallet) *model.Wallet {
		return d.toEntity()
	}), nil
}

func (w *walletRepository) Count(ctx context.Context, userID string) (int, error) {
	tx := cx.GetTx(ctx)
	var total int

	err := tx.GetContext(ctx, &total, tx.Rebind(
		`SELECT COUNT(*) FROM wallets WHERE user_id = ?`),
		userID,
	)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (w *walletRepository) Get(ctx context.Context, address string) (*model.Wallet, error) {
	tx := cx.GetTx(ctx)
	var wallet dbWallet

	err := tx.GetContext(ctx, &wallet, tx.Rebind(
		`SELECT * FROM wallets WHERE address = ?`),
		address,
	)
	if err != nil {
		return nil, err
	}
	return wallet.toEntity(), nil
}

func (w *walletRepository) Update(ctx context.Context, address string, patch model.WalletPatch) error {
	tx := cx.GetTx(ctx)

	_, err := tx.ExecContext(ctx, tx.Rebind(
		`UPDATE wallets SET balance = ?, updated_time = ? 
			WHERE address = ?`),
		patch.Balance, patch.UpdatedTime, address,
	)
	if err != nil {
		return err
	}
	return nil
}

type dbWallet struct {
	Address     string  `db:"address"`
	UserID      string  `db:"user_id"`
	Balance     float64 `db:"balance"`
	CreatedTime int64   `db:"created_time"`
	UpdatedTime int64   `db:"updated_time"`
}

func (d dbWallet) toEntity() *model.Wallet {
	return &model.Wallet{
		Address:     d.Address,
		UserID:      d.UserID,
		Balance:     d.Balance,
		CreatedTime: d.CreatedTime,
		UpdatedTime: d.UpdatedTime,
	}
}

func NewWalletRepository() WalletRepository {
	return &walletRepository{}
}
