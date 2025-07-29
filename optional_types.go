package goptional

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

type WithFunc[T any] func(value T) (T, bool)

type EnumType interface {
	String() string
	ValueFromStr(str string) error
}

type Optional[C any] struct {
	Value C
	isSet bool
}

func NewOptionalValue[C any](v C) Optional[C] {
	return Optional[C]{
		Value: v,
		isSet: true,
	}
}

func NewOptional[C any]() Optional[C] {
	return Optional[C]{}
}

func NewOptionalConditional[C any](v C, condFunc func(v C) bool) Optional[C] {
	if condFunc(v) {
		return NewOptionalValue(v)
	}
	return NewOptional[C]()
}

func (m Optional[C]) Set(v C) Optional[C] {
	return Optional[C]{
		Value: v,
		isSet: true,
	}
}

func (m Optional[C]) UnSet() Optional[C] {
	return Optional[C]{}
}

func (m Optional[C]) Get() (C, bool) {
	return m.Value, m.isSet
}

func (m Optional[C]) IfSetThenDo(fn WithFunc[C]) Optional[C] {
	if newValue, changed := fn(m.Value); changed {
		m.Value = newValue
		m.isSet = true
		return m
	}
	return m
}

func (m Optional[C]) IsSet() bool {
	return m.isSet
}

func (v Optional[C]) IsZero() bool {
	return v.isSet == false
}

func (v *Optional[C]) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		v.isSet = false
		return nil
	}

	err := json.Unmarshal(data, &v.Value)
	if err != nil {
		return err
	}

	v.isSet = true
	return nil
}

func (v Optional[C]) MarshalJSON() ([]byte, error) {
	if value, isSet := v.Get(); isSet {
		return json.Marshal(value)
	} else {
		return []byte("null"), nil
	}
}

func (o *Optional[C]) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode && value.Tag == "!!null" {
		var zero C
		o.Value = zero
		o.isSet = false
		return nil
	}

	var actual C
	if err := value.Decode(&actual); err != nil {
		return err
	}

	o.Value = actual
	o.isSet = true
	return nil
}

func (o Optional[C]) MarshalYAML() (any, error) {
	if value, isSet := o.Get(); isSet {
		return value, nil
	}
	return nil, nil
}

type OptionalEnum[C EnumType] struct {
	Value C
	isSet bool
}

func NewOptionalEnumValue[C EnumType](v C) OptionalEnum[C] {
	return OptionalEnum[C]{
		Value: v,
		isSet: true,
	}
}

func NewOptionalEnum[C EnumType]() OptionalEnum[C] {
	return OptionalEnum[C]{}
}

func (m OptionalEnum[C]) Set(v C) OptionalEnum[C] {
	return OptionalEnum[C]{
		Value: v,
		isSet: true,
	}
}

func (m OptionalEnum[C]) UnSet() OptionalEnum[C] {
	return OptionalEnum[C]{}
}

func (m OptionalEnum[C]) Get() (C, bool) {
	return m.Value, m.isSet
}

func (m OptionalEnum[C]) IfSetThenDo(fn WithFunc[C]) OptionalEnum[C] {
	if newValue, changed := fn(m.Value); changed {
		m.Value = newValue
		m.isSet = true
		return m
	}
	return m
}

func (m *OptionalEnum[C]) IsSet() bool {
	return m.isSet
}

func (v OptionalEnum[C]) IsZero() bool {
	return v.isSet == false
}

func (v *OptionalEnum[C]) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	var s string
	err := json.Unmarshal(data, &s)

	if err != nil {
		return err
	}
	err = v.Value.ValueFromStr(s)
	if err != nil {
		return err
	}
	v.isSet = true
	return nil
}

func (v OptionalEnum[C]) MarshalJSON() ([]byte, error) {
	if value, isSet := v.Get(); isSet {
		return json.Marshal(value)
	} else {
		return []byte("null"), nil
	}
}

func (o *OptionalEnum[C]) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode && value.Tag == "!!null" {
		var zero C
		o.Value = zero
		o.isSet = false
		return nil
	}

	var actual C
	if err := value.Decode(&actual); err != nil {
		return err
	}

	o.Value = actual
	o.isSet = true
	return nil
}

func (o OptionalEnum[C]) MarshalYAML() (any, error) {
	if value, isSet := o.Get(); isSet {
		return value, nil
	}
	return nil, nil
}
