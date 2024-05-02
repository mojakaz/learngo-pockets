// Package repository accesses the habits data.
package repository

import (
	"context"
	"fmt"
	"learngo-pockets/habits/internal/habit"
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

func (hr *HabitRepository) ListAll(_ context.Context) ([]habit.Habit, error) {
	hr.mutex.RLock()
	defer hr.mutex.RUnlock()
	var habits []habit.Habit
	for _, h := range hr.storage {
		habits = append(habits, h)
	}
	return habits, nil
}
