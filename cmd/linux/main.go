package main

import (
	"fmt"

	"github.com/OlegLaban/geo-flag/internal/app"
)

var configPath = "./configs/config.yaml"

func main() {
	config, err := app.LoadConfig(configPath)
	if err != nil {
		panic(fmt.Sprintf("can`t load config by path - %s", configPath))
	}
	app.RunApp(config, )
}
