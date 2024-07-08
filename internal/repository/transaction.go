package repository

import (
	"context"
	"fmt"

	"github.com/truongnqse05461/ewallet/internal/cx"
	"github.com/truongnqse05461/ewallet/internal/model"
	"github.com/truongnqse05461/ewallet/pkg/slice"
)

type TransactionRepository interface {
	Create(ctx context.Context, transaction model.Transaction) error
	Count(ctx context.Context, userID string, condition model.TransactionSearchCondition) (int, error)
	List(ctx context.Context, userID string, condition model.TransactionSearchCondition) ([]*model.Transaction, error)
}

type transactionRepository struct {
}

func (w *transactionRepository) Create(ctx context.Context, transaction model.Transaction) error {
	tx := cx.GetTx(ctx)

	_, err := tx.ExecContext(ctx, tx.Rebind(
		`INSERT INTO transactions (tx, user_id, from_address, to_address, amount, nonce, timestamp)
			VALUES (?, ?, ?, ?, ?, ?, ?)`),
		transaction.Tx, transaction.UserID, transaction.From, transaction.To, transaction.Amount, transaction.Nonce, transaction.Timestamp,
	)
	if err != nil {
		return err
	}
	return nil
}

// Count implements TransactionRepository.
func (w *transactionRepository) Count(ctx context.Context, userID string, condition model.TransactionSearchCondition) (int, error) {
	tx := cx.GetTx(ctx)

	var cond = "user_id = ?"
	values := []any{userID}
	if condition.From != nil {
		cond += " AND from = ?"
		values = append(values, *condition.From)
	}
	if condition.To != nil {
		cond += " AND to = ?"
		values = append(values, *condition.To)
	}
	if condition.Start != nil {
		cond += " AND timestamp > ?"
		values = append(values, *condition.Start)
	}
	if condition.End != nil {
		cond += " AND timestamp < ?"
		values = append(values, *condition.End)
	}
	var total int

	stmt := fmt.Sprintf(`SELECT COUNT(*) FROM transactions WHERE %s`, cond)

	err := tx.GetContext(ctx, &total, tx.Rebind(stmt), values...)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (w *transactionRepository) List(ctx context.Context, userID string, condition model.TransactionSearchCondition) ([]*model.Transaction, error) {
	tx := cx.GetTx(ctx)

	var cond = "user_id = ?"
	values := []any{userID}
	if condition.From != nil {
		cond += " AND from = ?"
		values = append(values, *condition.From)
	}
	if condition.To != nil {
		cond += " AND to = ?"
		values = append(values, *condition.To)
	}
	if condition.Start != nil {
		cond += " AND timestamp > ?"
		values = append(values, *condition.Start)
	}
	if condition.End != nil {
		cond += " AND timestamp < ?"
		values = append(values, *condition.End)
	}
	values = append(values, condition.Limit, condition.Offset)
	var transactions []dbTransaction

	stmt := fmt.Sprintf(`SELECT * FROM transactions WHERE %s LIMIT ? OFFSET ?`, cond)

	err := tx.SelectContext(ctx, &transactions, tx.Rebind(stmt), values...)
	if err != nil {
		return nil, err
	}
	return slice.Map(transactions, func(d dbTransaction) *model.Transaction {
		return d.toEntity()
	}), nil
}

type dbTransaction struct {
	Tx        string  `db:"tx"`
	UserID    string  `db:"user_id"`
	From      string  `db:"from_address"`
	To        string  `db:"to_address"`
	Amount    float64 `db:"amount"`
	Nonce     string  `db:"nonce"`
	Timestamp int64   `db:"timestamp"`
}

func (d dbTransaction) toEntity() *model.Transaction {
	return &model.Transaction{
		Tx:        d.Tx,
		UserID:    d.UserID,
		From:      d.From,
		To:        d.To,
		Amount:    d.Amount,
		Nonce:     d.Nonce,
		Timestamp: d.Timestamp,
	}
}

func NewTransactionRepository() TransactionRepository {
	return &transactionRepository{}
}
