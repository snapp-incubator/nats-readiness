package main

import "github.com/go-resty/resty/v2"

func main() {
	client := resty.New()

	client.R().Get("127.0.0.1:8222/healthz")
}
