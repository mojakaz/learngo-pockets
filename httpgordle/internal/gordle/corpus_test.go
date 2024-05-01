package gordle

import (
	_ "embed"
	"errors"
	"testing"
)

//go:embed testdata/english.txt
var englishCorpus string

//go:embed testdata/empty.txt
var emptyCorpus string

func TestReadCorpus(t *testing.T) {
	tt := map[string]struct {
		data   string
		length int
		err    error
	}{
		"English corpus": {
			data:   englishCorpus,
			length: 156,
			err:    nil,
		},
		"empty corpus": {
			data:   emptyCorpus,
			length: 0,
			err:    ErrCorpusIsEmpty,
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			words, err := ReadCorpus(tc.data)
			if len(words) != tc.length {
				t.Errorf("different words length: got %d, want %d", len(words), tc.length)
			}
			if !errors.Is(err, tc.err) {
				t.Errorf("different error: got %v, want %v", err, tc.err)
			}
		})
	}
}
