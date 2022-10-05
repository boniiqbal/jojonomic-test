package model

import "time"

type Harga struct {
	ID           string     `json:"id"`
	TopupPrice   int64      `json:"topup_price"`
	BuybackPrice int64      `json:"buyback_price"`
	UserID       string     `json:"user_id"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}
