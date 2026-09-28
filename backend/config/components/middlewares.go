package components

import (
	"github.com/go-raptor/middlewares/csrf"
	"github.com/go-raptor/middlewares/logger"
	"github.com/go-raptor/raptor/v4"
	"github.com/h00s-go/tsd-timetable/app/middlewares"
)

func Middlewares() raptor.Middlewares {
	return raptor.Middlewares{
		raptor.Use(&logger.LoggerMiddleware{}),
		raptor.Use(&csrf.CSRFMiddleware{}),
		raptor.UseOnly(&middlewares.UpdateTokenMiddleware{}, "Schedules.UpdateCurrent"),
	}
}
