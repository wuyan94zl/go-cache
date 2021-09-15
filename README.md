### go-cache
***
distributed cache running in project  
在项目中运行的分布式缓存
***
#### 安装
`go get github.com/wuyan94zl/go-cache`
***
#### 实例化
```
    // 引入缓存包
    import "github.com/wuyan94zl/go-cache"
    
    // 实例化单机版缓存服务
    cache.Default(nil)

    // 实例化分布式版缓存服务
    cache.Default(&cache.Config{
        MaxLen:         100000, // 缓存长度 默认10000
        BackupInterval: 1,      // 缓存备份间隔时间 1分钟备份一次（默认60分钟），备份文件保存在根目录 db 文件中，启动或重启时会自动把备份的db数据同步到内存缓存中
        // 远程节点基于grpc方式交互的分布式缓存配置
        Grpc: &cache.GrpcConfig{
            Port:        "8888",                                       // 当前服务监听端口
            CurrentHost: "localhost:8888",                             // 当前服务Grpc地址
            AllHosts:    []string{"localhost:8811", "localhost:8821"}, // 所有服务Grpc地址 直连需要手动配置
        },
    })
    
    // 操作使用
    
    // 设置一个键为test值为value的缓存，60秒后过期
    cache.Instance.Set("test", "value", 60)
    
    // 获取键为test的缓存值
    cache.Instance.Get("test")
    
    // 设置一个键为test_nx值为value的缓存，60秒后过期，返回true。如果键test_nx缓存存在则不操作返回false
    cache.Instance.SetNX("test_nx", "value", 60)
    
    // 获取一个键为cache_key的缓存值，不存在则调用回调函数f()获取数据并添加缓存。存在则直接返回，不会触发f()回调函数
    var f cache.CallBack = cache.CallBackFunc(func(key string, params map[string]interface{}) ([]byte, error) {
        return json.Marshal("cache write data")
    })
    // 回调函数的参数对应Cache函数的key,map[string]interface{}{}
    cache.Instance.CallBackFunc(f).Cache("cache_key", map[string]interface{}{}, 600)
```

### 项目中使用

#### 单机版

[单机版示例](/example/single.go)

#### 分布式版

[分布式版示例1](/example/cluster1.go)

[分布式版示例2](/example/cluster2.go)

[分布式版示例3](/example/cluster3.go)

3个示例仅 `rpcPort, httpPort := "8811", ":8810"` 端口不同  
保证 AllHosts 列表一致且包含所有的缓存节点信息

***