package login

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"eyes/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestController_RegisterRoutes(t *testing.T) {
}

type mockLoginService struct{}

func (l mockLoginService) LoginByPass(ctx context.Context, article domain.User) (int64, error) {
	return 1, nil
}

func (l mockLoginService) LoginByMobile(ctx context.Context, article domain.User) error {
	return nil
}

func TestNewController(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockLogin := mockLoginService{}
	logger, _ := zap.NewDevelopment()
	controller := NewController(mockLogin, logger)
	controller.RegisterRoutes(r)

	req, _ := http.NewRequest(http.MethodGet, "/v1.api.login/login-by-password", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
