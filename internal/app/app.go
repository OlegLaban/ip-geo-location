package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/OlegLaban/geo-flag/internal/config"
	"github.com/OlegLaban/geo-flag/internal/domain"
	"github.com/OlegLaban/geo-flag/internal/repository/geodata"
	"github.com/OlegLaban/geo-flag/internal/usecase"
	"github.com/OlegLaban/geo-flag/pkg/cache"
	"github.com/OlegLaban/geo-flag/pkg/http/client"
	"github.com/OlegLaban/geo-flag/pkg/ipservice"
	"github.com/OlegLaban/geo-flag/pkg/locationdata"
	"github.com/OlegLaban/geo-flag/pkg/logger"
	"github.com/getlantern/systray"
)

type FlagServiceI interface {
	CountryCodeToEmoji(ctx context.Context, code string) string
	CountryCodeToPng(ctx context.Context, code string) ([]byte, error)
}

type TrayDataI interface {
	GetIcon() []byte
	GetTitle() string
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

func RunApp(config *config.Config) {
	logger := logger.SetupLogger(config.Env)
	httpClient := client.NewClient(logger)
	IPService := ipservice.NewIPService(httpClient, logger)
	cache := cache.NewCacheService(logger)
	geoPkgService := locationdata.NewGeoService(IPService, httpClient, cache, logger)
	geoService := geodata.NewGeoDataService(geoPkgService, logger)
	flagService := locationdata.NewFlagService(httpClient, cache, logger)
	app := NewApp(flagService, geoService)
	err := app.LoadData(logger, config)
	if err != nil {
		logger.Error("can`t load data", "err", err)
	}
	systray.Run(app.run(config, logger), app.exit(logger))
}

func (a *App) LoadData(logger *slog.Logger, config *config.Config) error {
	geoData, err := a.geoService.GetCountryData(a.Ctx)
	if err != nil {
		logger.Error("can`t load geodata", "err", err)
		return errors.Join(ErrLoadGeoData, err)
	}
	a.GeoData = geoData

	flagBytes, err := a.flagService.CountryCodeToPng(a.Ctx, a.GeoData.CountryCode)
	if err != nil {
		logger.Error("can`t load flag", "err", err)
		return errors.Join(ErrLoadFlag, err)
	}
	a.GeoData.Flag = flagBytes
	geoData.Flag = flagBytes

	trayData := usecase.NewTrayData(config, geoData)

	systray.SetIcon(trayData.GetIcon())
	systray.SetTitle(trayData.GetTitle())

	return nil
}

func (a *App) exit(logger *slog.Logger) func() {
	return func() {
		logger.Info("Exit...")
		a.Cancel()
	}
}

func (a *App) run(config *config.Config, logger *slog.Logger) func() {
	return func() {
		ticker := time.NewTicker(time.Duration(config.Round) * time.Second)
		logger.Info(fmt.Sprintf("ticker was setted success round - %d s", config.Round))

		go func() {
			defer ticker.Stop()
			for range ticker.C {
				select {
				case <-a.Ctx.Done():
					return
				default:
				}
				err := a.LoadData(logger, config)
				if err != nil {
					logger.Error("can`t load data", "err", err)
				}
			}
		}()

		mQuit := systray.AddMenuItem("Exit", "Close app")
		mSettings := systray.AddMenuItem("Settings", "Open settings window")
		go func() {
			<-mQuit.ClickedCh
			logger.Debug("exit button was clicked")
			systray.Quit()
		}()

		go func() {
			for range mSettings.ClickedCh {
				fmt.Println("Setting window open")
			}
		}()
	}
}
