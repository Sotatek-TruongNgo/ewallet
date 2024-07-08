package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/rs/xid"
	"github.com/truongnqse05461/ewallet/internal/model"
)

func generateTransaction() model.Transaction {
	tx := uuid.New().String()
	userID := uuid.New().String()
	from := uuid.New().String()
	to := uuid.New().String()
	amount := rand.Float64() * 1000 // Random amount between 0 and 1000
	nonce := xid.New().String()
	timestamp := time.Now().UnixMilli()

	return model.Transaction{
		Tx:        tx,
		UserID:    userID,
		From:      from,
		To:        to,
		Amount:    amount,
		Nonce:     nonce,
		Timestamp: timestamp,
	}
}

func DummyTransaction() {
	rand.Seed(time.Now().UnixNano())
	numTransactions := 5000000
	filePath := "transactions.csv"

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Failed to create file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"tx", "user_id", "from_wallet", "to_wallet", "amount", "nonce", "timestamp"}
	writer.Write(header)

	for i := 0; i < numTransactions; i++ {
		transaction := generateTransaction()

		record := []string{
			transaction.Tx,
			transaction.UserID,
			transaction.From,
			transaction.To,
			fmt.Sprintf("%.2f", transaction.Amount),
			transaction.Nonce,
			fmt.Sprintf("%d", transaction.Timestamp),
		}
		writer.Write(record)
	}

	fmt.Printf("CSV file with %d transactions generated: %s\n", numTransactions, filePath)
}
