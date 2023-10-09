package bookworms

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// loadBookworms reads the file and returns the list of bookworms, and their beloved books, found therein.
func loadBookworms(filePath string) ([]Bookworm, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)

	// Initialise the type in which the file will be decoded.
	var bookworms []Bookworm

	// Decode the file and store the content in the value pointed to by bookworms
	buffedReader := bufio.NewReaderSize(f, 1024*1024)
	decoder := json.NewDecoder(buffedReader)
	err = decoder.Decode(&bookworms)
	if err != nil {
		return nil, err
	}

	return bookworms, nil
}

// A Bookworm contains the list of books on a bookworm's shelf.
type Bookworm struct {
	Name  string `json:"name"`
	Books []Book `json:"books"`
}

// Book describes a book on a bookworm's shelf.
type Book struct {
	Author string `json:"author"`
	Title  string `json:"title"`
}

// booksCount registers all the books and their occurrences from the bookworms' shelves.
func booksCount(bookworms []Bookworm) map[Book]uint {
	count := make(map[Book]uint)

	for _, bookworm := range bookworms {
		for _, book := range bookworm.Books {
			count[book]++
		}
	}

	return count
}

// findCommonBooks returns books that are on more than one bookworm's shelf.
func findCommonBooks(bookworms []Bookworm) []Book {
	booksOnShelves := booksCount(bookworms)
	var commonBooks []Book
	for book, count := range booksOnShelves {
		if count > 1 {
			commonBooks = append(commonBooks, book)
		}
	}
	return sortBooks2(commonBooks)
}

// sortBooks sorts the books by Author and then Title.
func sortBooks(books []Book) []Book {
	sort.Slice(books, func(i, j int) bool {
		if books[i].Author != books[j].Author {
			return books[i].Author < books[j].Author
		}
		return books[i].Title < books[j].Title
	})
	return books
}

// displayBooks prints out the titles and authors of a list of books.
func displayBooks(books []Book) {
	fmt.Println("Here are the common books:")
	for _, book := range books {
		fmt.Println("-", book.Title, "by", book.Author)
	}
}

// byAuthor is a list of Books. Defining a custom type to implement sort.Interface
type byAuthor []Book

// Len implements sort.Interface by returning the length of the collection.
func (b byAuthor) Len() int { return len(b) }

// Swap implements sort.Interface and swaps two books.
func (b byAuthor) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

// Less implements sort.Interface and returns books sorted by Author and then Title.
func (b byAuthor) Less(i, j int) bool {
	if b[i].Author != b[j].Author {
		return b[i].Author < b[j].Author
	}
	return b[i].Title < b[j].Title
}

// sortBooks2 sorts the books by Author and then Title in alphabetical order.
func sortBooks2(books []Book) []Book {
	sort.Sort(byAuthor(books))
	return books
}
