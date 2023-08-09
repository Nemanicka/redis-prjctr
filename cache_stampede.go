package main

import (
    "github.com/redis/go-redis/v9"
    "github.com/go-redis/cache/v9"
    "context"
    "fmt"
    "net/http"
    "log"
    "time"
    "strings"
    "math"
    "math/rand"
)

var ctx = context.Background()
var rdb *redis.Client
var server *redis.Client
var lc *cache.Cache

type Object struct {
	Value    string
	Expiry   time.Time
	Delta    time.Duration
}

func handleCGet(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Get")
	key := r.URL.Query().Get("key")

    var item Object
    keyMiss := false
    oerr := lc.Get(ctx, key, &item)
    if oerr != nil {
        fmt.Println("err", oerr)
        keyMiss = strings.Contains(oerr.Error(), "key is missing")
        if !keyMiss  {
            panic(oerr)
        }
    }
    
    forceUpdate := false
    if !keyMiss {
        randF := math.Log(rand.Float64())
        beta := 5.0
        subTime := time.Duration(float64(item.Delta) * randF * beta)
        fmt.Printf("randF = %v, delta = %v, subTime = %v\n", randF, item.Delta, subTime.Nanoseconds())
        tte := time.Now().Add( subTime )
        now := time.Now()
        exp := item.Expiry
        fmt.Printf("tte = %02d:%02d:%02d.%09d, expiry = %02d:%02d:%02d.%09d, now =  %02d:%02d:%02d.%09d\n", 
        tte.Hour(), tte.Minute(), tte.Second(), tte.Nanosecond(),
        exp.Hour(), exp.Minute(), exp.Second(), exp.Nanosecond(),
        now.Hour(), now.Minute(), now.Second(), now.Nanosecond())
        forceUpdate = tte.After(exp)
    }

    fmt.Printf("Cache miss = %v, forceUpdate = %v \n", keyMiss, forceUpdate)
    if keyMiss || forceUpdate {
            start := time.Now()
            value, verr := server.Get(ctx, key).Result()
            if verr != nil {
                w.Header().Set("Content-Type", "text/plain")
                fmt.Fprint(w, verr.Error())
                return
            }
            delta := time.Now().Sub(start)
            fmt.Printf("Key = %s, value = %s, delta = %d\n", key, value, delta)
            writeCache(key, value, delta)
            return    
    }

    fmt.Println("Cache hit")
    w.Header().Set("Content-Type", "text/plain")
    fmt.Fprint(w, item.Value)    
}


func writeCache(key, value string, delta time.Duration) {
    ttl := time.Second*5
    obj := &Object{
        Value: value,
        Expiry: time.Now().Add(ttl),
        Delta: delta,
    }

    if err := lc.Set(&cache.Item{
        Ctx:   ctx,
        Key:   key,
        Value: obj,
        TTL:   ttl,
    }); err != nil {
        panic(err)
    }
}

func handleCSet(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Set")
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")

    writeCache(key, value, 0) 

    w.WriteHeader(http.StatusOK)
}

func handleGet(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Get")
	key := r.URL.Query().Get("key")

    val, err := server.Get(ctx, key).Result()
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
    }

    w.Header().Set("Content-Type", "text/plain")
    fmt.Fprint(w, val)    
}

func handleSet(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Set")
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")

    err := server.Set(ctx, key, value, 0).Err()
    if err != nil {
        fmt.Println(err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    w.WriteHeader(http.StatusOK)
}

func init() {
    rdb = redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })
    server = redis.NewClient(&redis.Options{
        Addr:     "localhost:6397",
        Password: "", // no password set
        DB:       0,  // use default DB
    })

    lc = cache.New(&cache.Options{
        Redis:      rdb,
        LocalCache: cache.NewTinyLFU(1000, time.Minute),
    })

    key := "mykey"
    obj := &Object{
        Value: "mystring",
        Expiry: time.Now().Add(time.Second*5),
        Delta: time.Second,
    }

    if err := lc.Set(&cache.Item{
        Ctx:   ctx,
        Key:   key,
        Value: obj,
        TTL:   time.Hour,
    }); err != nil {
        panic(err)
    }

    var wanted Object
    if err := lc.Get(ctx, key, &wanted); err == nil {
        fmt.Println(wanted)
    }
}

func checkRedis(c *redis.Client) {
    fmt.Println("Check redis...")
    err := c.Set(ctx, "key", "value", 0).Err()
    if err != nil {
        panic(err)
    }

    val, err := c.Get(ctx, "key").Result()
    if err != nil {
        panic(err)
    }
    fmt.Println("key", val)
    if val == "value" {
        fmt.Println("Redis ok")
    }
}

func main() {
    fmt.Println("Check redis...")
    checkRedis(rdb)    
    checkRedis(server)    

    http.HandleFunc("/get", handleGet)
	http.HandleFunc("/set", handleSet)
    http.HandleFunc("/cget", handleCGet)
	http.HandleFunc("/cset", handleCSet)
    fmt.Printf("Server listening on port 8080...\n")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
