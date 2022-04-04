package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	fmt.Println("starting server")

	RedisCache := NewRedisCache()

	http.HandleFunc("/RedisCache", RedisCache.Handler)

	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil)
}

type RedisCacheResponse struct {
	Cache bool                `json:"cache"`
	Data  []NominatimResponse `json:"data"`
}
