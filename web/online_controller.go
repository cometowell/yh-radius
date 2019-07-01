package web

import (
	"github.com/gin-gonic/gin"
	"go-rad/common"
	"go-rad/database"
	"go-rad/model"
	"go-rad/radius"
	"net/http"
)

// -------------------------- online --------------------------------
func listOnline(c *gin.Context) {
	var online model.RadOnlineUser
	c.ShouldBindJSON(&online)

	whereSql := "1=1 "
	whereArgs := make([]interface{}, 0)
	if online.UserName != "" {
		whereSql += "and ol.username like ? "
		whereArgs = append(whereArgs, "%"+online.UserName+"%")
	}

	if online.IpAddr != "" {
		whereSql += "and ol.ip_addr = ? "
		whereArgs = append(whereArgs, online.IpAddr)
	}

	if online.RealName != "" {
		whereSql += "and ru.real_name like ? "
		whereArgs = append(whereArgs, "%"+online.RealName+"%")
	}

	var onlines []model.Online
	count, _ := database.DataBaseEngine.Table("online_user").Alias("ol").
		Join("INNER", []string{"rad_user", "ru"}, "ol.username = ru.username").
		Where(whereSql, whereArgs...).
		Limit(online.PageSize, online.PageSize*(online.Page-1)).
		FindAndCount(&onlines)

	pagination := model.NewPagination(onlines, count, online.Page, online.PageSize)
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "success", Data: pagination})
}

func offOnline(c *gin.Context) {
	var online model.RadOnlineUser
	c.ShouldBindJSON(&online)

	var dst model.RadOnlineUser
	database.DataBaseEngine.ID(online.Id).Get(&dst)

	if dst.Id == 0 {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: "no records were found", Data: nil})
		return
	}

	//下线用户
	err := radius.OfflineUser(dst)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, common.NewSuccessJsonResult("下线成功", nil))
}

func deleteOnline(c *gin.Context) {
	var online model.RadOnlineUser
	c.ShouldBindJSON(&online)
	count, e := database.DataBaseEngine.ID(online.Id).Delete(&online)
	if e != nil || count == 0 {
		c.JSON(http.StatusOK, common.NewErrorJsonResult("清理在线用户失败"))
		return
	}
	c.JSON(http.StatusOK, common.NewSuccessJsonResult("清理在线用户成功", nil))
}
