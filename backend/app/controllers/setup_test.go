package controllers_test

import (
	"os"
	"testing"

	"github.com/go-raptor/raptor/v4"
	"github.com/h00s-go/tsd-timetable/config"
	"github.com/h00s-go/tsd-timetable/config/components"
)

var app *raptor.Raptor

// .raptor.test.yaml selects the stub schedule source, so no test reaches Edupage.
func TestMain(m *testing.M) {
	app = raptor.NewTestApp(components.New(), config.Routes())
	os.Exit(m.Run())
}
