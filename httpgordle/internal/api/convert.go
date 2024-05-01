package api

import "learngo-pockets/httpgordle/internal/session"

// ToGameResponse converts a session.Game into an GameResponse.
func ToGameResponse(g session.Game) GameResponse {
	apiGame := GameResponse{
		ID:           string(g.ID),
		AttemptsLeft: g.AttemptsLeft,
		Guesses:      make([]session.Guess, len(g.Guesses)),
		Status:       string(g.Status),
		Solution:     g.Solution,
		WordLength:   byte(len(g.Solution)),
	}

	for index := 0; index < len(g.Guesses); index++ {
		apiGame.Guesses[index].Word = g.Guesses[index].Word
		apiGame.Guesses[index].Feedback = g.Guesses[index].Feedback
	}

	if g.AttemptsLeft != 0 {
		apiGame.Solution = ""
	}

	return apiGame
}
