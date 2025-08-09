package geodata

import (
	"context"
	"log"

	"github.com/OlegLaban/geo-flag/internal/domain"
	"github.com/OlegLaban/geo-flag/internal/usecase"
	"github.com/OlegLaban/geo-flag/pkg/locationdata"
)

type GeoService interface {
	GetCountryData(ctx context.Context) (locationdata.GeoData, error)	
}

type GeoDataService struct {
	geoService GeoService
}

func NewGeoDataService(gs GeoService) *GeoDataService {
	return &GeoDataService{geoService: gs}
}

func (gd *GeoDataService) GetCountryData(ctx context.Context) (domain.GeoData, error) {
	geoData, err := gd.geoService.GetCountryData(ctx)
	if err != nil {
		log.Println("Can`t get country data", err)
		return domain.GeoData{}, err
	}

	return usecase.MapModelToDomain(geoData), nil
} 