package model

type Wallet struct {
	Address     string  `json:"address"`
	UserID      string  `json:"userID"`
	Balance     float64 `json:"balance"`
	CreatedTime int64   `json:"createdTime"`
	UpdatedTime int64   `json:"updatedTime"`
}

type WalletPatch struct {
	Balance     float64
	UpdatedTime int64
}

type WalletPage struct {
	Limit   int       `json:"limit"`
	Offset  int       `json:"offset"`
	Total   int       `json:"total"`
	Wallets []*Wallet `json:"wallets"`
}
