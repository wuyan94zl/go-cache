package main

import (
	"github.com/wuyan94zl/go-cache"
	"net/http"
)

func main() {
	rpcPort, httpPort := "8821", ":8820"
	cache.Default(&cache.Config{
		MaxLen:         100000, // 缓存长度 默认10000
		BackupInterval: 1,      // 缓存备份间隔时间 1分钟备份一次（默认60分钟），备份文件保存在根目录 db 文件中，启动或重启时会自动把备份的db数据同步到内存缓存中
		Grpc: &cache.GrpcConfig{
			Port:        rpcPort,
			CurrentHost: "localhost:" + rpcPort,
			AllHosts:    []string{"localhost:8811", "localhost:8821", "localhost:8831"},
		},
	})

	http.HandleFunc("/api/cache", func(w http.ResponseWriter, r *http.Request) {
		// 使用缓存，根据key值计算hash，获取对应的节点（AllHosts 列表中的一个），和当前节点相同则直接本地操作，不同则以rpc协议访问远程节点操作
		cache.Instance.Set("test", "value", 60)
		cache.Instance.Get("test")
		cache.Instance.SetNX("test_nx", "value", 60)
	})
	http.ListenAndServe(httpPort, nil)
}
