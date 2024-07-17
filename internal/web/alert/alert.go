package alert

import (
	"encoding/json"

	"eyes/internal/domain"
	"eyes/internal/service/alert"
	C "eyes/internal/web/common"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MonitorController struct {
	service alert.MonitorService
	logger  *zap.Logger
}

func NewMonitorController(service alert.MonitorService, logger *zap.Logger) *MonitorController {
	return &MonitorController{
		service: service,
		logger:  logger,
	}
}

func (a *MonitorController) RegisterRoutes(server *gin.Engine) {
	api := server.Group("/v1/api/alert")

	api.GET("/count", a.Count)
	api.GET("/detail", a.Detail)
	api.GET("/list", a.List)
	api.POST("/query", a.Query)
}

func (a *MonitorController) Count(c *gin.Context) {
	if params, ok := c.Get("_params"); ok {
		a.logger.Info("获取params: ", zap.String("x-request-id", c.GetHeader("X-Request-ID")), zap.Any("params", params), zap.String("x-request-id", c.GetHeader("X-Request-ID")))
	}

	var _monitor domain.Monitor

	err := c.Bind(&_monitor)
	if err != nil {
		a.logger.Error("绑定错误", zap.String("x-request-id", c.GetHeader("X-Request-ID")), zap.Error(err), zap.String("x-request-id", c.GetHeader("X-Request-ID")))
	}

	cnt, err := a.service.Count(c.Request.Context(), _monitor)
	if err != nil {
		a.logger.Error("统计总数:", zap.String("x-request-id", c.GetHeader("X-Request-ID")), zap.Error(err), zap.String("x-request-id", c.GetHeader("X-Request-ID")))
		C.RespErrorWithMsg(c, 40001, "索引不存在")
		return
	}

	C.RespSuccess(c, []gin.H{{"count": cnt}}, C.Page{Total: 1, Page: 1, Size: 10})
	return
}

func (a *MonitorController) Detail(c *gin.Context) {
	if params, ok := c.Get("_params"); ok {
		// a.logger.Info("获取params: ", zap.Any("params", params), zap.String("x-request-id", c.GetHeader("X-Request-ID")))
		a.logger.Info("获取params", zap.String("x-request-id", c.GetHeader("X-Request-ID")), zap.Any("params", params), zap.String("x-request-id", c.GetHeader("X-Request-ID")))
	}

	ID := c.Request.URL.Query().Get("id")

	alertDetail, err := a.service.Detail(c.Request.Context(), domain.Monitor{ID: ID})
	if err != nil {
		a.logger.Error("获取alertDetail异常:", zap.String("x-request-id", c.GetHeader("X-Request-ID")), zap.Error(err), zap.String("x-request-id", c.GetHeader("X-Request-ID")))
		C.RespErrorWithMsg(c, 404, "未查询到数据")
		return
	}

	C.Response(c, []*domain.Monitor{alertDetail}, C.Page{Total: 1, Page: 1, Size: 10})
	return
}

func (a *MonitorController) List(ctx *gin.Context) {
	if params, ok := ctx.Get("_params"); ok {
		a.logger.Info("获取params: ", zap.String("x-request-id", ctx.GetHeader("X-Request-ID")), zap.Any("params", params), zap.String("x-request-id", ctx.GetHeader("X-Request-ID")))
	}

	page, size := C.GetPageInfo(ctx)
	monitorList, total, err := a.service.List(ctx.Request.Context(), page, size, domain.Monitor{})
	if err != nil {
		a.logger.Error("alert.List:", zap.String("x-request-id", ctx.GetHeader("X-Request-ID")), zap.Error(err))
		C.RespErrorWithMsg(ctx, 404, "未查询到数据")
		return
	}

	C.Response(ctx, monitorList, C.Page{Page: page, Size: size, Total: total})
	return
}

func (a *MonitorController) Query(c *gin.Context) {
	if params, ok := c.Get("_params"); ok {
		a.logger.Info("获取params:", zap.String("x-request-id", c.GetHeader("X-Request-ID")), zap.Any("params", params))
	}

	page, size := C.GetPageInfo(c)

	reqParams := make(map[string]any)
	reqParams["from"] = (page - 1) * size
	reqParams["size"] = size

	filterMust := make([]map[string]any, 0)
	// 源端口
	srcPort, ok := c.GetPostForm("src_port")
	if ok && srcPort != "" {
		filterMust = append(filterMust, map[string]any{"match": map[string]string{"src_port": srcPort}})
	}

	// 目的端口
	destPort, ok := c.GetPostForm("dest_port")
	if ok && destPort != "" {
		filterMust = append(filterMust, map[string]any{"match": map[string]string{"dest_port": destPort}})
	}

	// 协议
	proto, ok := c.GetPostForm("proto")
	if ok && proto != "" {
		filterMust = append(filterMust, map[string]any{"match": map[string]string{"proto": proto}})
	}

	Action, ok := c.GetPostForm("action")
	if ok && Action != "" {
		filterMust = append(filterMust, map[string]any{"match": map[string]string{"alert.action": Action}})
	}

	reqParams["query"] = map[string]any{"bool": map[string]any{"must": filterMust}}
	filter, err := json.Marshal(reqParams)
	if err != nil {
		a.logger.Error("json.Marshal(reqParams):", zap.String("x-request-id", c.GetHeader("X-Request-ID")), zap.Error(err))
		C.RespErrorWithMsg(c, 400, "参数错误")
		return
	}

	a.logger.Info("ES查询:", zap.String("x-request-id", c.GetHeader("X-Request-ID")), zap.String("filter", string(filter)))

	alertList, total, err := a.service.Query(c.Request.Context(), page, size, string(filter), domain.Monitor{})
	if err != nil {
		a.logger.Error("alert.Query(string(filter))", zap.String("x-request-id", c.GetHeader("X-Request-ID")), zap.Error(err))

		C.RespErrorWithMsg(c, 404, "未查询到数据")
		return
	}

	C.Response(c, alertList, C.Page{Page: page, Size: size, Total: total})
	return
}
