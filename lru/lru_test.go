package lru

import (
	"testing"
)

func TestCacheGet(t *testing.T) {
	lru := New(2, nil)

	lru.Add("test1", "test1", 1)
	lru.Add("test2", "test2", 1)

	// 添加后能正确获取到值
	if _, ok := lru.Get("test2"); !ok {
		t.Fatalf("cache hit test2")
	}
	// 添加后能正确获取到值
	if _, ok := lru.Get("test1"); !ok {
		t.Fatalf("cache hit test1")
	}

	lru.Add("test3", "test3", 1)
	// 添加test3时超过容量，删除test2(最久被使用的数据)，此时test1 和 test3 有数据
	if _, ok := lru.Get("test2"); ok {
		t.Fatalf("cache hit test1 is error")
	}
}
