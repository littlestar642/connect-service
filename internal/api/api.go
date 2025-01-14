package api

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Requester interface{
	SendPostRequest(endpoint string, count int)
}

type api struct {
	client *http.Client
}

func New() Requester {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	return &api{
		client: client,
	}
}

func (a *api) SendPostRequest(endpoint string, count int) {
	url := url.Values{}
	url.Add("count", fmt.Sprint(count))
	endpoint += "?" + url.Encode()

	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v\n", err)
	}

	log.Printf("Response Code: %d\n", resp.StatusCode)
}
