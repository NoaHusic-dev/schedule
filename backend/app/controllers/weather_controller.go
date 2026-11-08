package controllers

import (
	"github.com/go-raptor/raptor/v4"
	"github.com/go-raptor/raptor/v4/errs"
	"github.com/h00s-go/tsd-timetable/app/services"
)

type WeatherController struct {
	raptor.Controller

	Weather *services.WeatherService
}

func (c *WeatherController) Current(ctx *raptor.Context) error {
	weather, err := c.Weather.Current(ctx.Request().Context())
	if err != nil {
		return errs.NewErrorBadGateway("Current weather is unavailable")
	}
	return ctx.Data(weather)
}
