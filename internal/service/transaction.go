package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/xid"
	"github.com/truongnqse05461/ewallet/internal/model"
	"github.com/truongnqse05461/ewallet/internal/repository"
)

type TransactionService interface {
	Transfer(ctx context.Context, from, to string, amount float64) (*model.Transaction, error)
	List(ctx context.Context, userID string, condition model.TransactionSearchCondition) (model.TransactionPage, error)
}

type transactionService struct {
	transactionRepo repository.TransactionRepository
	walletRepo      repository.WalletRepository
	notificationSvc NotificationService
}

// Transfer implements TransactionService.
func (t *transactionService) Transfer(ctx context.Context, from, to string, amount float64) (*model.Transaction, error) {
	fromWallet, err := t.walletRepo.Get(ctx, from)
	if err != nil {
		return nil, err
	}

	toWallet, err := t.walletRepo.Get(ctx, to)
	if err != nil {
		return nil, err
	}
	if fromWallet.Balance < amount {
		return nil, errors.New("insufficient funds")
	}
	tx := &model.Transaction{
		Tx:        uuid.New().String(),
		UserID:    fromWallet.UserID,
		From:      from,
		To:        to,
		Amount:    amount,
		Nonce:     xid.New().String(),
		Timestamp: time.Now().UnixMilli(),
	}

	err = t.walletRepo.Update(ctx, from, model.WalletPatch{
		UpdatedTime: time.Now().UnixMilli(),
		Balance:     fromWallet.Balance - amount,
	})
	if err != nil {
		return nil, err
	}
	err = t.walletRepo.Update(ctx, to, model.WalletPatch{
		UpdatedTime: time.Now().UnixMilli(),
		Balance:     toWallet.Balance + amount,
	})
	if err != nil {
		return nil, err
	}
	err = t.transactionRepo.Create(ctx, *tx)
	if err != nil {
		return nil, err
	}
	go t.notificationSvc.SendNotification(ctx, model.Notification{
		UserID:  fromWallet.UserID,
		Message: fmt.Sprintf("You have successfully transferred %.2f to user %s", amount, toWallet.UserID),
	})
	go t.notificationSvc.SendNotification(ctx, model.Notification{
		UserID:  toWallet.UserID,
		Message: fmt.Sprintf("You have received %.2f from wallet %s", amount, fromWallet.UserID),
	})
	return tx, nil
}

func (t *transactionService) List(ctx context.Context, userID string, condition model.TransactionSearchCondition) (model.TransactionPage, error) {
	txs, err := t.transactionRepo.List(ctx, userID, condition)
	if err != nil {
		return model.TransactionPage{}, err
	}
	page := model.TransactionPage{
		Transactions: txs,
		Limit:        condition.Limit,
		Offset:       condition.Offset,
	}
	return page, nil
}

func NewTransactionService(notificationSvc NotificationService) TransactionService {
	return &transactionService{
		transactionRepo: repository.NewTransactionRepository(),
		walletRepo:      repository.NewWalletRepository(),
		notificationSvc: notificationSvc,
	}
}
