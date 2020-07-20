package entities

type Password struct {
	Value string `json:"value"`
}

type Response struct {
	Valid bool `json:"valid"`
}
