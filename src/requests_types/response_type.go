package req_types

type RespImp struct {
	Id     uint    `json:"id"`
	Width  uint    `json:"width"`
	Height uint    `json:"height"`
	Tile   string  `json:"tile"`
	Url    string  `json:"url"`
	Price  float64 `json:"price"`
}

// Response from ad partners
type SuccesResponse struct {
	Id  string    `json:"id"`
	Imp []RespImp `json:"imp"`
}
