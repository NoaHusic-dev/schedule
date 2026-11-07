package services

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/go-raptor/raptor/v4"
	"github.com/h00s-go/tsd-timetable/app/models"
	"github.com/h00s/goopenmeteo"
)

const (
	daruvarLatitude  = 45.59
	daruvarLongitude = 17.23
	weatherCacheTTL  = 10 * time.Minute
)

type weatherSource interface {
	Fetch(context.Context) (models.WeatherResponse, error)
}

type WeatherService struct {
	raptor.Service

	source    weatherSource
	refreshMu sync.Mutex
	mu        sync.RWMutex
	current   models.WeatherResponse
	fetchedAt time.Time
}

func (s *WeatherService) Setup() error {
	s.source = &openMeteoSource{client: goopenmeteo.NewOpenMeteo()}
	return nil
}

func (s *WeatherService) Current(ctx context.Context) (models.WeatherResponse, error) {
	if current, fresh := s.cached(time.Now()); fresh {
		return current, nil
	}

	s.refreshMu.Lock()
	defer s.refreshMu.Unlock()

	if current, fresh := s.cached(time.Now()); fresh {
		return current, nil
	}

	current, err := s.source.Fetch(ctx)
	if err != nil {
		s.Log.Warn("Failed to refresh current weather", "error", err)
		if previous, ok := s.lastKnown(); ok {
			return previous, nil
		}
		return models.WeatherResponse{}, err
	}

	s.mu.Lock()
	s.current = current
	s.fetchedAt = time.Now()
	s.mu.Unlock()
	return current, nil
}

func (s *WeatherService) cached(now time.Time) (models.WeatherResponse, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.fetchedAt.IsZero() || now.Sub(s.fetchedAt) >= weatherCacheTTL {
		return models.WeatherResponse{}, false
	}
	return s.current, true
}

func (s *WeatherService) lastKnown() (models.WeatherResponse, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.current, !s.fetchedAt.IsZero()
}

type openMeteoSource struct {
	client *goopenmeteo.OpenMeteo
}

func (o *openMeteoSource) Fetch(ctx context.Context) (models.WeatherResponse, error) {
	forecast, err := o.client.ForecastContext(ctx, goopenmeteo.ForecastOptions{
		Latitude:        daruvarLatitude,
		Longitude:       daruvarLongitude,
		Current:         goopenmeteo.WeatherVariables{goopenmeteo.Temperature2M, goopenmeteo.WeatherCode},
		TemperatureUnit: goopenmeteo.TemperatureUnitCelsius,
		Timezone:        "Europe/Zagreb",
	})
	if err != nil {
		return models.WeatherResponse{}, fmt.Errorf("Open-Meteo request: %w", err)
	}

	temperature, ok := forecast.Current.Data[goopenmeteo.Temperature2M]
	if !ok || math.IsNaN(temperature) || math.IsInf(temperature, 0) {
		return models.WeatherResponse{}, errors.New("Open-Meteo response is missing a valid temperature")
	}
	weatherCode, ok := forecast.Current.Data[goopenmeteo.WeatherCode]
	if !ok || math.IsNaN(weatherCode) || math.IsInf(weatherCode, 0) || math.Trunc(weatherCode) != weatherCode {
		return models.WeatherResponse{}, errors.New("Open-Meteo response is missing a valid weather code")
	}

	return models.WeatherResponse{
		Temperature: int(math.Round(temperature)),
		WeatherCode: int(weatherCode),
	}, nil
}
