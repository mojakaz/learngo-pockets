// Package repository accesses the habits data.
package repository

import (
	"context"
	"fmt"
	"learngo-pockets/habits/internal/habit"
	"sort"
	"sync"
)

// HabitRepository holds all the current habits.
type HabitRepository struct {
	mutex   sync.RWMutex
	storage map[habit.ID]habit.Habit
	lgr     Logger
}

type Logger interface {
	Logf(format string, args ...any)
}

func New(lgr Logger) *HabitRepository {
	return &HabitRepository{
		storage: make(map[habit.ID]habit.Habit),
		lgr:     lgr,
	}
}

func (hr *HabitRepository) Add(_ context.Context, habit habit.Habit) error {
	hr.lgr.Logf("Creating habit: %v", habit)
	hr.mutex.Lock()
	defer hr.mutex.Unlock()
	if _, ok := hr.storage[habit.ID]; ok {
		return fmt.Errorf("habit already exists: %v", habit)
	}
	hr.storage[habit.ID] = habit
	return nil
}

// FindAll returns all habits sorted by creation time.
func (hr *HabitRepository) FindAll(_ context.Context) ([]habit.Habit, error) {
	hr.lgr.Logf("Listing habits, sorted by creation time...")

	// Lock the reading and the writing of the habits.
	hr.mutex.Lock()
	defer hr.mutex.Unlock()

	habits := make([]habit.Habit, 0)
	for _, h := range hr.storage {
		habits = append(habits, h)
	}

	// Ensure the output is deterministic by sorting the habits.
	sort.Slice(habits, func(i, j int) bool {
		return habits[i].CreationTime.Before(habits[j].CreationTime)
	})

	return habits, nil
}
