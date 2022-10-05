package model

type Transaksi struct {
	ID           string  `json:"id"`
	TopupPrice   int64   `json:"topup_price"`
	BuybackPrice int64   `json:"buyback_price"`
	RekeningID   string  `json:"rekening_id"`
	Gram         float64 `json:"gram"`
	Type         string  `json:"type"`
	CreatedAt    int     `json:"created_at"`
}
