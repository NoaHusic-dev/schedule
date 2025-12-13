package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ScheduleResponse struct {
	Classes []ClassSchedule `json:"classes"`
}

type ClassSchedule struct {
	ClassName string     `json:"className"`
	Schedule  []TimeSlot `json:"schedule"`
}

type TimeSlot struct {
	Start   time.Time      `json:"start"`
	End     time.Time      `json:"end"`
	Lessons []LessonDetail `json:"lessons"`
}

type LessonDetail struct {
	Subject Subject `json:"subject"`
	Teacher Teacher `json:"teacher"`
	Room    string  `json:"room"`
}

type Subject struct {
	Name  string `json:"name"`
	Short string `json:"short"`
}

type Teacher struct {
	Firstname string `json:"firstName"`
	Lastname  string `json:"lastName"`
	Short     string `json:"short"`
}

// --- Client Definition ---

type EdupageClient struct {
	Client     *http.Client
	Subdomain  string
	Username   string
	Password   string
	GsecHash   string
	SchoolYear int
	Dbi        DBI
}

// --- Internal Data ---

type DBI struct {
	Subjects   map[string]DBItem `json:"subjects"`
	Teachers   map[string]DBItem `json:"teachers"`
	Classrooms map[string]DBItem `json:"classrooms"`
	Classes    map[string]DBItem `json:"classes"`
	Periods    map[string]DBItem `json:"periods"`
}

type DBItem struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Short     string `json:"short"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	StartTime string `json:"starttime"`
	EndTime   string `json:"endtime"`
}

type rawTimetableResponse struct {
	R struct {
		TTItems []rawLesson `json:"ttitems"`
		Error   string      `json:"error"`
	} `json:"r"`
}

type rawLesson struct {
	SubjectID    string   `json:"subjectid"`
	TeacherIDs   []string `json:"teacherids"`
	ClassroomIDs []string `json:"classroomids"`
	Period       string   `json:"period"`
	StartTime    string   `json:"starttime"`
	EndTime      string   `json:"endtime"`
	Type         string   `json:"type"`
}

func NewEdupageClient(subdomain, username, password string) (*EdupageClient, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	return &EdupageClient{
		Client: &http.Client{
			Jar:     jar,
			Timeout: 20 * time.Second,
		},
		Subdomain: subdomain,
		Username:  username,
		Password:  password,
	}, nil
}

func (e *EdupageClient) loginInternal() error {
	loginURL := fmt.Sprintf("https://%s.edupage.org/login/edubarLogin.php", e.Subdomain)
	initURL := fmt.Sprintf("https://%s.edupage.org/login/?cmd=MainLogin", e.Subdomain)

	// 1. Fetch CSRF
	resp, err := e.Client.Get(initURL)
	if err != nil {
		return fmt.Errorf("init request failed: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	csrfMatch := regexp.MustCompile(`"csrftoken":"([^"]+)"`).FindSubmatch(body)
	if len(csrfMatch) < 2 {
		return fmt.Errorf("CSRF token not found")
	}

	// 2. Post Login
	form := url.Values{}
	form.Set("csrfauth", string(csrfMatch[1]))
	form.Set("username", e.Username)
	form.Set("password", e.Password)

	resp, err = e.Client.PostForm(loginURL, form)
	if err != nil {
		return fmt.Errorf("login post failed: %w", err)
	}
	defer resp.Body.Close()

	// 3. Extract User Data & DBI
	userBody, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(userBody), "userhome(") {
		respUser, err := e.Client.Get(fmt.Sprintf("https://%s.edupage.org/user", e.Subdomain))
		if err != nil {
			return fmt.Errorf("user dashboard fetch failed: %w", err)
		}
		defer respUser.Body.Close()
		userBody, _ = io.ReadAll(respUser.Body)
	}

	return e.parseLoginData(string(userBody))
}

func (e *EdupageClient) parseLoginData(html string) error {
	hashMatch := regexp.MustCompile(`ASC\.gsechash="([^"]+)"`).FindStringSubmatch(html)
	if len(hashMatch) >= 2 {
		e.GsecHash = hashMatch[1]
	}

	parts := strings.SplitN(html, "userhome(", 2)
	if len(parts) < 2 {
		return fmt.Errorf("userhome data not found")
	}
	jsonRaw := parts[1]

	segments := strings.Split(jsonRaw, ");")
	if len(segments) >= 3 {
		jsonRaw = strings.Join(segments[:len(segments)-2], ");")
	} else {
		lastIdx := strings.LastIndex(jsonRaw, ");")
		if lastIdx != -1 {
			jsonRaw = jsonRaw[:lastIdx]
		}
	}

	jsonRaw = strings.ReplaceAll(jsonRaw, "\t", "")
	jsonRaw = strings.ReplaceAll(jsonRaw, "\n", "")
	jsonRaw = strings.ReplaceAll(jsonRaw, "\r", "")

	var data struct {
		DP  interface{} `json:"dp"`
		DBI DBI         `json:"dbi"`
	}
	if err := json.Unmarshal([]byte(jsonRaw), &data); err != nil {
		return fmt.Errorf("json parse error: %w", err)
	}

	e.SchoolYear = 0
	if dpMap, ok := data.DP.(map[string]interface{}); ok {
		if year, ok := dpMap["year"].(float64); ok {
			e.SchoolYear = int(year)
		}
	}
	if e.SchoolYear == 0 {
		e.SchoolYear = time.Now().Year()
	}

	e.Dbi = data.DBI

	if e.GsecHash == "" {
		return fmt.Errorf("security hash (gsec) not found")
	}

	return nil
}

func (e *EdupageClient) GetFullSchedule(date time.Time) (ScheduleResponse, error) {
	if len(e.Dbi.Classes) == 0 {
		if err := e.loginInternal(); err != nil {
			return ScheduleResponse{}, fmt.Errorf("auto-login failed: %w", err)
		}
	}

	response := ScheduleResponse{
		Classes: []ClassSchedule{},
	}
	var mu sync.Mutex
	var wg sync.WaitGroup

	sem := make(chan struct{}, 5)
	dateStr := date.Format("2006-01-02")

	for classID, classItem := range e.Dbi.Classes {
		wg.Add(1)
		go func(cid, cname string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			rawLessons, err := e.fetchRawClassTimetable(cid, dateStr)
			if err != nil {
				return
			}

			// Pass the requested date to correctly form time.Time
			classSched := e.processClassLessons(cname, rawLessons, date)

			if len(classSched.Schedule) > 0 {
				mu.Lock()
				response.Classes = append(response.Classes, classSched)
				mu.Unlock()
			}
		}(classID, classItem.Name)
	}

	wg.Wait()

	sort.Slice(response.Classes, func(i, j int) bool {
		return response.Classes[i].ClassName < response.Classes[j].ClassName
	})

	return response, nil
}

func (e *EdupageClient) fetchRawClassTimetable(classID, dateStr string) ([]rawLesson, error) {
	endpoint := fmt.Sprintf("https://%s.edupage.org/timetable/server/currenttt.js?__func=curentttGetData", e.Subdomain)

	payload := map[string]interface{}{
		"__args": []interface{}{
			nil,
			map[string]interface{}{
				"year":       e.SchoolYear,
				"datefrom":   dateStr,
				"dateto":     dateStr,
				"table":      "classes",
				"id":         classID,
				"showColors": true,
				"showOrig":   true,
				"log_module": "CurrentTTView",
			},
		},
		"__gsh": e.GsecHash,
	}

	bodyBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := e.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp rawTimetableResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}
	if apiResp.R.Error != "" {
		return nil, fmt.Errorf("error: %s", apiResp.R.Error)
	}

	return apiResp.R.TTItems, nil
}

func (e *EdupageClient) processClassLessons(className string, raw []rawLesson, baseDate time.Time) ClassSchedule {
	grouped := make(map[string][]LessonDetail)
	timeKeys := make(map[string]time.Time) // Map string key to start time object

	for _, r := range raw {
		if r.SubjectID == "" || r.Type == "period" {
			continue
		}

		subject := Subject{
			Name:  "-",
			Short: "-",
		}
		if s, ok := e.Dbi.Subjects[r.SubjectID]; ok {
			subject.Short = strings.TrimSpace(s.Short)
			subject.Name = strings.TrimSpace(s.Name)
		}

		teacher := Teacher{
			Firstname: "",
			Lastname:  "",
			Short:     "",
		}
		for _, tid := range r.TeacherIDs {
			if t, ok := e.Dbi.Teachers[tid]; ok {
				teacher.Firstname = strings.TrimSpace(t.Firstname)
				teacher.Lastname = strings.TrimSpace(t.Lastname)
				teacher.Short = strings.TrimSpace(t.Short)
				break
			}
		}

		var room string
		for _, rid := range r.ClassroomIDs {
			if rm, ok := e.Dbi.Classrooms[rid]; ok {
				room = rm.Short
			}
		}

		startStr, endStr := r.StartTime, r.EndTime
		if startStr == "" {
			if p, ok := e.Dbi.Periods[r.Period]; ok {
				startStr, endStr = p.StartTime, p.EndTime
			}
		}

		// Convert HH:MM to time.Time
		startTime := parseEduTime(startStr, baseDate)
		endTime := parseEduTime(endStr, baseDate)

		// Use a unique key for grouping lessons in the same slot
		slotKey := fmt.Sprintf("%s|%s", startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))
		timeKeys[slotKey] = startTime

		grouped[slotKey] = append(grouped[slotKey], LessonDetail{
			Subject: subject,
			Teacher: teacher,
			Room:    room,
		})
	}

	var schedule []TimeSlot
	for key, lessons := range grouped {
		// Split key back to times or rely on stored times (requires passing End around)
		// Simpler: re-parse key or better yet, store the Time objects in a struct map?
		// Let's just re-parse the key which is standard format now.
		parts := strings.Split(key, "|")
		t1, _ := time.Parse(time.RFC3339, parts[0])
		t2, _ := time.Parse(time.RFC3339, parts[1])

		schedule = append(schedule, TimeSlot{
			Start:   t1,
			End:     t2,
			Lessons: lessons,
		})
	}

	sort.Slice(schedule, func(i, j int) bool {
		return schedule[i].Start.Before(schedule[j].Start)
	})

	return ClassSchedule{
		ClassName: className,
		Schedule:  schedule,
	}
}

// Helper to combine "HH:MM" with the date's Year/Month/Day
func parseEduTime(timeStr string, date time.Time) time.Time {
	if timeStr == "" {
		return time.Time{}
	}
	parts := strings.Split(timeStr, ":")
	if len(parts) < 2 {
		return time.Time{}
	}
	h, _ := strconv.Atoi(parts[0])
	m, _ := strconv.Atoi(parts[1])

	return time.Date(date.Year(), date.Month(), date.Day(), h, m, 0, 0, date.Location())
}
