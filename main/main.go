package main

import (
	tiny_cache "7daysgo/tiny-cache"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":        "Tom",
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

func createGroup() *tiny_cache.Group {
	return tiny_cache.MakeGroup("scores", 2<<10, tiny_cache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

func startCacheServer(addr string, addrs []string, gee *tiny_cache.Group) {
	server := tiny_cache.MakeHTTPPool(addr)
	server.Set(addrs...)
	gee.RegisterPeers(server)
	log.Println("tiny-cache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], server))
}

func startAPIServer(apiAddr string, gee *tiny_cache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := gee.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())
		}))
	log.Println("frontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func startCacheServerGRPC(addr string, addrs []string, gee *tiny_cache.Group) {
	server := tiny_cache.MakeGrpcPool(addr)
	server.Set(addrs...)
	gee.RegisterPeers(server)
	log.Println("tiny-cache is running at", addr)
	server.Run()
}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8004, "tiny-cache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "localhost:8001", //HTTP VERSION: http://localhost:8001
		8002: "localhost:8002", //HTTP VERSION: http://localhost:8002
		8003: "localhost:8003", //HTTP VERSION: http://localhost:8003
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	gee := createGroup()
	if api {
		go startAPIServer(apiAddr, gee)
	}
	//HTTP VERSION: startCacheServer(addrMap[port], addrs, gee)
	startCacheServerGRPC(addrMap[port], addrs, gee)
}
