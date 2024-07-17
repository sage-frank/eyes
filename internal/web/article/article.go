package article

import (
	"fmt"
	"net/http"

	"eyes/internal/domain"
	"eyes/internal/service/article"
	C "eyes/internal/web/common"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type (
	Controller struct {
		service article.ArtService
		logger  *zap.Logger
	}

	// VOArticle VO => view object
	VOArticle struct {
		Id      int64  `json:"id"`
		Title   string `json:"title" form:"title"`
		Content string `json:"content" form:"content"`
	}
)

func NewArticleController(service article.ArtService, logger *zap.Logger) *Controller {
	return &Controller{
		service: service,
		logger:  logger,
	}
}

func (a *Controller) RegisterRoutes(server *gin.Engine) {
	art := server.Group("v1/api/article")
	art.POST("/new", a.Save)
	art.POST("/publish", a.Publish)
	art.GET("/new", a.SavePage)
	art.GET("/publish-success", func(c *gin.Context) {
		c.String(http.StatusOK, "发表成功")
	})

	art.GET("/publish-failed", func(c *gin.Context) {
		c.String(http.StatusOK, "发表失败")
	})
}

func (a *Controller) Save(c *gin.Context) {
	var art VOArticle
	err := c.Bind(&art)
	if err != nil {
		a.logger.Error("c.Bind(&art): %v", zap.Error(err))
		return
	}

	if art.Title == "" || art.Content == "" {
		a.logger.Info("标题或者内容不能为空")
		C.RespErrorWithMsg(c, 400, "标题或者内容为空")
		return
	}

	a.logger.Info("标题", zap.String("x-request-id", c.GetHeader("X-Request-ID")), zap.String("title", art.Title), zap.String("content:", art.Content))

	id, err := a.service.Save(c.Request.Context(), domain.Article{
		ID:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Author:  123,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "系统错误")
		return
	}
	c.String(http.StatusOK, fmt.Sprintf("%d", id))
}

func (a *Controller) Publish(c *gin.Context) {
	var art VOArticle
	err := c.BindJSON(&art)
	if err != nil {
		a.logger.Error("c.BindJSON(&art)", zap.Error(err))
		c.Redirect(http.StatusTemporaryRedirect, "/v1.api.article/publish-failed")
		return
	}
	err = a.service.Publish(c.Request.Context(), domain.Article{
		ID:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Author:  123,
	})
	if err != nil {
		a.logger.Error("a.service.Publish", zap.Error(err))
		c.Redirect(http.StatusTemporaryRedirect, "/v1.api.article/publish-failed")
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, "/v1.api.article/publish-success")
}

func (a *Controller) SavePage(c *gin.Context) {
	c.HTML(http.StatusOK, "write_article.gohtml", nil)
}
