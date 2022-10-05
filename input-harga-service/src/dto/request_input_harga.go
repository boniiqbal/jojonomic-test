package dto

type RequestInputHarga struct {
	AdminID      string `json:"norek"`
	HargaTopup   int    `json:"harga_topup"`
	HargaBuyback int    `json:"harga_buyback"`
}
