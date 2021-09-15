package main

import (
	"encoding/json"
	"github.com/wuyan94zl/go-cache"
	"log"
	"net/http"
	"sync"
)

func main() {
	cache.Default(nil)
	http.HandleFunc("/api/cache", myHandler)
	http.ListenAndServe(":8888", nil)
}
func myHandler(res http.ResponseWriter, req *http.Request) {
	var f cache.CallBack = cache.CallBackFunc(func(key string, params map[string]interface{}) ([]byte, error) {
		return json.Marshal("cache write data")
	})

	// Cache
	r, _ := cache.Instance.CallBackFunc(f).Cache("cache_key", map[string]interface{}{}, 600)
	log.Println("cache：", string(r.B))

	// Set
	log.Println("set：", cache.Instance.Set("set_key", "set_val", 100))

	// Get
	log.Println(cache.Instance.Get("set_key"))

	// setNX 使用
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			log.Println("setNX：", i, cache.Instance.SetNX("keyNx", "1", 10))
			wg.Done()
		}(i, &wg)
	}
	wg.Wait()
}
