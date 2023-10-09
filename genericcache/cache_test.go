package genericcache

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestCache_Upsert(t *testing.T) {
	tt := map[string]struct {
		key   string
		value int
	}{
		"nominal": {
			key:   "count",
			value: 1,
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			c := New[string, int](3, time.Millisecond*100)
			err := c.Upsert(tc.key, tc.value)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestCache_Upsert2(t *testing.T) {
	tt := map[string]struct {
		key    string
		value  int
		value2 int
	}{
		"nominal": {
			key:    "count",
			value:  1,
			value2: 2,
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			c := New[string, int](3, time.Millisecond*100)
			err := c.Upsert(tc.key, tc.value)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			err = c.Upsert(tc.key, tc.value2)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestCache_Read(t *testing.T) {
	tt := map[string]struct {
		key      string
		value    int
		expected int
	}{
		"nominal": {
			key:      "count",
			value:    1,
			expected: 1,
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			c := New[string, int](3, time.Millisecond*100)
			_ = c.Upsert(tc.key, tc.value)
			got, found := c.Read(tc.key)
			assert.Truef(t, found, "expected the key %s is found, but it is not found.", tc.key)
			assert.Equalf(t, tc.expected, got, "expected %d, got %d", tc.expected, got)
		})
	}
}

func TestCache_Read_KeyNotFound(t *testing.T) {
	tt := map[string]struct {
		key   string
		value int
	}{
		"no such key": {
			key:   "non",
			value: 1,
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			c := New[string, int](3, time.Millisecond*100)
			_, found := c.Read(tc.key)
			if found {
				t.Errorf("expected the key %s is not found, but it is found.", tc.key)
			}
		})
	}
}

func TestCache_Delete(t *testing.T) {
	tt := map[string]struct {
		key   string
		value int
		found bool
	}{
		"nominal": {
			key:   "count",
			value: 1,
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			c := New[string, int](3, time.Millisecond*100)
			_ = c.Upsert(tc.key, tc.value)
			c.Delete(tc.key)
			_, found := c.Read(tc.key)
			if found {
				t.Error("expected the key is not found, but it is found.")
			}
		})
	}
}

func TestCache_Delete_KeyNonExist(t *testing.T) {
	tt := map[string]struct {
		key   string
		value int
		found bool
	}{
		"nominal": {
			key: "count",
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			c := New[string, int](3, time.Millisecond*100)
			c.Delete(tc.key)
			_, found := c.Read(tc.key)
			if found {
				t.Error("expected the key is not found, but it is found.")
			}
		})
	}
}

func TestCache_Parallel_goroutines(t *testing.T) {
	c := New[int, string](3, time.Millisecond*100)

	const parallelTasks = 10
	wg := sync.WaitGroup{}
	wg.Add(parallelTasks)

	for i := 0; i < parallelTasks; i++ {
		go func(j int) {
			defer wg.Done()
			_ = c.Upsert(4, fmt.Sprint(j))
		}(i)
	}

	wg.Wait()
}

func TestCache_Parallel(t *testing.T) {
	c := New[int, string](3, time.Millisecond*100)

	t.Run("write six", func(t *testing.T) {
		t.Parallel()
		_ = c.Upsert(6, "six")
	})

	t.Run("write kuus", func(t *testing.T) {
		t.Parallel()
		_ = c.Upsert(6, "kuus")
	})
}

func TestCache_TTL(t *testing.T) {
	t.Parallel()

	c := New[string, string](3, time.Millisecond*100)
	_ = c.Upsert("Norwegian", "Blue")

	// Check the item is there.
	got, found := c.Read("Norwegian")
	assert.True(t, found)
	assert.Equal(t, "Blue", got)

	time.Sleep(time.Millisecond * 200)

	got, found = c.Read("Norwegian")

	assert.False(t, found)
	assert.Equal(t, "", got)
}

// TestCache_MaxSize tests the maximum capacity feature of a cache.
// It checks that update items are properly requeued as "new" items,
// and that we make room by removing the most ancient item for the new ones.
func TestCache_MaxSize(t *testing.T) {
	t.Parallel()

	// Give it a TTL long enough to survive this test
	c := New[int, int](3, time.Minute)

	_ = c.Upsert(1, 1)
	_ = c.Upsert(2, 2)
	_ = c.Upsert(3, 3)

	got, found := c.Read(1)
	assert.True(t, found)
	assert.Equal(t, 1, got)

	// Update 1, which will no longer make it the oldest
	_ = c.Upsert(1, 10)

	// Adding a fourth element will discard the oldest - 2 in this case.
	_ = c.Upsert(4, 4)

	// Trying to retrieve an element that should've been discarded by now.
	got, found = c.Read(2)
	assert.False(t, found)
	assert.Equal(t, 0, got)
}
