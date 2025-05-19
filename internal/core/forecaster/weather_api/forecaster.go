package weather_api

import (
	"context"
	"github.com/EduardMikhrin/forecaster/internal/core/forecaster"
	"github.com/EduardMikhrin/forecaster/internal/core/mailer"
	"github.com/EduardMikhrin/forecaster/internal/data"
	"gitlab.com/distributed_lab/logan/v3"
	"sync"
	"time"
)

const name = "https://api.weatherapi.com/v1/current.json"

type weatherApi struct {
	apiKey          string
	url             string
	pollingInterval time.Duration

	db     data.MasterQ
	mailer mailer.Mailer

	log *logan.Entry
}

func (f *weatherApi) Run(ctx context.Context) {
	ticker := time.NewTicker(f.pollingInterval)
	wg := &sync.WaitGroup{}
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			f.log.Debug("service stopped: exiting forecaster")
			return

		case <-ticker.C:
			cities, err := f.db.CitiesQ().GetAll()
			if err != nil {
				f.log.WithError(err).Error("error getting cities")
				continue
			}

			if err == nil && cities == nil {
				f.log.Debug("no cities available for forecasting")
			}

			for _, city := range cities {
				wg.Add(1)
				go func(city data.City) {
					defer wg.Done()

					subs, err := f.db.New().SubscriptionQ().FilterByCityId(int(city.Id)).GetAll()
					if err != nil {
						f.log.WithError(err).Error("error getting subscriptions")
						return
					}
					if subs == nil {
						f.log.Debug("no emails available for forecasting")
						return
					}

					payload, err := f.getWeather(ctx, city.Name)
					if err != nil {
						f.log.WithError(err).Error("error getting weather")
						return
					}

					if err := f.mailer.SendInfoEmail(toEmailsList(subs), payload); err != nil {
						f.log.WithError(err).Error("error sending email")
						return
					}

				}(city)
			}

			wg.Wait()
		}
	}
}

func NewForecaster(apiKey, url string, interval int) forecaster.WeatherForecaster {
	return &weatherApi{
		apiKey:          apiKey,
		url:             url,
		pollingInterval: time.Duration(interval) * time.Second,
	}
}
