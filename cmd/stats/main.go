package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

type report map[string][]string

func main() {
	https := flag.Bool("https", false, "whether backends support HTTPs")

	flag.Parse()

	scheme := "http"
	if *https {
		scheme = "http"
	}

	serversPool := []string{
		"localhost:8080",
		"localhost:8081",
		"localhost:8082",
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	res := make([]report, len(serversPool))

	for i, s := range serversPool {
		resp, err := client.Get(fmt.Sprintf("%s://%s/report", scheme, s))
		if err == nil {
			var data report

			if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
				log.Printf("error parsing report %s: %s", s, err)
			} else {
				for k, v := range data {
					l := len(v)
					if l > 5 {
						l = 5
					}
					data[k] = v[len(v)-l:]
				}
				res[i] = data
			}
		} else {
			log.Printf("error %s %s", s, err)
		}

		if err = resp.Body.Close(); err != nil {
			continue
		}

		log.Println("=========================")
		log.Println("SERVER", i, serversPool[i])
		log.Println("=========================")

		data, _ := json.MarshalIndent(res[i], "", "  ")

		log.Println(string(data))
	}
}