package gordle

import (
	"fmt"
	"golang.org/x/exp/slices"
	"os"
	"strings"
)

// Game holds all the information we need to play a game of gordle.
type Game struct {
	solution []rune
}

// New returns a Game, which can be used to Play!
func New(corpus []string) (*Game, error) {
	if len(corpus) == 0 {
		return nil, ErrCorpusIsEmpty
	}

	g := &Game{
		splitToUppercaseCharacters(pickWord(corpus)), // pick a random word from the corpus
	}

	return g, nil
}

// Play runs the game.
func (g *Game) Play(guess string) (Feedback, error) {
	sr := splitToUppercaseCharacters(guess)

	fb, err := computeFeedback(sr, g.solution)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		return nil, err
	}

	if slices.Equal(sr, g.solution) {
		fmt.Printf("🎉 You won! You found it! The word was: %s.\n", string(g.solution))
		return fb, err
	}
	fmt.Printf("😞 Wrong guess! Try again!\n")
	return fb, err
}

func (g *Game) ShowAnswer() string {
	return string(g.solution)
}

func (g *Game) ShowLength() int {
	return len(g.solution)
}

// splitToUppercaseCharacters is a naive implementation to turn a string into a list of characters.
func splitToUppercaseCharacters(input string) []rune {
	return []rune(strings.ToUpper(input))
}

// ErrInvalidGuess is returned when the guess has the wrong number of characters.
var ErrInvalidGuess = fmt.Errorf("invalid guess, word doesn't have the same number of characters as the solution")

// validateGuess ensures the guess is valid enough.
func validateGuess(guess, solution []rune) error {
	if len(guess) != len(solution) {
		return fmt.Errorf("expected %d, got %d, %w", len(solution), len(guess), ErrInvalidGuess)
	}

	return nil
}

// computeFeedback verifies every character of the guess against the solution.
// It assumes that the length of the guess and the solution are the same.
func computeFeedback(guess, solution []rune) (Feedback, error) {
	err := validateGuess(guess, solution)
	if err != nil {
		return nil, err
	}
	fb := make(Feedback, len(guess))
	seen := make([]bool, len(solution))

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
	return fb, err
}
