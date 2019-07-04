package web

import (
	"github.com/gin-gonic/gin"
	"go-rad/common"
	"go-rad/database"
	"go-rad/model"
	"net/http"
)

func loadControllers(router *gin.Engine) {
	router.POST("/login", login)
	router.POST("/logout", logout)

	router.POST("/system/user/session/info", sessionManagerInfo)
	router.POST("/system/user/list", managerList)
	router.POST("/system/user/info", managerById)
	router.POST("/system/user/add", addManager)
	router.POST("/system/user/update", updateManager)
	router.POST("/system/user/delete", delManager)
	router.POST("/system/user/change/password", changeManagerPassword)

	router.POST("/product/add", addProduct)
	router.POST("/product/list", listProduct)
	router.POST("/product/update", updateProduct)
	router.POST("/fetch/product", fetchProductList)
	router.POST("/product/info", getProductInfo)
	router.POST("/product/delete", deleteProduct)

	router.POST("/user/list", listUser)
	router.POST("/user/add", addUser)
	router.POST("/user/info", fetchUser)
	router.POST("/user/update", updateUser)
	router.POST("/user/order/record", fetchUserOrderRecord)
	router.POST("/user/continue", continueProduct)
	router.POST("/user/delete", deleteUser)

	router.POST("/fetch/department", fetchDepartments)
	router.POST("/department/info", getDepartmentInfo)
	router.POST("/department/list", listDepartments)
	router.POST("/department/add", addDepartment)
	router.POST("/department/update", updateDepartment)
	router.POST("/department/delete", deleteDepartment)

	router.POST("/nas/info", getNasInfo)
	router.POST("/nas/list", listNas)
	router.POST("/nas/add", addNas)
	router.POST("/nas/update", updateNas)
	router.POST("/nas/delete", deleteNas)

	router.POST("/resource/list", listRes)
	router.POST("/session/resource", getSessionResource)

	router.POST("/role/info", getRoleInfo)
	router.POST("/role/list", listRole)
	router.POST("/role/add", addRole)
	router.POST("/role/update", updateRole)
	router.POST("/role/delete", deleteRole)
	router.POST("/role/resources", empowerRole)
	router.POST("/role/empower/:roleId", doEmpowerRole)

	router.POST("/online/list", listOnline)
	router.POST("/online/off", offOnline)
	router.POST("/online/delete", deleteOnline)

}

func login(c *gin.Context) {
	var manager model.SysUser
	err := c.ShouldBindJSON(&manager)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	manager.Password = common.Encrypt(manager.Password)
	ok, _ := database.DataBaseEngine.Get(&manager)

	if !ok {
		c.JSON(http.StatusOK, common.NewErrorJsonResult("username or password is incorrect"))
		return
	}

	session := GlobalSessionManager.CreateSession(c)
	manager.Password = ""
	session.SetAttr("manager", manager)

	var resources []model.SysResource
	err = database.DataBaseEngine.Table("sys_resource").Alias("sr").
		Join("LEFT", []string{"sys_role_resource_rel", "srr"}, "sr.id = srr.resource_id").
		Join("LEFT", []string{"sys_role", "r"}, "srr.role_id = r.id").
		Join("LEFT", []string{"sys_user_role_rel", "smr"}, "smr.role_id = r.id").
		Join("LEFT", []string{"sys_user", "m"}, "smr.sys_user_id = m.id").
		Where("m.id = ? or sr.should_perm_control = 0", manager.Id).
		Find(&resources)

	if err != nil {
		c.JSON(http.StatusOK, common.NewErrorJsonResult("find user error: "+err.Error()))
		return
	}

	if len(resources) == 0 {
		c.JSON(http.StatusForbidden, common.NewErrorJsonResult("user's permission is poor"))
		return
	}

	session.SetAttr("resources", resources)

	buttonPermissions := make([]int64, 0)
	for _, res := range resources {
		if res.Level == 3 {
			buttonPermissions = append(buttonPermissions, res.Id)
		}
	}

	c.JSON(http.StatusOK, common.NewSuccessJsonResult("success", gin.H{
		"sessionId":         session.SessionId(),
		"buttonPermissions": buttonPermissions,
	}))
}

func logout(c *gin.Context) {
	GlobalSessionManager.DestroySession(c)
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "success"})
}
