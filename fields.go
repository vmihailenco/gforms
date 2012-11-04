package gforms

import (
	"errors"
	"fmt"
	"html/template"
	"mime/multipart"
	"reflect"
	"strconv"
)

var (
	ErrRequired = errors.New("This field is required")
)

//------------------------------------------------------------------------------

type Field interface {
	Validator

	HasName() bool
	SetName(string)
	Name() string

	HasLabel() bool
	SetLabel(string)
	Label() string

	SetWidget(Widget)
	Widget() Widget

	SetIsMulti(bool)
	IsMulti() bool
	SetIsMultipart(bool)
	IsMultipart() bool
	SetIsRequired(bool)
	IsRequired() bool

	AddValidator(Validator)
	ApplyValidators(interface{}) error

	HasValidationError() bool
	SetValidationError(error)
	ValidationError() error

	Reset()
	Render(...string) template.HTML
}

func isEmpty(value interface{}) bool {
	v := reflect.ValueOf(value)

	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

func IsFieldValid(f Field, rawValue interface{}) bool {
	f.Reset()

	if rawValue == nil || isEmpty(rawValue) {
		if f.IsRequired() {
			f.SetValidationError(ErrRequired)
			return false
		} else {
			return true
		}
	}

	if err := f.Validate(rawValue); err != nil {
		f.SetValidationError(err)
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

//------------------------------------------------------------------------------

type BaseField struct {
	hasName  bool
	name     string
	hasLabel bool
	label    string
	widget   Widget

	isMulti     bool
	isMultipart bool
	isRequired  bool

	validators      []Validator
	validationError error
	iValue          interface{}
}

func (f *BaseField) HasName() bool {
	return f.hasName
}

func (f *BaseField) SetName(name string) {
	f.hasName = true
	f.name = name
	attrs := f.Widget().Attrs()
	attrs.Set("id", name)
	attrs.Set("name", name)
}

func (f *BaseField) Name() string {
	return f.name
}

func (f *BaseField) HasLabel() bool {
	return f.hasLabel
}

func (f *BaseField) SetLabel(label string) {
	f.hasLabel = true
	f.label = label
}

func (f *BaseField) Label() string {
	return f.label
}

func (f *BaseField) SetWidget(widget Widget) {
	f.widget = widget
}

func (f *BaseField) Widget() Widget {
	return f.widget
}

func (f *BaseField) SetIsMulti(flag bool) {
	f.isMulti = flag
}

func (f *BaseField) IsMulti() bool {
	return f.isMulti
}

func (f *BaseField) SetIsMultipart(flag bool) {
	f.isMultipart = flag
}

func (f *BaseField) IsMultipart() bool {
	return f.isMultipart
}

func (f *BaseField) SetIsRequired(flag bool) {
	f.isRequired = flag
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

func (f *BaseField) ApplyValidators(rawValue interface{}) error {
	for _, validator := range f.validators {
		if err := validator.Validate(rawValue); err != nil {
			return err
		}
	}
	return nil
}

func (f *BaseField) validate(rawValue interface{}) error {
	panic("not implemented.")
}

func (f *BaseField) HasValidationError() bool {
	return f.validationError != nil
}

func (f *BaseField) SetValidationError(err error) {
	f.validationError = err
}

func (f *BaseField) ValidationError() error {
	return f.validationError
}

func (f *BaseField) StringValue() string {
	if f.iValue == nil {
		return ""
	}
	return fmt.Sprint(f.iValue)
}

func (f *BaseField) Reset() {
	f.iValue = nil
	f.validationError = nil
}

func (f *BaseField) Render(attrs ...string) template.HTML {
	panic("not implemented")
}

//------------------------------------------------------------------------------

type StringField struct {
	*BaseField
	MinLen, MaxLen int
}

func (f *StringField) Value() string {
	if f.iValue == nil {
		return ""
	}
	return f.iValue.(string)
}

func (f *StringField) Validate(rawValue interface{}) error {
	value := fmt.Sprint(rawValue)

	valueLen := len(value)
	if f.MinLen > 0 && valueLen < f.MinLen {
		return fmt.Errorf("This field should have at least %d symbols", f.MinLen)
	}
	if f.MaxLen > 0 && valueLen > f.MaxLen {
		return fmt.Errorf("This field should have less than %d symbols", f.MaxLen)
	}

	if err := f.ApplyValidators(value); err != nil {
		return err
	}

	f.iValue = value
	return nil
}

func (f *StringField) SetInitial(initial string) {
	f.iValue = initial
}

func (f *StringField) Render(attrs ...string) template.HTML {
	return f.Widget().Render(attrs, f.StringValue())
}

func NewStringField() *StringField {
	return &StringField{
		BaseField: &BaseField{
			widget: NewTextWidget(),
		},
	}
}

type TextareaStringField struct {
	*StringField
}

func NewTextareaStringField() *TextareaStringField {
	return &TextareaStringField{
		StringField: &StringField{
			BaseField: &BaseField{
				widget: NewTextareaWidget(),
			},
		},
	}
}

//------------------------------------------------------------------------------

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
				widget: NewSelectWidget(),
			},
		},
	}
}

func NewRadioStringField() *StringChoiceField {
	return &StringChoiceField{
		&StringField{
			BaseField: &BaseField{
				widget: NewRadioWidget(),
			},
		},
	}
}

// ----------------------------------------------------------------------------

type Int64Field struct {
	*BaseField
}

func (f *Int64Field) Value() int64 {
	if f.iValue == nil {
		return 0
	}
	return f.iValue.(int64)
}

func (f *Int64Field) Validate(rawValue interface{}) error {
	value, err := strconv.ParseInt(fmt.Sprint(rawValue), 10, 64)
	if err != nil {
		return err
	}

	if err := f.ApplyValidators(value); err != nil {
		return err
	}

	f.iValue = value
	return nil
}

func (f *Int64Field) SetInitial(initial int64) {
	f.iValue = initial
}

func (f *Int64Field) Render(attrs ...string) template.HTML {
	return f.Widget().Render(attrs, f.StringValue())
}

func NewInt64Field() *Int64Field {
	return &Int64Field{
		BaseField: &BaseField{
			widget: NewTextWidget(),
		},
	}
}

//------------------------------------------------------------------------------

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
				widget: NewSelectWidget(),
			},
		},
	}
}

func NewRadioInt64Field() *Int64ChoiceField {
	return &Int64ChoiceField{
		&Int64Field{
			BaseField: &BaseField{
				widget: NewRadioWidget(),
			},
		},
	}
}

//------------------------------------------------------------------------------

type BoolField struct {
	*BaseField
}

func (f *BoolField) Value() bool {
	if f.iValue == nil {
		return false
	}
	return f.iValue.(bool)
}

func (f *BoolField) Validate(rawValue interface{}) error {
	value := fmt.Sprint(rawValue) == "true"

	if err := f.ApplyValidators(value); err != nil {
		return err
	}

	f.iValue = value
	return nil
}

func (f *BoolField) SetInitial(initial bool) {
	f.iValue = initial
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
			widget: NewCheckboxWidget(),
		},
	}
}

//------------------------------------------------------------------------------

type MultiStringChoiceField struct {
	*StringChoiceField
}

func (f *MultiStringChoiceField) Value() []string {
	if f.iValue == nil {
		return nil
	}
	return f.iValue.([]string)
}

func (f *MultiStringChoiceField) Validate(rawValue interface{}) error {
	valuesI, ok := rawValue.([]interface{})
	if !ok {
		return fmt.Errorf("Type %T is not supported", rawValue)
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

	f.iValue = values
	return nil
}

func (f *MultiStringChoiceField) SetInitial(initial []string) {
	f.iValue = initial
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
					widget:  NewMultiSelectWidget(),
					isMulti: true,
				},
			},
		},
	}
}

//------------------------------------------------------------------------------

type MultiInt64ChoiceField struct {
	*Int64ChoiceField
}

func (f *MultiInt64ChoiceField) Value() []int64 {
	if f.iValue == nil {
		return nil
	}
	return f.iValue.([]int64)
}

func (f *MultiInt64ChoiceField) Validate(rawValue interface{}) error {
	valuesI, ok := rawValue.([]interface{})
	if !ok {
		return fmt.Errorf("Type %T is not supported", rawValue)
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

	f.iValue = values
	return nil
}

func (f *MultiInt64ChoiceField) SetInitial(initial []int64) {
	f.iValue = initial
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
					widget:  NewMultiSelectWidget(),
					isMulti: true,
				},
			},
		},
	}
}

//------------------------------------------------------------------------------

type FileField struct {
	*BaseField
}

func (f *FileField) Value() *multipart.FileHeader {
	if f.iValue == nil {
		return nil
	}
	return f.iValue.(*multipart.FileHeader)
}

func (f *FileField) Validate(rawValue interface{}) error {
	value, ok := rawValue.(*multipart.FileHeader)
	if !ok {
		return fmt.Errorf("Type %T is not supported", rawValue)
	}

	if err := f.ApplyValidators(value); err != nil {
		return err
	}

	f.iValue = value
	return nil
}

func (f *FileField) SetInitial(initial *multipart.FileHeader) {
	f.iValue = initial
}

func (f *FileField) Render(attrs ...string) template.HTML {
	return f.widget.Render(attrs)
}

func NewFileField() *FileField {
	return &FileField{
		BaseField: &BaseField{
			widget:      NewFileWidget(),
			isMultipart: true,
		},
	}
}
