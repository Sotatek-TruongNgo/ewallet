package model

type Transaction struct {
	Tx        string  `json:"tx"`
	UserID    string  `json:"user_id"`
	From      string  `json:"from"`
	To        string  `json:"to"`
	Amount    float64 `json:"amount"`
	Nonce     string  `json:"nonce"`
	Timestamp int64   `json:"timestamp"`
}

type TransactionSearchCondition struct {
	From   *string
	To     *string
	Start  *int64
	End    *int64
	Offset int
	Limit  int
}

type TransactionPage struct {
	Limit        int            `json:"limit"`
	Offset       int            `json:"offset"`
	Total        int            `json:"total"`
	Transactions []*Transaction `json:"transactions"`
}
