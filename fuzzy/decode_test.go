package fuzzy

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFuzzyDecodeRequest(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		outputType  any
		expectError bool
	}{
		{
			name:       "simple struct with string field",
			input:      []byte(`{"name": "John"}`),
			outputType: &struct{ Name string }{},
		},
		{
			name:       "struct with int field from string",
			input:      []byte(`{"age": "30"}`),
			outputType: &struct{ Age int }{},
		},
		{
			name:       "struct with bool field from string",
			input:      []byte(`{"active": "true"}`),
			outputType: &struct{ Active bool }{},
		},
		{
			name:       "struct with *bool field from string",
			input:      []byte(`{"active": "true"}`),
			outputType: &struct{ Active *bool }{},
		},
		{
			name:       "struct with float field from string",
			input:      []byte(`{"price": "19.99"}`),
			outputType: &struct{ Price float64 }{},
		},
		{
			name:  "nested struct with fuzzy values",
			input: []byte(`{"user": {"id": "123", "premium": "1"}}`),
			outputType: &struct {
				User struct {
					ID      int
					Premium bool
				}
			}{},
		},
		{
			name:       "empty input",
			input:      []byte(`{}`),
			outputType: &struct{}{},
		},
		{
			name:       "pointer fields",
			input:      []byte(`{"ptr": "123"}`),
			outputType: &struct{ Ptr *int }{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decode, err := Decode(tt.input, tt.outputType)
			if err != nil {
				return
			}
			t.Log(tt.outputType)
			err = json.Unmarshal(decode, tt.outputType)
			assert.NoError(t, err)
		})
	}
}

func TestFuzzyDecodeRequest_EdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		input      []byte
		outputType any
		check      func(t *testing.T, output any)
	}{
		{
			name:       "null to zero value",
			input:      []byte(`{"value": null}`),
			outputType: &struct{ Value int }{},
			check: func(t *testing.T, output any) {
				t.Helper()
				assert.Equal(t, 0, output.(*struct{ Value int }).Value)
			},
		},
		{
			name:       "empty string to zero value",
			input:      []byte(`{"value": ""}`),
			outputType: &struct{ Value int }{},
			check: func(t *testing.T, output any) {
				t.Helper()
				assert.Equal(t, 0, output.(*struct{ Value int }).Value)
			},
		},
		{
			name:  "string true/false to bool",
			input: []byte(`{"trueVal": "true", "falseVal": "false"}`),
			outputType: &struct {
				TrueVal  bool
				FalseVal bool
			}{},
			check: func(t *testing.T, output any) {
				t.Helper()
				out := output.(*struct {
					TrueVal  bool
					FalseVal bool
				})
				assert.True(t, out.TrueVal)
				assert.False(t, out.FalseVal)
			},
		},
		{
			name:  "string numbers to bool",
			input: []byte(`{"trueVal": "1", "falseVal": "0"}`),
			outputType: &struct {
				TrueVal  bool
				FalseVal bool
			}{},
			check: func(t *testing.T, output any) {
				t.Helper()
				out := output.(*struct {
					TrueVal  bool
					FalseVal bool
				})
				assert.True(t, out.TrueVal)
				assert.False(t, out.FalseVal)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := tt.outputType
			_, err := Decode(tt.input, output)
			assert.NoError(t, err)
			tt.check(t, output)
		})
	}
}
