package gordle

import (
	"bufio"
	"fmt"
	"golang.org/x/exp/slices"
	"io"
	"os"
	"strings"
)

// Game holds all the information we need to play a game of gordle.
type Game struct {
	reader      *bufio.Reader
	solution    []rune
	maxAttempts int
}

// New returns a Game, which can be used to Play!
func New(playerInput io.Reader, corpus []string, maxAttempts int) (*Game, error) {
	if len(corpus) == 0 {
		return nil, ErrCorpusIsEmpty
	}

	g := &Game{
		bufio.NewReader(playerInput),
		splitToUppercaseCharacters(pickWord(corpus)), // pick a random word from the corpus
		maxAttempts,
	}

	return g, nil
}

// Play runs the game.
func (g *Game) Play() {
	fmt.Println("Welcome to Gordle!")

	for currentAttempt := 1; currentAttempt <= g.maxAttempts; currentAttempt++ {
		guess := g.ask()

		fb := computeFeedback(guess, g.solution)
		fmt.Println(fb)

		if slices.Equal(guess, g.solution) {
			fmt.Printf(
				"🎉 You won! You found it in %d guess(es)! The word was: %s.\n",
				currentAttempt,
				string(g.solution),
			)
			return
		}
	}
	fmt.Printf("😞 You've lost! The solution was: %s. \n", string(g.solution))
}

// splitToUppercaseCharacters is a naive implementation to turn a string into a list of characters.
func splitToUppercaseCharacters(input string) []rune {
	return []rune(strings.ToUpper(input))
}

// ask reads input until a valid suggestion is made (and returned).
func (g *Game) ask() []rune {
	fmt.Printf("Enter a %d-character guess:\n", len(g.solution))

	for {
		playerInput, _, err := g.reader.ReadLine()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Gordle failed to read your guess: %s\n", err.Error())
			continue
		}

		guess := splitToUppercaseCharacters(string(playerInput))

		err = validateGuess(guess, g.solution)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Your attempt is invalid with Gordle's solution: %s\n", err.Error())
		}
		return guess
	}
}

// errInvalidWordLength is returned when the guess has the wrong number of characters.
var errInvalidWordLength = fmt.Errorf("invalid guess, word doesn't have the same number of characters as the solution")

// validateGuess ensures the guess is valid enough.
func validateGuess(guess, solution []rune) error {
	if len(guess) != len(solution) {
		return fmt.Errorf("expected %d, got %d, %w", len(solution), len(guess), errInvalidWordLength)
	}

	return nil
}

// computeFeedback verifies every character of the guess against the solution.
// It assumes that the length of the guess and the solution are the same.
func computeFeedback(guess, solution []rune) feedback {
	fb := make(feedback, len(guess))
	seen := make([]bool, len(solution))
	if len(guess) != len(solution) {
		//_, _ = fmt.Fprintf(os.Stderr, "Internal error! Guess and solution have different lengths: %d vs %d", len(guess), len(solution))
		return nil
	}
	for i, char := range guess {
		if char == solution[i] {
			fb[i] = correctPosition
			seen[i] = true
		}
	}
	for i, char := range guess {
		if fb[i] != absentCharacter {
			continue
		}
		for j, target := range solution {
			if seen[j] {
				continue
			}
			if char == target {
				fb[i] = wrongPosition
				seen[j] = true
				break
			}
		}
	}
	return fb
}
