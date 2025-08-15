package geodata

import (
	"context"
	"log/slog"

	"github.com/OlegLaban/geo-flag/internal/domain"
	"github.com/OlegLaban/geo-flag/internal/usecase"
	"github.com/OlegLaban/geo-flag/pkg/locationdata"
)

type GeoService interface {
	GetCountryData(ctx context.Context) (locationdata.GeoData, error)
}

type GeoDataService struct {
	geoService GeoService
	logger     *slog.Logger
}

func NewGeoDataService(gs GeoService, logger *slog.Logger) *GeoDataService {
	return &GeoDataService{geoService: gs, logger: logger}
}

func (gd *GeoDataService) GetCountryData(ctx context.Context) (domain.GeoData, error) {
	geoData, err := gd.geoService.GetCountryData(ctx)
	if err != nil {
		gd.logger.Error("can`t get country data", "err", err)
		return domain.GeoData{}, err
	}
	gd.logger.Info("geo data was got success")
	return usecase.MapModelToDomain(geoData), nil
}
