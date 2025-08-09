package ipservice

import (
	"context"
	"errors"
	"io"
)

type HttpClient interface{
	Get(context.Context, string) (io.ReadCloser, error)
}

type IPService struct {
	httpClient HttpClient
}

func NewIPService(httpClient HttpClient) *IPService {
	return &IPService{httpClient: httpClient}
}

func (ips *IPService) GetPublicIP(ctx context.Context) (string, error) {
	rc, err := ips.httpClient.Get(ctx, "https://api.ipify.org/")
	if err != nil {
		return "", errors.Join(ErrGetIP, err)
	}

	defer rc.Close()
	body, err := io.ReadAll(rc)	
	if err != nil {
		return "", errors.Join(ErrDecodeBody, err)
	}

	return string(body), nil
}
