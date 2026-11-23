package controllers

import (
	"encoding/json/v2"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/go-raptor/raptor/v4"
	"github.com/go-raptor/raptor/v4/errs"
	"github.com/h00s-go/tsd-timetable/app/models"
	"github.com/h00s-go/tsd-timetable/app/services"
)

const publicScheduleURL = "https://tsd-timetable.husak.me/api/v1/schedules/current"

var publicScheduleCache = struct {
	sync.RWMutex
	body      []byte
	fetchedAt time.Time
}{}

type SchedulesController struct {
	raptor.Controller

	Schedules *services.SchedulesService
}

func (c *SchedulesController) Current(ctx *raptor.Context) error {
	publicScheduleCache.RLock()
	cachedBody := append([]byte(nil), publicScheduleCache.body...)
	cachedAt := publicScheduleCache.fetchedAt
	publicScheduleCache.RUnlock()

	if len(cachedBody) > 0 && time.Since(cachedAt) < 60*time.Second {
		ctx.Response().Header().Set("Cache-Control", "public, max-age=60")
		_, _ = ctx.Response().Write(cachedBody)
		return nil
	}

	client := &http.Client{Timeout: 8 * time.Second}
	res, err := client.Get(publicScheduleURL)
	if err == nil {
		defer res.Body.Close()
		if res.StatusCode == http.StatusOK {
			body, readErr := io.ReadAll(io.LimitReader(res.Body, 4<<20))
			if readErr == nil {
				var schedule models.ScheduleResponse
				if err := json.Unmarshal(body, &schedule); err == nil {
					publicScheduleCache.Lock()
					publicScheduleCache.body = append([]byte(nil), body...)
					publicScheduleCache.fetchedAt = time.Now()
					publicScheduleCache.Unlock()
					ctx.Response().Header().Set("Cache-Control", "public, max-age=60")
					_, _ = ctx.Response().Write(body)
					return nil
				}
			}
		}
	}

	schedule := c.Schedules.Current()
	return ctx.Data(models.NewScheduleResponse(&schedule))
}

// UpdateCurrent refetches today's timetable on demand. The service logs the cause of a failure.
func (c *SchedulesController) UpdateCurrent(ctx *raptor.Context) error {
	if err := c.Schedules.Refresh(); err != nil {
		return errs.NewErrorBadGateway("Failed to refresh the schedule from Edupage")
	}
	return ctx.NoContent()
}
