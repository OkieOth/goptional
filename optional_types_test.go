package goptional_test

import (
	"fmt"
	"testing"

	"github.com/okieoth/goptional"
)

func TestOptionalString(t *testing.T) {
	var s goptional.Optional[string]
	if s.IsSet() {
		t.Errorf("uninitialized value is set after creation")
	}
}

type DummyEnum int64

const (
	value1 DummyEnum = iota
	value2
	value3
	value4
	value5
	value6
)

func (s DummyEnum) String() string {
	switch s {
	case value1:
		return "value1"
	case value2:
		return "value2"
	case value3:
		return "value3"
	case value4:
		return "value4"
	case value5:
		return "value5"
	default:
		return "value6"
	}
}

func (s DummyEnum) ValueFromStr(str string) error {
	switch str {
	case "value1":
		s = value1
	case "value2":
		s = value2
	case "value3":
		s = value3
	case "value4":
		s = value4
	case "value5":
		s = value5
	case "value6":
		s = value6
	default:
		return fmt.Errorf("input not part of the enum: %v", str)
	}
	return nil
}

func TestOptionalEnum(t *testing.T) {
	var s goptional.OptionalEnum[DummyEnum]
	if s.IsSet() {
		t.Errorf("uninitialized value is set after creation")
	}
}

type ListCont struct {
	L1 []string
	L2 goptional.Optional[[]string]
}

type ListContBuilder struct {
	l1 []string
	l2 goptional.Optional[[]string]
}

func (b *ListContBuilder) L1(v []string) *ListContBuilder {
	b.l1 = v
	return b
}

func (b *ListContBuilder) L2(v []string) *ListContBuilder {
	b.l2.Set(v)
	return b
}

func (b *ListContBuilder) Build() ListCont {
	var ret ListCont
	ret.L1 = b.l1
	if b.l2.IsSet() {
		ret.L2 = b.l2
	}
	return ret
}

func NewListContBuilder() *ListContBuilder {
	var lb ListContBuilder
	return &lb
}

func TestList(t *testing.T) {
	var a []string
	t1 := NewListContBuilder().
		L1(a).
		Build()

	if len(t1.L1) != 0 {
		t.Errorf("t1 doesn't have len 0")
	}

	t1.L1 = append(t1.L1, "test")

	if len(t1.L1) != 1 {
		t.Errorf("t1 doesn't contain a value")
	}

	if t1.L2.IsSet() {
		t.Errorf("t1.L2 wrongly set")
	}

	t1.L2.Set(a)

	if v, isSet := t1.L2.Get(); isSet {
		v := append(v, "test2")
		if v2, isSet := t1.L2.Get(); isSet {
			if len(v2) != 0 {
				t.Errorf("v2 has wrong value")
			}
		} else {
			t.Errorf("t1.L2 couldn't get value (1)")
		}
		t1.L2.Set(v)
		if v2, isSet := t1.L2.Get(); isSet {
			if len(v2) != 1 {
				t.Errorf("v2 has wrong value")
			}
		} else {
			t.Errorf("t1.L2 couldn't get value (1)")
		}
	} else {
		t.Errorf("t1.L2 couldn't get value (2)")
	}

	if !t1.L2.IsSet() {
		t.Errorf("t1.L2 should be set")
	}

}
