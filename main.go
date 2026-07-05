package main

import (
	"context"
	"fmt"
	"weatherCLI/internal/provider/openmeteo"
)

func main() {
	ctx := context.Background()
	c := openmeteo.NewClient()
	today, err := c.GetTodayFormatted(ctx, "Minsk")
	if err != nil {
		fmt.Println("error in getting todays forecast: ", err)
	}
	fmt.Println(today)
}
