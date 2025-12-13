package services

import (
	"time"

	"github.com/go-raptor/raptor/v4"
	"github.com/h00s-go/tsd-timetable/internal"
)

type SchedulesService struct {
	raptor.Service

	edupageClient *internal.EdupageClient
	Schedule      internal.ScheduleResponse
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

		s.Schedule, err = s.edupageClient.GetFullSchedule(time.Now())

		return err
	})

	return s
}

func (s *SchedulesService) Current() internal.ScheduleResponse {
	return s.Schedule
}
