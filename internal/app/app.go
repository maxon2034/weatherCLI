package app

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"
	"weatherCLI/internal/cache"
	"weatherCLI/internal/config"
	"weatherCLI/internal/domain"
	"weatherCLI/internal/provider"
	"weatherCLI/internal/provider/openmeteo"
	"weatherCLI/internal/ui"
)

type application struct {
	config   config.Config
	cache    *cache.TTLCache
	provider provider.WeatherProvider
	city     string
}

func Run() error {
	app := application{}
	var err error

	app.cache = cache.New()

	app.config, err = config.Load()
	if err != nil {
		return fmt.Errorf("Error in loading config: %w", err)
	}

	if app.config.DefaultCity == "" {
		fmt.Println("Input a default city")
		fmt.Scanln(&app.city)
	} else {
		app.city = app.config.DefaultCity
	}

	ctx := context.Background()
	client := openmeteo.NewClient()

	reader := bufio.NewReader(os.Stdin)
	today := domain.Today{}
	if todayCached, _, ok := app.cache.Get("today:" + app.city); ok {
		today = todayCached.(domain.Today)
		app.cache.Set("today:"+app.city, today, 5*time.Minute)
	} else {
		today, err = client.GetToday(ctx, app.city)
		if err != nil {
			return fmt.Errorf("Error in generating today's forecast with %s city: %w", app.city, err)
		}
		app.cache.Set("today:"+app.city, today, 5*time.Minute)
	}

	fmt.Println(ui.RenderToday(today))
	fmt.Println(ui.RenderMenu())
	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch strings.ToLower(input) {
		case "1":
			list := []domain.HourlyEntry{}
			if hourlyCached, _, ok := app.cache.Get("hourly:" + app.city); ok {
				list = hourlyCached.([]domain.HourlyEntry)
				app.cache.Set("hourly:"+app.city, list, 15*time.Minute)
			} else {
				list, err := client.GetHourly(ctx, app.city, 12)
				if err != nil {
					return fmt.Errorf("Error in generating hourly forecast: %w", err)
				}
				app.cache.Set("hourly:"+app.city, list, 15*time.Minute)
			}
			fmt.Println(ui.RenderHourly(list))
			fmt.Println(ui.RenderMenu())
		case "2":
			list := []domain.DailyEntry{}
			if dailyCached, _, ok := app.cache.Get("daily:" + app.city); ok {
				list = dailyCached.([]domain.DailyEntry)
				app.cache.Set("daily:"+app.city, list, 30*time.Minute)
			} else {
				list, err := client.GetDaily(ctx, app.city, 7)
				if err != nil {
					return fmt.Errorf("Error in generating daily forecast: %w", err)
				}
				app.cache.Set("daily:"+app.city, list, 30*time.Minute)
			}
			fmt.Println(ui.RenderDaily(list))
			fmt.Println(ui.RenderMenu())
		case "c":
			fmt.Println("Input your city")
			fmt.Print("> ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			app.city = input
			forecast := domain.Today{}
			if todayCached, _, ok := app.cache.Get("today:" + app.city); ok {
				forecast = todayCached.(domain.Today)
				app.cache.Set("today:"+app.city, forecast, 5*time.Minute)
			} else {
				forecast, err = client.GetToday(ctx, app.city)
				if err != nil {
					return fmt.Errorf("Error in generating today's forecast with %s city: %w", app.city, err)
				}
				app.cache.Set("today:"+app.city, forecast, 5*time.Minute)
			}
			app.config.DefaultCity = app.city
			err = config.Save(app.config)
			if err != nil {
				return fmt.Errorf("Error in saving new configuration: %w", err)
			}
			fmt.Println(ui.RenderToday(forecast))
			fmt.Println(ui.RenderMenu())
		case "r":
			today := domain.Today{}
			if todayCached, _, ok := app.cache.Get("today:" + app.city); ok {
				today = todayCached.(domain.Today)
				app.cache.Set("today:"+app.city, today, 5*time.Minute)
			} else {
				today, err = client.GetToday(ctx, app.city)
				if err != nil {
					return fmt.Errorf("Error in generating today's forecast with %s city: %w", app.city, err)
				}
				app.cache.Set("today:"+app.city, today, 5*time.Minute)
			}
			fmt.Println(ui.RenderToday(today))
			fmt.Println(ui.RenderMenu())
			continue
		case "q":
			return nil
		}
	}
}
