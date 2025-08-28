package cache

import (
	"myapp/src/internal/order_model"
	"container/list"
	"sync"
)

var maxSize int = 50

type Cache struct{
	mu sync.RWMutex
	orders map[string]*model.Order
	current_size int
	max_size int
	orderList *list.List
}

func NewCache() *Cache {
	return &Cache{
		orders: make(map[string]*model.Order),
		current_size: 0,
		max_size: maxSize,
		orderList: list.New(),
	}
}

func (c *Cache) Add(order *model.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.current_size >= c.max_size {
		c.removeOldest()
	}
	c.orderList.PushBack(order)
	c.orders[order.TrackNumber] = order
	c.current_size++
}

func (c *Cache) Get(trackNumber string) (*model.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, exists := c.orders[trackNumber]
	return order, exists
}


func (c *Cache) removeOldest(){
	oldest := c.orderList.Front()
	if oldest != nil {
		order := oldest.Value.(*model.Order)
		delete(c.orders, order.TrackNumber)
		c.orderList.Remove(oldest)
		c.current_size--
	}
}

func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.current_size
}

func (c *Cache) MaxSize() int {
	return c.max_size
}