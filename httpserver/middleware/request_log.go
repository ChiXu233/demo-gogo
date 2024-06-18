package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"io/ioutil"
	"strings"
)

var (
	IgnorePATH = []string{
		"api/upload",
		"api/point_cloud",
		"api/big_file",
		"api/nerf_camera",
		"api/nerf_model/shot",
		"api/images/image_upload",
	}
)

// RequestInfo 添加接口调用前日志
func RequestInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Next()
		method := c.Request.Method
		uri := c.Request.RequestURI
		contentType := c.Request.Header.Get("Content-Type")
		clientIP := c.ClientIP()
		// mock请求body为nil，实际接口请求不为nil
		if c.Request.Body == nil || ignoreUri(uri) || strings.Contains(contentType, "multipart/form-data") {
			logger.Debug("RequestInfo context: ClientIP:[%s] Method:[%s] URI:[%s]", clientIP, method, uri)
		} else {
			if data, err := ioutil.ReadAll(c.Request.Body); err == nil {
				// 新建缓冲区并替换原有Request.body
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
				logger.Debug("RequestInfo context: ClientIP:[%s] Method:[%s] URI:[%s] RequestBody:[%+v]", clientIP, method, uri, string(data))
			} else {
				logger.Debug("RequestInfo context: ClientIP:[%s] Method:[%s] URI:[%s] GetRequestBodyError:[%+v]", clientIP, method, uri, err)
			}
		}
	}
}

func ignoreUri(uri string) bool {
	for _, item := range IgnorePATH {
		if strings.Contains(uri, item) {
			return true
		}
	}
	return false
}
