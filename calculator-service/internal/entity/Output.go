package entity

type Output struct {
	Result float64 `json:"result"`
	Error  string  `json:"error,omitempty"`
}
