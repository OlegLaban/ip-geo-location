package locationdata

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
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
	useCache  bool
}

func NewGeoService(IPService IPService, http HttpClient, cache CacheInterface, useCache bool) *GeoService {
	return &GeoService{IPService: IPService, http: http, cache: cache, useCache: useCache}
}

func (gs *GeoService) GetCountryData(ctx context.Context) (GeoData, error) {
	ip, err := gs.IPService.GetPublicIP(ctx)
	var readCloser io.ReadCloser
	if err != nil {
		return GeoData{}, errors.Join(ErrCantGetIP, err)
	}

	bytesData, err := gs.cache.GetWithCallback(ip, func() ([]byte, error) {
		rc, err := gs.loadViaAPI(ctx)
		if err != nil {
			return []byte{}, err
		}
		return io.ReadAll(rc)
	})
	if err != nil {
		return GeoData{}, errors.Join(ErrCantGetDataFromCache, err)
	}

	readCloser = io.NopCloser(bytes.NewReader(bytesData))

	defer readCloser.Close()

	var data GeoResponse
	err = json.NewDecoder(readCloser).Decode(&data)
	if err != nil {
		return GeoData{}, errors.Join(ErrCantDecodeGeoData, err)
	}

	return GeoData{Country: data.CountryCode, CountryCode: data.CountryCode}, nil
}

func (gs *GeoService) loadViaAPI(ctx context.Context) (io.ReadCloser, error) {
	url := "https://ipinfo.io/json"

	return gs.http.Get(ctx, url)
}

func (gs *GeoService) GetRc(ctx context.Context) (io.ReadCloser, error) {
	readCloser, err := gs.loadViaAPI(ctx)
	if err != nil {
		return nil, errors.Join(ErrCantGetGeoData, err)
	}

	return readCloser, nil
}
