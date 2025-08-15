package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type Client struct {
	logger *slog.Logger
}

func NewClient(logger *slog.Logger) *Client {
	return &Client{logger: logger}
}

func (c *Client) Get(ctx context.Context, url string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	c.logger.Debug("Request was created for method GET")
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
	c.logger.Debug("Client was created with timeout and DisableKeepAlives")
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		resp, err := client.Do(req)
		if err == nil {
			return resp.Body, nil
		}
		c.logger.Debug(fmt.Sprintf("Trying %d was fail", i+1), "err", err)
		time.Sleep(time.Second * 2)
	}
	return nil, errors.Join(ErrDoRequest, fmt.Errorf("can`t do request via url %s", url))

}
