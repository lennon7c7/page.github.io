package middleware

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/log"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"page.github.io/pkg/util"
)

var e *casbin.Enforcer

func CasbinSetup(db *gorm.DB, modelText string) {
	m, _ := model.NewModelFromString(modelText)
	a, _ := gormadapter.NewAdapterByDBUseTableName(db, "", "casbin_rule")
	e, _ := casbin.NewEnforcer(m, a)

	logger := &log.DefaultLogger{}
	logger.EnableLog(true)
	logger.IsEnabled()
	e.SetLogger(logger)

	err := e.LoadPolicy()
	if err != nil {
		fmt.Println(err)
	}

	//// Check the permission.
	//fmt.Println(e.Enforce("alice", "data1", "read"))
	//
	//// Modify the policy.
	//e.AddPolicy("qlp", "book_policy", "read")
	////e.RemovePolicy(...)
	//e.AddRoleForUser("qlp", "admin_role", "com.lninl")
	//e.AddRoleForUser("cjz", "editor_role", "com.lninl")
	//e.AddPermissionForUser("qlp", "book_permission", "write")
	//
	//fmt.Println("\"admin\", \"data1\", \"read\"")
	//fmt.Println(e.Enforce("admin", "data1", "read"))
	//fmt.Println(e.Enforce("alice", "admin", "read"))
	//fmt.Println(e.GetAllRoles())
	//fmt.Println(e.GetAllActions())
	//fmt.Println(e.GetAllObjects())
	//fmt.Println(e.GetGroupingPolicy())
	//// Save the policy back to DB.
	////e.SavePolicy()
	//e.LoadPolicy()
}

// AuthCheck 认证验证
func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.Keys[util.JWTClaims]
		result, err := e.Enforce(role, c.Request.URL.Path, c.Request.Method)
		if err != nil {
			util.ErrorSystem(c, err.Error())
			return
		}
		if !result {
			util.Resp(c, util.CodeErrorForbidden, "无权限访问", "")
			return
		}
		c.Next()
	}
}
