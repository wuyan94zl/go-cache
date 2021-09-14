### go-cache
***
distributed cache running in project  
在项目中运行的分布式缓存
***
#### 安装
`go get github.com/wuyan94zl/go-cache`
***
#### 单机版使用
```
import "github.com/wuyan94zl/go-cache"

func main(){
    // 实例化单机版缓存
    cache.Default(nil)
    
    // 运行web服务
    http.HandleFunc("/api/cache", myHandler)
	http.ListenAndServe(":8888", nil)
}

func myHandler(res http.ResponseWriter, req *http.Request) {
	// 定义CallBack 在没有缓存数据时 获取新数据
	var f cache.CallBack = cache.CallBackFunc(func(key string, params map[string]interface{}) ([]byte, error) {
	    time.Sleep(5 * time.Second) //模拟耗时操作
		return json.Marshal("test cache data")
	})
	// cache.Instance 为缓存实例 CallBackFunc() 传入定义好的数据来源 Get(key,params,60) key和params为CallBack回调函数参数， 600为缓存时间：600秒
	// 当key值缓存不存在时调用f()获取数据并写入缓存，存在时直接拉去缓存数据，不会执行f()
	r, _ := cache.Instance.CallBackFunc(f).Get(key, map[string]interface{}{}, 600)
	res.Write(r.ByteSlice())
}
```
运行程序 `go run main`  
多次访问 localhost:8888/api/cache  
第一次访问会执行5秒，后面再次访问就会很快啦
***
### 分布式集群使用

```
import "github.com/wuyan94zl/go-cache"
func main(){
    // 实例化分布式缓存
    cache.Default(&cache.Config{
        MaxLen:         100000, // 缓存长度 默认10000
        BackupInterval: 1,      // 缓存备份间隔时间 1分钟备份一次（默认60分钟），备份文件保存在根目录 db 文件中，启动或重启时会自动把备份的db数据同步到内存缓存中
        // 基于grpc方式交互 的 分布式缓存配置
        Grpc: &cache.GrpcConfig{
            Port:        "8881",                                       // 当前服务监听端口
            CurrentHost: "localhost:8881",                             // 当前服务Grpc地址
            AllHosts:    []string{"localhost:8882", "localhost:8881"}, // 所有服务Grpc地址 直连需要手动配置
        },
    })
    // 运行web服务省略同上
}
```


8881 和 8882 为2个程序缓存的相关节点信息  
在程序1 或 程序2 中调用 cache.Instance.CallBackFunc(f).Get(key, map[string]interface{}{}, 600)
程序会根据key值 计算出相应的缓存节点
如果是本节点缓存，直接获取程序内存缓存数据，缓存不存在，调用回调获取数据并写入本节点缓存。
如果非本节点缓存，根据节点远程获取缓存数据，缓存不存在，调用回调获取数据并写入本节点缓存，10秒过期（）。
  


```
