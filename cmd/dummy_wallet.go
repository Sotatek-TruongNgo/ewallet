package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/truongnqse05461/ewallet/internal/model"
)

func generateWallet() model.Wallet {
	address := uuid.New().String()
	userID := faker.UUIDDigit()
	balance := rand.Float64() * 1000 // Random balance between 0 and 1000
	createdTime := time.Now().Unix()
	updatedTime := createdTime

	return model.Wallet{
		Address:     address,
		UserID:      userID,
		Balance:     balance,
		CreatedTime: createdTime,
		UpdatedTime: updatedTime,
	}
}

func DummyWallet() {
	rand.Seed(time.Now().UnixNano())
	numWallets := 1000 // Change this to the desired number of wallets
	filePath := "wallets.csv"

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Failed to create file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"address", "user_id", "balance", "created_time", "updated_time"}
	writer.Write(header)

	for i := 0; i < numWallets; i++ {
		wallet := generateWallet()

		record := []string{
			wallet.Address,
			wallet.UserID,
			fmt.Sprintf("%.2f", wallet.Balance),
			fmt.Sprintf("%d", wallet.CreatedTime),
			fmt.Sprintf("%d", wallet.UpdatedTime),
		}
		writer.Write(record)
	}

	fmt.Printf("CSV file with %d wallets generated: %s\n", numWallets, filePath)
}
