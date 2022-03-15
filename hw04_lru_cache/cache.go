package hw04lrucache

import "fmt"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
	PrintCache()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	val, ok := c.items[key]

	if !ok {
		c.items[key] = c.queue.PushFront(cacheItem{key, value})

		if c.queue.Len() > c.capacity {
			item := c.queue.Back().Value.(cacheItem)
			delete(c.items, item.key)
			c.queue.Remove(c.queue.Back())
		}
	} else {
		val.Value = cacheItem{key, value}
		c.queue.MoveToFront(val)
	}

	return ok
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	val, ok := c.items[key]
	if !ok {
		return nil, false
	}
	c.queue.MoveToFront(val)
	return val.Value.(cacheItem).value, ok
}

func (c *lruCache) Clear() {
	c.items = make(map[Key]*ListItem, c.capacity)
	c.queue.Clear()
}

/* Used for a debug */
func (c *lruCache) PrintCache() {
	fmt.Println("================================================")
	c.queue.PrintList(
		func(x interface{}) {
			fmt.Printf("[%v]=%v ", x.(cacheItem).key, x.(cacheItem).value)
		})
	fmt.Printf("\n{ ")

	for k := range c.items {
		fmt.Printf("%v ", k)
	}
	fmt.Printf("}\n")
}
