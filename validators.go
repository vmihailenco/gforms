package gforms

import (
	"fmt"
)

type Validator interface {
	Validate(interface{}) error
}

type StringChoicesValidator struct {
	Choices []StringChoice
}

func (v *StringChoicesValidator) Validate(rawValue interface{}) error {
	value, ok := rawValue.(string)
	if !ok {
		return fmt.Errorf("Type %T is not supported", rawValue)
	}
	for _, choice := range v.Choices {
		if choice.Value == value {
			return nil
		}
	}
	return fmt.Errorf("%v is invalid choice", value)
}

func NewStringChoicesValidator(choices []StringChoice) *StringChoicesValidator {
	return &StringChoicesValidator{Choices: choices}
}

type Int64ChoicesValidator struct {
	Choices []Int64Choice
}

func (v *Int64ChoicesValidator) Validate(rawValue interface{}) error {
	value, ok := rawValue.(int64)
	if !ok {
		return fmt.Errorf("Type %T is not supported", rawValue)
	}
	for _, choice := range v.Choices {
		if choice.Value == value {
			return nil
		}
	}
	return fmt.Errorf("%v is invalid choice", value)
}

func NewInt64ChoicesValidator(choices []Int64Choice) *Int64ChoicesValidator {
	return &Int64ChoicesValidator{Choices: choices}
}
