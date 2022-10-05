package dto

type RequestTransaction struct {
	Norek     string `json:"norek"`
	StartDate int `json:"start_date"`
	EndDate   int `json:"end_date"`
}
