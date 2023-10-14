package gordle

import (
	"errors"
	"testing"
)

func TestReadCorpus(t *testing.T) {
	tt := map[string]struct {
		file   string
		length int
		err    error
	}{
		"English corpus": {
			file:   "../corpus/english.txt",
			length: 156,
			err:    nil,
		},
		"empty corpus": {
			file:   "../corpus/empty.txt",
			length: 0,
			err:    ErrCorpusIsEmpty,
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			words, err := ReadCorpus(tc.file)
			if len(words) != tc.length {
				t.Errorf("different words length: got %d, want %d", len(words), tc.length)
			}
			if !errors.Is(err, tc.err) {
				t.Errorf("different error: got %v, want %v", err, tc.err)
			}
		})
	}
}
