package util

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strconv"
)

// Request 请求对象
type Request struct {
	Page   int                    `json:"page"`   // 页码从1开始
	Size   int                    `json:"size"`   // 每页项目数
	Sort   string                 `json:"sort"`   // 排序
	With   string                 `json:"with"`   // 需要额外加载的字段
	UserID int64                  `json:"userID"` // 用户ID
	Params map[string]interface{} `json:"params"` // 所有参数
}

type IDList struct {
	ID []int64 // ID列表
}

// Offset 获取分页偏移量
func (request *Request) Offset() int {
	return (request.Page - 1) * request.Size
}

func (request *Request) Value(key string) interface{} {
	return request.Params[key]
}

// NewDefaultRequest 返回默认Request对象
func NewDefaultRequest() *Request {
	return &Request{
		Page: 1,
		Size: 20,
	}
}

func BindPage(c *gin.Context) (request *Request) {
	request = NewDefaultRequest()
	_ = c.ShouldBind(&request)
	if Page := c.Query("page"); Page != "" {
		request.Page, _ = strconv.Atoi(Page)
	} else if request.Params["page"] != nil {
		request.Page = request.Params["page"].(int)
	}
	if Size := c.Query("size"); Size != "" {
		request.Size, _ = strconv.Atoi(Size)
	} else if request.Params["size"] != nil {
		request.Size = request.Params["size"].(int)
	}
	if Sort := c.Query("sort"); Sort != "" {
		request.Sort = Sort
	} else if request.Params["sort"] != nil {
		request.Sort = request.Params["sort"].(string)
	}
	if value, exists := c.Get(JWTClaims); exists {
		userClaims := value.(*jwt.StandardClaims)
		request.UserID, _ = strconv.ParseInt(userClaims.Issuer, 10, 0)
	}
	return
}

func BindParams(c *gin.Context) (request *Request) {
	request = NewDefaultRequest()
	_ = c.ShouldBindJSON(&request.Params)
	if Page := c.Query("page"); Page != "" {
		request.Page, _ = strconv.Atoi(Page)
	} else if request.Params["page"] != nil {
		request.Page = int(request.Params["page"].(float64))
	}
	if Size := c.Query("size"); Size != "" {
		request.Size, _ = strconv.Atoi(Size)
	} else if request.Params["size"] != nil {
		request.Size = int(request.Params["size"].(float64))
	}
	if Sort := c.Query("sort"); Sort != "" {
		request.Sort = Sort
	} else if request.Params["sort"] != nil {
		request.Sort = request.Params["sort"].(string)
	}
	if value, exists := c.Get(JWTClaims); exists {
		userClaims := value.(*jwt.StandardClaims)
		request.UserID, _ = strconv.ParseInt(userClaims.Issuer, 10, 0)
	}
	return
}
