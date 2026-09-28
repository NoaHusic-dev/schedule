package models_test

import (
	"encoding/json/v2"
	"testing"
	"time"

	"github.com/h00s-go/tsd-timetable/app/models"
	"github.com/h00s-go/tsd-timetable/internal"
)

func sampleClasses(loc *time.Location) []internal.ClassSchedule {
	return []internal.ClassSchedule{{
		Class: "1 EL",
		Schedule: []internal.TimeSlot{{
			Start: time.Date(2026, 9, 28, 8, 0, 0, 0, loc),
			End:   time.Date(2026, 9, 28, 8, 45, 0, 0, loc),
			Lessons: []internal.LessonDetail{
				{
					Subject: internal.Subject{Name: "Matematika", Short: "MAT"},
					Teacher: internal.Teacher{Firstname: "Ana", Lastname: "Anić", Short: "A. Anić"},
					Room:    "ZR 2",
					Groups:  []string{"1. grupa"},
				},
				{
					Subject: internal.Subject{Name: "Fizika", Short: "FIZ"},
					Teacher: internal.Teacher{Firstname: "Ivo", Lastname: "Ivić", Short: "I. Ivić"},
					Room:    "12",
				},
			},
		}},
	}}
}

// The response must serialize exactly like the parser's own structs did before the DTO layer,
// plus the new date, so a frontend still running on the TV keeps working.
func TestNewScheduleResponseKeepsTheParserShape(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Zagreb")
	if err != nil {
		t.Fatal(err)
	}
	classes := sampleClasses(loc)

	legacy, err := json.Marshal(internal.ScheduleResponse{Classes: classes})
	if err != nil {
		t.Fatal(err)
	}
	got, err := json.Marshal(models.NewScheduleResponse(&models.Schedule{
		Date:    time.Date(2026, 9, 28, 5, 0, 0, 0, loc),
		Classes: classes,
	}))
	if err != nil {
		t.Fatal(err)
	}

	want := `{"date":"2026-09-28",` + string(legacy[1:])
	if string(got) != want {
		t.Errorf("response JSON differs from the parser's\n got: %s\nwant: %s", got, want)
	}
}

func TestNewScheduleResponseBeforeFirstFetch(t *testing.T) {
	got, err := json.Marshal(models.NewScheduleResponse(&models.Schedule{}))
	if err != nil {
		t.Fatal(err)
	}
	if want := `{"date":"","classes":[]}`; string(got) != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
