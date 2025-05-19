package weather_api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/EduardMikhrin/forecaster/internal"
	"github.com/EduardMikhrin/forecaster/internal/data"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

func (f *weatherApi) getWeather(ctx context.Context, city string) (*internal.WeatherPayload, error) {
	client := http.Client{Timeout: 5 * time.Second}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s?key=%s&q=%s", f.url,
		f.apiKey, city), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build weather request")
	}

	resp, err := client.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Wrap(err, "failed to fetch weather")
	}

	var apiResp weatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, errors.Wrap(err, "failed to decode weather response")
	}

	payload := &internal.WeatherPayload{
		City:        apiResp.Location.Name,
		Temperature: fmt.Sprintf("%.1f°C", apiResp.Current.TempC),
		Humidity:    fmt.Sprintf("%d%%", apiResp.Current.Humidity),
		WindSpeed:   fmt.Sprintf("%.1f kph", apiResp.Current.WindKph),
		Condition:   apiResp.Current.Condition.Text,
	}

	return payload, nil
}

func toEmailsList(subs []data.Subscription) []string {
	var emails []string
	for _, sub := range subs {
		emails = append(emails, sub.Email)
	}

	return emails
}
