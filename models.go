package coffeezone

type LocationModel struct {
	Address   string  `json:"address"`
	Longitude float64 `json:"lon"`
	Latitude  float64 `json:"lat"`
}

type CafeModel struct {
	Code     string        `json:"code"`
	Title    string        `json:"title"`
	Location LocationModel `json:"location,omitempty"`
}
