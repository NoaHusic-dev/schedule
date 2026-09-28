package controllers

import (
	"github.com/go-raptor/raptor/v4"
	"github.com/go-raptor/raptor/v4/errs"
	"github.com/h00s-go/tsd-timetable/app/models"
	"github.com/h00s-go/tsd-timetable/app/services"
)

type SchedulesController struct {
	raptor.Controller

	Schedules *services.SchedulesService
}

func (c *SchedulesController) Current(ctx *raptor.Context) error {
	schedule := c.Schedules.Current()
	return ctx.Data(models.NewScheduleResponse(&schedule))
}

// UpdateCurrent refetches today's timetable on demand. The service logs the cause of a failure.
func (c *SchedulesController) UpdateCurrent(ctx *raptor.Context) error {
	if err := c.Schedules.Refresh(); err != nil {
		return errs.NewErrorBadGateway("Failed to refresh the schedule from Edupage")
	}
	return ctx.NoContent()
}
