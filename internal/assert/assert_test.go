package assert_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/kyuff/htmx/internal/assert"
)

func TestAsserts(t *testing.T) {
	var testCases = []struct {
		name   string
		assert func(t *testing.T)
		failed bool
	}{
		{
			name: "Equal success",
			assert: func(t *testing.T) {
				assert.Equal(t, 1, 1)
			},
			failed: false,
		},
		{
			name: "Equal failed",
			assert: func(t *testing.T) {
				assert.Equal(t, 1, 2)
			},
			failed: true,
		},
		{
			name: "GreaterOrEqual success",
			assert: func(t *testing.T) {
				assert.GreaterOrEqual(t, 5, 5)
			},
			failed: false,
		},
		{
			name: "GreaterOrEqual success greater",
			assert: func(t *testing.T) {
				assert.GreaterOrEqual(t, 6, 5)
			},
			failed: false,
		},
		{
			name: "GreaterOrEqual failed",
			assert: func(t *testing.T) {
				assert.GreaterOrEqual(t, 4, 5)
			},
			failed: true,
		},
		{
			name: "LessOrEqual success",
			assert: func(t *testing.T) {
				assert.LessOrEqual(t, 5, 5)
			},
			failed: false,
		},
		{
			name: "LessOrEqual success less",
			assert: func(t *testing.T) {
				assert.LessOrEqual(t, 4, 5)
			},
			failed: false,
		},
		{
			name: "LessOrEqual failed",
			assert: func(t *testing.T) {
				assert.LessOrEqual(t, 6, 5)
			},
			failed: true,
		},
		{
			name: "Len success slice",
			assert: func(t *testing.T) {
				assert.Len(t, []int{1, 2, 3}, 3)
			},
			failed: false,
		},
		{
			name: "Len failed slice",
			assert: func(t *testing.T) {
				assert.Len(t, []int{1, 2, 3}, 2)
			},
			failed: true,
		},
		{
			name: "Len success string",
			assert: func(t *testing.T) {
				assert.Len(t, "hello", 5)
			},
			failed: false,
		},
		{
			name: "Len success map",
			assert: func(t *testing.T) {
				assert.Len(t, map[string]int{"a": 1}, 1)
			},
			failed: false,
		},
		{
			name: "Len success array",
			assert: func(t *testing.T) {
				assert.Len(t, [3]int{1, 2, 3}, 3)
			},
			failed: false,
		},
		{
			name: "Len success chan",
			assert: func(t *testing.T) {
				ch := make(chan int, 2)
				ch <- 1
				assert.Len(t, ch, 1)
			},
			failed: false,
		},
		{
			name: "Len failed unsupported type",
			assert: func(t *testing.T) {
				assert.Len(t, 42, 1)
			},
			failed: true,
		},
		{
			name: "NoError success",
			assert: func(t *testing.T) {
				assert.NoError(t, nil)
			},
			failed: false,
		},
		{
			name: "NoError failed",
			assert: func(t *testing.T) {
				assert.NoError(t, errors.New("oops"))
			},
			failed: true,
		},
		{
			name: "Error success",
			assert: func(t *testing.T) {
				assert.Error(t, errors.New("oops"))
			},
			failed: false,
		},
		{
			name: "Error failed",
			assert: func(t *testing.T) {
				assert.Error(t, nil)
			},
			failed: true,
		},
		{
			name: "Truef success",
			assert: func(t *testing.T) {
				assert.Truef(t, true, "should be true")
			},
			failed: false,
		},
		{
			name: "Truef failed",
			assert: func(t *testing.T) {
				assert.Truef(t, false, "should be true")
			},
			failed: true,
		},
		{
			name: "False success",
			assert: func(t *testing.T) {
				assert.False(t, false)
			},
			failed: false,
		},
		{
			name: "False failed",
			assert: func(t *testing.T) {
				assert.False(t, true)
			},
			failed: true,
		},
		{
			name: "Contains success",
			assert: func(t *testing.T) {
				assert.Contains(t, "hello world", "world")
			},
			failed: false,
		},
		{
			name: "Contains failed",
			assert: func(t *testing.T) {
				assert.Contains(t, "hello world", "missing")
			},
			failed: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// arrange
			var (
				x  = &testing.T{}
				wg sync.WaitGroup
			)

			// act
			// Run in a goroutine so that Fatalf's runtime.Goexit() is contained.
			wg.Add(1)
			go func() {
				defer wg.Done()
				testCase.assert(x)
			}()
			wg.Wait()

			// assert
			assert.Equal(t, testCase.failed, x.Failed())
		})
	}
}
