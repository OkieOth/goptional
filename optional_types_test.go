package goptional_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/okieoth/goptional/v3"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestOptionalString(t *testing.T) {
	var s goptional.Optional[string]
	require.False(t, s.IsSet())

	s = s.Set("test")
	v, b := s.Get()
	require.True(t, b)
	require.Equal(t, "test", v)
	s1 := s.IfSetThenDo(func(val string) (string, bool) {
		require.Equal(t, "test", v, "setting a value failed")
		return val, false
	})
	require.Equal(t, s, s1)
	s2 := s1.IfSetThenDo(func(val string) (string, bool) {
		require.Equal(t, "test", v, "setting a value failed")
		return "newTest", true
	})
	require.NotEqual(t, s1, s2)
	v, b = s2.Get()
	require.True(t, b)
	require.Equal(t, "newTest", v)
}

type Demo struct {
	Id            int64                                `yaml:"id,omitempty" json:"id,omitempty"`
	List          []DummyEnum                          `yaml:"list,omitempty" json:"list,omitempty"`
	OptionalMap   goptional.Optional[map[string]int32] `yaml:"omap,omitempty" json:"omap,omitzero"`
	OptionalList  goptional.Optional[[]string]         `yaml:"olist,omitempty" json:"olist,omitzero"`
	OptionalList2 goptional.Optional[[][]DummyEnum]    `yaml:"olist2,omitempty" json:"olist2,omitzero"`
	OptionalEnum  goptional.OptionalEnum[DummyEnum]    `yaml:"oenum,omitempty" json:"oenum,omitzero"`
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
	require.False(t, s.IsSet())

	s = s.Set(value3)
	v, b := s.Get()
	require.True(t, b)
	require.Equal(t, value3, v)
	s1 := s.IfSetThenDo(func(val DummyEnum) (DummyEnum, bool) {
		require.Equal(t, value3, v, "setting a value failed")
		return val, false
	})
	require.Equal(t, s, s1)
	s2 := s1.IfSetThenDo(func(val DummyEnum) (DummyEnum, bool) {
		require.Equal(t, value3, v, "setting a value failed")
		return value5, true
	})
	require.NotEqual(t, s1, s2)
	v, b = s2.Get()
	require.True(t, b)
	require.Equal(t, value5, v)
}

func TestJson(t *testing.T) {
	tests := []struct {
		value    Demo
		destFile string
	}{
		{
			value: Demo{
				Id:   13,
				List: []DummyEnum{value1, value3},
			},
			destFile: "temp/demo_01.json",
		},
		{
			value: Demo{
				Id:           13,
				List:         []DummyEnum{value1, value3},
				OptionalList: goptional.NewOptionalValue([]string{"eins", "zwei", "drei"}),
			},
			destFile: "temp/demo_02.json",
		},
	}
	for _, test := range tests {
		if _, err := os.Stat(test.destFile); err != nil {
			os.Remove(test.destFile)
			bOut, err := json.MarshalIndent(test.value, "", "  ")
			require.Nil(t, err)
			err = os.WriteFile(test.destFile, bOut, 0644)
			require.Nil(t, err)
			bIn, err := os.ReadFile(test.destFile)
			require.Nil(t, err)
			var d Demo
			err = json.Unmarshal(bIn, &d)
			require.Nil(t, err)
			require.Equal(t, test.value, d)
		}
	}

	//    Id            int64                                `yaml:"id,omitempty" json:"id,omitempty"`
	// List          []DummyEnum                          `yaml:"list,omitempty" json:"list,omitempty"`
	// OptionalMap   goptional.Optional[map[string]int32] `yaml:"omap,omitempty" json:"omap,omitempty"`
	// OptionalList  goptional.Optional[[]string]         `yaml:"olist,omitempty" json:"olist,omitempty"`
	// OptionalList2 goptional.Optional[[][]DummyEnum]    `yaml:"olist2,omitempty" json:"olist2,omitempty"`
	// OptionalEnum  goptional.OptionalEnum[DummyEnum]    `yaml:"oenum,omitempty" json:"oenum,omitempty"`

}

func TestYaml(t *testing.T) {
	tests := []struct {
		value    Demo
		destFile string
	}{
		{
			value: Demo{
				Id:   13,
				List: []DummyEnum{value1, value3},
			},
			destFile: "temp/demo_01.yaml",
		},
		{
			value: Demo{
				Id:           13,
				List:         []DummyEnum{value1, value3},
				OptionalList: goptional.NewOptionalValue([]string{"eins", "zwei", "drei"}),
			},
			destFile: "temp/demo_02.yaml",
		},
	}
	for _, test := range tests {
		if _, err := os.Stat(test.destFile); err != nil {
			os.Remove(test.destFile)
			bOut, err := yaml.Marshal(test.value)
			require.Nil(t, err)
			err = os.WriteFile(test.destFile, bOut, 0644)
			require.Nil(t, err)
			bIn, err := os.ReadFile(test.destFile)
			require.Nil(t, err)
			var d Demo
			err = yaml.Unmarshal(bIn, &d)
			require.Nil(t, err)
			require.Equal(t, test.value, d)
		}
	}
}

func TestNewOptionalConditional(t *testing.T) {
	o1 := goptional.NewOptionalConditional(10, func(i int) bool { return i != 10 })
	_, isSet := o1.Get()
	require.False(t, isSet)
	o2 := goptional.NewOptionalConditional(10, func(i int) bool { return i != 0 })
	i2, isSet := o2.Get()
	require.True(t, isSet)
	require.Equal(t, 10, i2)
	o3 := goptional.NewOptionalConditional("", func(i string) bool { return i != "" })
	_, isSet = o3.Get()
	require.False(t, isSet)
	o4 := goptional.NewOptionalConditional("xxx", func(i string) bool { return i != "" })
	s, isSet := o4.Get()
	require.Equal(t, "xxx", s)
}
