package guess

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"learngo-pockets/httpgordle/internal/api"
	"learngo-pockets/httpgordle/internal/session"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler(t *testing.T) {
	handleFunc := Handler(gameFinderUpdaterStab{})
	body := strings.NewReader(`{"guess":"pocket"}`)
	req, err := http.NewRequest(http.MethodPut, "/games/", body)
	require.NoError(t, err)

	// add path parameters
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(api.GameID, "123456")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	recorder := httptest.NewRecorder()

	handleFunc(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

type gameFinderUpdaterStab struct {
	err error
}

func (g gameFinderUpdaterStab) Find(gameID session.GameID) (session.Game, error) {
	return session.Game{}, g.err
}

func (g gameFinderUpdaterStab) Update(game session.Game) error {
	return g.err
}
