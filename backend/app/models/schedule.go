package models

import (
	"time"

	"github.com/h00s-go/tsd-timetable/internal"
)

// Schedule is one day's timetable for every class, as parsed from Edupage.
type Schedule struct {
	Date    time.Time // the day the timetable is for; zero until the first successful fetch
	Classes []internal.ClassSchedule
}

type (
	ScheduleResponse struct {
		Date    string                  `json:"date"` // YYYY-MM-DD, or "" before the first successful fetch
		Classes []ClassScheduleResponse `json:"classes"`
	}

	ClassScheduleResponse struct {
		Class    string             `json:"class"`
		Schedule []TimeSlotResponse `json:"schedule"`
	}

	TimeSlotResponse struct {
		Start   time.Time        `json:"start"`
		End     time.Time        `json:"end"`
		Lessons []LessonResponse `json:"lessons"`
	}

	LessonResponse struct {
		Subject SubjectResponse `json:"subject"`
		Teacher TeacherResponse `json:"teacher"`
		Room    string          `json:"room"`
		Groups  []string        `json:"groups"`
	}

	SubjectResponse struct {
		Name  string `json:"name"`
		Short string `json:"short"`
	}

	TeacherResponse struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Short     string `json:"short"`
	}
)

func NewScheduleResponse(s *Schedule) ScheduleResponse {
	res := ScheduleResponse{Classes: make([]ClassScheduleResponse, len(s.Classes))}
	if !s.Date.IsZero() {
		res.Date = s.Date.Format(time.DateOnly)
	}
	for i := range s.Classes {
		res.Classes[i] = newClassScheduleResponse(&s.Classes[i])
	}
	return res
}

func newClassScheduleResponse(c *internal.ClassSchedule) ClassScheduleResponse {
	res := ClassScheduleResponse{Class: c.Class, Schedule: make([]TimeSlotResponse, len(c.Schedule))}
	for i := range c.Schedule {
		res.Schedule[i] = newTimeSlotResponse(&c.Schedule[i])
	}
	return res
}

func newTimeSlotResponse(t *internal.TimeSlot) TimeSlotResponse {
	res := TimeSlotResponse{Start: t.Start, End: t.End, Lessons: make([]LessonResponse, len(t.Lessons))}
	for i := range t.Lessons {
		res.Lessons[i] = newLessonResponse(&t.Lessons[i])
	}
	return res
}

func newLessonResponse(l *internal.LessonDetail) LessonResponse {
	return LessonResponse{
		Subject: SubjectResponse{Name: l.Subject.Name, Short: l.Subject.Short},
		Teacher: TeacherResponse{FirstName: l.Teacher.Firstname, LastName: l.Teacher.Lastname, Short: l.Teacher.Short},
		Room:    l.Room,
		Groups:  l.Groups,
	}
}
