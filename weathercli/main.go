package main

import (
	"fmt"
	"weatherCLI/internal/app"
)

func main() {
	err := app.Run()
	if err != nil {
		fmt.Println(err)
	}
}
