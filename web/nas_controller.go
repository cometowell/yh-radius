package web

import (
	"github.com/gin-gonic/gin"
	"go-rad/common"
	"go-rad/database"
	"go-rad/model"
	"net/http"
)

// -------------------------- nas start -----------------------------

func getNasInfo(c *gin.Context) {
	var nas model.RadNas
	err := c.ShouldBindJSON(&nas)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	database.DataBaseEngine.Id(nas.Id).Get(&nas)
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "success", Data: nas})
}

func addNas(c *gin.Context) {
	var nas model.RadNas
	err := c.ShouldBindJSON(&nas)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	session := database.DataBaseEngine.NewSession()
	defer session.Close()
	session.Begin()
	count, _ := session.Table("rad_nas").Where("ip_addr = ?", nas.IpAddr).Count()
	if count > 0 {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: "错误：IP地址重复"})
		session.Rollback()
		return
	}
	session.InsertOne(&nas)
	session.Commit()
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "success"})
}

func updateNas(c *gin.Context) {
	var nas model.RadNas
	err := c.ShouldBindJSON(&nas)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	session := database.DataBaseEngine.NewSession()
	defer session.Close()
	session.Begin()
	count, _ := session.Table("rad_nas").Where("ip_addr = ? and id != ?", nas.IpAddr, nas.Id).Count()
	if count > 0 {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: "错误：IP地址重复"})
		session.Rollback()
		return
	}
	session.ID(nas.Id).Update(&nas)
	session.Commit()
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "修改成功"})
}

func listNas(c *gin.Context) {
	var nas model.RadNas
	err := c.ShouldBindJSON(&nas)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	c.Set("current", nas.Page)
	c.Set("pageSize", nas.PageSize)
	var nasList []model.RadNas
	PageByConditions(c, &nasList, &nas)
}

func deleteNas(c *gin.Context) {
	var nas model.RadNas
	err := c.ShouldBindJSON(&nas)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	database.DataBaseEngine.Id(nas.Id).Delete(&nas)
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "已删除"})
}

// -------------------------- nas end -----------------------------
