package article_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"eyes/internal/web/article"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"eyes/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockArtService struct {
	mock.Mock
}

func (m *MockArtService) Save(ctx context.Context, art domain.Article) (int64, error) {
	args := m.Called(ctx, art)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockArtService) Publish(ctx context.Context, art domain.Article) error {
	args := m.Called(ctx, art)
	return args.Error(0)
}

func TestSave(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockArtService)
	logger, _ := zap.NewDevelopment()
	controller := article.NewArticleController(mockService, logger)

	router := gin.Default()
	controller.RegisterRoutes(router)

	t.Run("successful save", func(t *testing.T) {
		mockService.On("Save", mock.Anything, mock.AnythingOfType("domain.Article")).Return(int64(1), nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/v1.api.article/new", strings.NewReader(`{"title": "Test Title", "content": "Test Content"}`))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "1", w.Body.String())

		mockService.AssertExpectations(t)
	})

	t.Run("missing title or content", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/v1.api.article/new", strings.NewReader(`{"title": "", "content": "Test Content"}`))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "标题或者内容为空")
	})
}

func TestPublish(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockArtService)
	logger, _ := zap.NewDevelopment()
	controller := article.NewArticleController(mockService, logger)

	router := gin.Default()
	controller.RegisterRoutes(router)

	t.Run("successful publish", func(t *testing.T) {
		mockService.On("Publish", mock.Anything, mock.AnythingOfType("domain.Article")).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/v1.api.article/publish",
			strings.NewReader(`{"title": "Test Title", "content": "Test Content"}`))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		assert.Equal(t, "/v1/api/article/publish-success", w.Header().Get("Location"))

		mockService.AssertExpectations(t)
	})

	t.Run("publish failed", func(t *testing.T) {
		mockService.On("Publish", mock.Anything, mock.AnythingOfType("domain.Article")).Return(assert.AnError)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/v1/api/article/publish",
			strings.NewReader(`{"title": "Test Title1", "content": "Test Content"}`))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		assert.Equal(t, "/v1/api/article/publish-failed", w.Header().Get("Location"))
	})
}
