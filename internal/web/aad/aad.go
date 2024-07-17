package aad

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	neturl "net/url"
	"strings"

	"eyes/internal/domain"
	"eyes/internal/service/aad"
	C "eyes/internal/web/common"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type (
	VOAzure struct {
		ID      int64
		Author  int64
		Title   string
		Content string
	}

	AzureController struct {
		service aad.AzureService
		logger  *zap.Logger
	}

	ContextKey string

	azureRequest struct {
		ClientID     string `json:"client_id"`
		ResponseType string `json:"response_type"`
		RedirectURI  string `json:"redirect_uri"`
		Scope        string `json:"scope"`
		ResponseMode string `json:"response_mode"`
		State        string `json:"state"`
		Nonce        string `json:"nonce"`
		Prompt       string `json:"prompt"`
	}
)

func NewAzureAzureController(service aad.AzureService, logger *zap.Logger) *AzureController {
	return &AzureController{
		service: service,
		logger:  logger,
	}
}

func newAzureRequest(url string) *azureRequest {
	state, _ := generateSecureRandomString(6)
	nonce, _ := generateSecureRandomString(6)

	return &azureRequest{
		ClientID:     viper.GetString("aad.client_id"),
		ResponseType: "id_token",
		RedirectURI:  url,
		Scope:        neturl.QueryEscape(viper.GetString("aad.scope")), // python urllib.parse.quote
		ResponseMode: "form_post",
		State:        state,
		Nonce:        nonce,
		Prompt:       "select_account",
	}
}

func (a *AzureController) RegisterRoutes(server *gin.Engine) {
	azure := server.Group("/api/v1/azure")
	{
		azure.POST("/read-token", a.CallBackToken)
		azure.GET("/read-token", a.CallBackToken)
	}
}

// CallBackToken controller 微软回调本机接口
func (a *AzureController) CallBackToken(c *gin.Context) {
	var _payLoad string
	a.logger.Info("微软回调本机接口 获取token")
	idToken := c.PostForm("id_token")
	state := c.PostForm("state")
	a.logger.Info("获取到的参数", zap.String("id-token", idToken), zap.String("state", state))

	IDTokens := strings.Split(idToken, ".")
	if len(IDTokens) == 3 {
		_payLoad = IDTokens[1]
	} else {
		a.logger.Error("id-token格式不符合预期")
		C.RespSuccess(c, []map[string]string{}, C.Page{Total: 0, Page: 1, Size: 10})
		return
	}

	payLoad, err := base64.RawURLEncoding.DecodeString(_payLoad)
	if err != nil {
		a.logger.Error("base64.RawURLEncoding.DecodeString(payLoad)", zap.Error(err))
		C.RespSuccess(c, []map[string]string{}, C.Page{Total: 0, Page: 1, Size: 10})
		return
	}

	var payloadMap map[string]interface{}
	err = json.Unmarshal(payLoad, &payloadMap)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		a.logger.Error("解析payload失败.json.Unmarshal")

		C.RespSuccess(c, []map[string]string{}, C.Page{Total: 0, Page: 1, Size: 10})
		return
	}

	a.logger.Info("payload", zap.Any("payload", payloadMap))
	email, ok := payloadMap["email"]
	if !ok {
		a.logger.Error("id-token中没有email这个key")
		C.RespSuccess(c, []map[string]string{}, C.Page{Total: 0, Page: 1, Size: 10})
		return
	}

	session := sessions.Default(c)
	originalURL := session.Get("original_url")
	stateSession := session.Get("state")

	if state != stateSession {
		a.logger.Error("微软跳转回来携带的state不是本次传递过去的state")
		c.Redirect(http.StatusMovedPermanently, "/login")
	}

	userInfo, err := a.JwtGetInfo(c.Request.Context(), email.(string))
	if err != nil {
		a.logger.Error("JwtGetInfo", zap.Any("email", email), zap.Error(err))
		c.Redirect(http.StatusMovedPermanently, "/login")
	}

	a.logger.Info("JwtGetInfo:", zap.Any("UserInfo", userInfo))
	// 登录成功,把userid放入到session中
	session.Set("user_id", userInfo.ID)
	session.Delete("state")

	c.Redirect(http.StatusMovedPermanently, originalURL.(string))
}

// AzureLogin controller 前端调用
func (a *AzureController) AzureLogin(c *gin.Context) {
	originalURL, ok := c.GetQuery("initial_url") // 前端请求地址

	if !ok {
		a.logger.Error("获取跳转URL失败")
		c.Redirect(http.StatusMovedPermanently, "/login")
		return
	}

	aq := newAzureRequest(originalURL)

	params := aq.toURLParams()
	a.logger.Info("参数:", zap.String("params", params))

	requestURL := fmt.Sprintf("%s/%s/oauth2/v2.0/authorize?%s",
		viper.GetString("aad.authority"),
		viper.GetString("aad.tenant_id"),
		params,
	)

	session := sessions.Default(c)
	session.Set("original_url", originalURL)
	session.Set("state", aq.State)

	a.logger.Info("跳转到微软的地址", zap.String("requestURL", requestURL))
	c.Redirect(http.StatusMovedPermanently, requestURL)
}

// help function
func (ar *azureRequest) toURLParams() string {
	v := neturl.Values{}
	v.Set("client_id", ar.ClientID)
	v.Set("response_type", ar.ResponseType)
	v.Set("redirect_uri", ar.RedirectURI)
	v.Set("scope", ar.Scope)
	v.Set("response_mode", ar.ResponseMode)
	v.Set("state", ar.State)
	v.Set("nonce", ar.Nonce)
	v.Set("prompt", ar.Prompt)
	return v.Encode()
}

func (a *AzureController) JwtGetInfo(ctx context.Context, email string) (domain.UserInfo, error) {
	var emailKey ContextKey = "email"
	return a.service.Select(context.WithValue(ctx, emailKey, email))
}

func generateSecureRandomString(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	for i, v := range b {
		b[i] = charset[v%byte(len(charset))]
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
