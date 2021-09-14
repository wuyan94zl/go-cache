package lru

import (
	"container/list"
	"sync"
	"time"
)

type Cache struct {
	maxLen         int
	currLen        int
	backupInterval int
	ll             *list.List
	cache          map[string]*list.Element
	wg             sync.WaitGroup
	OnEvicted      func(key string, val interface{})
}

type entry struct {
	key string
	val interface{}
	ttl int64
}

func New(maxLen int, backupInterval int, onEvicted func(string2 string, value interface{})) *Cache {
	c := &Cache{
		maxLen:         maxLen,
		backupInterval: backupInterval,
		ll:             list.New(),
		cache:          make(map[string]*list.Element),
		OnEvicted:      onEvicted,
	}
	c.syncDisk()
	c.writeDisk()
	return c
}

func (c *Cache) Get(key string) (value interface{}, ok bool) {
	c.wg.Wait()
	if ele, ok := c.cache[key]; ok {
		kv := ele.Value.(*entry)
		if kv.ttl > time.Now().Unix() {
			c.ll.MoveToFront(ele) // 将ele 移动到对尾
			return kv.val, true
		} else {
			c.ll.Remove(ele)
			delete(c.cache, kv.key)
		}
	}
	return
}

func (c *Cache) RemoveOldest() {
	c.wg.Wait()
	ele := c.ll.Back() // 取首节点，并删除
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.val)
		}
	}
}

func (c *Cache) Add(key string, value interface{}, ttl int64) {
	c.wg.Wait()
	if ele, ok := c.cache[key]; ok {
		kv := ele.Value.(*entry)
		kv.val = value
		kv.ttl = time.Now().Unix() + ttl
		c.ll.MoveToFront(ele)
	} else {
		ele := c.ll.PushFront(&entry{key: key, val: value, ttl: time.Now().Unix() + ttl})
		c.cache[key] = ele
	}
	for c.maxLen < c.Len() {
		c.RemoveOldest()
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}

func (c *Cache) writeDisk() {
	c.write()
}

func (c *Cache) syncDisk() {
	c.syncData()
}
