package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/OlegLaban/geo-flag/internal/domain"
	"github.com/OlegLaban/geo-flag/internal/repository/geodata"
	"github.com/OlegLaban/geo-flag/pkg/cache"
	"github.com/OlegLaban/geo-flag/pkg/http/client"
	"github.com/OlegLaban/geo-flag/pkg/ipservice"
	"github.com/OlegLaban/geo-flag/pkg/locationdata"
	"github.com/getlantern/systray"
)

type FlagServiceI interface {
	CountryCodeToEmoji(code string) string
	CountryCodeToPng(code string) ([]byte, error)
}

type GeoServiceI interface {
	GetCountryData(ctx context.Context) (domain.GeoData, error)
}

type App struct {
	Ctx         context.Context
	Cancel      context.CancelFunc
	flagService FlagServiceI
	geoService  GeoServiceI
	GeoData     domain.GeoData
}

func NewApp(flagSerice FlagServiceI, geoService GeoServiceI) App {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	return App{Ctx: ctx, Cancel: cancel, flagService: flagSerice, geoService: geoService}
}

func RunApp() {
	httpClient := client.NewClient()
	IPService := ipservice.NewIPService(httpClient)
	cache := cache.NewCacheService()
	geoPkgService := locationdata.NewGeoService(IPService, httpClient, cache, true)
	geoService := geodata.NewGeoDataService(geoPkgService)
	flagService := locationdata.NewFlagService()
	app := NewApp(flagService, geoService)
	err := app.LoadData()
	if err != nil {
		log.Println(err)
	}
	systray.Run(app.run(), app.exit())
}

func (a *App) LoadData() error {
	geoData, err := a.geoService.GetCountryData(a.Ctx)
	if err != nil {
		return errors.Join(ErrLoadGeoData, err)
	}
	a.GeoData = geoData

	a.GeoData.Flag, err = a.flagService.CountryCodeToPng(a.GeoData.CountryCode)
	if err != nil {
		return errors.Join(ErrLoadFlag, err)
	}

	systray.SetIcon(a.GeoData.Flag)
	systray.SetTitle(a.GeoData.CountryName)

	return nil
}

func (a *App) exit() func() {
	return func() {
		a.Cancel()
	}
}

func (a *App) run() func() {
	return func() {
		ticker := time.NewTicker(10 * time.Second)

		go func() {
			defer ticker.Stop()
			for range ticker.C {
				select {
				case <-a.Ctx.Done():
					return
				default:
				}
				err := a.LoadData()
				if err != nil {
					log.Println("can`t load data")
				}
				systray.SetIcon(a.GeoData.Flag)
				systray.SetTitle(a.GeoData.CountryName)
			}
		}()

		mQuit := systray.AddMenuItem("Exit", "Close app")
		mSettings := systray.AddMenuItem("Settings", "Open settings window")
		go func() {
			<-mQuit.ClickedCh
			systray.Quit()
		}()

		go func() {
			for range mSettings.ClickedCh {
				fmt.Println("Setting window open")
			}
		}()
	}
}
