package util

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	Version = "1.0.0"

	CodeOK                = 20000
	CodeErrorBusiness     = 30000
	CodeErrorForm         = 40000
	CodeErrorForbidden    = 40403
	CodeErrorSystem       = 50000
	CodeErrorToken        = 50008
	CodeErrorTokenExpired = 50014
)

// Response 返回
type Response struct {
	Code int         `json:"code"` // 错误码：20000正常,30000业务异常,40000表单异常,40403无权限,50000系统异常,50008Token异常,50014Token过期
	Msg  string      `json:"msg"`  // 错误信息
	Data interface{} `json:"data"` // 正常返回内容
}

// List 列表
type List struct {
	Total int64       `json:"total"` // 总页码
	Items interface{} `json:"items"` // 项目列表
}

// ListTotal 单列表多总数量
type ListTotal struct {
	Total interface{} `json:"total"` // 总数量
	Items interface{} `json:"items"` // 项目列表
}

// 返回
func Resp(c *gin.Context, code int, msg string, data interface{}) {
	c.Header("Version", Version)
	c.JSONP(http.StatusOK, Response{Code: code, Msg: msg, Data: data})
	c.Abort()
}

// OK 正常返回
func OK(c *gin.Context, msg string, data interface{}) {
	c.Header("Version", Version)
	c.JSONP(http.StatusOK, Response{Code: CodeOK, Msg: msg, Data: data})
	c.Abort()
}

// OK 正常返回，不返回提示信息，也不返回数据，只返回正常状态码
func OKToo(c *gin.Context) {
	c.Header("Version", Version)
	c.JSONP(http.StatusOK, Response{Code: CodeOK})
	c.Abort()
}

// OK 正常返回，不返回数据，只返回正常状态码、提示信息
func OKMsg(c *gin.Context, msg string) {
	c.Header("Version", Version)
	c.JSONP(http.StatusOK, Response{Code: CodeOK, Msg: msg})
	c.Abort()
}

// OK 正常返回，不返回提示信息，只返回正常状态码、data
func OKData(c *gin.Context, data interface{}) {
	c.Header("Version", Version)
	c.JSONP(http.StatusOK, Response{Code: CodeOK, Data: data})
	c.Abort()
}

// ErrorSystem 系统错误，比如数据库链接错误等
func ErrorSystem(c *gin.Context, err string) {
	c.Header("Version", Version)
	c.JSONP(http.StatusOK, Response{Code: CodeErrorSystem, Msg: err})
	c.Abort()
}

// Error 系统错误，比如数据库链接错误等
func Error(c *gin.Context, msg string, data interface{}) {
	c.Header("Version", Version)
	c.JSONP(http.StatusOK, Response{Code: CodeErrorSystem, Msg: msg, Data: data})
	c.Abort()
}

// ErrorForm 表单错误返回，data一般包含错误项
func ErrorForm(c *gin.Context, data interface{}) {
	c.Header("Version", Version)
	c.JSONP(http.StatusOK, Response{Code: CodeErrorForm, Data: data})
	c.Abort()
}

// ErrorBusiness 业务错误返回，msg中包含业务错误内容
func ErrorBusiness(c *gin.Context, msg string) {
	c.Header("Version", Version)
	c.JSONP(http.StatusOK, Response{Code: CodeErrorBusiness, Msg: msg})
	c.Abort()
}
