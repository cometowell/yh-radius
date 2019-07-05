package web

import (
	"github.com/gin-gonic/gin"
	"go-rad/common"
	"go-rad/database"
	"go-rad/model"
	"net/http"
)

// -------------------------- manager start ----------------------------
func sessionManagerInfo(c *gin.Context) {
	token := c.GetHeader(common.SessionName)
	session := GlobalSessionManager.GetSession(token)

	managerInfo := session.GetAttr("manager")
	resources := session.GetAttr("resources").([]model.SysResource)

	buttonPermissions := make([]int64, 0)
	for _, res := range resources {
		if res.Level == 3 {
			buttonPermissions = append(buttonPermissions, res.Id)
		}
	}

	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "success", Data: gin.H{
		"manager":           managerInfo,
		"buttonPermissions": buttonPermissions,
	}})
}

func managerList(c *gin.Context) {
	var params model.SysUser
	err := c.ShouldBindJSON(&params)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	c.Set("current", params.Page)
	c.Set("pageSize", params.PageSize)
	var managers []model.SysUser
	whereSql := "1=1 "
	whereArgs := make([]interface{}, 0)
	if params.Username != "" {
		whereSql += "and username like ? "
		whereArgs = append(whereArgs, "%"+params.Username+"%")
	}

	if params.RealName != "" {
		whereSql += "and real_name like ? "
		whereArgs = append(whereArgs, "%"+params.RealName+"%")
	}

	if params.Status != 0 {
		whereSql += "and status = ?"
		whereArgs = append(whereArgs, params.Status)
	}

	PageByWhereSql(c, &managers, whereSql, whereArgs)
}

func managerById(c *gin.Context) {
	var manager model.SysUser
	err := c.ShouldBindJSON(&manager)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	database.DataBaseEngine.Id(manager.Id).Get(&manager)
	manager.Password = ""
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "success", Data: manager})
}

func addManager(c *gin.Context) {
	var manager model.SysUser
	err := c.ShouldBindJSON(&manager)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	count, _ := database.DataBaseEngine.Table(&model.SysUser{}).Where("username=?", manager.Username).Count()
	if count > 0 {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: "用户名已存在!"})
		return
	}
	manager.Status = 1
	manager.Password = common.Encrypt(manager.Password)
	manager.CreateTime = model.NowTime()
	database.DataBaseEngine.InsertOne(&manager)
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "管理员信息添加成功!"})
}

func updateManager(c *gin.Context) {
	var manager model.SysUser
	err := c.ShouldBindJSON(&manager)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	if manager.Password != "" {
		manager.Password = common.Encrypt(manager.Password)
	}

	count, _ := database.DataBaseEngine.Table(&model.SysUser{}).Where("username=? and id != ?", manager.Username, manager.Id).Count()
	if count > 0 {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: "用户名已存在!"})
		return
	}
	manager.UpdateTime = model.NowTime()
	database.DataBaseEngine.Id(manager.Id).Update(&manager)
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "管理员信息更新成功!"})
}

func delManager(c *gin.Context) {
	var manager model.SysUser
	err := c.ShouldBindJSON(&manager)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	manager.Status = 3 // 标记为已删除
	manager.UpdateTime = model.NowTime()
	database.DataBaseEngine.Id(manager.Id).Cols("status").Update(&manager)
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "删除成功!"})
}

func changeManagerPassword(c *gin.Context) {
	var managerPassword model.SysUserPassword
	err := c.ShouldBindJSON(&managerPassword)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}

	var manager model.SysUser
	manager.Password = common.Encrypt(managerPassword.NewPassword)
	_, err = database.DataBaseEngine.ID(managerPassword.Id).Cols("password").Update(&manager)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "密码修改成功!"})
}

// -------------------------- manager end ----------------------------
