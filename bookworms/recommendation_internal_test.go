package bookworms

import (
	"math"
	"testing"
)

var (
	emmy   = Reader{Name: "Emmy", Books: []Book{handmaidsTale, janeEyre}}
	tom    = Reader{Name: "Tom", Books: []Book{oryxAndCrake, janeEyre}}
	horny  = Reader{Name: "horny", Books: []Book{theBellJar, oryxAndCrake, handmaidsTale}}
	michel = Reader{Name: "Michel", Books: []Book{janeEyre}}
	goro   = Reader{Name: "Goro", Books: []Book{theBellJar, oryxAndCrake}}
)

func TestRecommend(t *testing.T) {
	tt := map[string]struct {
		allReaders []Reader
		target     Reader
		n          int
		expected   []Recommendation
	}{
		"two readers and a target have one book in common " +
			"and the target does not have a book that the others have": {
			allReaders: []Reader{emmy, tom, horny},
			target:     horny,
			n:          3,
			expected:   []Recommendation{{Book: janeEyre, Score: 2}},
		},
		"a reader and a target hove a book in common, another reader and the target have no books in common": {
			allReaders: []Reader{emmy, horny, michel},
			target:     horny,
			n:          3,
			expected:   []Recommendation{{Book: janeEyre, Score: 1}},
		},
		"a reader and a target have two books in common, and the reader has a book that the target doesn't have": {
			allReaders: []Reader{horny, goro},
			target:     goro,
			n:          3,
			expected:   []Recommendation{{Book: handmaidsTale, Score: math.Log(2) + 1}},
		},
		"no book in common": {
			allReaders: []Reader{horny, michel, goro},
			target:     michel,
			n:          3,
			expected:   []Recommendation{},
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got := recommend(tc.allReaders, tc.target, tc.n)
			if !equalRecommendation(got, tc.expected) {
				t.Errorf("got %v, expected %v", got, tc.expected)
			}
		})
	}
}

func equalRecommendation(got []Recommendation, expected []Recommendation) bool {
	if len(got) != len(expected) {
		return false
	}
	for i := range got {
		if got[i].Score != expected[i].Score {
			return false
		}
	}
	return true
}
