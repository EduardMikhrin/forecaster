package forecaster

import "context"

type WeatherForecaster interface {
	Run(ctx context.Context)
}
