package gforms

import (
	"errors"
	"fmt"
	"html/template"
	"mime/multipart"
	"strconv"
)

// ----------------------------------------------------------------------------
// Field
// ----------------------------------------------------------------------------

type Field interface {
	Validator

	AddValidator(Validator)
	ApplyValidators(interface{}) error
	HasValidationError() bool
	Render(...string) template.HTML
	ToBaseField() *BaseField
}

func IsFieldValid(field Field, rawValue interface{}) bool {
	bf := field.ToBaseField()

	if rawValue == nil {
		if bf.IsRequired {
			bf.ValidationError = errors.New("This field is required.")
			return false
		} else {
			return true
		}
	}

	if err := field.Validate(rawValue); err != nil {
		bf.ValidationError = err
		return false
	}

	return true
}

type SingleValueField interface {
	StringValue() string
}

type MultiValueField interface {
	StringValue() []string
}

// ----------------------------------------------------------------------------
// BaseField
// ----------------------------------------------------------------------------

type BaseField struct {
	Name   string
	Label  string
	Widget Widget

	IsMulti     bool
	IsMultipart bool
	IsRequired  bool

	Validators      []Validator
	ValidationError error
	IValue          interface{}
}

func (f *BaseField) AddValidator(validator Validator) {
	if f.Validators == nil {
		f.Validators = make([]Validator, 0)
	}
	f.Validators = append(f.Validators, validator)
}

func (f *BaseField) ApplyValidators(rawValue interface{}) error {
	for _, validator := range f.Validators {
		if err := validator.Validate(rawValue); err != nil {
			return err
		}
	}
	return nil
}

func (f *BaseField) Validate(rawValue interface{}) error {
	panic("not implemented.")
}

func (f *BaseField) HasValidationError() bool {
	return f.ValidationError != nil
}

func (f *BaseField) StringValue() string {
	if f.IValue == nil {
		return ""
	}
	return fmt.Sprint(f.IValue)
}

func (f *BaseField) Render(attrs ...string) template.HTML {
	panic("not implemented.")
	return template.HTML("")
}

func (f *BaseField) ToBaseField() *BaseField {
	return f
}

// ----------------------------------------------------------------------------
// StringField
// ----------------------------------------------------------------------------

type StringField struct {
	*BaseField
	MinLen, MaxLen int
}

func (f *StringField) Value() string {
	if f.IValue == nil {
		return ""
	}
	return f.IValue.(string)
}

func (f *StringField) Validate(rawValue interface{}) error {
	value := fmt.Sprint(rawValue)

	valueLen := len(value)
	if f.MinLen > 0 && valueLen < f.MinLen {
		return fmt.Errorf("This field should have at least %d symbols.", f.MinLen)
	}
	if f.MaxLen > 0 && valueLen > f.MaxLen {
		return fmt.Errorf("This field should have less than %d symbols.", f.MaxLen)
	}

	if err := f.ApplyValidators(value); err != nil {
		return err
	}

	f.IValue = value
	return nil
}

func (f *StringField) SetInitial(initial string) {
	f.IValue = initial
}

func (f *StringField) Render(attrs ...string) template.HTML {
	return f.Widget.Render(attrs, f.StringValue())
}

func NewStringField() *StringField {
	return &StringField{
		BaseField: &BaseField{
			Widget:     NewTextWidget(),
			IsRequired: true,
		},
	}
}

func NewTextareaStringField() *StringField {
	return &StringField{
		BaseField: &BaseField{
			Widget:     NewTextareaWidget(),
			IsRequired: true,
		},
	}
}

// ----------------------------------------------------------------------------
// StringChoiceField
// ----------------------------------------------------------------------------

type StringChoice struct {
	Value string
	Label string
}

type StringChoiceField struct {
	*StringField
}

func (f *StringChoiceField) SetChoices(choices []StringChoice) {
	strChoices := make([][2]string, 0, len(choices))
	for _, choice := range choices {
		strChoices = append(strChoices, [2]string{choice.Value, choice.Label})
	}

	f.Widget.(ChoiceWidget).SetChoices(strChoices)
	f.AddValidator(NewStringChoicesValidator(choices))
}

func NewSelectStringField() *StringChoiceField {
	return &StringChoiceField{
		StringField: &StringField{
			BaseField: &BaseField{
				Widget:     NewSelectWidget(),
				IsRequired: true,
			},
		},
	}
}

func NewRadioStringField() *StringChoiceField {
	return &StringChoiceField{
		&StringField{
			BaseField: &BaseField{
				Widget:     NewRadioWidget(),
				IsRequired: true,
			},
		},
	}
}

// ----------------------------------------------------------------------------
// Int64Field
// ----------------------------------------------------------------------------

type Int64Field struct {
	*BaseField
}

func (f *Int64Field) Value() int64 {
	if f.IValue == nil {
		return 0
	}
	return f.IValue.(int64)
}

func (f *Int64Field) Validate(rawValue interface{}) error {
	value, err := strconv.ParseInt(fmt.Sprint(rawValue), 10, 64)
	if err != nil {
		return err
	}

	if err := f.ApplyValidators(value); err != nil {
		return err
	}

	f.IValue = value
	return nil
}

func (f *Int64Field) SetInitial(initial int64) {
	f.IValue = initial
}

func (f *Int64Field) Render(attrs ...string) template.HTML {
	return f.Widget.Render(attrs, f.StringValue())
}

func NewInt64Field() *Int64Field {
	return &Int64Field{
		BaseField: &BaseField{
			Widget:     NewTextWidget(),
			IsRequired: true,
		},
	}
}

// ----------------------------------------------------------------------------
// Int64ChoiceField
// ----------------------------------------------------------------------------

type Int64ChoiceField struct {
	*Int64Field
}

type Int64Choice struct {
	Value int64
	Label string
}

func (f *Int64ChoiceField) SetChoices(choices []Int64Choice) {
	strChoices := make([][2]string, 0, len(choices))
	for _, choice := range choices {
		strChoice := [2]string{strconv.FormatInt(choice.Value, 10), choice.Label}
		strChoices = append(strChoices, strChoice)
	}

	f.Widget.(ChoiceWidget).SetChoices(strChoices)
	f.AddValidator(NewInt64ChoicesValidator(choices))
}

func NewSelectInt64Field() *Int64ChoiceField {
	return &Int64ChoiceField{
		Int64Field: &Int64Field{
			BaseField: &BaseField{
				Widget:     NewSelectWidget(),
				IsRequired: true,
			},
		},
	}
}

func NewRadioInt64Field() *Int64ChoiceField {
	return &Int64ChoiceField{
		&Int64Field{
			BaseField: &BaseField{
				Widget:     NewRadioWidget(),
				IsRequired: true,
			},
		},
	}
}

// ----------------------------------------------------------------------------
// BoolField
// ----------------------------------------------------------------------------

type BoolField struct {
	*BaseField
}

func (f *BoolField) Value() bool {
	if f.IValue == nil {
		return false
	}
	return f.IValue.(bool)
}

func (f *BoolField) Validate(rawValue interface{}) error {
	value := fmt.Sprint(rawValue) == "true"

	if err := f.ApplyValidators(value); err != nil {
		return err
	}

	f.IValue = value
	return nil
}

func (f *BoolField) SetInitial(initial bool) {
	f.IValue = initial
}

func (f *BoolField) Render(attrs ...string) template.HTML {
	if f.StringValue() == "true" {
		attrs = append(attrs, "checked", "checked")
	}
	return f.Widget.Render(attrs, "true")
}

func NewBoolField() *BoolField {
	return &BoolField{
		BaseField: &BaseField{
			Widget:     NewCheckboxWidget(),
			IsRequired: true,
		},
	}
}

// ----------------------------------------------------------------------------
// MultiStringChoiceField
// ----------------------------------------------------------------------------

type MultiStringChoiceField struct {
	*StringChoiceField
}

func (f *MultiStringChoiceField) Value() []string {
	if f.IValue == nil {
		return nil
	}
	return f.IValue.([]string)
}

func (f *MultiStringChoiceField) Validate(rawValue interface{}) error {
	valuesI, ok := rawValue.([]interface{})
	if !ok {
		return fmt.Errorf("Type %T is not supported.", rawValue)
	}

	values := make([]string, 0)
	for _, valueI := range valuesI {
		values = append(values, fmt.Sprint(valueI))
	}

	for _, value := range values {
		if err := f.ApplyValidators(value); err != nil {
			return err
		}
	}

	f.IValue = values
	return nil
}

func (f *MultiStringChoiceField) SetInitial(initial []string) {
	f.IValue = initial
}

func (f *MultiStringChoiceField) StringValue() []string {
	return f.Value()
}

func (f *MultiStringChoiceField) Render(attrs ...string) template.HTML {
	return f.Widget.Render(attrs, f.StringValue()...)
}

func NewMultiSelectStringField() *MultiStringChoiceField {
	return &MultiStringChoiceField{
		StringChoiceField: &StringChoiceField{
			StringField: &StringField{
				BaseField: &BaseField{
					Widget:     NewMultiSelectWidget(),
					IsRequired: true,
					IsMulti:    true,
				},
			},
		},
	}
}

// ----------------------------------------------------------------------------
// MultiInt64ChoiceField
// ----------------------------------------------------------------------------

type MultiInt64ChoiceField struct {
	*Int64ChoiceField
}

func (f *MultiInt64ChoiceField) Value() []int64 {
	if f.IValue == nil {
		return nil
	}
	return f.IValue.([]int64)
}

func (f *MultiInt64ChoiceField) Validate(rawValue interface{}) error {
	valuesI, ok := rawValue.([]interface{})
	if !ok {
		return fmt.Errorf("Type %T is not supported.", rawValue)
	}

	values := make([]int64, 0)
	for _, valueI := range valuesI {
		value, err := strconv.ParseInt(fmt.Sprint(valueI), 10, 64)
		if err != nil {
			return err
		}
		values = append(values, value)
	}

	for _, value := range values {
		if err := f.ApplyValidators(value); err != nil {
			return err
		}
	}

	f.IValue = values
	return nil
}

func (f *MultiInt64ChoiceField) SetInitial(initial []int64) {
	f.IValue = initial
}

func (f *MultiInt64ChoiceField) StringValue() []string {
	values := make([]string, 0, len(f.Value()))
	for _, value := range f.Value() {
		values = append(values, fmt.Sprint(value))
	}
	return values
}

func (f *MultiInt64ChoiceField) Render(attrs ...string) template.HTML {
	return f.Widget.Render(attrs, f.StringValue()...)
}

func NewMultiSelectInt64Field() *MultiInt64ChoiceField {
	return &MultiInt64ChoiceField{
		Int64ChoiceField: &Int64ChoiceField{
			Int64Field: &Int64Field{
				BaseField: &BaseField{
					Widget:     NewMultiSelectWidget(),
					IsRequired: true,
					IsMulti:    true,
				},
			},
		},
	}
}

// ----------------------------------------------------------------------------
// FileField
// ----------------------------------------------------------------------------

type FileField struct {
	*BaseField
}

func (f *FileField) Value() *multipart.FileHeader {
	if f.IValue == nil {
		return nil
	}
	return f.IValue.(*multipart.FileHeader)
}

func (f *FileField) Validate(rawValue interface{}) error {
	value, ok := rawValue.(*multipart.FileHeader)
	if !ok {
		return fmt.Errorf("Type %T is not supported.", rawValue)
	}

	if err := f.ApplyValidators(value); err != nil {
		return err
	}

	f.IValue = value
	return nil
}

func (f *FileField) SetInitial(initial *multipart.FileHeader) {
	f.IValue = initial
}

func (f *FileField) Render(attrs ...string) template.HTML {
	return f.Widget.Render(attrs)
}

func NewFileField() *FileField {
	return &FileField{
		BaseField: &BaseField{
			Widget:      NewFileWidget(),
			IsMultipart: true,
		},
	}
}
