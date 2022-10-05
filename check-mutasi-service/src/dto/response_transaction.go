package dto

type DetailTransaction struct {
	Date         int     `json:"date"`
	Type         string  `json:"type"`
	Gram         float64 `json:"gram"`
	HargaTopup   int64   `json:"harga_topup"`
	HargaBuyback int64   `json:"harga_buyback"`
	Saldo        float64 `json:"saldo"`
}

type ResponseTransaction struct {
	Data []DetailTransaction `json:"data"`
}
