package main

import (
	"fmt"

	"github.com/OlegLaban/geo-flag/internal/app"
)

func main() {
	config, err := app.LoadConfig(app.ConfigPath)
	if err != nil {
		panic(fmt.Sprintf("can`t load config by path - %s", app.ConfigPath))
	}
	app.RunApp(config)
}
