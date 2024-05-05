package habit

import (
	"context"
	"fmt"
	"time"
)

//go:generate minimock -i tickFinder -s "_mock.go" -o "mocks"
type tickFinder interface {
	FindTick(ctx context.Context, id ID, t time.Time) ([]time.Time, error)
}

type TicksCount int32

type StatusResponse struct {
	ID              ID
	Name            Name
	WeeklyFrequency WeeklyFrequency
	TicksCount      TicksCount
}

func Status(ctx context.Context, habitDB habitFinder, tickDB tickFinder, id ID, t time.Time) (StatusResponse, error) {
	// Check if the habit exists.
	h, err := habitDB.Find(ctx, id)
	if err != nil {
		return StatusResponse{}, fmt.Errorf("cannot find habit %q: %w", id, err)
	}

	ticks, err := tickDB.FindTick(ctx, id, t)
	if err != nil {
		return StatusResponse{}, fmt.Errorf("unexpected error while finding ticks for habit %q: %w", id, err)
	}

	return StatusResponse{
		ID:              h.ID,
		Name:            h.Name,
		WeeklyFrequency: h.WeeklyFrequency,
		TicksCount:      TicksCount(len(ticks)),
	}, nil
}
