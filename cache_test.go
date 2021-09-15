package cache

import (
	"encoding/json"
	"log"
	"sync"
	"testing"
	"time"
)

func init() {
	// 初始化缓存服务
	Default(nil)
}

func TestCache(t *testing.T) {
	var hit CallBack = CallBackFunc(func(key string, params map[string]interface{}) ([]byte, error) {
		t.Fatalf("cache hit err") // 应该缓存命中 却执行了函数 测试不通过
		return json.Marshal(params)
	})

	var add CallBack = CallBackFunc(func(key string, params map[string]interface{}) ([]byte, error) {
		return json.Marshal(params)
	})

	// 缓存不存在调用 CallBackFunc 并设置缓存  时效1秒
	Instance.CallBackFunc(add).Cache("test", map[string]interface{}{}, 1)

	time.Sleep(2 * time.Second)
	// sleep 2秒 缓存过期 调用CallBackFunc并设置缓存  时效1秒
	Instance.CallBackFunc(add).Cache("test", map[string]interface{}{}, 1)

	// 未sleep 获取缓存数据 cache hit
	Instance.CallBackFunc(hit).Cache("test", map[string]interface{}{}, 1)

	log.Println("cache test pass")
}

func TestSetGet(t *testing.T) {
	Instance.Set("set_key", "set_value", 1)
	time.Sleep(2 * time.Second)
	if _, err := Instance.Get("set_key"); err == nil {
		t.Fatalf("set_key 不存在 却返回数据")
	}
	Instance.Set("set_key", "set_value", 1)
	if _, err := Instance.Get("set_key"); err != nil {
		t.Fatalf("set_key 存在 却无返回数据")
	}
	log.Println("set get test pass")
}

func TestSetNX(t *testing.T) {
	wg := sync.WaitGroup{}
	// setNX 并发测试
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			log.Println("setNX 并发测试：", i, Instance.SetNX("keyNx", "1", 10))
			wg.Done()
		}(i, &wg)
	}
	wg.Wait()
	log.Println("setNX test pass")
}
