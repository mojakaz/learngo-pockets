package ecbank

import (
	"errors"
	"testing"
)

func TestDumpResponseBody(t *testing.T) {

}

func TestCheckStatusCode(t *testing.T) {
	tt := map[string]struct {
		statusCode int
		err        error
	}{
		"100": {
			statusCode: 100,
			err:        ErrUnknownStatusCode,
		},
		"200": {
			statusCode: 200,
			err:        nil,
		},
		"300": {
			statusCode: 300,
			err:        ErrUnknownStatusCode,
		},
		"400": {
			statusCode: 400,
			err:        ErrClientSide,
		},
		"500": {
			statusCode: 500,
			err:        ErrServerSide,
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			err := checkStatusCode(tc.statusCode)
			if !errors.Is(err, tc.err) {
				t.Errorf("expected %s, got %s", tc.err, err)
			}
		})
	}
}

func TestHttpStatusClass(t *testing.T) {
	tt := map[string]struct {
		statusCode int
		expected   int
	}{
		"100": {
			statusCode: 100,
			expected:   1,
		},
		"200": {
			statusCode: 200,
			expected:   2,
		},
		"300": {
			statusCode: 300,
			expected:   3,
		},
		"400": {
			statusCode: 400,
			expected:   4,
		},
		"500": {
			statusCode: 500,
			expected:   5,
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got := httpStatusClass(tc.statusCode)
			if got != tc.expected {
				t.Errorf("expected %d, got %d", tc.expected, got)
			}
		})
	}

}
