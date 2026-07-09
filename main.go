package main

import (
	"context"
	"fmt"
	"weatherCLI/internal/config"
	"weatherCLI/internal/provider/openmeteo"
	"weatherCLI/internal/ui"
)

func main() {
	c := openmeteo.NewClient()
	ctx := context.Background()
	conf, err := config.Load()
	if err != nil {
		fmt.Println(err)
		return
	}
	list, err := c.GetDaily(ctx, conf.DefaultCity, 12)
	if err != nil {
		fmt.Println(err)
		return
	}
	str := ui.RenderDaily(list)
	fmt.Println(str)
}
