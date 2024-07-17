package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
{
	"code": 10000, // 程序中的错误码
	"msg": xx,     // 提示信息
	"data": {},    // 数据
	"page": {"total":0, "size":10, "page":1}
}
*/

type (
	Page struct {
		Total int64 `json:"total"`
		Page  int64 `json:"page"`
		Size  int64 `json:"size"`
	}

	RespData struct {
		Code ResCode `json:"code"`
		Msg  any     `json:"msg"`
		Data any     `json:"data"`
		Page Page    `json:"page"`
	}
)

func RespError(c *gin.Context, code ResCode) {
	c.JSON(http.StatusOK, &RespData{
		Code: code,
		Msg:  code.Msg(),
		Data: []string{},
		Page: Page{Total: 0, Page: 0, Size: 10},
	})
}

func RespErrorWithMsg(c *gin.Context, code ResCode, msg any) {
	c.JSON(http.StatusOK, &RespData{
		Code: code,
		Msg:  msg,
		Data: []string{},
		Page: Page{Total: 0, Page: 0, Size: 10},
	})
}

func RespSuccess(c *gin.Context, data any, page Page) {
	c.JSON(http.StatusOK, &RespData{
		Code: CodeSuccess,
		Msg:  CodeSuccess.Msg(),
		Data: data,
		Page: page,
	})
}

func Response(c *gin.Context, data any, page Page) {
	c.JSON(http.StatusOK, &RespData{
		Code: CodeSuccess,
		Msg:  CodeSuccess.Msg(),
		Data: data,
		Page: page,
	})
}
