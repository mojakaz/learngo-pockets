package server

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"learngo-pockets/habits/api"
	"learngo-pockets/habits/internal/repository"
	"net"
	"sync"
	"testing"
)

func TestIntegration(t *testing.T) {
	// Skip this test when running lightweight suites
	if testing.Short() {
		t.Skip()
	}

	grpcServ := newServer(t)
	listener, err := net.Listen("tcp", "")
	require.NoError(t, err)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = grpcServ.Serve(listener)
		require.NoError(t, err)
	}()
	defer func() {
		// terminate the GRPC server
		grpcServ.Stop()
		// when that is done, and no error were caught, we can end this test
		wg.Wait()
	}()

	// create client
	habitsCli, err := newClient(t, listener.Addr().String())
	require.NoError(t, err)

	// add 2 habits
	addHabit(t, habitsCli, nil, "walk in the forest")
	addHabit(t, habitsCli, ptr(3), "read a few pages")
	addHabitWithError(t, habitsCli, ptr(5), "        ", codes.InvalidArgument)

	// check that the 2 habits are present
	listHabitsMatches(t, habitsCli, []*api.Habit{
		{
			Id:              "",
			Name:            "walk in the forest",
			WeeklyFrequency: 1,
		},
		{
			Id:              "",
			Name:            "read a few pages",
			WeeklyFrequency: 3,
		},
	})

	// ...
}

func newServer(t *testing.T) *grpc.Server {
	t.Helper()
	s := New(repository.New(t), t)
	grpcServer := grpc.NewServer()
	api.RegisterHabitsServer(grpcServer, s)
	return grpcServer
}

func newClient(t *testing.T, serverAddress string) (api.HabitsClient, error) {
	t.Helper()
	creds := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.Dial(serverAddress, creds)
	if err != nil {
		return nil, err
	}

	return api.NewHabitsClient(conn), nil
}

func listHabitsMatches(t *testing.T, habitsCli api.HabitsClient, expected []*api.Habit) {
	list, err := habitsCli.ListHabits(context.Background(), &api.ListHabitsRequest{})
	require.NoError(t, err)

	for i := range list.Habits {
		assert.NotEqual(t, "", list.Habits[i].Id)
		list.Habits[i].Id = ""
	}
	assert.Equal(t, expected, list.Habits)
}

func ptr(i int32) *int32 {
	return &i
}

func addHabit(t *testing.T, habitsCli api.HabitsClient, freq *int32, name string) {
	_, err := habitsCli.CreateHabit(context.Background(), &api.CreateHabitRequest{
		Name:            name,
		WeeklyFrequency: freq,
	})
	assert.NoError(t, err)
}

func addHabitWithError(t *testing.T, habitsCli api.HabitsClient, freq *int32, name string, expectedErrCode codes.Code) {
	_, err := habitsCli.CreateHabit(context.Background(), &api.CreateHabitRequest{
		Name:            name,
		WeeklyFrequency: freq,
	})
	assert.Equal(t, expectedErrCode, status.Code(err))
}
