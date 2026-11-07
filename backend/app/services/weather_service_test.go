package services

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-raptor/raptor/v4"
	"github.com/h00s-go/tsd-timetable/app/models"
	"github.com/h00s/goopenmeteo"
)

type fakeWeatherSource struct {
	current models.WeatherResponse
	err     error
	calls   int
}

func (f *fakeWeatherSource) Fetch(context.Context) (models.WeatherResponse, error) {
	f.calls++
	return f.current, f.err
}

func newTestWeatherService(t *testing.T, source weatherSource) *WeatherService {
	t.Helper()
	s := &WeatherService{source: source}
	if err := s.Init(raptor.NewTestResources()); err != nil {
		t.Fatal(err)
	}
	return s
}

func TestWeatherCurrentCachesSuccessfulFetch(t *testing.T) {
	source := &fakeWeatherSource{current: models.WeatherResponse{Temperature: 13, WeatherCode: 2}}
	s := newTestWeatherService(t, source)

	for range 2 {
		got, err := s.Current(context.Background())
		if err != nil {
			t.Fatalf("Current: %v", err)
		}
		if got != source.current {
			t.Errorf("Current = %+v, want %+v", got, source.current)
		}
	}
	if source.calls != 1 {
		t.Errorf("Fetch called %d times, want 1", source.calls)
	}
}

func TestWeatherCurrentReturnsLastKnownDataWhenRefreshFails(t *testing.T) {
	source := &fakeWeatherSource{current: models.WeatherResponse{Temperature: 13, WeatherCode: 2}}
	s := newTestWeatherService(t, source)
	s.current = source.current
	s.fetchedAt = time.Now().Add(-weatherCacheTTL)
	source.err = errors.New("Open-Meteo is unavailable")

	got, err := s.Current(context.Background())
	if err != nil {
		t.Fatalf("Current: %v", err)
	}
	if got != source.current {
		t.Errorf("Current = %+v, want cached weather %+v", got, source.current)
	}
}

func TestWeatherCurrentReturnsErrorWithoutCachedData(t *testing.T) {
	source := &fakeWeatherSource{err: errors.New("Open-Meteo is unavailable")}
	s := newTestWeatherService(t, source)

	if _, err := s.Current(context.Background()); !errors.Is(err, source.err) {
		t.Fatalf("Current error = %v, want %v", err, source.err)
	}
}

func TestOpenMeteoSourceFetchesDaruvarCurrentWeather(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/forecast" {
			t.Errorf("request path = %q, want /forecast", r.URL.Path)
		}
		query := r.URL.Query()
		if query.Get("latitude") != "45.59" || query.Get("longitude") != "17.23" {
			t.Errorf("coordinates = %q, %q, want Daruvar", query.Get("latitude"), query.Get("longitude"))
		}
		if query.Get("timezone") != "Europe/Zagreb" || query.Get("temperature_unit") != "celsius" {
			t.Errorf("timezone/unit = %q/%q, want Europe/Zagreb/celsius", query.Get("timezone"), query.Get("temperature_unit"))
		}
		_, _ = w.Write([]byte(`{"current":{"time":"2026-11-07T10:00","interval":900,"temperature_2m":13.4,"weather_code":2}}`))
	}))
	defer server.Close()

	client := goopenmeteo.NewOpenMeteo()
	client.BaseURL = server.URL
	source := &openMeteoSource{client: client}

	got, err := source.Fetch(context.Background())
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	want := models.WeatherResponse{Temperature: 13, WeatherCode: 2}
	if got != want {
		t.Errorf("Fetch = %+v, want %+v", got, want)
	}
}
