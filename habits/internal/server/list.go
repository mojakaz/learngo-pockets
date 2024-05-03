package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"learngo-pockets/habits/api"
	"learngo-pockets/habits/internal/habit"
)

func (s *Server) ListHabits(ctx context.Context, request *api.ListHabitsRequest) (*api.ListHabitsResponse, error) {
	listedHabits, err := habit.List(ctx, s.db)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var responseHabits []*api.Habit
	for _, lh := range listedHabits {
		responseHabits = append(responseHabits, &api.Habit{
			Id:              string(lh.ID),
			Name:            string(lh.Name),
			WeeklyFrequency: int32(lh.WeeklyFrequency),
		})
	}
	return &api.ListHabitsResponse{
		Habits: responseHabits,
	}, nil
}
