package dto

type CurrencyResponse struct {
	Valute map[string]Valute `json:"Valute"`
}

type Valute struct {
	ID       string  `json:"ID"`
	NumCode  string  `json:"NumCode"`
	CharCode string  `json:"CharCode"`
	Nominal  int     `json:"Nominal"`
	Name     string  `json:"Name"`
	Value    float32 `json:"Value"`
	Previous float32 `json:"Previous"`
}
