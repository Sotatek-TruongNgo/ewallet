package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/truongnqse05461/ewallet/internal/model"
)

func DummyUser() {
	numUsers := 50000000 // Change this to 50000000 for 50 million users
	filePath := "users.csv"

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Failed to create file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"id", "name", "created_time", "updated_time"}
	writer.Write(header)

	for i := 0; i < numUsers; i++ {
		id := uuid.New().String()
		name := faker.Name()
		createdTime := time.Now().Unix()
		updatedTime := createdTime

		user := model.User{
			ID:          id,
			Name:        name,
			CreatedTime: createdTime,
			UpdatedTime: updatedTime,
		}

		record := []string{user.ID, user.Name, fmt.Sprintf("%d", user.CreatedTime), fmt.Sprintf("%d", user.UpdatedTime)}
		writer.Write(record)
	}

	fmt.Printf("CSV file with %d users generated: %s\n", numUsers, filePath)
}
