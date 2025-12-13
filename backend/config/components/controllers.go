package components

import (
	"github.com/go-raptor/raptor/v4"
	"github.com/h00s-go/tsd-timetable/app/controllers"
)

func Controllers() raptor.Controllers {
	return raptor.Controllers{
		&controllers.SchedulesController{},
		controllers.NewSPAController("public", "index.html"),
	}
}
