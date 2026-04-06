package components

import (
	"github.com/go-raptor/raptor/v4"
	"github.com/h00s-go/tsd-timetable/app/services"
)

func Services() raptor.Services {
	return raptor.Services{
		&services.SchedulesService{},
	}
}
