package validator

import (
	"testing"
)

func TestValidate(t *testing.T) {
	req := struct {
		Name string `json:"name" validate:"required" label:"姓名"`
		Age  int    `json:"age" validate:"gte=10" label:"年龄"`
	}{
		Name: "jaronnie",
		Age:  9,
	}

	err := Validate(req)
	if err != nil {
		t.Error(err)
	}

	req2 := struct {
		Name string `path:"name" validate:"required"`
		Age  int    `json:"age" validate:"gte=10"`
	}{
		Name: "",
		Age:  9,
	}

	err = Validate(req2, "zh_Hans_CN")
	if err != nil {
		t.Error(err)
	}
}
