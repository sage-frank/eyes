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

// NewAzureController creates a new AzureController instance
func NewAzureController(service aad.AzureService, logger *zap.Logger) *AzureController {
	return &AzureController{
		service: service,
		logger:  logger,
	}
}

// newAzureRequest creates a new azureRequest with the given URL
func newAzureRequest(url string) (*azureRequest, error) {
	state, err := generateSecureRandomString(6)
	if err != nil {
		return nil, err
	}
	nonce, err := generateSecureRandomString(6)
	if err != nil {
		return nil, err
	}

	return &azureRequest{
		ClientID:     viper.GetString("aad.client_id"),
		ResponseType: "id_token",
		RedirectURI:  url,
		Scope:        neturl.QueryEscape(viper.GetString("aad.scope")),
		ResponseMode: "form_post",
		State:        state,
		Nonce:        nonce,
		Prompt:       "select_account",
	}, nil
}

// RegisterRoutes registers the routes for AzureController
func (a *AzureController) RegisterRoutes(server *gin.Engine) {
	azure := server.Group("/api/v1/azure")
	{
		azure.POST("/read-token", a.CallBackToken)
		azure.GET("/read-token", a.CallBackToken)
	}
}

func isAllowed(url string, allowedList []string) bool {
	for _, allowed := range allowedList {
		if url == allowed {
			return true
		}
	}
	return false
}

// CallBackToken handles the callback from Azure and processes the token
func (a *AzureController) CallBackToken(c *gin.Context) {
	a.logger.Info("微软回调本机接口 获取token")
	idToken := c.PostForm("id_token")
	state := c.PostForm("state")
	a.logger.Info("获取到的参数", zap.String("id-token", idToken), zap.String("state", state))

	// 检查 Referer 和 Origin
	referer := c.Request.Referer()
	origin := c.Request.Header.Get("Origin")
	a.logger.Info("Referer 和 Origin", zap.String("Referer", referer), zap.String("Origin", origin))

	// 获取预期的 Referer 和 Origin 列表
	expectedReferrers := viper.GetStringSlice("aad.expected_referers")
	expectedOrigins := viper.GetStringSlice("aad.expected_origins")

	// 检查 Referer 和 Origin 是否在允许列表中
	if !isAllowed(referer, expectedReferrers) || !isAllowed(origin, expectedOrigins) {
		a.logger.Error("非法的 Referer 或 Origin", zap.String("referer", referer), zap.String("origin", origin))
		c.Redirect(http.StatusMovedPermanently, "/login")
		return
	}

	payload, err := decodeIDToken(idToken)
	if err != nil {
		a.logger.Error("解析ID Token失败", zap.Error(err))
		C.RespSuccess(c, []map[string]string{}, C.Page{Total: 0, Page: 1, Size: 10})
		return
	}

	a.logger.Info("payload", zap.Any("payload", payload))
	email, ok := payload["email"]
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
		return
	}

	userInfo, err := a.JwtGetInfo(c.Request.Context(), email.(string))
	if err != nil {
		a.logger.Error("JwtGetInfo", zap.Any("email", email), zap.Error(err))
		c.Redirect(http.StatusMovedPermanently, "/login")
		return
	}

	a.logger.Info("JwtGetInfo:", zap.Any("UserInfo", userInfo))
	session.Set("user_id", userInfo.ID)
	session.Delete("state")
	session.Save()

	c.Redirect(http.StatusMovedPermanently, originalURL.(string))
}

// AzureLogin handles the login request and redirects to Azure for authentication
func (a *AzureController) AzureLogin(c *gin.Context) {
	originalURL, ok := c.GetQuery("initial_url")
	if !ok {
		a.logger.Error("获取跳转URL失败")
		c.Redirect(http.StatusMovedPermanently, "/login")
		return
	}

	aq, err := newAzureRequest(originalURL)
	if err != nil {
		a.logger.Error("创建Azure请求失败", zap.Error(err))
		c.Redirect(http.StatusMovedPermanently, "/login")
		return
	}

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
	session.Save()

	a.logger.Info("跳转到微软的地址", zap.String("requestURL", requestURL))
	c.Redirect(http.StatusMovedPermanently, requestURL)
}

// decodeIDToken decodes the ID token and returns the payload
func decodeIDToken(idToken string) (map[string]interface{}, error) {
	IDTokens := strings.Split(idToken, ".")
	if len(IDTokens) != 3 {
		return nil, fmt.Errorf("id-token格式不符合预期")
	}

	payload, err := base64.RawURLEncoding.DecodeString(IDTokens[1])
	if err != nil {
		return nil, fmt.Errorf("base64.RawURLEncoding.DecodeString: %w", err)
	}

	var payloadMap map[string]interface{}
	if err := json.Unmarshal(payload, &payloadMap); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return payloadMap, nil
}

// toURLParams converts azureRequest to URL parameters
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

// JwtGetInfo retrieves user information based on the email from the JWT token
func (a *AzureController) JwtGetInfo(ctx context.Context, email string) (domain.UserInfo, error) {
	var emailKey ContextKey = "email"
	return a.service.Select(context.WithValue(ctx, emailKey, email))
}

// generateSecureRandomString generates a secure random string of the given length
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
