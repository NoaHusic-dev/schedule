package controllers

import (
	"github.com/go-raptor/raptor/v4"
	"github.com/go-raptor/raptor/v4/errs"
	"github.com/h00s-go/tsd-timetable/app/services"
)

type SchedulesController struct {
	raptor.Controller

	Schedules *services.SchedulesService
}

func (c *SchedulesController) Current(ctx *raptor.Context) error {
	return ctx.Data(c.Schedules.Current())
}

func (c *SchedulesController) UpdateCurrent(ctx *raptor.Context) error {
	if err := c.Schedules.UpdateCurrent(); err != nil {
		return errs.NewErrorInternal("Failed to update schedule", "error", err)
	}

	return ctx.NoContent()
}
