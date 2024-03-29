package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(3)

		require.False(t, c.Set("a", 1))
		require.False(t, c.Set("b", 2))
		require.False(t, c.Set("c", 3))

		c.Clear()

		require.False(t, c.Set("a", 4))

		val, ok := c.Get("a")
		require.True(t, ok)
		require.Equal(t, 4, val)

		_, ok = c.Get("b")
		require.False(t, ok)

		_, ok = c.Get("c")
		require.False(t, ok)
	})

	t.Run("displace", func(t *testing.T) {
		c := NewCache(3)

		require.False(t, c.Set("a", 1)) // [a]
		require.False(t, c.Set("b", 2)) // [b, a]
		require.False(t, c.Set("c", 3)) // [c, b, a]
		require.False(t, c.Set("d", 4)) // [d, c, b]

		val, ok := c.Get("a")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("displace 2", func(t *testing.T) {
		c := NewCache(3)

		require.False(t, c.Set("a", 1)) // [a]
		require.False(t, c.Set("b", 2)) // [b, a]
		require.False(t, c.Set("c", 3)) // [c, b, a]
		require.True(t, c.Set("a", 10)) // [a, c, b]

		require.False(t, c.Set("d", 4)) // [d, a, c]

		val, ok := c.Get("b") // [d, a, c]
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("c") // [c, d, a]
		require.True(t, ok)
		require.Equal(t, 3, val)

		val, ok = c.Get("a") // [a, c, d]
		require.True(t, ok)
		require.Equal(t, 10, val)

		val, ok = c.Get("d") // [d, a, c]
		require.True(t, ok)
		require.Equal(t, 4, val)

		c.Get("a")                       // [a, d, c]
		c.Get("c")                       // [c, a, d]
		require.False(t, c.Set("b", 20)) // [b, c, a]

		_, ok = c.Get("d") // displaced
		require.False(t, ok)
		require.False(t, c.Set("d", 40)) // [d, b, c]
		val, ok = c.Get("d")
		require.True(t, ok)
		require.Equal(t, 40, val)
		val, ok = c.Get("b")
		require.True(t, ok)
		require.Equal(t, 20, val)
		val, ok = c.Get("c")
		require.True(t, ok)
		require.Equal(t, 3, val)
		_, ok = c.Get("a") // displaced
		require.False(t, ok)
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
