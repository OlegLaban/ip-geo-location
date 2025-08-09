package locationdata

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type FlagService struct{}

func NewFlagService() *FlagService {
	return &FlagService{}
}

func (fs *FlagService) CountryCodeToEmoji(code string) string {
	runes := []rune{}
	for _, char := range code {
		if char >= 'A' && char <= 'Z' {
			runes = append(runes, rune(127397+char))
		} else if char >= 'a' && char <= 'z' {
			runes = append(runes, rune(127397+char-32))
		}
	}
	return string(runes)
}

func (fs *FlagService) CountryCodeToPng(code string) ([]byte, error) {
	url := fmt.Sprintf("https://flagcdn.com/64x48/%s.png", strings.ToLower(code))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err
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
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)
			return body, nil
		}
		log.Printf("Попытка %d не удалась: %v", i+1, err)
		time.Sleep(time.Second * 2)
	}

	return []byte{}, errors.New("can`t get flag")
}
