package habit_test

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"learngo-pockets/habits/internal/habit"
	"learngo-pockets/habits/internal/habit/mocks"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	h := habit.Habit{
		ID:              "",
		Name:            "Test",
		WeeklyFrequency: 0,
		CreationTime:    time.Time{},
	}
	var dbErr error
	tt := map[string]struct {
		db          func(ctl minimock.Tester) *mocks.HabitCreatorMock
		expectedErr error
	}{
		"nominal": {
			db: func(ctl minimock.Tester) *mocks.HabitCreatorMock {
				db := mocks.NewHabitCreatorMock(ctl)
				db.AddMock.Return(nil)
				return db
			},
			expectedErr: nil,
		},
		"error case": {
			db: func(ctl minimock.Tester) *mocks.HabitCreatorMock {
				db := mocks.NewHabitCreatorMock(ctl)
				db.AddMock.Return(dbErr)
				return db
			},
			expectedErr: dbErr,
		},
		"db timeout": {
			db: func(ctl minimock.Tester) *mocks.HabitCreatorMock {
				db := mocks.NewHabitCreatorMock(ctl)
				db.AddMock.Set(
					func(ctx context.Context, habit habit.Habit) error {
						select {
						// This tick is longer than a database call
						case <-time.Tick(2 * time.Second):
							return nil
						case <-ctx.Done():
							return ctx.Err()
						}
					})
				return db
			},
			expectedErr: context.DeadlineExceeded,
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			ctl := minimock.NewController(t)
			db := tc.db(ctl)

			got, err := habit.Create(ctx, db, h)
			assert.ErrorIs(t, err, tc.expectedErr)
			if tc.expectedErr == nil {
				assert.Equal(t, h.Name, got.Name)
			}
		})
	}
}
