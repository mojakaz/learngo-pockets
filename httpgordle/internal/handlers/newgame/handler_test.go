package newgame

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"learngo-pockets/httpgordle/internal/session"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestHandler(t *testing.T) {
	handleFunc := Handler(gameAdderStub{})
	req, err := http.NewRequest(http.MethodPost, "/games", nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()

	handleFunc(recorder, req)

	assert.Equal(t, http.StatusCreated, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
	// idFinderRegexp is a regular expression that will ensure the body contains an id field with a value that contains
	// only letters (uppercase and/or lowercase) and/or digits.
	idFinderRegexp := regexp.MustCompile(`.*"id":"(\w+)".+`)
	body := recorder.Body.String()
	id := idFinderRegexp.FindStringSubmatch(body)
	if len(id) != 2 {
		t.Fatalf("cannot find one id in the json output: %s", body)
	}
	body = strings.Replace(body, id[1], "123456", 1)

	assert.JSONEq(t, `{"id":"123456","attempts_left":5,"guesses":null,"status":"Playing","word_length":5}`, body)
}

type gameAdderStub struct {
	err error
}

func (g gameAdderStub) Add(_ session.Game) error {
	return g.err
}
