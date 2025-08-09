package usecase

import (
	"github.com/OlegLaban/geo-flag/internal/domain"
	"github.com/OlegLaban/geo-flag/pkg/locationdata"
)

func MapModelToDomain(geodata locationdata.GeoData) domain.GeoData {
	return domain.GeoData{
		CountryCode: geodata.CountryCode,
		CountryName: geodata.Country,
		Flag: geodata.Flag,
	}
}