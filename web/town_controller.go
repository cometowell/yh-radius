package web

import (
	"github.com/gin-gonic/gin"
	"go-rad/common"
	"go-rad/database"
	"go-rad/model"
	"net/http"
)

func fetchTowns(c *gin.Context) {
	var town model.RadTown
	err := c.ShouldBindJSON(&town)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	var towns []model.RadTown
	database.DataBaseEngine.Table(&model.RadTown{}).Alias("t").Select(`t.*, a.id as area_id, a.name as area_name`).
		Join("LEFT", []interface{}{&model.RadArea{}, "a"}, "t.area_id = a.id").Where("area_id = ?", town.AreaId).Find(&towns)
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "success", Data: towns})
}

func listTowns(c *gin.Context) {
	var town model.RadTown
	err := c.ShouldBindJSON(&town)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}

	whereSql := "1=1 "
	whereArgs := make([]interface{}, 0)
	if town.Code != "" {
		whereSql += "and t.code like ? "
		whereArgs = append(whereArgs, "%"+town.Code+"%")
	}

	if town.AreaId != 0 {
		whereSql += "and t.area_id = ? "
		whereArgs = append(whereArgs, town.AreaId)
	}

	if town.Name != "" {
		whereSql += "and t.name like ? "
		whereArgs = append(whereArgs, "%"+town.Name+"%")
	}

	var towns []model.RadTown
	count, _ := database.DataBaseEngine.Table(&model.RadTown{}).Alias("t").Select(`t.*, a.id as area_id, a.name as area_name`).
		Join("LEFT", []interface{}{&model.RadArea{}, "a"}, "t.area_id = a.id").
		Where(whereSql, whereArgs...).
		Limit(town.PageSize, town.PageSize*(town.Page-1)).
		FindAndCount(&towns)

	pagination := model.NewPagination(towns, count, town.Page, town.PageSize)
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "success", Data: pagination})

}

func getTownInfo(c *gin.Context) {
	var town model.RadTown
	err := c.ShouldBindJSON(&town)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	database.DataBaseEngine.Table(&model.RadTown{}).
		Alias("t").Select(`t.*, a.id as area_id, a.name as area_name`).
		Join("LEFT", []interface{}{&model.RadArea{}, "a"}, "t.area_id = a.id").
		Get(&town)
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "success", Data: town})
}

func addTown(c *gin.Context) {
	var town model.RadTown
	err := c.ShouldBindJSON(&town)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	session := database.DataBaseEngine.NewSession()
	defer session.Close()
	session.Begin()
	count, _ := session.Table(&model.RadTown{}).Where("area_id = ? and (code = ? or name = ?)", town.AreaId, town.Code, town.Name).Count()
	if count > 0 {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: "错误：编码或者名称重复"})
		session.Rollback()
		return
	}
	town.Status = 1
	town.CreateTime = model.NowTime()
	session.InsertOne(&town)
	session.Commit()
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "添加成功"})
}

func updateTown(c *gin.Context) {
	var town model.RadTown
	err := c.ShouldBindJSON(&town)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	session := database.DataBaseEngine.NewSession()
	defer session.Close()
	session.Begin()
	count, _ := session.Table(&model.RadTown{}).
		Where("area_id = ? and (code = ? or name = ?) and id != ?", town.AreaId, town.Code, town.Name, town.Id).Count()
	if count > 0 {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: "错误：编码或者名称重复"})
		session.Rollback()
		return
	}
	town.UpdateTime = model.NowTime()
	session.Cols("name", "update_time", "description", "status").ID(town.Id).Update(&town)
	session.Commit()
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "修改成功"})
}

func deleteTown(c *gin.Context) {
	var town model.RadTown
	err := c.ShouldBindJSON(&town)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	count, err := database.DataBaseEngine.Table(&model.RadUser{}).Where("town_id = ?", town.Id).Count()

	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}

	if count > 0 {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: "村镇/街道已经被用户关联，不允许删除，可以修改为停用"})
		town.Status = 2 // 标记为停用
		database.DataBaseEngine.Id(town.Id).Cols("status").Update(&town)
		return
	}
	database.DataBaseEngine.Id(town.Id).Delete(&town)
	c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: "删除成功!"})

}
