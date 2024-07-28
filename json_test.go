package resteasy_test

import (
	"testing"

	"github.com/kvlach/resteasy"
)

func unmarshall(json string, v any) error {
	return resteasy.StrictUnmarshalJSON([]byte(json), v)
}

type ErrMissingJsonTag struct {
	NoTag int
}
type NoErrMissingJsonTag struct {
	NoTag int `json:"no_tag"`
}

func TestMissingJsonTag(t *testing.T) {
	var errRet ErrMissingJsonTag
	var noErrRet NoErrMissingJsonTag

	json := `{"no_tag": 123}`

	if err := unmarshall(json, &errRet); err == nil {
		t.Error("expected error")
	}
	if err := unmarshall(json, &noErrRet); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

type ErrMissingStructField struct {
	Field1 string `json:"field_1"`
}
type NoErrMissingStructField struct {
	Field1        string `json:"field_1"`
	JsonOnlyField int    `json:"json_only_field"`
}

func TestMissingStructField(t *testing.T) {
	var errRet ErrMissingStructField
	var noErrRet NoErrMissingStructField

	json := `{"field_1": "test_val", "json_only_field": 123}`

	if err := unmarshall(json, &errRet); err == nil {
		t.Error("expected error")
	}
	if err := unmarshall(json, &noErrRet); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

type MissingJsonField struct {
	Field1          string `json:"field_1"`
	StructOnlyField int    `json:"struct_only_field"`
}

func TestMissingJsonField(t *testing.T) {
	var ret MissingJsonField

	errJson := `{"field_1": "test_val"}`
	noErrJson := `{"field_1": "test_val", "struct_only_field": 123}`

	if err := unmarshall(errJson, &ret); err == nil {
		t.Error("expected error")
	}
	if err := unmarshall(noErrJson, &ret); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

type MissingSliceOfStructsJsonField struct {
	Field1         string `json:"field_1"`
	SliceOfStructs []struct {
		Field2          int    `json:"field_2"`
		StructOnlyField string `json:"struct_only_field"`
	} `json:"slice_of_structs"`
}

func TestMissingSliceOfStructsJsonField(t *testing.T) {
	var ret MissingSliceOfStructsJsonField

	errJson := `{
		"field_1": "test_val",
		"slice_of_structs": [
			{
				"field_2": 123
			}
		]
	}`
	noErrJson := `{
		"field_1": "test_val",
		"slice_of_structs": [
			{
				"field_2": 123,
				"struct_only_field": "val"
			}
		]
	}`

	if err := unmarshall(errJson, &ret); err == nil {
		t.Error("expected error")
	}
	if err := unmarshall(noErrJson, &ret); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}
