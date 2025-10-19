package cfg

import (
	v "github.com/go-ozzo/ozzo-validation/v4"
)

type BaseConfigModel struct {
	Redis struct {
		Host string `koanf:"host" validate:"nonzero"`
		Port string `koanf:"port" validate:"nonzero"`
	} `koanf:"redis" validate:"nonnil"`
}

func (c *BaseConfigModel) Validate() error {
	return v.ValidateStruct(&c) // Digits cannot be empty, and the length must be 4
	//validation.Field(&d.Digits, validation.Required, validation.Length(4, 4)),
}
