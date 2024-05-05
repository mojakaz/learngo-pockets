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

func (s *Server) TickHabit(ctx context.Context, request *api.TickHabitRequest) (*api.TickHabitResponse, error) {
	err := habit.Tick(ctx, s.db, s.db, habit.ID(request.HabitId), time.Now())
	if err != nil {
		switch {
		case errors.Is(err, repository.DBError{}):
			return nil, status.Errorf(codes.NotFound, "couldn't find habit %q in repository", request.HabitId)
		default:
			return nil, status.Errorf(codes.Internal, "cannot tick habit %q: %s", request.HabitId, err.Error())
		}
	}
	return &api.TickHabitResponse{}, nil
}
