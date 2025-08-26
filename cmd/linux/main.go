package main

import (
	"fmt"

	"github.com/OlegLaban/geo-flag/internal/app"
	"github.com/OlegLaban/geo-flag/internal/config"
)

func main() {
	config, err := config.LoadConfig(app.ConfigPath)
	if err != nil {
		panic(fmt.Sprintf("can`t load config by path - %s, err - %s", app.ConfigPath, err.Error()))
	}
	app.RunApp(config)
}
