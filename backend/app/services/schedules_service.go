package services

import (
	"sync"
	"time"
	_ "time/tzdata" // the distroless image has no zoneinfo to rely on

	"github.com/go-raptor/raptor/v4"
	"github.com/h00s-go/tsd-timetable/app/models"
)

const (
	refreshHour     = 5 // Edupage has the day's substitutions by then
	firstRetryDelay = 5 * time.Minute
	maxRetryDelay   = time.Hour
	defaultTimezone = "Europe/Zagreb"
)

type SchedulesService struct {
	raptor.Service

	source   scheduleSource
	location *time.Location

	refreshMu sync.Mutex // one refresh at a time: the Edupage client keeps session state
	mu        sync.RWMutex
	current   models.Schedule

	quit chan struct{}
}

func (s *SchedulesService) Setup() error {
	source, err := newScheduleSource(s.Config.AppConfig)
	if err != nil {
		return err
	}
	s.source = source

	s.location, err = loadLocation(s.Config.AppConfig["timezone"])
	if err != nil {
		s.Log.Warn("Invalid timezone, using the default", "timezone", s.Config.AppConfig["timezone"], "default", defaultTimezone)
		s.location, _ = loadLocation("")
	}

	// A failed first fetch must not keep the app from booting: the scheduler retries it.
	failures := 0
	if err := s.Refresh(); err != nil {
		failures = 1
	}

	s.quit = make(chan struct{})
	go s.schedule(failures)
	return nil
}

func (s *SchedulesService) Cleanup() error {
	close(s.quit)
	return nil
}

func (s *SchedulesService) Current() models.Schedule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.current
}

// Refresh fetches today's timetable. An empty result (a weekend, a holiday, or a parse that
// found nothing) keeps the previous schedule; its date tells the frontend it isn't today's.
func (s *SchedulesService) Refresh() error {
	s.refreshMu.Lock()
	defer s.refreshMu.Unlock()

	day := time.Now().In(s.location)
	schedule, err := s.source.Fetch(day)
	if err != nil {
		s.Log.Error("Failed to refresh the schedule", "error", err)
		return err
	}
	if len(schedule.Classes) == 0 {
		s.Log.Info("No lessons found, keeping the previous schedule", "day", day.Format(time.DateOnly))
		return nil
	}

	s.mu.Lock()
	s.current = models.Schedule{Date: day, Classes: schedule.Classes}
	s.mu.Unlock()
	s.Log.Info("Schedule updated", "day", day.Format(time.DateOnly), "classes", len(schedule.Classes))
	return nil
}

func (s *SchedulesService) schedule(failures int) {
	for {
		timer := time.NewTimer(time.Until(nextRefresh(time.Now().In(s.location), failures)))
		select {
		case <-s.quit:
			timer.Stop()
			return
		case <-timer.C:
		}

		if err := s.Refresh(); err != nil {
			failures++
		} else {
			failures = 0
		}
	}
}

// nextRefresh is 05:00 the next morning after a success. After consecutive failures it backs off
// from 5 minutes up to an hour, so a morning outage still gets today's timetable on the TV.
func nextRefresh(now time.Time, failures int) time.Time {
	if failures > 0 {
		delay := firstRetryDelay
		for i := 1; i < failures && delay < maxRetryDelay; i++ {
			delay *= 2
		}
		return now.Add(min(delay, maxRetryDelay))
	}

	next := time.Date(now.Year(), now.Month(), now.Day(), refreshHour, 0, 0, 0, now.Location())
	if !next.After(now) {
		next = next.AddDate(0, 0, 1)
	}
	return next
}

func loadLocation(name string) (*time.Location, error) {
	if name == "" {
		name = defaultTimezone
	}
	return time.LoadLocation(name)
}
