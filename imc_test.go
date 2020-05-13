package imc //nolint:golint,stylecheck

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetOrSet(t *testing.T) {
	testCache := NewInMemoryCache()

	t.Run("simple", func(t *testing.T) {
		val := testCache.GetOrSet("test1", func() Value { return "test1" })
		require.Equal(t, "test1", val)
		require.Equal(t, 1, len(testCache.data))

		val = testCache.GetOrSet("test2", func() Value { return "test2" })
		require.Equal(t, "test2", val)
		require.Equal(t, 2, len(testCache.data))

		val = testCache.GetOrSet("", func() Value { return "test3" })
		require.Equal(t, "test3", val)
		require.Equal(t, 3, len(testCache.data))

		val = testCache.GetOrSet("тест1", func() Value { return "тест1" })
		require.Equal(t, "тест1", val)
		require.Equal(t, 4, len(testCache.data))
	})

	t.Run("replay", func(t *testing.T) {
		val := testCache.GetOrSet("test1", func() Value { return "test4" })
		require.Equal(t, "test1", val)
		require.Equal(t, 4, len(testCache.data))

		val = testCache.GetOrSet("test2", func() Value { return "test5" })
		require.Equal(t, "test2", val)
		require.Equal(t, 4, len(testCache.data))
	})
}

func TestGetOrSetMultithreading(t *testing.T) {
	testCache := NewInMemoryCache()
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		for i := 0; i < 1_000_000; i++ {
			go testCache.GetOrSet(strconv.Itoa(rand.Intn(1_000)), func() Value { return strconv.Itoa(rand.Intn(1_000)) })
		}
	}()

	wg.Wait()
}
