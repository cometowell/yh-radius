package web

import (
	"github.com/gin-gonic/gin"
	"go-rad/common"
	"go-rad/database"
	"go-rad/model"
	"net/http"
)

// -------------------------- role start -------------------------------

func getRoleInfo(c *gin.Context) {
	var role model.SysRole
	err := c.ShouldBindJSON(&role)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	database.DataBaseEngine.Id(role.Id).Get(&role)
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "success", Data: role})
}

func addRole(c *gin.Context) {
	var role model.SysRole
	err := c.ShouldBindJSON(&role)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	session := database.DataBaseEngine.NewSession()
	defer session.Close()
	session.Begin()
	count, _ := session.Table("sys_role").Where("code = ?", role.Code).Count()
	if count > 0 {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: "错误：编码已存在，不能重复!"})
		session.Rollback()
		return
	}
	role.CreateTime = model.NowTime()
	role.Enable = 1
	session.InsertOne(&role)
	session.Commit()
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "添加成功"})
}

func updateRole(c *gin.Context) {
	var role model.SysRole
	err := c.ShouldBindJSON(&role)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	session := database.DataBaseEngine.NewSession()
	defer session.Close()
	session.Begin()
	count, _ := session.Table("sys_role").Where("code = ? and id != ?", role.Code, role.Id).Count()
	if count > 0 {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: "错误：编码重复!"})
		session.Rollback()
		return
	}
	session.ID(role.Id).Update(&role)
	session.Commit()
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "修改成功"})
}

func listRole(c *gin.Context) {
	var role model.SysRole
	err := c.ShouldBindJSON(&role)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	c.Set("current", role.Page)
	c.Set("pageSize", role.PageSize)
	var roleList []model.SysRole

	whereSql := "1=1 "
	whereArgs := make([]interface{}, 0)
	if role.Name != "" {
		whereSql += "and name like ? "
		whereArgs = append(whereArgs, "%"+role.Name+"%")
	}

	if role.Code != "" {
		whereSql += "and code like ? "
		whereArgs = append(whereArgs, "%"+role.Code+"%")
	}
	PageByWhereSql(c, &roleList, whereSql, whereArgs)
}

func deleteRole(c *gin.Context) {
	var role model.SysRole
	err := c.ShouldBindJSON(&role)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	database.DataBaseEngine.Id(role.Id).Delete(&role)
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "已删除"})
}

// 角色赋权
func empowerRole(c *gin.Context) {
	var role model.SysRole
	err := c.ShouldBindJSON(&role)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	var selectedResList []model.SysResource
	database.DataBaseEngine.Table("sys_resource").Alias("sr").
		Join("INNER", []string{"sys_role_resource_rel", "srr"}, "sr.id = srr.resource_id").
		Join("INNER", []string{"sys_role", "r"}, "srr.role_id = r.id").
		Where("r.id = ?", role.Id).
		Find(&selectedResList)

	var resources []model.SysResource
	database.DataBaseEngine.Where("should_perm_control = ?", 1).Find(&resources)
	for index, item := range resources {
		for _, r := range selectedResList {
			if r.Id == item.Id {
				resources[index].Selected = true
			}
		}
	}
	c.JSON(http.StatusOK, common.DefaultSuccessJsonResult(getResLevel(resources)))
}

func doEmpowerRole(c *gin.Context) {
	roleId := c.Param("roleId")
	var roleResourceRels []model.SysRoleResourceRel
	err := c.ShouldBindJSON(&roleResourceRels)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	session := database.DataBaseEngine.NewSession()
	defer session.Close()
	session.Begin()
	session.Where("role_id = ?", roleId).Delete(&model.SysRoleResourceRel{})
	if roleResourceRels != nil && len(roleResourceRels) > 0 {
		session.Insert(&roleResourceRels)
	}
	session.Commit()
	c.JSON(http.StatusOK, common.NewSuccessJsonResult("赋权成功", nil))
}

// -------------------------- role end ---------------------------------
