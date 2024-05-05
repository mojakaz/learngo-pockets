package server

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"learngo-pockets/habits/api"
	"learngo-pockets/habits/internal/habit"
	"learngo-pockets/habits/internal/repository"
	"time"
)

func (s *Server) GetHabitStatus(ctx context.Context, request *api.GetHabitStatusRequest) (*api.GetHabitStatusResponse, error) {
	timestamp := request.Timestamp.AsTime()
	if year, _ := timestamp.ISOWeek(); year == 1970 {
		timestamp = time.Now()
	}
	res, err := habit.Status(ctx, s.db, s.db, habit.ID(request.HabitId), timestamp)
	if err != nil {
		switch {
		case errors.Is(err, repository.DBError{}):
			return nil, status.Errorf(codes.NotFound, "couldn't find habit %q: %s", request.HabitId, err.Error())
		default:
			return nil, status.Errorf(codes.Internal, "error while getting status for habit %q, timestamp %q: %s", request.HabitId, timestamp, err.Error())
		}
	}

	return &api.GetHabitStatusResponse{
		Habit: &api.Habit{
			Id:              string(res.ID),
			Name:            string(res.Name),
			WeeklyFrequency: int32(res.WeeklyFrequency),
		},
		TicksCount: int32(res.TicksCount),
	}, nil
}
