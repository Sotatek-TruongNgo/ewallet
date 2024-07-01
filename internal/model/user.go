package model

type User struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	CreatedTime int64  `json:"createdTime"`
	UpdatedTime int64  `json:"updatedTime"`
}
