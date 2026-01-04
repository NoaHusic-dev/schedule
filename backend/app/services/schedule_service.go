package services

import (
	"sync"
	"time"

	"github.com/go-raptor/raptor/v4"
	"github.com/h00s-go/tsd-timetable/internal"
)

type SchedulesService struct {
	raptor.Service

	edupageClient *internal.EdupageClient

	mu              sync.RWMutex
	CurrentSchedule internal.ScheduleResponse
}

func NewSchedulesService(c *raptor.Config) *SchedulesService {
	s := &SchedulesService{}

	s.OnInit(func() error {
		var err error
		s.edupageClient, err = internal.NewEdupageClient(
			c.AppConfig["edupage_subdomain"],
			c.AppConfig["edupage_username"],
			c.AppConfig["edupage_password"],
		)
		if err != nil {
			return err
		}

		err = s.UpdateCurrent()
		if err != nil {
			return err
		}

		go s.startScheduler()

		return nil
	})

	return s
}

func (s *SchedulesService) Current() internal.ScheduleResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.CurrentSchedule
}

func (s *SchedulesService) UpdateCurrent() error {
	schedule, err := s.edupageClient.GetFullSchedule(time.Now())
	if err != nil {
		s.Log.Error("Failed to update schedule", "error", err)
		return err
	}
	if len(schedule.Classes) > 0 {
		s.mu.Lock()
		s.CurrentSchedule = schedule
		s.mu.Unlock()
		s.Log.Info("Schedule updated", "classes", len(schedule.Classes))
	}
	return nil
}

func (s *SchedulesService) startScheduler() {
	for {
		now := time.Now()

		nextRun := time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, now.Location())
		if !nextRun.After(now) {
			nextRun = nextRun.AddDate(0, 0, 1)
		}

		time.Sleep(time.Until(nextRun))
		s.Log.Info("Updating schedule...")
		s.UpdateCurrent()
	}
}
