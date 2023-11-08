package getstatus

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"learngo-pockets/httpgordle/internal/api"
	"learngo-pockets/httpgordle/internal/session"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	handleFunc := Handler(gameFinderStab{})
	req, err := http.NewRequest(http.MethodPost, "/games/", nil)
	require.NoError(t, err)

	// add path parameters
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(api.GameID, "123456")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	recorder := httptest.NewRecorder()

	handleFunc(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.JSONEq(t, `{"id":"123456","attempts_left":0,"guesses":null,"word_length":0,"status":""}`, recorder.Body.String())
}

type gameFinderStab struct {
	err error
}

func (g gameFinderStab) Find(gameID session.GameID) (session.Game, error) {
	return session.Game{}, g.err
}
