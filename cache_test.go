package tiny_cache

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"github.com/go-redis/redis"

	"math/rand"
	"reflect"
	"testing"
	"time"
)

var db = map[string]string{
	"T":          "Fj",
	"Xw":         "mmLYaB0N",
	"6ZTSxCKs0":  "Byk5SH9pBF",
	"xnzyNnlR":   "_2wLrbN",
	"IbZo":       "Tn",
	"i":          "dj",
	"RIYV":       "2g",
	"MYpb3yx_":   "aNLWxS9QVO",
	"a_hMcdvZa":  "Fy47U",
	"oXhv":       "y-AlWqW31E",
	"-EyJ":       "9EXRWv1C",
	"A3T2kk":     "JOLK_Mr",
	"vI":         "nn3x2Slv",
	"9t4":        "d",
	"8XpMchWj":   "6x5YScYH",
	"Jm-XZHlLNz": "XoIR",
	"1976ki0":    "TySr99-tY",
	"Huj0O":      "2L6cO",
	"nexqQOmh":   "8DPCgjBhpk",
	"Ig0L88qZ":   "Rh-",
	"ZfYG9q":     "m_L7JskB",
	"TN4WB-4pum": "tDygi",
	"Q0F5068tG":  "nR",
	"N":          "t2Sc",
	"k0":         "gEdd",
	"r76PtW8Y":   "gjNQSEUX6",
	"91g2t":      "W",
	"h":          "ZX",
	"6j3-":       "Nq",
	"KuVBGU":     "VbzqQG",
	"1hDK":       "bjP-qjKcmx",
	"J0_7":       "X1q",
	"RMx9":       "Bp8Bo",
	"jRL0Elcy":   "My91drBiFM",
	"6jn":        "xsbvg2Ly9U",
	"Cdb8JaqKL":  "WjqUtNM",
	"FD":         "By-2",
	"NdYEqfP7T":  "G",
	"EKeWBzLKU":  "w_UgyI",
	"ASM":        "WJusuxXZr7",
	"acLizc8j":   "vxd0rOdw6",
	"VesjOr":     "0xPIo7Q3Dr",
	"_ez1I":      "pvHtC",
	"a":          "Fi5",
	"fTq6U68Zd5": "i2Vw_6C3",
	"_k46":       "m1ccJ",
	"PprAsAznO_": "JxHzII",
	"Syy9nCiH":   "O5ynQP",
	"k":          "FyX1D",
	"H":          "v8r8nhOGR",
}

func CacheTest(t *testing.T) {
	start := time.Now()
	for i := 0; i < 100; i++ {
		for k, v := range db {
			// tiny-cache hit
			resp, err := http.Get("http://localhost:9999/api?key=" + k)
			if err != nil {
				fmt.Println("请求失败：", err)
				return
			}
			go func(v string) {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					fmt.Println("读取响应失败：", err)
					return
				}
				if string(body) != v {
					panic("tiny-cache failed to get data")
				}
				resp.Body.Close()
			}(v)
		}
	}
	fmt.Println("tiny-cache::----------------", time.Since(start).Seconds(), "---------------------")
}
func RedisTest(t *testing.T) {
	start := time.Now()
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	client.ConfigSet("maxmemory-policy", "allkeys-lfu")
	client.ConfigSet("maxmemory", "10mb")
	for k, v := range db {
		client.Set(k, v, 0).Err()
	}

	for i := 0; i < 100; i++ {
		for k, v := range db {
			val, err := client.Get(k).Result()
			if err != nil {
				panic(err)
			}
			if val != v {
				t.Fatalf("Redis fail to get data")
			}
		}
	}
	fmt.Println("redis::----------------", time.Since(start).Seconds(), "---------------------")

}
func generateRandomString(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length]
}

func TestCompareWithRedis(t *testing.T) {
	CacheTest(t)
	RedisTest(t)
}

func TestGetter(t *testing.T) {
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expect := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Errorf("callback failed")
	}
}
