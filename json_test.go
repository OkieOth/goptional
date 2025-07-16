package goptional_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/okieoth/goptional"
)

type Demo2 struct {
	Id            int64                                `json:"id,omitempty"`
	S             goptional.Optional[string]           `json:"s,omitempty"`
	S2            string                               `json:"s2,omitempty"`
	List          []DummyEnum                          `json:"list,omitempty"`
	OptionalMap   goptional.Optional[map[string]int32] `json:"optionalMap,omitempty"`
	OptionalList  goptional.Optional[[]string]         `json:"optionalList,omitempty"`
	OptionalList2 goptional.Optional[[][]DummyEnum]    `json:"optionalList2,omitempty"`
	OptionalEnum  goptional.OptionalEnum[DummyEnum]    `json:"optionalEnum,omitempty"`
}

func TestJson(t *testing.T) {
	var d Demo2
	d.Id = 66
	d.S2 = "a string"

	if jsonData, err := json.MarshalIndent(d, "", "  "); err == nil {
		fmt.Println(jsonData)
	} else {
		t.Errorf("Error while marshal to JSON: %v", err)
	}
}
