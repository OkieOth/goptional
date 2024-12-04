package goptional

import (
	"encoding/json"
)

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
	return Optional[C]{
		isSet: false,
	}
}

func (m *Optional[C]) Set(v C) {
	m.Value = v
	m.isSet = true
}

func (m *Optional[C]) UnSet() {
	m.isSet = false
}

func (m *Optional[C]) Get() (C, bool) {
	return m.Value, m.isSet
}

func (m *Optional[C]) IsSet() bool {
	return m.isSet
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

func (v *Optional[C]) MarshalJSON() ([]byte, error) {
	if value, isSet := v.Get(); isSet {
		return json.Marshal(value)
	} else {
		return []byte("null"), nil
	}
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
	return OptionalEnum[C]{
		isSet: false,
	}
}

func (m *OptionalEnum[C]) Set(v C) {
	m.Value = v
	m.isSet = true
}

func (m *OptionalEnum[C]) UnSet() {
	m.isSet = false
}

func (m *OptionalEnum[C]) Get() (C, bool) {
	return m.Value, m.isSet
}

func (m *OptionalEnum[C]) IsSet() bool {
	return m.isSet
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

func (v *OptionalEnum[C]) MarshalJSON() ([]byte, error) {
	if value, isSet := v.Get(); isSet {
		return json.Marshal(value)
	} else {
		return []byte("null"), nil
	}
}
