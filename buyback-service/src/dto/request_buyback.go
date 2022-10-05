package dto

type RequestBuyback struct {
	Gram  float64 `json:"gram"`
	Harga int     `json:"harga"`
	Norek string  `json:"norek"`
}
