package guess

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"learngo-pockets/httpgordle/internal/api"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandle(t *testing.T) {
	body := strings.NewReader(`{"guess":"pocket"}`)
	req, err := http.NewRequest(http.MethodPut, "/games/", body)
	require.NoError(t, err)

	// add path parameters
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(api.GameID, "123456")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	recorder := httptest.NewRecorder()

	Handle(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
}
