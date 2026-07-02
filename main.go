package main

import (
	"fmt"
	"weatherCLI/internal/config"
)

func main() {
	confSave := config.Config{DefaultCity: "Brest"}
	err := config.Save(confSave)
	if err != nil {
		fmt.Println(err)
	}
	confLoad, err := config.Load()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(confLoad)
}
