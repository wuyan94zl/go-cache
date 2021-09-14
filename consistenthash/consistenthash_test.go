package consistenthash

import (
	"log"
	"testing"
)

func TestHashing(t *testing.T) {
	//hash := New(3, func(key []byte) uint32 {
	//	var i int
	//	for v,_ := range key{
	//		i += v
	//	}
	//	return uint32(i)
	//})
	hash := New(3, nil)

	hash.Add("localhost:9999","localhost:9998","localhost:9997")

	testCases := map[string]string{
		"2":"2",
		"11":"2",
		"23":"4",
		"27":"2",
	}
	log.Println(hash.Get("test123"))
	log.Println(hash.Get("test234"))
	log.Println(hash.Get("test345"))
	log.Println(hash.Get("test456"))
	log.Println(hash.Get("test567"))
	log.Println(hash.Get("test678"))
	//for k, v := range testCases {
	//	if hash.Get(k) != v {
	//		t.Errorf("Asking for %s, should have yielded %s", k, v)
	//	}
	//}

	// Adds 8, 18, 28
	hash.Add("8")

	// 27 should now map to 8.
	testCases["27"] = "8"

	//for k, v := range testCases {
	//	if hash.Get(k) != v {
	//		t.Errorf("Asking for %s, should have yielded %s", k, v)
	//	}
	//}

}
