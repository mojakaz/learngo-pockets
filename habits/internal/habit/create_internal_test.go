package habit

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_validateAndCompleteHabit(t *testing.T) {
	t.Parallel()

	t.Run("Full", testValidateAndCompleteHabitFull)
	t.Run("Partial", testValidateAndCompleteHabitPartial)
	t.Run("SpaceName", testValidateAndCompleteHabitSpaceName)
}

func testValidateAndCompleteHabitFull(t *testing.T) {
	t.Parallel()

	h := Habit{
		ID:              "1",
		Name:            "Test",
		WeeklyFrequency: 1,
		CreationTime:    time.Now(),
	}

	got, err := validateAndCompleteHabit(h)
	require.NoError(t, err)
	assert.Equal(t, h, got)
}

func testValidateAndCompleteHabitPartial(t *testing.T) {
	t.Parallel()

	h := Habit{
		ID:              "",
		Name:            "Test",
		WeeklyFrequency: 0,
		CreationTime:    time.Time{},
	}

	got, err := validateAndCompleteHabit(h)
	require.NoError(t, err)
	assert.Equal(t, h.Name, got.Name)
	assert.Equal(t, WeeklyFrequency(1), got.WeeklyFrequency)
	assert.NotEmpty(t, got.ID)
	assert.NotEmpty(t, got.CreationTime)
}

func testValidateAndCompleteHabitSpaceName(t *testing.T) {
	t.Parallel()

	h := Habit{
		ID:              "",
		Name:            "    ",
		WeeklyFrequency: 0,
		CreationTime:    time.Time{},
	}

	_, err := validateAndCompleteHabit(h)
	assert.ErrorAs(t, err, &InvalidInputError{})
}
