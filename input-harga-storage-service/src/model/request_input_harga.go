package model

type RequestInputHarga struct {
	AdminID      string `json:"admin_id"`
	HargaTopup   int    `json:"harga_topup"`
	HargaBuyback int    `json:"harga_buyback"`
}
