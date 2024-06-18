package app

import (
	"demo-gogo/httpserver/errcode"
	"github.com/gin-gonic/gin"
	log "github.com/wonderivan/logger"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// HandleNotFound 404处理
func HandleNotFound(c *gin.Context) {
	SendNotFoundErrorResponse(c, errcode.ErrorMsgHandleNotFound)
}

// MethodNotFound 404处理
func MethodNotFound(c *gin.Context) {
	SendNotFoundErrorResponse(c, errcode.ErrorMsgMethodNotFound)
}

// Success 请求成功，默认返回HTTP 200 + Business 0
func Success(c *gin.Context, data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	c.JSON(http.StatusOK, Response{
		Code:    errcode.SuccessCodeBusiness,
		Message: errcode.SuccessMsgBusiness,
		Data:    data,
	})
}

func Error(c *gin.Context, code int, message string, data interface{}) {
	businessCode := errcode.GetErrorCode(message)
	log.Error("Request [%s] [%s] Error.HttpCode[%d] BusinessCode[%d] Message[%s] ErrDetail[%v]", c.Request.Method, c.Request.RequestURI, code, businessCode, message, data)
	c.JSON(code, Response{
		Code:    businessCode,
		Message: message,
		Data:    gin.H{},
	})
}

// SendAuthorizedErrorResponse 鉴权有无
func SendAuthorizedErrorResponse(c *gin.Context, msg string) {
	Error(c, errcode.ErrorCodeUnauthorized, msg, nil)
}

// SendParameterErrorResponse 非法用户输入
func SendParameterErrorResponse(c *gin.Context, msg string) {
	if msg != "" || len(msg) > 0 {
		Error(c, http.StatusOK, msg, nil)
	} else {
		Error(c, http.StatusOK, errcode.ErrorMsgPrefixInvalidParameter, nil)
	}
}

// SendServerErrorResponse 服务端错误
func SendServerErrorResponse(c *gin.Context, msg string, err error) {
	if msg != "" || len(msg) > 0 {
		// 绕过errMsg收口，查询具体的业务错误信息，如果是预估的，取消收口
		if err == nil || errcode.GetErrorCode(err.Error()) == errcode.ErrorCodeBusiness {
			Error(c, http.StatusOK, msg, err)
			return
		}
		Error(c, http.StatusOK, err.Error(), err)
	} else {
		Error(c, http.StatusOK, errcode.ErrorMsgInternal, err)
	}
}

// SendNotFoundErrorResponse 路由错误
func SendNotFoundErrorResponse(c *gin.Context, msg string) {
	if msg != "" || len(msg) > 0 {
		Error(c, http.StatusOK, msg, nil)
	} else {
		Error(c, http.StatusOK, errcode.ErrorMsgNotfound, nil)
	}
}
