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

func TestList(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	h1 := habit.Habit{
		ID:              "1",
		Name:            "Test1",
		WeeklyFrequency: 1,
		CreationTime:    time.Time{},
	}
	h2 := habit.Habit{
		ID:              "2",
		Name:            "Test2",
		WeeklyFrequency: 1,
		CreationTime:    time.Time{},
	}
	var dbError error
	tt := map[string]struct {
		db             func(ctl minimock.Tester) *mocks.HabitListerMock
		expectedHabits []habit.Habit
		expectedErr    error
	}{
		"nominal": {
			db: func(ctl minimock.Tester) *mocks.HabitListerMock {
				db := mocks.NewHabitListerMock(ctl)
				db.FindAllMock.Expect(ctx).Return([]habit.Habit{h1, h2}, nil)
				return db
			},
			expectedHabits: []habit.Habit{h1, h2},
			expectedErr:    nil,
		},
		"empty": {
			db: func(ctl minimock.Tester) *mocks.HabitListerMock {
				db := mocks.NewHabitListerMock(ctl)
				db.FindAllMock.Expect(ctx).Return([]habit.Habit{}, nil)
				return db
			},
			expectedHabits: []habit.Habit{},
			expectedErr:    nil,
		},
		"error case": {
			db: func(ctl minimock.Tester) *mocks.HabitListerMock {
				db := mocks.NewHabitListerMock(ctl)
				db.FindAllMock.Expect(ctx).Return([]habit.Habit{h1, h2}, dbError)
				return db
			},
			expectedHabits: []habit.Habit{h1, h2},
			expectedErr:    dbError,
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			ctl := minimock.NewController(t)
			db := tc.db(ctl)

			got, err := habit.List(ctx, db)
			assert.ErrorIs(t, err, tc.expectedErr)
			if err == nil {
				assert.Equal(t, tc.expectedHabits, got)
			}
		})
	}
}
