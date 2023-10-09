package bookworms

import (
	"math"
	"sort"
)

type Recommendation struct {
	Book  Book
	Score float64
}

type Reader struct {
	Name  string
	Books []Book
}

type set map[Book]struct{}

// newSet returns a new Set with a list of unique books.
func newSet(books ...Book) set {
	uniqueBooks := make(map[Book]struct{})
	for _, book := range books {
		_, ok := uniqueBooks[book]
		if !ok {
			uniqueBooks[book] = struct{}{}
		}
	}
	return uniqueBooks
}

// Contains checks if the set contains the given book.
func (s set) Contains(b Book) bool {
	_, ok := s[b]
	return ok
}

func recommend(allReaders []Reader, target Reader, n int) []Recommendation {
	read := newSet(target.Books...)

	recommendations := map[Book]float64{}
	for _, reader := range allReaders {
		if reader.Name == target.Name {
			continue
		}

		var similarity float64
		for _, book := range reader.Books {
			if read.Contains(book) {
				similarity++
			}
			// you could also later extend to liked and dislike score
		}
		if similarity == 0 {
			continue
		}

		score := math.Log(similarity) + 1
		for _, book := range reader.Books {
			if !read.Contains(book) {
				recommendations[book] += score
			}
		}
	}
	// TODO: sort by score
	// TODO: only output a certain amount of recommendations (n)
	var result []Recommendation
	for book, score := range recommendations {
		result = append(result, Recommendation{Book: book, Score: score})
		if len(result) == n {
			break
		}
	}
	return sortRecommendation(result)
}

// sortRecommendation sorts the list of Recommendation by Score.
func sortRecommendation(sr []Recommendation) []Recommendation {
	sort.Slice(sr, func(i, j int) bool {
		return sr[i].Score < sr[j].Score
	})
	return sr
}
