package locationdata

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
)

type IPService interface {
	GetPublicIP(ctx context.Context) (string, error)
}

type HttpClient interface {
	Get(context.Context, string) (io.ReadCloser, error)
}

type GeoResponse struct {
	Country     string `json:"country_name"`
	CountryCode string `json:"country"`
}

type GeoService struct {
	http      HttpClient
	IPService IPService
	cache     CacheInterface
	logger    *slog.Logger
}

func NewGeoService(IPService IPService, http HttpClient, cache CacheInterface, logger *slog.Logger) *GeoService {
	return &GeoService{IPService: IPService, http: http, cache: cache, logger: logger}
}

func (gs *GeoService) GetCountryData(ctx context.Context) (GeoData, error) {
	var readCloser io.ReadCloser
	ip, err := gs.IPService.GetPublicIP(ctx)
	if err != nil {
		gs.logger.Error("can`t get public ip", "err", err)
		return GeoData{}, errors.Join(ErrCantGetIP, err)
	}
	gs.logger.Info("ip was got successfuly")
	bytesData, err := gs.cache.GetWithCallback(ip, func() ([]byte, error) {
		rc, err := gs.loadViaAPI(ctx)
		if err != nil {
			gs.logger.Error("can`t get geodata via API", "err", err)
			return []byte{}, err
		}
		gs.logger.Info("geodata was got successfuly via API")
		return io.ReadAll(rc)
	})
	if err != nil {
		gs.logger.Error("can`t get geodata from cache or API", "err", err)
		return GeoData{}, errors.Join(ErrCantGetDataFromCache, err)
	}
	gs.logger.Info("geodata was got successfuly via cache or API")
	readCloser = io.NopCloser(bytes.NewReader(bytesData))

	defer func () {
		if err := readCloser.Close(); err != nil {
			gs.logger.Error("can`t close reader with geodata", "err", err)
		}
	}()

	var data GeoResponse
	err = json.NewDecoder(readCloser).Decode(&data)
	if err != nil {
		gs.logger.Error("can`t decode geodata to json", "err", err)
		return GeoData{}, errors.Join(ErrCantDecodeGeoData, err)
	}
	gs.logger.Info("geodata was loaded successfuly")

	return GeoData{Country: data.CountryCode, CountryCode: data.CountryCode}, nil
}

func (gs *GeoService) loadViaAPI(ctx context.Context) (io.ReadCloser, error) {
	url := "https://ipinfo.io/json"

	return gs.http.Get(ctx, url)
}

func (gs *GeoService) GetRc(ctx context.Context) (io.ReadCloser, error) {
	readCloser, err := gs.loadViaAPI(ctx)
	if err != nil {
		gs.logger.Error("can`t get geodata via API", "err", err)
		return nil, errors.Join(ErrCantGetGeoData, err)
	}
	gs.logger.Debug("geodata was loaded successfuly")
	return readCloser, nil
}
