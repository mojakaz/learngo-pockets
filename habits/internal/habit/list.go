package habit

import (
	"context"
	"fmt"
)

//go:generate minimock -i habitLister -s "_mock.go" -o "mocks"
type habitLister interface {
	FindAll(ctx context.Context) ([]Habit, error)
}

// List lists all habits in the DB.
func List(ctx context.Context, db habitLister) ([]Habit, error) {
	habits, err := db.FindAll(ctx)
	if err != nil {
		return []Habit{}, fmt.Errorf("failed to list habits: %w", err)
	}
	return habits, nil
}
