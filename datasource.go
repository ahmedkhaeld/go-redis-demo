package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-redis/redis/v8"
)

// getData takes in query string that fetch a specifi url
// and returns slice of data from Nominatim
// getData here calling an external RedisCache, as a data source
func (a *RedisCache) getData(ctx context.Context, q string) ([]NominatimResponse, bool, error) {
	// is query cached?
	value, err := a.cache.Get(ctx, q).Result()
	// if key doesn't exist
	if err == redis.Nil {
		// we want call external data source and escape the url to safely put the query withing it
		escapedQ := url.PathEscape(q)
		address := fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json", escapedQ)

		resp, err := http.Get(address)
		if err != nil {
			return nil, false, err
		}

		// in order to read the reponse body, create a data container
		data := make([]NominatimResponse, 0)

		// decode the response data
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return nil, false, err
		}

		b, err := json.Marshal(data)
		if err != nil {
			return nil, false, err
		}

		// set the value
		err = a.cache.Set(ctx, q, bytes.NewBuffer(b).Bytes(), time.Second*15).Err()
		if err != nil {
			return nil, false, err
		}
		// return the response
		return data, false, nil

	} else if err != nil {
		fmt.Printf("error calling redis: %v\n", err)
		return nil, false, err
	} else {
		// cache hit
		data := make([]NominatimResponse, 0)

		// build response
		err := json.Unmarshal(bytes.NewBufferString(value).Bytes(), &data)
		if err != nil {
			return nil, false, err
		}

		// return response which is a cache hit
		return data, true, nil
	}
}
