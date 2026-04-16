package entity

type Calculation struct {
	Expression string  `json:"expression"`
	Result     float64 `json:"result"`
}

func (c *Calculation) Validate() error {
	if c.Expression == "" {
		return ErrEmptyExpression
	}
	return nil
}

type DomainError string

func (e DomainError) Error() string { return string(e) }

const (
	ErrEmptyExpression = DomainError("expression cannot be empty")
)
