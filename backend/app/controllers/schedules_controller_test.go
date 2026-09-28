package controllers_test

import (
	"encoding/json/v2"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/go-raptor/raptor/v4"
	"github.com/h00s-go/tsd-timetable/app/models"
)

const (
	currentPath = "/api/v1/schedules/current"
	updatePath  = "/api/v1/schedules/current/update"
	testToken   = "test-token" // .raptor.test.yaml
)

func TestCurrentReturnsTodaysSchedule(t *testing.T) {
	rec := app.TestGet(currentPath)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body %s", rec.Code, rec.Body)
	}

	var got models.ScheduleResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v; body %s", err, rec.Body)
	}
	zagreb, _ := time.LoadLocation("Europe/Zagreb")
	if today := time.Now().In(zagreb).Format(time.DateOnly); got.Date != today {
		t.Errorf("date = %q, want %q", got.Date, today)
	}
	if len(got.Classes) == 0 || got.Classes[0].Class != "1 EL" {
		t.Errorf("classes = %+v, want the stub's classes", got.Classes)
	}
	if !strings.Contains(rec.Body.String(), `"firstName"`) {
		t.Errorf("teacher keys must stay camelCase for the frontend; body %s", rec.Body)
	}
}

func TestUpdateRequiresTheToken(t *testing.T) {
	tests := []struct {
		name string
		opts []raptor.TestRequestOption
		want int
	}{
		{"no token", nil, http.StatusUnauthorized},
		{"wrong token", []raptor.TestRequestOption{raptor.WithHeader("Authorization", "Bearer guess")}, http.StatusUnauthorized},
		{"right token", []raptor.TestRequestOption{raptor.WithHeader("Authorization", "Bearer "+testToken)}, http.StatusNoContent},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := app.TestPost(updatePath, nil, tt.opts...)
			if rec.Code != tt.want {
				t.Errorf("status = %d, want %d; body %s", rec.Code, tt.want, rec.Body)
			}
		})
	}
}

func TestCrossOriginWriteIsRejected(t *testing.T) {
	rec := app.TestPost(updatePath, nil,
		raptor.WithHeader("Authorization", "Bearer "+testToken),
		raptor.WithHeader("Sec-Fetch-Site", "cross-site"),
		raptor.WithHeader("Origin", "https://evil.example"))
	if rec.Code != http.StatusForbidden {
		t.Errorf("status = %d, want 403, proving the csrf middleware is active; body %s", rec.Code, rec.Body)
	}
}

// The SPA catch-all on / would otherwise answer unknown API paths with an empty 404.
func TestUnknownAPIPathIsAJSON404(t *testing.T) {
	rec := app.TestGet("/api/v1/nope")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body %s", rec.Code, rec.Body)
	}
	var body struct {
		Code int `json:"code"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil || body.Code != http.StatusNotFound {
		t.Errorf("body = %q, want a JSON error envelope with code 404", rec.Body)
	}
}
