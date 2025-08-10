package locationdata

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
)

type BaseFlagService struct {
	logger *slog.Logger
}

type FlagService struct {
	BaseFlagService
	logger *slog.Logger
	http   HttpClient
}

func NewFlagService(http HttpClient, logger *slog.Logger) *FlagService {
	return &FlagService{http: http, logger: logger}
}

func (fs *BaseFlagService) CountryCodeToEmoji(ctx context.Context, code string) string {
	fs.logger.Debug("try get emoji from coutry code")
	runes := []rune{}
	for _, char := range code {
		if char >= 'A' && char <= 'Z' {
			runes = append(runes, rune(127397+char))
		} else if char >= 'a' && char <= 'z' {
			runes = append(runes, rune(127397+char-32))
		}
	}
	emoji := string(runes)
	fs.logger.Debug(fmt.Sprintf("emoji was generate - %s", emoji))

	return emoji
}

func (fs *FlagService) CountryCodeToPng(ctx context.Context, code string) ([]byte, error) {
	url := fmt.Sprintf("https://flagcdn.com/64x48/%s.png", strings.ToLower(code))

	data, err := fs.http.Get(ctx, url)
	if err != nil {
		fs.logger.Error("can`t flag image via http", err)
		return []byte{}, errors.Join(ErrGetImageViaHttp, err)
	}
	fs.logger.Debug("request was load success for flag image")

	bytesData, err := io.ReadAll(data)
	if err != nil {
		fs.logger.Error("can`t decode body to bytes data for image of flag", err)
		return []byte{}, err
	}
	fs.logger.Info("Flag was got success")

	return bytesData, nil
}
