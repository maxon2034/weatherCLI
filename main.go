package main

import (
	"context"
	"fmt"
	"weatherCLI/internal/provider/openmeteo"
	"weatherCLI/internal/ui"
)

func main() {
	client := openmeteo.NewClient()
	ctx := context.Background()
	city := "Minsk"
	today, err := client.GetToday(ctx, city)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ui.RenderToday(today))
}
