package repository

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"learngo-pockets/habits/internal/habit"
	"learngo-pockets/habits/isoweek"
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
			tc.hr.habits[tc.habit.ID] = tc.habit
			err := tc.hr.Add(context.Background(), tc.habit)
			assert.Error(t, err, "assume error but got no error")
		})
	}
}

func TestHabitRepository_FindAll(t *testing.T) {
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
			tc.hr.habits[tc.habit.ID] = tc.habit
			got, err := tc.hr.FindAll(context.Background())
			assert.NoError(t, err, "assume no error but got an error %w", err)
			assert.Equal(t, []habit.Habit{tc.habit}, got, "assume %v matches %v", tc.habit, got)
		})
	}
}

func TestHabitRepository_AddTick(t *testing.T) {
	t.Parallel()
	tt := map[string]struct {
		hr *HabitRepository
		id habit.ID
		t  time.Time
	}{
		"nominal": {
			hr: New(log.New(&bytes.Buffer{})),
			id: "test",
			t:  time.Now(),
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			err := tc.hr.AddTick(context.Background(), tc.id, tc.t)
			assert.NoError(t, err, "assume no error but got an error %w", err)
		})
	}
}

func TestHabitRepository_FindTick(t *testing.T) {
	t.Parallel()
	tt := map[string]struct {
		hr *HabitRepository
		id habit.ID
		t  time.Time
	}{
		"nominal": {
			hr: New(log.New(&bytes.Buffer{})),
			id: "test",
			t:  time.Now(),
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			year, week := time.Now().ISOWeek()
			ticks := []time.Time{tc.t}
			tpw := make(map[isoweek.ISO8601][]time.Time)
			tpw[isoweek.ISO8601{
				Year: year,
				Week: week,
			}] = ticks
			tc.hr.ticks[tc.id] = tpw

			got, err := tc.hr.FindTick(context.Background(), tc.id, tc.t)
			assert.NoError(t, err, "assume no error but got an error %w", err)
			assert.Equal(t, ticks, got, "assume %v matches %v", ticks, got)
		})
	}
}
