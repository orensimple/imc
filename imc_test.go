package imc //nolint:testpackage

import (
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
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
	var (
		count     int32
		sMap      sync.Map
		dataMutex sync.Mutex
	)

	wg := &sync.WaitGroup{}

	t.Run("load", func(t *testing.T) {
		testCache := NewInMemoryCache()
		for i := 0; i < 1_000_000; i++ {
			wg.Add(1)

			go func(i int, count *int32) {
				defer wg.Done()
				testCache.GetOrSet(strconv.Itoa(rand.Intn(1000)), func() Value { atomic.AddInt32(count, 1); return strconv.Itoa(i) })
			}(i, &count)
		}

		for i := 0; i < 1_000; i++ {
			wg.Add(1)

			go func(i int, count *int32) {
				defer wg.Done()
				testCache.GetOrSet(strconv.Itoa(i), func() Value { atomic.AddInt32(count, 1); return strconv.Itoa(i) })
			}(i, &count)
		}

		wg.Wait()
		require.Equal(t, int32(1000), count)
		require.Equal(t, 1000, len(testCache.data))
	})

	t.Run("value", func(t *testing.T) {
		testCache := NewInMemoryCache()

		for i := 0; i < 100_000; i++ {
			wg.Add(1)

			go func(sMap *sync.Map, dataMutex *sync.Mutex) {
				defer wg.Done()
				dataMutex.Lock()
				key := strconv.Itoa(rand.Intn(100))
				value := strconv.Itoa(rand.Intn(100))
				imcValue := testCache.GetOrSet(key, func() Value { return value })
				sMapValue, _ := sMap.LoadOrStore(key, value)
				dataMutex.Unlock()
				require.Equal(t, imcValue, sMapValue)
			}(&sMap, &dataMutex)
		}

		wg.Wait()
	})
}
