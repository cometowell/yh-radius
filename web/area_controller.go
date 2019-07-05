package web

import (
	"github.com/gin-gonic/gin"
	"go-rad/common"
	"go-rad/database"
	"go-rad/model"
	"net/http"
)

func fetchAreas(c *gin.Context) {
	var areas []model.RadArea
	database.DataBaseEngine.Where("status = 1").Find(&areas)
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "success", Data: areas})
}

func listAreas(c *gin.Context) {
	var area model.RadArea
	err := c.ShouldBindJSON(&area)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}

	whereSql := "1=1 "
	whereArgs := make([]interface{}, 0)
	if area.Code != "" {
		whereSql += "and a.code like ? "
		whereArgs = append(whereArgs, "%"+area.Code+"%")
	}

	if area.Name != "" {
		whereSql += "and a.name like ? "
		whereArgs = append(whereArgs, "%"+area.Name+"%")
	}

	var areas []model.RadArea
	count, _ := database.DataBaseEngine.Table(&model.RadArea{}).Alias("a").
		Where(whereSql, whereArgs...).
		Limit(area.PageSize, area.PageSize*(area.Page-1)).
		FindAndCount(&areas)

	pagination := model.NewPagination(areas, count, area.Page, area.PageSize)
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "success", Data: pagination})

}

func getAreaInfo(c *gin.Context) {
	var area model.RadArea
	err := c.ShouldBindJSON(&area)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	database.DataBaseEngine.Id(area.Id).Get(&area)
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "success", Data: area})
}

func addArea(c *gin.Context) {
	var area model.RadArea
	err := c.ShouldBindJSON(&area)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	session := database.DataBaseEngine.NewSession()
	defer session.Close()
	session.Begin()
	count, _ := session.Table(&model.RadArea{}).Where("code = ? or name = ?", area.Code, area.Name).Count()
	if count > 0 {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: "错误：编码或者名称重复"})
		session.Rollback()
		return
	}
	area.CreateTime = model.NowTime()
	area.Status = 1
	session.InsertOne(&area)
	session.Commit()
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "添加成功"})
}

func updateArea(c *gin.Context) {
	var area model.RadArea
	err := c.ShouldBindJSON(&area)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	session := database.DataBaseEngine.NewSession()
	defer session.Close()
	session.Begin()
	count, _ := session.Table(&model.RadArea{}).Where("(code = ? or name = ?) and id != ?", area.Code, area.Name, area.Id).Count()
	if count > 0 {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: "错误：编码或者名称重复"})
		session.Rollback()
		return
	}
	area.UpdateTime = model.NowTime()
	session.Cols("name", "update_time", "description", "status").ID(area.Id).Update(&area)
	session.Commit()
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "修改成功"})
}

func deleteArea(c *gin.Context) {
	var area model.RadArea
	err := c.ShouldBindJSON(&area)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}

	count, err := database.DataBaseEngine.Table(&model.RadTown{}).Where("area_id = ?", area.Id).Count()
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}

	if count > 0 {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: "片区已经被村镇/街道关联，不允许删除，可以修改为停用"})
		return
	}

	database.DataBaseEngine.Id(area.Id).Delete(&area)
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "删除成功!"})
}
