package services

import (
	"errors"
	"testing"
	"time"

	"github.com/go-raptor/raptor/v4"
	"github.com/h00s-go/tsd-timetable/internal"
)

type fakeSource struct {
	result internal.ScheduleResponse
	err    error
}

func (f *fakeSource) Fetch(time.Time) (internal.ScheduleResponse, error) {
	return f.result, f.err
}

func newTestService(t *testing.T, source scheduleSource) *SchedulesService {
	t.Helper()
	s := &SchedulesService{source: source, location: time.UTC}
	if err := s.Init(raptor.NewTestResources()); err != nil {
		t.Fatal(err)
	}
	return s
}

func oneClass(name string) internal.ScheduleResponse {
	return internal.ScheduleResponse{Classes: []internal.ClassSchedule{{Class: name}}}
}

func TestRefreshStoresTodaysSchedule(t *testing.T) {
	s := newTestService(t, &fakeSource{result: oneClass("1 EL")})

	if err := s.Refresh(); err != nil {
		t.Fatalf("Refresh: %v", err)
	}

	got := s.Current()
	if today := time.Now().In(time.UTC).Format(time.DateOnly); got.Date.Format(time.DateOnly) != today {
		t.Errorf("Date = %v, want %s", got.Date, today)
	}
	if len(got.Classes) != 1 || got.Classes[0].Class != "1 EL" {
		t.Errorf("Classes = %+v, want the fetched class", got.Classes)
	}
}

// Weekends and holidays come back empty; so does a parse that silently failed. Neither may wipe
// the last good schedule: the frontend tells stale data apart by its date.
func TestRefreshKeepsThePreviousScheduleWhenThereAreNoLessons(t *testing.T) {
	source := &fakeSource{result: oneClass("1 EL")}
	s := newTestService(t, source)
	if err := s.Refresh(); err != nil {
		t.Fatalf("first Refresh: %v", err)
	}
	before := s.Current()

	source.result = internal.ScheduleResponse{}
	if err := s.Refresh(); err != nil {
		t.Fatalf("empty Refresh: %v", err)
	}

	if got := s.Current(); len(got.Classes) != 1 || !got.Date.Equal(before.Date) {
		t.Errorf("Current() = %+v, want the previous schedule %+v", got, before)
	}
}

func TestRefreshReturnsTheFetchErrorAndKeepsTheSchedule(t *testing.T) {
	source := &fakeSource{result: oneClass("1 EL")}
	s := newTestService(t, source)
	if err := s.Refresh(); err != nil {
		t.Fatalf("first Refresh: %v", err)
	}

	source.err = errors.New("edupage is down")
	if err := s.Refresh(); !errors.Is(err, source.err) {
		t.Fatalf("Refresh error = %v, want %v", err, source.err)
	}
	if got := s.Current(); len(got.Classes) != 1 {
		t.Errorf("Current() lost the schedule after a failed refresh: %+v", got)
	}
}

func TestNextRefresh(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Zagreb")
	if err != nil {
		t.Fatal(err)
	}
	at := func(month time.Month, day, hour, minute int) time.Time {
		return time.Date(2026, month, day, hour, minute, 0, 0, loc)
	}

	tests := []struct {
		name     string
		now      time.Time
		failures int
		want     time.Time
	}{
		{"before 05:00 refreshes this morning", at(9, 28, 4, 30), 0, at(9, 28, 5, 0)},
		{"at 05:00 refreshes tomorrow", at(9, 28, 5, 0), 0, at(9, 29, 5, 0)},
		{"after 05:00 refreshes tomorrow", at(9, 28, 13, 0), 0, at(9, 29, 5, 0)},
		{"across the end of summer time", at(10, 24, 6, 0), 0, at(10, 25, 5, 0)},
		{"first failure retries in 5 minutes", at(9, 28, 5, 0), 1, at(9, 28, 5, 5)},
		{"second failure retries in 10 minutes", at(9, 28, 5, 0), 2, at(9, 28, 5, 10)},
		{"third failure retries in 20 minutes", at(9, 28, 5, 0), 3, at(9, 28, 5, 20)},
		{"fourth failure retries in 40 minutes", at(9, 28, 5, 0), 4, at(9, 28, 5, 40)},
		{"fifth failure is capped at an hour", at(9, 28, 5, 0), 5, at(9, 28, 6, 0)},
		{"a long outage stays at an hour", at(9, 28, 5, 0), 100, at(9, 28, 6, 0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nextRefresh(tt.now, tt.failures); !got.Equal(tt.want) {
				t.Errorf("nextRefresh(%v, %d) = %v, want %v", tt.now, tt.failures, got, tt.want)
			}
		})
	}
}

func TestNewScheduleSource(t *testing.T) {
	creds := map[string]string{"edupage_subdomain": "school", "edupage_username": "user", "edupage_password": "secret"}

	if source, err := newScheduleSource(map[string]string{"schedule_source": "stub"}); err != nil {
		t.Errorf("stub: %v", err)
	} else if _, ok := source.(stubSource); !ok {
		t.Errorf("stub: got %T, want stubSource", source)
	}

	if source, err := newScheduleSource(creds); err != nil {
		t.Errorf("edupage: %v", err)
	} else if _, ok := source.(*edupageSource); !ok {
		t.Errorf("edupage: got %T, want *edupageSource", source)
	}

	if _, err := newScheduleSource(map[string]string{"edupage_subdomain": "school"}); err == nil {
		t.Error("missing credentials: want an error, so a misconfigured deploy fails at startup")
	}

	if _, err := newScheduleSource(map[string]string{"schedule_source": "stubb"}); err == nil {
		t.Error("unknown schedule_source: want an error rather than a silent fallback to Edupage")
	}
}

func TestLoadLocation(t *testing.T) {
	loc, err := loadLocation("")
	if err != nil || loc.String() != "Europe/Zagreb" {
		t.Errorf(`loadLocation("") = %v, %v; want Europe/Zagreb`, loc, err)
	}
	if _, err := loadLocation("Nowhere/Atlantis"); err == nil {
		t.Error("an unknown zone must be an error")
	}
}

// The frontend only shows a schedule whose date is today, so the stub must produce lessons on the
// day it was asked for.
func TestStubSourceProducesLessonsForTheRequestedDay(t *testing.T) {
	loc, _ := loadLocation("")
	day := time.Date(2026, 9, 28, 5, 0, 0, 0, loc)

	schedule, err := stubSource{}.Fetch(day)
	if err != nil {
		t.Fatal(err)
	}
	if len(schedule.Classes) == 0 {
		t.Fatal("stub returned no classes")
	}
	for _, class := range schedule.Classes {
		for _, slot := range class.Schedule {
			if slot.Start.Format(time.DateOnly) != "2026-09-28" || slot.Start.Location() != loc || len(slot.Lessons) == 0 {
				t.Errorf("%s: slot %v is not a lesson on the requested day", class.Class, slot.Start)
			}
		}
	}
}
