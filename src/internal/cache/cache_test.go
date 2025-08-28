package cache

import (
	generator "myapp/src/testutils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCache_AddAndGet(t *testing.T) {
	cache := NewCache()
	order := generator.GenerateOrder("WB123456789")
	cache.Add(&order)
	retrievedOrder, exists := cache.Get("WB123456789")
	if !exists {
		t.Fatal("Expected order to exist in cache")
	}
	if retrievedOrder.TrackNumber != order.TrackNumber {
		t.Fatalf("Expected track number %s, got %s", order.TrackNumber, retrievedOrder.TrackNumber)
	}
}

func TestCacheMaxSize(t *testing.T) {
	cache := NewCache()
	orders := generator.GenerateOrders(100)
	for _, order := range orders {
		cache.Add(&order)
		assert.True(t, cache.Size() <= cache.MaxSize())
	}
}
