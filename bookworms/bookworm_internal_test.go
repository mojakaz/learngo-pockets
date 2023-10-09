package bookworms

import (
	"testing"
)

var (
	handmaidsTale = Book{Author: "Margaret Atwood", Title: "The Handmaid's Tale"}
	oryxAndCrake  = Book{Author: "Margaret Atwood", Title: "Oryx and Crake"}
	theBellJar    = Book{Author: "Sylvia Plath", Title: "The Bell Jar"}
	janeEyre      = Book{Author: "Charlotte Brontë", Title: "Jane Eyre"}
)

func TestLoadBookworms(t *testing.T) {
	tests := map[string]testCase{
		"file exists": {
			bookwormsFile: "testdata/bookworms.json",
			want: []Bookworm{
				{Name: "Fadi", Books: []Book{handmaidsTale, theBellJar}},
				{Name: "Peggy", Books: []Book{oryxAndCrake, handmaidsTale, janeEyre}},
			},
			wantErr: false,
		},
		"file doesn't exist": {
			bookwormsFile: "testdata/no_such_file.json",
			want:          nil,
			wantErr:       true,
		},
		"invalid JSON": {
			bookwormsFile: "testdata/invalid.json",
			want:          nil,
			wantErr:       true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := loadBookworms(tc.bookwormsFile)
			if err != nil && !tc.wantErr {
				t.Fatalf("expected no error, got one %v", err.Error())
			}
			if err == nil && tc.wantErr {
				t.Fatal("expected an error, got none")
			}
			if !equalBookworms(got, tc.want, t) {
				t.Fatalf("different result: got %v, expected %v", got, tc.want)
			}
		})
	}
}

type testCase struct {
	bookwormsFile string
	want          []Bookworm
	wantErr       bool
}

// equalBookworms is a helper to test the equality of two lists of Bookworms.
func equalBookworms(bookworms, target []Bookworm, t *testing.T) bool {
	t.Helper()
	if len(bookworms) != len(target) {
		return false
	}

	for i := range bookworms {
		if bookworms[i].Name != target[i].Name {
			return false
		}

		if !equalBooks(bookworms[i].Books, target[i].Books, t) {
			return false
		}
	}

	return true
}

// equalBooks is a helper to test the equality of two lists of Books.
func equalBooks(books, target []Book, t *testing.T) bool {
	t.Helper()
	if len(books) != len(target) {
		return false
	}

	for i := range books {
		if books[i] != target[i] {
			return false
		}
	}

	return true
}

func TestBooksCount(t *testing.T) {
	tests := map[string]struct {
		bookworms []Bookworm
		want      map[Book]uint
	}{
		"nominal use case": {
			bookworms: []Bookworm{
				{Name: "Fadi", Books: []Book{handmaidsTale, theBellJar}},
				{Name: "Peggy", Books: []Book{oryxAndCrake, handmaidsTale, janeEyre}},
			},
			want: map[Book]uint{handmaidsTale: 2, theBellJar: 1, oryxAndCrake: 1, janeEyre: 1},
		},
		"no bookworms": {
			bookworms: []Bookworm{},
			want:      map[Book]uint{},
		},
		"bookworm without books": {
			bookworms: []Bookworm{
				{Name: "Fadi", Books: []Book{handmaidsTale, theBellJar}},
				{Name: "Peggy", Books: []Book{}},
			},
			want: map[Book]uint{handmaidsTale: 1, theBellJar: 1},
		},
		"bookworm with twice the same book": {
			bookworms: []Bookworm{
				{Name: "Fadi", Books: []Book{handmaidsTale, handmaidsTale, theBellJar}},
				{Name: "Peggy", Books: []Book{oryxAndCrake, janeEyre}},
			},
			want: map[Book]uint{handmaidsTale: 2, theBellJar: 1, oryxAndCrake: 1, janeEyre: 1},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			count := booksCount(tc.bookworms)
			if !equalBooksCount(count, tc.want, t) {
				t.Fatalf("different counts: got %q, expected %q", count, tc.want)
			}
		})
	}
}

// equalBooksCount is a helper to test the equality of two maps of books count.
func equalBooksCount(got, want map[Book]uint, t *testing.T) bool {
	t.Helper()
	if len(got) != len(want) {
		return false
	}

	for book, targetCount := range want {
		count, ok := got[book]
		if !ok || targetCount != count {
			return false
		}
	}

	return true
}

func TestFindCommonBooks(t *testing.T) {
	tests := map[string]struct {
		bookworms []Bookworm
		want      []Book
	}{
		"three bookworms have the same books on their shelves": {
			bookworms: []Bookworm{
				{Name: "Fadi", Books: []Book{handmaidsTale, theBellJar}},
				{Name: "Peggy", Books: []Book{handmaidsTale, theBellJar}},
				{Name: "Tom", Books: []Book{handmaidsTale, theBellJar}},
			},
			want: []Book{
				handmaidsTale, theBellJar,
			},
		},
		"no common book": {
			bookworms: []Bookworm{
				{Name: "Fadi", Books: []Book{handmaidsTale, theBellJar}},
				{Name: "Peggy", Books: []Book{oryxAndCrake, janeEyre}},
			},
			want: []Book{},
		},
		"one common book": {
			bookworms: []Bookworm{
				{Name: "Fadi", Books: []Book{handmaidsTale, theBellJar}},
				{Name: "Peggy", Books: []Book{oryxAndCrake, handmaidsTale, janeEyre}},
			},
			want: []Book{
				handmaidsTale,
			},
		},
		"one bookworm has no book": {
			bookworms: []Bookworm{
				{Name: "Fadi", Books: []Book{handmaidsTale, theBellJar}},
				{Name: "Peggy", Books: []Book{}},
			},
			want: []Book{},
		},
		"nobody has any book": {
			bookworms: []Bookworm{
				{Name: "Fadi", Books: []Book{}},
				{Name: "Peggy", Books: []Book{}},
			},
			want: []Book{},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			commonBooks := findCommonBooks(tc.bookworms)
			if !equalCommonBooks(commonBooks, tc.want) {
				t.Fatalf("got different common books: got %q, expected %q", commonBooks, tc.want)
			}
		})
	}
}

func equalCommonBooks(got, want []Book) bool {
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}

func ExampleDisplayBooks() {
	displayBooks([]Book{handmaidsTale, oryxAndCrake})
	// Output
	// Here are the common books:
	// - The Handmaid's Tale by Margaret Atwood
	// - Oryx and Crake by Margaret Atwood
}
