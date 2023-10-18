package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value any) bool
	Get(key Key) (any, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mutex    sync.Mutex
}

type cacheItem struct {
	key   Key
	value any
}

// Создать новый кэш с ёмкостью capacity.
func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

// Добавить значение в кэш по ключу.
func (c *lruCache) Set(key Key, value any) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if element, exists := c.items[key]; exists {
		// Если такое значение есть
		// Переместить в начало
		c.queue.MoveToFront(element)
		// и обновить значение
		element.Value.(*cacheItem).value = value
		// выдать присутствие
		return true
	}

	// Если ёмкость превышена
	if c.queue.Len() >= c.capacity {
		c.removeOldestItem()
	}

	item := &cacheItem{
		key:   key,
		value: value,
	}

	element := c.queue.PushFront(item)
	c.items[item.key] = element

	return false
}

// Получить значение из кэша по ключу.
func (c *lruCache) Get(key Key) (any, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	element, exists := c.items[key]

	// Если элемента нет
	if !exists {
		// Вернуть пустое значение и false
		return nil, false
	}

	// Передвинуть в начало и венуть его
	c.queue.MoveToFront(element)
	return element.Value.(*cacheItem).value, true
}

// Очистить кэш.
func (c *lruCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

// Удалить самый старый элемент из кэша.
func (c *lruCache) removeOldestItem() {
	lastElement := c.queue.Back()                       // взять последний элемент
	c.queue.Remove(lastElement)                         // удалить из очереди
	delete(c.items, lastElement.Value.(*cacheItem).key) // удалить из словаря
}
