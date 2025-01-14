package api


import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type RequestBody struct {
	Count int `json:"count"`
}

func SendPostRequest(endpoint string, count int) {
	url := url.Values{}
	url.Add("count", fmt.Sprint(count))
	endpoint += "?" + url.Encode()

	_, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return
	}
}
