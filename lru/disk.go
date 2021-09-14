package lru

import (
	"encoding/json"
	"github.com/wuyan94zl/go-cache/byteview"
	"io"
	"io/ioutil"
	"os"
	"time"
)

type item struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	TTL   int64  `json:"ttl"`
}

const dbFile = "db"

func (c *Cache) write() {
	ticker := time.NewTicker(time.Duration(c.backupInterval) * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				D := make(map[string]item)
				c.wg.Add(1)
				for k, v := range c.cache {
					e := v.Value.(*entry)
					if e.ttl < time.Now().Unix() {
						c.ll.Remove(v)
						delete(c.cache, e.key)
					} else {
						D[k] = item{Key: e.key, Value: e.val.(byteview.ByteView).String(), TTL: e.ttl}
					}
				}
				c.wg.Done()
				res, _ := json.Marshal(D)
				os.Remove(dbFile)
				f, _ := os.OpenFile(dbFile, os.O_CREATE, 0666)
				io.WriteString(f, string(res))
				f.Close()
			}
		}
	}()
}

func (c *Cache) syncData() {
	str, _ := ioutil.ReadFile(dbFile)
	d := make(map[string]item)
	json.Unmarshal(str, &d)
	for _, v := range d {
		c.Add(v.Key, byteview.ByteView{B: []byte(v.Value)}, v.TTL-time.Now().Unix())
	}
}
