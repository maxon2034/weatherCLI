package main

import (
	"context"
	"fmt"
	"weatherCLI/internal/provider/openmeteo"
)

func main() {
	ctx := context.Background()
	c := openmeteo.NewClient()
	day, err := c.GetDaily(ctx, "Minsk", 5)
	if err != nil {
		fmt.Println("error in getting todays forecast: ", err)
	}
	fmt.Println(day)
}
