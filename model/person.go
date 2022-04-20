package model

import "developer.zopsmart.com/go/gofr/pkg/errors"

type Person struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Age     float64 `json:age`
	Address string  `json:address`
}

func (p *Person) Validate() error {
	if p.Name == "" {
		return errors.Error("invalid fileds")
	}

	if p.Address == "" {
		return errors.Error("invalid fileds")
	}

	return nil
}
