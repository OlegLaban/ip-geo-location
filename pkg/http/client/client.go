package client

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"time"
)

type Client struct {}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Get(ctx context.Context, url string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "curl/7.64.1")

	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		resp, err := client.Do(req)
		if err == nil {
			return resp.Body, nil
		}
		log.Printf("Попытка %d не удалась: %v", i+1, err)
		time.Sleep(time.Second * 2)
	}
	return nil, errors.New("can`t get IP")

}