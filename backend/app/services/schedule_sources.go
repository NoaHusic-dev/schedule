package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/h00s-go/tsd-timetable/internal"
)

// scheduleSource produces one day's timetable. The stub keeps tests and local development away
// from Edupage (app.schedule_source: "stub").
type scheduleSource interface {
	Fetch(day time.Time) (internal.ScheduleResponse, error)
}

func newScheduleSource(cfg map[string]string) (scheduleSource, error) {
	switch cfg["schedule_source"] {
	case "stub":
		return stubSource{}, nil
	case "", "edupage":
	default:
		return nil, fmt.Errorf("unknown schedule_source %q, use \"edupage\" or \"stub\"", cfg["schedule_source"])
	}

	source := &edupageSource{subdomain: cfg["edupage_subdomain"], username: cfg["edupage_username"], password: cfg["edupage_password"]}
	if source.subdomain == "" || source.username == "" || source.password == "" {
		return nil, errors.New("edupage_subdomain, edupage_username and edupage_password must be set")
	}
	return source, nil
}

type edupageSource struct {
	subdomain, username, password string
}

// Fetch logs in with a fresh client every time. The session doesn't survive a day, and logging in
// again while it is still alive fails: the login page has no CSRF token for a signed-in session.
func (e *edupageSource) Fetch(day time.Time) (internal.ScheduleResponse, error) {
	client, err := internal.NewEdupageClient(e.subdomain, e.username, e.password)
	if err != nil {
		return internal.ScheduleResponse{}, err
	}
	if _, err := client.Login(); err != nil {
		return internal.ScheduleResponse{}, fmt.Errorf("edupage login: %w", err)
	}
	return client.GetFullSchedule(day)
}

type stubSource struct{}

func (stubSource) Fetch(day time.Time) (internal.ScheduleResponse, error) {
	at := func(hour, minute int) time.Time {
		return time.Date(day.Year(), day.Month(), day.Day(), hour, minute, 0, 0, day.Location())
	}
	lesson := func(subject, short, teacher, room string, groups ...string) internal.LessonDetail {
		return internal.LessonDetail{
			Subject: internal.Subject{Name: subject, Short: short},
			Teacher: internal.Teacher{Short: teacher},
			Room:    room,
			Groups:  groups,
		}
	}
	slot := func(start, end time.Time, lessons ...internal.LessonDetail) internal.TimeSlot {
		return internal.TimeSlot{Start: start, End: end, Lessons: lessons}
	}

	return internal.ScheduleResponse{Classes: []internal.ClassSchedule{
		{Class: "1 EL", Schedule: []internal.TimeSlot{
			slot(at(8, 0), at(8, 45), lesson("Hrvatski jezik", "HJ", "A. Anić", "12")),
			slot(at(8, 50), at(9, 35), lesson("Matematika", "MAT", "I. Ivić", "14")),
			slot(at(9, 40), at(10, 25),
				lesson("Osnove elektrotehnike", "OE", "M. Marić", "ZR 2", "1. grupa"),
				lesson("Računalstvo", "RAČ", "P. Perić", "INF 1", "2. grupa")),
			slot(at(10, 40), at(12, 15), lesson("Praktična nastava", "PN", "K. Kovač", "RAD 1")),
		}},
		{Class: "2 EL", Schedule: []internal.TimeSlot{
			slot(at(8, 0), at(9, 35), lesson("Elektronika", "ELE", "M. Marić", "ZR 1")),
			slot(at(9, 40), at(10, 25), lesson("Engleski jezik", "EJ", "S. Horvat", "8")),
			slot(at(10, 40), at(11, 25), lesson("Tjelesna i zdravstvena kultura", "TZK", "D. Babić", "DVORANA")),
		}},
		{Class: "1 MT", Schedule: []internal.TimeSlot{
			slot(at(11, 30), at(12, 15), lesson("Tehničko crtanje", "TC", "B. Novak", "5")),
			slot(at(12, 20), at(13, 55), lesson("Mehatronika", "MEH", "T. Jurić", "MEH 1")),
		}},
	}}, nil
}
