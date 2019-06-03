package web

import (
	"github.com/gin-gonic/gin"
	"go-rad/common"
	"go-rad/database"
	"go-rad/model"
	"net/http"
)

// -------------------------- resource start ---------------------------

func listRes(c *gin.Context) {
	var resList []model.SysResource
	database.DataBaseEngine.Find(&resList)
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "success", Data: getResLevel(resList)})
}

func getSessionResource(c *gin.Context) {
	c.JSON(http.StatusOK, common.DefaultSuccessJsonResult(getManagerResources(c.GetHeader(common.SessionName))))
}

func getManagerResources(sessionId string) []model.SysResource {
	session := GlobalSessionManager.Provider.ReadSession(sessionId)
	resources := session.GetAttr("resources").([]model.SysResource)
	return getResLevel(resources)
}

// 菜单分层展示
func getResLevel(resList []model.SysResource) []model.SysResource {
	result := make([]model.SysResource, 0, 20)
	for _, res := range resList {
		if res.ParentId == 0 {
			r := res
			setChildren(&r, resList)
			result = append(result, r)
		}
	}
	return result
}

func setChildren(r *model.SysResource, resList []model.SysResource) {
	if r.Children == nil {
		r.Children = make([]model.SysResource, 0, 20)
	}
	for _, item := range resList {
		res := item
		if r.Id == res.ParentId {
			setChildren(&res, resList)
			r.Children = append(r.Children, res)
		}
	}
}

// -------------------------- resource end -----------------------------
