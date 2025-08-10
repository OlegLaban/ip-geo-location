package ipservice

import (
	"context"
	"errors"
	"io"
	"log/slog"
)

type HttpClient interface {
	Get(context.Context, string) (io.ReadCloser, error)
}

type IPService struct {
	httpClient HttpClient
	logger     *slog.Logger
}

func NewIPService(httpClient HttpClient, logger *slog.Logger) *IPService {
	return &IPService{httpClient: httpClient, logger: logger}
}

func (ips *IPService) GetPublicIP(ctx context.Context) (string, error) {
	ips.logger.Debug("Try get public IP via api.ipify")
	rc, err := ips.httpClient.Get(ctx, "https://api.ipify.org/")
	if err != nil {
		ips.logger.Error("can`t get public ip via api.ipify", err)
		return "", errors.Join(ErrGetIP, err)
	}
	ips.logger.Debug("request for getting public ip was got successfuly")
	defer rc.Close()
	body, err := io.ReadAll(rc)
	if err != nil {
		ips.logger.Error("can`t decode body for getting public ip", err)
		return "", errors.Join(ErrDecodeBody, err)
	}
	ip := string(body)
	ips.logger.Debug("ip was got success ip - " + ip)
	ips.logger.Info("ip was got successfuly")
	return ip, nil
}
