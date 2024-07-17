package login

import (
	"eyes/internal/domain"
	"eyes/internal/service/login"
	C "eyes/internal/web/common"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type (
	Controller struct {
		service login.LService
		logger  *zap.Logger
	}

	VOUser struct {
		Name     string `json:"name" form:"name"`
		PassWord string `json:"password" form:"password"`
	}
)

func NewController(service login.LService, logger *zap.Logger) *Controller {
	return &Controller{
		service: service,
		logger:  logger,
	}
}

func (a *Controller) RegisterRoutes(server *gin.Engine) {
	login := server.Group("v1/api/login")
	login.POST("/login-by-password", a.LoginByPassword)
}

func (a *Controller) LoginByPassword(c *gin.Context) {
	var user VOUser
	err := c.Bind(&user)
	if err != nil {
		a.logger.Error("c.Bind(&user)", zap.Error(err))
		return
	}

	userID, err := a.service.LoginByPass(c.Request.Context(), domain.User{
		Name:     user.Name,
		PassWord: user.PassWord,
	})
	if err != nil {
		a.logger.Error("a.service.LoginByPass", zap.Error(err))
		return
	}

	if userID > 0 {
		session := sessions.Default(c)
		session.Set("userID", userID)

		a.logger.Info("登录成功，Write the sessionID to session")

		C.RespSuccess(c, nil, C.Page{})
		return
	} else {
		a.logger.Info("登录失败，用户名或密码不存在")
		C.RespSuccess(c, "登录失败，用户名或密码不存在", C.Page{})
	}
}
