package req_types

type RespImp struct {
	Width  uint
	Height uint
	Tile   string
	Url    string
	Price  float64
}

type SuccesResponse struct {
	Id  string    `json:"id"`
	Imp []RespImp `json:"imp"`
}
