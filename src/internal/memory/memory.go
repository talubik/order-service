package memory

import (
	"log"
	"myapp/src/internal/cache"
	"myapp/src/internal/database"
	"myapp/src/internal/order_model"
	"sort"
)

type Memory struct{
	cache *cache.Cache
	repo *database.OrderRepository
}

func NewMemory(repo *database.OrderRepository) *Memory {
	mem := &Memory{
		cache: cache.NewCache(),
		repo: repo,
	}
	mem.loadCache()
	return mem
}

func (m *Memory) loadCache() {
	orders, err := m.repo.FindForCache(m.cache.MaxSize())
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].CreatedAt.Before(orders[j].CreatedAt)
	})
	if err != nil {
		log.Default().Println("Error loading cache:", err)
		return
	}
	for i := len(orders)-1; i >= 0; i-- {
		m.cache.Add(&orders[i])
	}
}

func (m *Memory) Add(order *model.Order) error {
	if err := m.repo.Create(order); err != nil {
		return err
	}
	m.cache.Add(order)
	return nil
}

func (m *Memory) Get(trackNumber string) (*model.Order, error) {
	if order, exists := m.cache.Get(trackNumber); exists {
		log.Println("Loaded from cache")
		return order, nil
	}
	order, err := m.repo.GetByTrackNumber(trackNumber)
	if err != nil {
		return nil, err
	}
	return order, nil
}
