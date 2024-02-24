package cache

import (
	"sync"
	"time"
	"errors"
)


type Cache struct {
	sync.RWMutex
	defaultExpiration time.Duration
    cleanupInterval   time.Duration
	orders map[string]Order
}

type Order struct {
    Value      interface{}
    Expiration int64
	Created    time.Time
}

func New(defaultExpiration, cleanupInterval time.Duration) *Cache {

    orders := make(map[string]Order)

    cache := Cache{
        orders:             orders,
        defaultExpiration: defaultExpiration,
        cleanupInterval:   cleanupInterval,
    }

    if cleanupInterval > 0 {
        cache.StartGC()
    }

    return &cache
}

func (c *Cache) Set(orderID string, value interface{}, duration time.Duration) {

	var expiration int64

	if duration == 0 {
        duration = c.defaultExpiration
    } else {
		expiration = time.Now().Add(duration).UnixNano()
	}

	c.Lock()
	defer c.Unlock()
	c.orders[orderID] = Order{
        Value:      value,
        Expiration: expiration,
        Created:    time.Now(),
    }
}

func (c *Cache) Get(orderID string) (interface{}, bool) {

    c.RLock()
    defer c.RUnlock()
    item, found := c.orders[orderID]

    if !found {
        return "orderID", false
    }

    if item.Expiration > 0 {
        if time.Now().UnixNano() > item.Expiration {
            delete(c.orders, orderID)
            return nil, false
        }

    }
    item.Expiration = time.Now().Add(c.defaultExpiration).UnixNano()
    c.orders[orderID] = item

    return item.Value, true
}


func (c *Cache) Delete(orderID string) error {

    c.Lock()
    defer c.Unlock()

    if _, found := c.orders[orderID]; !found {
        return errors.New("order ID not found")
    }
    delete(c.orders, orderID)
    return nil
}

func (c *Cache) Count() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.orders)
}


func (c *Cache) StartGC()  {
    go c.GC()
}

func (c *Cache) GC() {

    for {
        <-time.After(c.cleanupInterval)
        if c.orders == nil {
            return
        }
        if orderIDs := c.expiredIDs(); len(orderIDs) != 0 {
            c.clearOrders(orderIDs)

        }

    }

}

func (c *Cache) expiredIDs() (orderIDs []string) {

    c.RLock()
    defer c.RUnlock()
    for k, i := range c.orders {
        if time.Now().UnixNano() > i.Expiration && i.Expiration > 0 {
            orderIDs = append(orderIDs, k)
        }
    }
    return
}

func (c *Cache) clearOrders(orderIDs []string) {

    c.Lock()
    defer c.Unlock()
    for _, k := range orderIDs {
        delete(c.orders, k)
    }
}