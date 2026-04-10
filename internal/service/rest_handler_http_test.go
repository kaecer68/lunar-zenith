package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRestHandlerGetCalendarInvalidDateReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewRestHandler(NewAggregator(nil, nil)).RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/v1/calendar?date=2026/04/11", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d; want %d", rec.Code, http.StatusBadRequest)
	}

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if body["error"] != invalidCalendarDateMessage {
		t.Errorf("error = %q; want %q", body["error"], invalidCalendarDateMessage)
	}
}

func TestRestHandlerGetCalendarEmptyDateUsesTaipeiToday(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewRestHandler(NewAggregator(nil, nil)).RegisterRoutes(router)
	want, err := resolveCalendarQueryTime("")
	if err != nil {
		t.Fatalf("resolveCalendarQueryTime(\"\") error = %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/calendar", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; want %d", rec.Code, http.StatusOK)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if body["gregorian_date"] != want.Format("2006-01-02") {
		t.Errorf("gregorian_date = %v; want %q", body["gregorian_date"], want.Format("2006-01-02"))
	}
}

func TestRestHandlerGetCalendarLocksSolarTermAndSnakeCaseNestedFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewRestHandler(NewAggregator(nil, nil)).RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/v1/calendar?date=2026-03-18", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; want %d", rec.Code, http.StatusOK)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	solarTerm, ok := body["solar_term"].(map[string]any)
	if !ok {
		t.Fatalf("solar_term should be object, got %T", body["solar_term"])
	}
	if solarTerm["name"] == "" {
		t.Errorf("solar_term.name should not be empty")
	}
	if _, exists := body["solar_term_detail"]; exists {
		t.Errorf("solar_term_detail should be removed in v4")
	}

	shenSha, ok := body["shen_sha"].([]any)
	if !ok || len(shenSha) == 0 {
		t.Fatalf("shen_sha should contain at least one item, got %T", body["shen_sha"])
	}
	firstShenSha, ok := shenSha[0].(map[string]any)
	if !ok {
		t.Fatalf("shen_sha[0] should be object, got %T", shenSha[0])
	}
	if _, ok := firstShenSha["name"]; !ok {
		t.Errorf("shen_sha[0] should expose snake_case key name")
	}
	if _, ok := firstShenSha["Name"]; ok {
		t.Errorf("shen_sha[0] should not expose PascalCase key Name")
	}

	mansion, ok := body["mansion"].(map[string]any)
	if !ok {
		t.Fatalf("mansion should be object, got %T", body["mansion"])
	}
	if _, ok := mansion["name"]; !ok {
		t.Errorf("mansion should expose snake_case key name")
	}
	if _, ok := mansion["Name"]; ok {
		t.Errorf("mansion should not expose PascalCase key Name")
	}
	if _, ok := mansion["full_name"]; !ok {
		t.Errorf("mansion should expose snake_case key full_name")
	}

	dailyDeity, ok := body["daily_deity"].(map[string]any)
	if !ok {
		t.Fatalf("daily_deity should be object, got %T", body["daily_deity"])
	}
	if _, ok := dailyDeity["name"]; !ok {
		t.Errorf("daily_deity should expose snake_case key name")
	}
	if _, ok := dailyDeity["Name"]; ok {
		t.Errorf("daily_deity should not expose PascalCase key Name")
	}

	fetalGod, ok := body["fetal_god"].(map[string]any)
	if !ok {
		t.Fatalf("fetal_god should be object, got %T", body["fetal_god"])
	}
	if _, ok := fetalGod["position"]; !ok {
		t.Errorf("fetal_god should expose snake_case key position")
	}
	if _, ok := fetalGod["Position"]; ok {
		t.Errorf("fetal_god should not expose PascalCase key Position")
	}

	clashSha, ok := body["clash_sha"].(map[string]any)
	if !ok {
		t.Fatalf("clash_sha should be object, got %T", body["clash_sha"])
	}
	if _, ok := clashSha["clash_zodiac"]; !ok {
		t.Errorf("clash_sha should expose snake_case key clash_zodiac")
	}
	if _, ok := clashSha["ClashZodiac"]; ok {
		t.Errorf("clash_sha should not expose PascalCase key ClashZodiac")
	}
}
