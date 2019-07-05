package web

import (
	"github.com/gin-gonic/gin"
	"go-rad/common"
	"go-rad/database"
	"go-rad/model"
	"net/http"
)

// department
func fetchDepartments(c *gin.Context) {
	var departments []model.SysDepartment
	database.DataBaseEngine.Find(&departments)
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "success", Data: departments})
}

func listDepartments(c *gin.Context) {
	var department model.SysDepartment
	err := c.ShouldBindJSON(&department)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}

	whereSql := "1=1 "
	whereArgs := make([]interface{}, 0)
	if department.Code != "" {
		whereSql += "and sd.code like ? "
		whereArgs = append(whereArgs, "%"+department.Code+"%")
	}

	if department.Name != "" {
		whereSql += "and sd.name like ? "
		whereArgs = append(whereArgs, "%"+department.Name+"%")
	}

	if department.ParentId != 0 {
		whereSql += "and sd.parent_id = ? "
		whereArgs = append(whereArgs, department.ParentId)
	}

	var departments []model.Department
	count, _ := database.DataBaseEngine.Cols("sd.*, d.name").Table(&model.SysDepartment{}).Alias("sd").
		Join("LEFT", []interface{}{&model.SysDepartment{}, "d"}, "sd.parent_id = d.id").
		Where(whereSql, whereArgs...).
		Limit(department.PageSize, department.PageSize*(department.Page-1)).
		FindAndCount(&departments)

	pagination := model.NewPagination(departments, count, department.Page, department.PageSize)
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "success", Data: pagination})

}

func getDepartmentInfo(c *gin.Context) {
	var department model.SysDepartment
	err := c.ShouldBindJSON(&department)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	database.DataBaseEngine.Id(department.Id).Get(&department)
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "success", Data: department})
}

func addDepartment(c *gin.Context) {
	var department model.SysDepartment
	err := c.ShouldBindJSON(&department)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	session := database.DataBaseEngine.NewSession()
	defer session.Close()
	session.Begin()
	count, _ := session.Table(&model.SysDepartment{}).Where("code = ? or name = ?", department.Code, department.Name).Count()
	if count > 0 {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: "错误：编码或者名称重复"})
		session.Rollback()
		return
	}
	department.Status = 1
	department.CreateTime = model.NowTime()
	session.InsertOne(&department)
	session.Commit()
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "添加成功"})
}

func updateDepartment(c *gin.Context) {
	var department model.SysDepartment
	err := c.ShouldBindJSON(&department)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	session := database.DataBaseEngine.NewSession()
	defer session.Close()
	session.Begin()
	count, _ := session.Table(&model.SysDepartment{}).
		Where("(code = ? or name = ?) and id != ?", department.Code, department.Name, department.Id).Count()
	if count > 0 {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: "错误：编码或者名称重复"})
		session.Rollback()
		return
	}
	department.UpdateTime = model.NowTime()
	session.Cols("name", "update_time", "description", "status").ID(department.Id).Update(&department)
	session.Commit()
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "修改成功"})
}

func deleteDepartment(c *gin.Context) {
	var department model.SysDepartment
	err := c.ShouldBindJSON(&department)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	department.Status = 2 // 标记为停用
	database.DataBaseEngine.Id(department.Id).Cols("status").Update(&department)
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "已停用!"})
}
