// Package repository accesses the habits data.
package repository

import (
	"context"
	"fmt"
	"learngo-pockets/habits/internal/habit"
	"learngo-pockets/habits/isoweek"
	"sort"
	"sync"
	"time"
)

// ticksPerWeek holds all the timestamps for a given week number.
type ticksPerWeek map[isoweek.ISO8601][]time.Time

// HabitRepository holds all the current habits.
type HabitRepository struct {
	mutex  sync.RWMutex
	habits map[habit.ID]habit.Habit
	ticks  map[habit.ID]ticksPerWeek
	lgr    Logger
}

type Logger interface {
	Logf(format string, args ...any)
}

func New(lgr Logger) *HabitRepository {
	return &HabitRepository{
		habits: make(map[habit.ID]habit.Habit),
		ticks:  make(map[habit.ID]ticksPerWeek),
		lgr:    lgr,
	}
}

// Add stores new habit to the storage.
func (hr *HabitRepository) Add(_ context.Context, habit habit.Habit) error {
	hr.lgr.Logf("Creating habit: %v", habit)
	hr.mutex.Lock()
	defer hr.mutex.Unlock()
	if _, ok := hr.habits[habit.ID]; ok {
		return DBError{reason: fmt.Sprintf("habit already exists: %v", habit)}
	}
	hr.habits[habit.ID] = habit
	return nil
}

// FindAll returns all habits sorted by creation time.
func (hr *HabitRepository) FindAll(_ context.Context) ([]habit.Habit, error) {
	hr.lgr.Logf("Listing habits, sorted by creation time...")

	// Lock the reading and the writing of the habits.
	hr.mutex.Lock()
	defer hr.mutex.Unlock()

	habits := make([]habit.Habit, 0)
	for _, h := range hr.habits {
		habits = append(habits, h)
	}

	// Ensure the output is deterministic by sorting the habits.
	sort.Slice(habits, func(i, j int) bool {
		return habits[i].CreationTime.Before(habits[j].CreationTime)
	})

	return habits, nil
}

// Find checks if the habit exists in the storage.
func (hr *HabitRepository) Find(_ context.Context, id habit.ID) (habit.Habit, error) {
	hr.lgr.Logf("Finding habit %q", id)

	// Lock the reading of the habits.
	hr.mutex.RLock()
	defer hr.mutex.RUnlock()

	if h, ok := hr.habits[id]; ok {
		return h, nil
	}
	return habit.Habit{}, DBError{reason: fmt.Sprintf("habit %q not found", id)}
}

// AddTick adds a new tick for the habit.
func (hr *HabitRepository) AddTick(_ context.Context, id habit.ID, t time.Time) error {
	hr.lgr.Logf("Adding tick for habit %q", id)

	// Lock the reading and writing of ticks.
	hr.mutex.Lock()
	defer hr.mutex.Unlock()

	tpw := make(map[isoweek.ISO8601][]time.Time)
	year, week := t.ISOWeek()

	// Key id not in ticks storage.
	if _, ok := hr.ticks[id]; !ok {
		tpw[isoweek.ISO8601{
			Year: year,
			Week: week,
		}] = []time.Time{t}
		hr.ticks[id] = tpw
		return nil
	}

	ticks := hr.ticks[id][isoweek.ISO8601{
		Year: year,
		Week: week,
	}]
	ticks = append(ticks, t)
	tpw[isoweek.ISO8601{
		Year: year,
		Week: week,
	}] = ticks
	hr.ticks[id] = tpw

	return nil
}

// FindTick returns all ticks associated with given id and timestamp.
func (hr *HabitRepository) FindTick(_ context.Context, id habit.ID, t time.Time) ([]time.Time, error) {
	hr.lgr.Logf("Finding tick for habit %q", id)

	// Lock the reading and writing of ticks.
	hr.mutex.RLock()
	defer hr.mutex.RUnlock()

	year, week := t.ISOWeek()
	ticks := hr.ticks[id][isoweek.ISO8601{
		Year: year,
		Week: week,
	}]
	return ticks, nil
}
