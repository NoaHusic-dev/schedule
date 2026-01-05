package internal

import (
	"bytes"
	"compress/flate"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html"
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

// --- Custom Errors ---

var (
	ErrTwoFactorNeeded = errors.New("two-factor authentication required")
	ErrCaptchaNeeded   = errors.New("captcha required")
	ErrBadCredentials  = errors.New("bad credentials")
)

// --- Structs for Schedule ---

type ScheduleResponse struct {
	Classes []ClassSchedule `json:"classes"`
}

type ClassSchedule struct {
	Class    string     `json:"class"`
	Schedule []TimeSlot `json:"schedule"`
}

type TimeSlot struct {
	Start   time.Time      `json:"start"`
	End     time.Time      `json:"end"`
	Lessons []LessonDetail `json:"lessons"`
}

type LessonDetail struct {
	Subject Subject  `json:"subject"`
	Teacher Teacher  `json:"teacher"`
	Room    string   `json:"room"`
	Groups  []string `json:"groups"`
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

// --- Controllers ---

type TwoFactorController struct {
	client    *EdupageClient
	authEp    string
	authToken string
	csrfToken string
	code      string
}

type CaptchaController struct {
	client    *EdupageClient
	Image     []byte
	inputName string
	csrfToken string
}

// --- Internal API Structs ---

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
	GroupNames   []string `json:"groupnames"`
	Period       string   `json:"period"`
	StartTime    string   `json:"starttime"`
	EndTime      string   `json:"endtime"`
	Type         string   `json:"type"`
}

// --- Client Implementation ---

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

func (e *EdupageClient) Login() (interface{}, error) {
	return e.loginInternal(nil)
}

func (e *EdupageClient) ReloadData(subdomain, sessionID, username string) error {
	u, err := url.Parse(fmt.Sprintf("https://%s.edupage.org", subdomain))
	if err != nil {
		return fmt.Errorf("invalid subdomain: %w", err)
	}

	cookie := &http.Cookie{
		Name:   "PHPSESSID",
		Value:  sessionID,
		Path:   "/",
		Domain: fmt.Sprintf("%s.edupage.org", subdomain),
	}
	e.Client.Jar.SetCookies(u, []*http.Cookie{cookie})

	e.Subdomain = subdomain
	e.Username = username

	reqURL := fmt.Sprintf("https://%s.edupage.org/user", subdomain)
	resp, err := e.Client.Get(reqURL)
	if err != nil {
		return fmt.Errorf("failed to reload user data: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	return e.parseLoginData(string(bodyBytes))
}

func (e *EdupageClient) loginInternal(captchaData map[string]string) (interface{}, error) {
	requestURL := fmt.Sprintf("https://%s.edupage.org/login/?cmd=MainLogin", e.Subdomain)
	resp, err := e.Client.Get(requestURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)
	data := string(bodyBytes)

	csrfToken := extractCsrfToken(data)
	if csrfToken == "" {
		csrfToken = extractInputValue(data, "csrfauth")
		if csrfToken == "" {
			csrfToken = extractInputValue(data, "requestid")
		}
	}
	if csrfToken == "" {
		return nil, fmt.Errorf("could not locate `csrftoken`, `csrfauth`, or `requestid`")
	}

	loginURL := fmt.Sprintf("https://%s.edupage.org/login/edubarLogin.php", e.Subdomain)
	form := url.Values{}
	form.Set("csrfauth", csrfToken)
	form.Set("username", e.Username)
	form.Set("password", e.Password)

	for k, v := range captchaData {
		form.Set(k, v)
	}

	resp, err = e.Client.PostForm(loginURL, form)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respData, _ := io.ReadAll(resp.Body)
	respStr := string(respData)
	finalURL := resp.Request.URL.String()

	if strings.Contains(finalURL, "cap=1") || strings.Contains(finalURL, "lerr=b43b43") {
		return e.handleCaptchaPage(respStr, finalURL)
	}

	if strings.Contains(finalURL, "bad=1") {
		return nil, ErrBadCredentials
	}

	if strings.Contains(finalURL, "twofactor") {
		if e.Subdomain == "login1" {
			parts := strings.Split(respStr, "-->")
			if len(parts) > 0 {
				subparts := strings.Split(parts[0], " ")
				e.Subdomain = subparts[len(subparts)-1]
			}
		}

		twoFactorURL := fmt.Sprintf("https://%s.edupage.org/login/twofactor?sn=1", e.Subdomain)
		tfResp, err := e.Client.Get(twoFactorURL)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch 2fa page: %w", err)
		}
		defer tfResp.Body.Close()
		tfBody, _ := io.ReadAll(tfResp.Body)
		tfData := string(tfBody)

		tfController := &TwoFactorController{
			client:    e,
			csrfToken: extractInputValue(tfData, "csrfauth"),
			authToken: extractInputValue(tfData, "au"),
			authEp:    extractInputValue(tfData, "gu"),
		}
		return tfController, ErrTwoFactorNeeded
	}

	if e.Subdomain == "login1" {
		parts := strings.Split(respStr, "-->")
		if len(parts) > 0 {
			subparts := strings.Split(parts[0], " ")
			e.Subdomain = subparts[len(subparts)-1]
		}
	}

	return nil, e.parseLoginData(respStr)
}

func (e *EdupageClient) handleCaptchaPage(htmlContent, pageURL string) (*CaptchaController, error) {
	imgRe := regexp.MustCompile(`(?i)<img[^>]+src="([^"]*captcha[^"]*)"`)
	matches := imgRe.FindStringSubmatch(htmlContent)
	if len(matches) < 2 {
		return nil, fmt.Errorf("captcha detected but image URL not found in HTML")
	}

	imgRelPath := html.UnescapeString(matches[1])

	u, _ := url.Parse(pageURL)
	rel, _ := url.Parse(imgRelPath)
	imgAbsURL := u.ResolveReference(rel).String()

	inputName := "captcha"
	inputRe := regexp.MustCompile(`(?i)<input[^>]+name="(captcha[^"]*)"`)
	inputMatches := inputRe.FindStringSubmatch(htmlContent)
	if len(inputMatches) >= 2 {
		inputName = inputMatches[1]
	}

	imgResp, err := e.Client.Get(imgAbsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download captcha image: %w", err)
	}
	defer imgResp.Body.Close()
	imgData, err := io.ReadAll(imgResp.Body)
	if err != nil {
		return nil, err
	}

	return &CaptchaController{
		client:    e,
		Image:     imgData,
		inputName: inputName,
	}, ErrCaptchaNeeded
}

func (cc *CaptchaController) Submit(answer string) (interface{}, error) {
	data := map[string]string{
		cc.inputName: answer,
	}
	return cc.client.loginInternal(data)
}

// --- 2FA Methods ---

func (tf *TwoFactorController) IsConfirmed() (bool, error) {
	reqURL := fmt.Sprintf("https://%s.edupage.org/login/twofactor?akcia=checkIfConfirmed", tf.client.Subdomain)
	resp, err := tf.client.Client.Post(reqURL, "application/json", nil)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var res map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return false, err
	}

	status, _ := res["status"].(string)
	if status == "fail" {
		return false, nil
	}
	if status != "ok" {
		return false, fmt.Errorf("invalid response from server: %v", res)
	}

	if data, ok := res["data"].(string); ok {
		tf.code = data
		return true, nil
	}
	return true, nil
}

func (tf *TwoFactorController) ResendNotifications() error {
	reqURL := fmt.Sprintf("https://%s.edupage.org/login/twofactor?akcia=resendNotifs", tf.client.Subdomain)
	resp, err := tf.client.Client.Post(reqURL, "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (tf *TwoFactorController) Finish(manualCode ...string) error {
	code := tf.code
	if len(manualCode) > 0 {
		code = manualCode[0]
	}
	if code == "" {
		return errors.New("no confirmation code found (did you call IsConfirmed?)")
	}

	if tf.csrfToken != "" && tf.authEp != "" && tf.authToken != "" {
		return tf.finishLegacy(code)
	}
	return tf.finishEQ(code)
}

func (tf *TwoFactorController) finishLegacy(code string) error {
	reqURL := fmt.Sprintf("https://%s.edupage.org/login/edubarLogin.php", tf.client.Subdomain)
	form := url.Values{}
	form.Set("csrfauth", tf.csrfToken)
	form.Set("t2fasec", code)
	form.Set("2fNoSave", "y")
	form.Set("2fform", "1")
	form.Set("gu", tf.authEp)
	form.Set("au", tf.authToken)

	resp, err := tf.client.Client.PostForm(reqURL, form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)

	if strings.Contains(string(bodyBytes), "window.location = gu;") {
		return tf.reloadAfterTwoFactor()
	}

	return errors.New("second factor failed (legacy)")
}

func (tf *TwoFactorController) finishEQ(code string) error {
	reqURL := fmt.Sprintf("https://%s.edupage.org/login/twofactor?akcia=checkIfConfirmed", tf.client.Subdomain)

	params := map[string]interface{}{
		"t2fasec":  code,
		"2fNoSave": "y",
		"2fform":   "1",
		"tu":       nil,
		"gu":       tf.authEp,
		"au":       tf.authToken,
	}

	payload, err := buildEqPayload(params)
	if err != nil {
		return err
	}

	postData := url.Values{}
	for k, v := range payload {
		postData.Set(k, v)
	}

	resp, err := tf.client.Client.PostForm(reqURL, postData)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var res map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return fmt.Errorf("invalid response json: %w", err)
	}

	if status, ok := res["status"].(string); !ok || status != "ok" {
		return fmt.Errorf("second factor failed: %v", res)
	}

	return tf.reloadAfterTwoFactor()
}

func (tf *TwoFactorController) reloadAfterTwoFactor() error {
	u, _ := url.Parse(fmt.Sprintf("https://%s.edupage.org", tf.client.Subdomain))
	cookies := tf.client.Client.Jar.Cookies(u)
	var sessionID string
	for _, c := range cookies {
		if c.Name == "PHPSESSID" {
			sessionID = c.Value
			break
		}
	}

	if sessionID == "" {
		return errors.New("missing PHPSESSID after 2FA")
	}

	return tf.client.ReloadData(tf.client.Subdomain, sessionID, tf.client.Username)
}

// --- Helper Functions ---

func buildEqPayload(params map[string]interface{}) (map[string]string, error) {
	rpcPayloadBytes, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	v := url.Values{}
	v.Set("rpcparams", string(rpcPayloadBytes))
	encoded := v.Encode()

	var b bytes.Buffer
	fw, _ := flate.NewWriter(&b, 9)
	if _, err := fw.Write([]byte(encoded)); err != nil {
		return nil, err
	}
	if err := fw.Close(); err != nil {
		return nil, err
	}
	compressed := b.Bytes()

	b64Compressed := base64.StdEncoding.EncodeToString(compressed)
	cs0 := "dz:" + b64Compressed

	hasher := sha1.New()
	hasher.Write([]byte(cs0))
	sha1Hash := hex.EncodeToString(hasher.Sum(nil))

	return map[string]string{
		"eqap":  cs0,
		"eqacs": sha1Hash,
		"eqaz":  "1",
	}, nil
}

func extractCsrfToken(html string) string {
	var dummy map[string]interface{}
	if err := json.Unmarshal([]byte(html), &dummy); err == nil {
		if val, ok := dummy["csrftoken"].(string); ok {
			return val
		}
	}
	re := regexp.MustCompile(`"csrftoken"\s*:\s*"([^"]+)"`)
	matches := re.FindStringSubmatch(html)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func extractInputValue(htmlContent, key string) string {
	patterns := []string{
		fmt.Sprintf(`%s"?\s*value="([^"]+)"`, regexp.QuoteMeta(key)),
		fmt.Sprintf(`name=["']%s["'][^>]*value=["']([^"']+)["']`, regexp.QuoteMeta(key)),
		fmt.Sprintf(`id=["']%s["'][^>]*value=["']([^"']+)["']`, regexp.QuoteMeta(key)),
		fmt.Sprintf(`%s\s*[:=]\s*"([^"]+)"`, regexp.QuoteMeta(key)),
		fmt.Sprintf(`%s\s*[:=]\s*'([^']+)'`, regexp.QuoteMeta(key)),
	}

	for _, p := range patterns {
		re := regexp.MustCompile("(?i)" + p)
		match := re.FindStringSubmatch(htmlContent)
		if len(match) > 1 {
			return html.UnescapeString(match[1])
		}
	}
	return ""
}

func (e *EdupageClient) parseLoginData(htmlStr string) error {
	hashMatch := regexp.MustCompile(`ASC\.gsechash="([^"]+)"`).FindStringSubmatch(htmlStr)
	if len(hashMatch) >= 2 {
		e.GsecHash = hashMatch[1]
	}

	parts := strings.SplitN(htmlStr, "userhome(", 2)
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
		return fmt.Errorf("json parse error in userhome: %w", err)
	}

	// 1. Extract School Year
	e.SchoolYear = 0
	if dpMap, ok := data.DP.(map[string]interface{}); ok {
		if year, ok := dpMap["year"].(float64); ok {
			e.SchoolYear = int(year)
		}
	}

	// 2. Fix Year for Early Months
	if e.SchoolYear == 0 {
		now := time.Now()
		e.SchoolYear = now.Year()
		if now.Month() < 8 {
			e.SchoolYear = now.Year() - 1
		}
	}

	e.Dbi = data.DBI

	if e.GsecHash == "" {
		return fmt.Errorf("security hash (gsec) not found")
	}

	return nil
}

// --- Schedule Logic ---

func (e *EdupageClient) GetFullSchedule(date time.Time) (ScheduleResponse, error) {
	if len(e.Dbi.Classes) == 0 {
		// Attempt auto-login if classes are empty
		if _, err := e.Login(); err != nil {
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
		return response.Classes[i].Class < response.Classes[j].Class
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

func (e *EdupageClient) processClassLessons(class string, raw []rawLesson, baseDate time.Time) ClassSchedule {
	grouped := make(map[string][]LessonDetail)

	for _, r := range raw {
		if r.SubjectID == "" || r.Type == "period" {
			continue
		}

		subject := Subject{Name: "-", Short: "-"}
		if s, ok := e.Dbi.Subjects[r.SubjectID]; ok {
			subject.Short = strings.TrimSpace(s.Short)
			subject.Name = strings.TrimSpace(s.Name)
		}

		teacher := Teacher{Firstname: "", Lastname: "", Short: ""}
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

		var groups []string
		for _, g := range r.GroupNames {
			trimmed := strings.TrimSpace(g)
			if trimmed != "" {
				groups = append(groups, trimmed)
			}
		}

		startStr, endStr := r.StartTime, r.EndTime
		if startStr == "" {
			if p, ok := e.Dbi.Periods[r.Period]; ok {
				startStr, endStr = p.StartTime, p.EndTime
			}
		}

		startTime := parseEduTime(startStr, baseDate)
		endTime := parseEduTime(endStr, baseDate)

		slotKey := fmt.Sprintf("%s|%s", startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))

		grouped[slotKey] = append(grouped[slotKey], LessonDetail{
			Subject: subject,
			Teacher: teacher,
			Room:    room,
			Groups:  groups,
		})
	}

	var schedule []TimeSlot
	for key, lessons := range grouped {
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
		Class:    class,
		Schedule: schedule,
	}
}

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
