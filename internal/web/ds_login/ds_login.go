package ds_login

import (
	"errors"
	"fmt"
	"image/color"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"eyes/internal/service/ds_login"
	"github.com/fogleman/gg"
	"github.com/gin-contrib/sessions"

	"eyes/internal/domain"
	C "eyes/internal/web/common"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type (
	Controller struct {
		service ds_login.DsLoginServer
		logger  *zap.Logger
	}

	DsLogin struct {
		Phone    string `json:"phone,omitempty"`
		Username string `json:"username,omitempty"`
		Email    string `json:"email,omitempty"`
		Password string `json:"password"`
	}

	DsRegisterUser struct {
		Nickname    string `json:"nickname,omitempty"`    // 用户昵称
		Username    string `json:"username,omitempty"`    // 用户名
		Password    string `json:"password,omitempty"`    // 密码
		Phone       string `json:"phone,omitempty"`       // 用户电话号码
		Email       string `json:"email,omitempty"`       // 用户电子邮箱
		Avatar      string `json:"avatar,omitempty"`      // 用户头像URL
		Description string `json:"description,omitempty"` // 用户描述或简介
	}
)

func NewDsLoginController(service ds_login.DsLoginServer, logger *zap.Logger) *Controller {
	return &Controller{
		service: service,
		logger:  logger,
	}
}

func (a *Controller) RegisterRoutes(server *gin.Engine) {
	ds := server.Group("v1/api/ds_login")
	ds.GET("/log-in", a.HtmlCaptcha)
	ds.GET("/captcha", a.Captcha)
	ds.POST("/verify", a.Verify)
	ds.GET("/log-out", a.Logout)
	ds.POST("/register-user", a.registerUser)
}

func generateCaptcha1() (string, string, int) {
	rand.Seed(time.Now().UnixNano())

	// 生成随机初速度 (u)、加速度 (a) 和时间 (t)
	u := rand.Intn(20) + 1 // 保证不是0
	a := rand.Intn(10) + 1 // 保证不是0
	t := rand.Intn(5) + 1  // 保证不是0

	// 计算位移 s = ut + 0.5at^2
	s := u*t + (a*t*t)/2

	// 生成验证码公式字符串
	captchaExpression := fmt.Sprintf("%d * %d + 0.5 * %d * %d^2", u, t, a, t)

	// 生成验证码ID
	captchaID := strconv.Itoa(u) + "u" + strconv.Itoa(a) + "a" + strconv.Itoa(t) + "t"

	return captchaID, captchaExpression, s
}

func generateCaptcha() (string, string, int) {
	rand.Seed(time.Now().UnixNano())

	// 生成两个随机数字
	num1 := rand.Intn(9) + 1 // 保证不是0
	num2 := rand.Intn(9) + 1

	// 生成一个随机运算符
	operators := []string{"+", "-", "*"}
	operator := operators[rand.Intn(len(operators))]

	// 计算结果
	var result int
	switch operator {
	case "+":
		result = num1 + num2
	case "-":
		result = num1 - num2
	case "*":
		result = num1 * num2
	}

	// 生成5个干扰字符
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	interference := make([]byte, 5)
	for i := range interference {
		interference[i] = chars[rand.Intn(len(chars))]
	}

	// 将干扰字符分为两部分
	firstHalf := interference[:2]
	secondHalf := interference[2:]

	// 创建特定结构的验证码组件
	captchaComponents := make([]string, 7)
	captchaComponents[0] = strconv.Itoa(num1)
	for i, c := range firstHalf {
		captchaComponents[i+1] = string(c)
	}
	captchaComponents[3] = operator
	for i, c := range secondHalf {
		captchaComponents[i+4] = string(c)
	}
	captchaComponents[6] = strconv.Itoa(num2)

	// 打乱干扰字符
	rand.Shuffle(len(firstHalf), func(i, j int) {
		firstHalf[i], firstHalf[j] = firstHalf[j], firstHalf[i]
	})
	rand.Shuffle(len(secondHalf), func(i, j int) {
		secondHalf[i], secondHalf[j] = secondHalf[j], secondHalf[i]
	})

	// 生成验证码ID
	captchaID := strings.Join(captchaComponents, "")
	captchaExpression := fmt.Sprintf("%d %s %d", num1, operator, num2)

	return captchaID, captchaExpression, result
}

func createCaptchaImage(captchaID string) *gg.Context {
	const width = 240
	const height = 80

	dc := gg.NewContext(width, height)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	if err := dc.LoadFontFace("D:\\code\\golang\\eyes\\fonts\\font.ttf", 64); err != nil {
		panic(err)
	}
	dc.SetColor(color.Black)
	dc.DrawStringAnchored(captchaID, width/2, height/2, 0.5, 0.5)

	// 生成更多干扰线
	for i := 0; i < 10; i++ { // 增加干扰线数量
		x1 := rand.Float64() * width
		y1 := rand.Float64() * height
		x2 := rand.Float64() * width
		y2 := rand.Float64() * height
		dc.DrawLine(x1, y1, x2, y2)
		dc.Stroke()
	}

	return dc
}

func (a *Controller) registerUser(c *gin.Context) {
	var registerUserInfo DsRegisterUser
	err := c.Bind(&registerUserInfo)
	if err != nil {
		a.logger.Error("c.Bind(&registerUserInfo): %w", zap.Error(err))
		return
	}

	dUser := domain.DSUser{
		Nickname:    registerUserInfo.Nickname,
		Username:    registerUserInfo.Username,
		Password:    registerUserInfo.Password,
		Phone:       registerUserInfo.Phone,
		Email:       registerUserInfo.Email,
		Avatar:      registerUserInfo.Avatar,
		Description: registerUserInfo.Description,
	}
	r, err := a.service.Select(c.Request.Context(), dUser)
	a.logger.Info("查询是否存在", zap.Any("user", r))
	if err != nil {
		a.logger.Error("a.service.Select", zap.Error(err))
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			C.RespErrorWithMsg(c, 400, "信息错误")
			return
		}
	} else if r.ID != "" {
		a.logger.Error("a.service.Select(: %w", zap.Error(err))
		C.RespErrorWithMsg(c, 400, "信息已经存在")
		return
	}

	id, err := a.service.Save(c.Request.Context(), dUser)
	a.logger.Info("保存数据:", zap.Any("save-user-id", id))
	if err != nil {
		a.logger.Error("a.service.Save(c.Request.Context(), dUser)", zap.Error(err))
		C.RespErrorWithMsg(c, 400, "保存信息错误")
		return
	}
	C.RespSuccess(c, []gin.H{{"id": id}}, C.Page{Page: 0, Size: 10, Total: 0})
}

func (a *Controller) Verify(c *gin.Context) {
	type CaptchaData struct {
		Captcha string `json:"captcha" binding:"required"`
	}
	var userResult CaptchaData
	err := c.BindJSON(&userResult)
	a.logger.Info("userResult", zap.Any("userResult", userResult))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "input captcha error"})
		return
	}

	session := sessions.Default(c)
	captchaResult := session.Get("captcha_result")
	a.logger.Info("session captcha_result", zap.Any("captchaResult", captchaResult), zap.Any("userResult", userResult))

	if userResult.Captcha == captchaResult {
		c.JSON(http.StatusOK, gin.H{"data": "success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"data": "failed"})
	}
}

func (a *Controller) Captcha(c *gin.Context) {
	captchaID, _, captchaResult := generateCaptcha()
	session := sessions.Default(c)
	session.Set("captcha_id", captchaID)
	session.Set("captcha_result", strconv.Itoa(captchaResult))
	session.Save()
	img := createCaptchaImage(captchaID)
	// 设置响应头，告知内容类型为图片
	c.Header("Content-Type", "image/png")
	img.EncodePNG(c.Writer)
}

func (a *Controller) HtmlCaptcha(c *gin.Context) {
	c.HTML(http.StatusOK, "captcha.html", gin.H{})
}

func (a *Controller) Login(c *gin.Context) {
	var ds DsLogin
	err := c.Bind(&ds)
	if err != nil {
		a.logger.Error("c.Bind(&ds): %v", zap.Error(err))
		C.RespErrorWithMsg(c, 400, "账号密码不能为空")
		return
	}

	if ds.Phone == "" && ds.Username == "" && ds.Email == "" {
		a.logger.Info("账号或者手机号不能为空")
		C.RespErrorWithMsg(c, 400, "账号密码不能为空")
		return
	}

	if ds.Password == "" {
		a.logger.Info("密码不能为空")
		C.RespErrorWithMsg(c, 400, "密码不能为空")
		return
	}

	a.logger.Info("账号", zap.String("x-request-id", c.GetHeader("X-Request-ID")),
		zap.String("Phone", ds.Phone),
		zap.String("Password:", ds.Password))

	id, err := a.service.LoginByPass(c.Request.Context(), domain.DSUser{
		Phone:    ds.Phone,
		Username: ds.Username,
		Email:    ds.Email,
		Password: ds.Password,
	})

	if err != nil {
		a.logger.Error("err:", zap.Error(err))
		if errors.Is(err, C.ErrPassword) {
			C.RespErrorWithMsg(c, 400, "账号或密码不正确")
			return
		}
		c.String(400, "未知错误")
		return

	} else if id != "" {

		// 登录成功 把用户id 写入session中
		session := sessions.Default(c)
		session.Set("userID", id)

		a.logger.Info("登录成功")
		C.RespSuccess(c, []gin.H{{"id": id}}, C.Page{Page: 1, Size: 10, Total: 1})
		return

	} else {
		a.logger.Info("账号不存在")
		C.RespSuccess(c, "账号不存在", C.Page{Page: 1, Size: 10, Total: 1})
		return
	}
}

func (a *Controller) Logout(c *gin.Context) {
}
