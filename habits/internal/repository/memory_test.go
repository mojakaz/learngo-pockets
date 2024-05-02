package repository

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"learngo-pockets/habits/internal/habit"
	"learngo-pockets/habits/log"
	"testing"
	"time"
)

func TestHabitRepository_Add_Success(t *testing.T) {
	t.Parallel()
	tt := map[string]struct {
		hr    *HabitRepository
		habit habit.Habit
	}{
		"nominal": {
			hr: New(log.New(&bytes.Buffer{})),
			habit: habit.Habit{
				ID:              "1",
				Name:            "Test",
				WeeklyFrequency: 1,
				CreationTime:    time.Now(),
			},
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			err := tc.hr.Add(context.Background(), tc.habit)
			assert.NoError(t, err, "assume no error but got an error %w", err)
		})
	}
}

func TestHabitRepository_Add_Fail(t *testing.T) {
	t.Parallel()
	tt := map[string]struct {
		hr    *HabitRepository
		habit habit.Habit
	}{
		"key exists": {
			hr: New(log.New(&bytes.Buffer{})),
			habit: habit.Habit{
				ID:              "a",
				Name:            "Test",
				WeeklyFrequency: 0,
				CreationTime:    time.Time{},
			},
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			tc.hr.storage[tc.habit.ID] = tc.habit
			err := tc.hr.Add(context.Background(), tc.habit)
			assert.Error(t, err, "assume error but got no error")
		})
	}
}

func TestHabitRepository_ListAll(t *testing.T) {
	t.Parallel()
	tt := map[string]struct {
		hr    *HabitRepository
		habit habit.Habit
	}{
		"nominal": {
			hr: New(log.New(&bytes.Buffer{})),
			habit: habit.Habit{
				ID:              "b",
				Name:            "Test",
				WeeklyFrequency: 0,
				CreationTime:    time.Time{},
			},
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			_, err := tc.hr.ListAll(context.Background())
			assert.NoError(t, err, "assume no error but got an error %w", err)
		})
	}
}
