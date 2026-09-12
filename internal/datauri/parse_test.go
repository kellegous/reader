package datauri

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		Name     string
		Input    string
		Expected *URI
	}{
		{
			Name:  "image data",
			Input: "image/png;base64,aGVsbG8=",
			Expected: &URI{
				MediaType: "image/png",
				Params:    map[string]string{},
				Data:      []byte("hello"),
			},
		},
		{
			Name:  "case insensitive base64 marker",
			Input: "image/svg+xml;charset=utf-8;BASE64,PHN2Zy8+",
			Expected: &URI{
				MediaType: "image/svg+xml",
				Params:    map[string]string{"charset": "utf-8"},
				Data:      []byte("<svg/>"),
			},
		},
		{
			Name:  "default media type",
			Input: "base64,SGVsbG8=",
			Expected: &URI{
				MediaType: "text/plain",
				Params:    map[string]string{"charset": "US-ASCII"},
				Data:      []byte("Hello"),
			},
		},
		{
			Name:  "percent encoded metadata and plus encoded data",
			Input: "text/plain%3Bcharset%3Dutf-8%3Bbase64,+w%3D%3D",
			Expected: &URI{
				MediaType: "text/plain",
				Params:    map[string]string{"charset": "utf-8"},
				Data:      []byte{0xfb},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			got, err := Parse(test.Input)
			if err != nil {
				t.Fatalf("Parse(%q) returned an unexpected error: %v", test.Input, err)
			}
			if got.MediaType != test.Expected.MediaType {
				t.Errorf("MediaType = %q, want %q", got.MediaType, test.Expected.MediaType)
			}
			if !reflect.DeepEqual(got.Params, test.Expected.Params) {
				t.Errorf("Params = %#v, want %#v", got.Params, test.Expected.Params)
			}
			if !bytes.Equal(got.Data, test.Expected.Data) {
				t.Errorf("Data = %q, want %q", got.Data, test.Expected.Data)
			}
		})
	}
}

func TestParseErrors(t *testing.T) {
	tests := []struct {
		Name          string
		Input         string
		ExpectedError string
	}{
		{
			Name:          "missing comma",
			Input:         "image/png;base64",
			ExpectedError: "missing a comma",
		},
		{
			Name:          "not base64",
			Input:         "image/png,aGVsbG8=",
			ExpectedError: "not base64",
		},
		{
			Name:          "invalid escaped metadata",
			Input:         "%zz;base64,aGVsbG8=",
			ExpectedError: "invalid URL escape",
		},
		{
			Name:          "invalid media type",
			Input:         "text/plain;charset=\"unterminated;base64,aGVsbG8=",
			ExpectedError: "parse metadata",
		},
		{
			Name:          "invalid escaped data",
			Input:         "image/png;base64,%zz",
			ExpectedError: "invalid URL escape",
		},
		{
			Name:          "invalid base64 data",
			Input:         "image/png;base64,not-base64",
			ExpectedError: "illegal base64 data",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			got, err := Parse(test.Input)
			if err == nil {
				t.Fatalf("Parse(%q) = %#v, want an error", test.Input, got)
			}
			if !strings.Contains(err.Error(), test.ExpectedError) {
				t.Errorf("Parse(%q) error = %q, want it to contain %q", test.Input, err, test.ExpectedError)
			}
		})
	}
}
