package cache

import (
	"encoding/json"
	"testing"
	"time"
)

func TestSingle(t *testing.T) {
	var hit CallBack = CallBackFunc(func(key string, params map[string]interface{}) ([]byte, error) {
		t.Fatalf("cache hit err") // 应该缓存命中 却执行了函数 测试不通过
		return json.Marshal(params)
	})

	var add CallBack = CallBackFunc(func(key string, params map[string]interface{}) ([]byte, error) {
		return json.Marshal(params)
	})

	// 初始化缓存服务
	Default(nil)

	// 缓存不存在调用 CallBackFunc 并设置缓存  时效1秒
	p := make(map[string]interface{})
	p["data"] = 123456789
	Instance.CallBackFunc(add).Get("test", p, 1)

	time.Sleep(2 * time.Second)
	// sleep 2秒 缓存过期 调用CallBackFunc并设置缓存  时效1秒
	Instance.CallBackFunc(add).Get("test", p, 1)

	// 未sleep 获取缓存数据 cache hit
	Instance.CallBackFunc(hit).Get("test", p, 1)
}

func TestCluster(t *testing.T) {
	//http.HandleFunc("/api/cache")
}

func cluster1(){

}

func cluster2(){

}
