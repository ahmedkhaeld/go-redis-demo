package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (a *RedisCache) Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in the handler")

	// extract the query string from the requested url
	q := r.URL.Query().Get("q")
	data, cacheHit, err := a.getData(r.Context(), q)
	if err != nil {
		fmt.Printf("error calling data source: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := RedisCacheResponse{
		Cache: cacheHit,
		Data:  data,
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		fmt.Printf("error encoding response: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
