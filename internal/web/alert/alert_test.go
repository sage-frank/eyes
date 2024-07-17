package alert

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"eyes/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type MockMonitorService struct{}

func (m *MockMonitorService) Count(ctx context.Context, monitor domain.Monitor) (int64, error) {
	return 10, nil
}

func (m *MockMonitorService) Detail(ctx context.Context, monitor domain.Monitor) (*domain.Monitor, error) {
	return &domain.Monitor{ID: monitor.ID}, nil
}

func (m *MockMonitorService) List(ctx context.Context, page, size int64, monitor domain.Monitor) ([]*domain.Monitor, int64, error) {
	return []*domain.Monitor{{ID: "1"}, {ID: "2"}}, 2, nil
}

func (m *MockMonitorService) Query(ctx context.Context, page, size int64, filter string, monitor domain.Monitor) ([]*domain.Monitor, int64, error) {
	return []*domain.Monitor{{ID: "1"}, {ID: "2"}}, 2, nil
}

func TestMonitorController_Count(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := &MockMonitorService{}
	logger, _ := zap.NewDevelopment()
	controller := NewMonitorController(mockService, logger)
	controller.RegisterRoutes(r)

	req, _ := http.NewRequest(http.MethodGet, "/v1.api.alert/count", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestMonitorController_Detail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := &MockMonitorService{}
	logger, _ := zap.NewDevelopment()
	controller := NewMonitorController(mockService, logger)
	controller.RegisterRoutes(r)

	req, _ := http.NewRequest(http.MethodGet, "/v1/api/alert/detail?id=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestMonitorController_List(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := &MockMonitorService{}
	logger, _ := zap.NewDevelopment()
	controller := NewMonitorController(mockService, logger)
	controller.RegisterRoutes(r)

	req, _ := http.NewRequest(http.MethodGet, "/v1/api/alert/list", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestMonitorController_Query(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := &MockMonitorService{}
	logger, _ := zap.NewDevelopment()
	controller := NewMonitorController(mockService, logger)
	controller.RegisterRoutes(r)

	queryParams := map[string]string{
		"src_port":  "8080",
		"dest_port": "9090",
		"proto":     "TCP",
		"action":    "ALLOW",
	}

	jsonValue, _ := json.Marshal(queryParams)
	req, _ := http.NewRequest(http.MethodPost, "/v1/api/alert/query", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
