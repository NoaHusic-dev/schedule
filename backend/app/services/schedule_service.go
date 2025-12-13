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

	mu       sync.RWMutex
	Schedule internal.ScheduleResponse
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

		initialSchedule, err := s.edupageClient.GetFullSchedule(time.Now())
		if err != nil {
			return err
		}
		s.Schedule = initialSchedule

		go s.startScheduler()

		return nil
	})

	return s
}

func (s *SchedulesService) Current() internal.ScheduleResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Schedule
}

func (s *SchedulesService) startScheduler() {
	for {
		now := time.Now()
		nextRun := time.Date(now.Year(), now.Month(), now.Day(), 6, 0, 0, 0, now.Location())

		if now.After(nextRun) {
			nextRun = nextRun.Add(24 * time.Hour)
		}

		timer := time.NewTimer(time.Until(nextRun))

		<-timer.C
		newSchedule, err := s.edupageClient.GetFullSchedule(time.Now())
		if err == nil {
			s.mu.Lock()
			s.Schedule = newSchedule
			s.mu.Unlock()
		} else {
			s.Log.Error("Failed to update schedule", "error", err)
		}
	}
}
