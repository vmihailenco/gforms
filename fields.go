package gforms

import (
	"errors"
	"fmt"
	"html/template"
	"strconv"
)

type Field interface {
	Validator

	SetName(string)
	Name() string
	SetLabel(string)
	Label() string
	IsMulti() bool
	SetIsRequired(bool)
	IsRequired() bool
	AddValidator(Validator)
	Render(...string) template.HTML
	SetValidationError(error)
	ValidationError() error
	HasValidationError() bool
	SetWidget(Widget)
	Widget() Widget
}

func IsFieldValid(field Field, rawValue interface{}) bool {
	if rawValue == nil {
		if field.IsRequired() {
			field.SetValidationError(errors.New("This field is required."))
			return false
		} else {
			return true
		}
	}

	if err := field.Validate(rawValue); err != nil {
		field.SetValidationError(err)
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

type BaseField struct {
	name       string
	label      string
	widget     Widget
	isMulti    bool
	isRequired bool
	validators []Validator
	error      error
	value      interface{}
}

func (f *BaseField) SetName(name string) {
	f.name = name
}

func (f *BaseField) Name() string {
	return f.name
}

func (f *BaseField) SetLabel(label string) {
	f.label = label
}

func (f *BaseField) Label() string {
	return f.label
}

func (f *BaseField) IsMulti() bool {
	return f.isMulti
}

func (f *BaseField) SetIsRequired(isRequired bool) {
	f.isRequired = isRequired
}

func (f *BaseField) IsRequired() bool {
	return f.isRequired
}

func (f *BaseField) AddValidator(validator Validator) {
	if f.validators == nil {
		f.validators = make([]Validator, 0)
	}
	f.validators = append(f.validators, validator)
}

func (f *BaseField) applyValidators(rawValue interface{}) error {
	for _, validator := range f.validators {
		if err := validator.Validate(rawValue); err != nil {
			return err
		}
	}
	return nil
}

func (f *BaseField) SetValidationError(err error) {
	f.error = err
}

func (f *BaseField) ValidationError() error {
	return f.error
}

func (f *BaseField) HasValidationError() bool {
	return f.error != nil
}

func (f *BaseField) StringValue() string {
	if f.value == nil {
		return ""
	}
	return fmt.Sprint(f.value)
}

func (f *BaseField) Render(attrs ...string) template.HTML {
	panic("not implemented.")
	return template.HTML("")
}

func (f *BaseField) SetWidget(w Widget) {
	f.widget = w
}

func (f *BaseField) Widget() Widget {
	return f.widget
}

type StringField struct {
	*BaseField
	minLen, maxLen int
}

func (f *StringField) SetMinLen(minLen int) {
	f.minLen = minLen
}

func (f *StringField) SetMaxLen(maxLen int) {
	f.maxLen = maxLen
}

func (f *StringField) Value() string {
	if f.value == nil {
		return ""
	}
	return f.value.(string)
}

func (f *StringField) Validate(rawValue interface{}) error {
	value := fmt.Sprint(rawValue)

	valueLen := len(value)
	if f.minLen > 0 && valueLen < f.minLen {
		return fmt.Errorf("This field should have at least %d symbols.", f.minLen)
	}
	if f.maxLen > 0 && valueLen > f.maxLen {
		return fmt.Errorf("This field should have less than %d symbols.", f.maxLen)
	}

	if err := f.BaseField.applyValidators(value); err != nil {
		return err
	}

	f.value = value
	return nil
}

func (f *StringField) SetInitial(initial string) {
	f.value = initial
}

func (f *StringField) Render(attrs ...string) template.HTML {
	return f.Widget().Render(attrs, f.StringValue())
}

func NewStringField() *StringField {
	return &StringField{
		BaseField: &BaseField{
			widget:     NewTextWidget(),
			isRequired: true,
		},
	}
}

func NewTextareaStringField() *StringField {
	return &StringField{
		BaseField: &BaseField{
			widget:     NewTextareaWidget(),
			isRequired: true,
		},
	}
}

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

	f.Widget().(ChoiceWidget).SetChoices(strChoices)
	f.AddValidator(NewStringChoicesValidator(choices))
}

func NewSelectStringField() *StringChoiceField {
	return &StringChoiceField{
		StringField: &StringField{
			BaseField: &BaseField{
				widget:     NewSelectWidget(),
				isRequired: true,
			},
		},
	}
}

func NewRadioStringField() *StringChoiceField {
	return &StringChoiceField{
		&StringField{
			BaseField: &BaseField{
				widget:     NewRadioWidget(),
				isRequired: true,
			},
		},
	}
}

type Int64Field struct {
	*BaseField
}

func (f *Int64Field) Value() int64 {
	if f.value == nil {
		return 0
	}
	return f.value.(int64)
}

func (f *Int64Field) Validate(rawValue interface{}) error {
	value, err := strconv.ParseInt(fmt.Sprint(rawValue), 10, 64)
	if err != nil {
		return err
	}

	if err := f.BaseField.applyValidators(value); err != nil {
		return err
	}

	f.value = value
	return nil
}

func (f *Int64Field) SetInitial(initial int64) {
	f.value = initial
}

func (f *Int64Field) Render(attrs ...string) template.HTML {
	return f.Widget().Render(attrs, f.StringValue())
}

func NewInt64Field() *Int64Field {
	return &Int64Field{
		BaseField: &BaseField{
			widget:     NewTextWidget(),
			isRequired: true,
		},
	}
}

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

	f.Widget().(ChoiceWidget).SetChoices(strChoices)
	f.AddValidator(NewInt64ChoicesValidator(choices))
}

func NewSelectInt64Field() *Int64ChoiceField {
	return &Int64ChoiceField{
		Int64Field: &Int64Field{
			BaseField: &BaseField{
				widget:     NewSelectWidget(),
				isRequired: true,
			},
		},
	}
}

func NewRadioInt64Field() *Int64ChoiceField {
	return &Int64ChoiceField{
		&Int64Field{
			BaseField: &BaseField{
				widget:     NewRadioWidget(),
				isRequired: true,
			},
		},
	}
}

type BoolField struct {
	*BaseField
}

func (f *BoolField) Value() bool {
	if f.value == nil {
		return false
	}
	return f.value.(bool)
}

func (f *BoolField) Validate(rawValue interface{}) error {
	value := fmt.Sprint(rawValue) == "true"

	if err := f.BaseField.applyValidators(value); err != nil {
		return err
	}

	f.value = value
	return nil
}

func (f *BoolField) SetInitial(initial bool) {
	f.value = initial
}

func (f *BoolField) Render(attrs ...string) template.HTML {
	if f.StringValue() == "true" {
		attrs = append(attrs, "checked", "checked")
	}
	return f.Widget().Render(attrs, "true")
}

func NewBoolField() *BoolField {
	return &BoolField{
		BaseField: &BaseField{
			widget:     NewCheckboxWidget(),
			isRequired: true,
		},
	}
}

type MultiStringChoiceField struct {
	*StringChoiceField
}

func (f *MultiStringChoiceField) Value() []string {
	if f.value == nil {
		return nil
	}
	return f.value.([]string)
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
		if err := f.BaseField.applyValidators(value); err != nil {
			return err
		}
	}

	f.value = values
	return nil
}

func (f *MultiStringChoiceField) SetInitial(initial []string) {
	f.value = initial
}

func (f *MultiStringChoiceField) StringValue() []string {
	return f.Value()
}

func (f *MultiStringChoiceField) Render(attrs ...string) template.HTML {
	return f.Widget().Render(attrs, f.StringValue()...)
}

func NewMultiSelectStringField() *MultiStringChoiceField {
	return &MultiStringChoiceField{
		StringChoiceField: &StringChoiceField{
			StringField: &StringField{
				BaseField: &BaseField{
					widget:     NewMultiSelectWidget(),
					isRequired: true,
					isMulti:    true,
				},
			},
		},
	}
}

type MultiInt64ChoiceField struct {
	*Int64ChoiceField
}

func (f *MultiInt64ChoiceField) Value() []int64 {
	if f.value == nil {
		return nil
	}
	return f.value.([]int64)
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
		if err := f.BaseField.applyValidators(value); err != nil {
			return err
		}
	}

	f.value = values
	return nil
}

func (f *MultiInt64ChoiceField) SetInitial(initial []int64) {
	f.value = initial
}

func (f *MultiInt64ChoiceField) StringValue() []string {
	values := make([]string, 0, len(f.Value()))
	for _, value := range f.Value() {
		values = append(values, fmt.Sprint(value))
	}
	return values
}

func (f *MultiInt64ChoiceField) Render(attrs ...string) template.HTML {
	return f.Widget().Render(attrs, f.StringValue()...)
}

func NewMultiSelectInt64Field() *MultiInt64ChoiceField {
	return &MultiInt64ChoiceField{
		Int64ChoiceField: &Int64ChoiceField{
			Int64Field: &Int64Field{
				BaseField: &BaseField{
					widget:     NewMultiSelectWidget(),
					isRequired: true,
					isMulti:    true,
				},
			},
		},
	}
}
