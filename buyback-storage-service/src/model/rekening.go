package model

type Rekening struct {
	ID     string  `json:"id"`
	UserID string  `json:"user_id"`
	Norek  string  `json:"norek"`
	Saldo  float64 `json:"saldo"`
}
