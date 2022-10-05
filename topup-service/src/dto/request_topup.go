package dto

type RequestTopup struct {
	Gram  float64 `json:"gram"`
	Harga int     `json:"harga"`
	Norek string  `json:"norek"`
}
