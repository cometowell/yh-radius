package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-rad/common"
	"go-rad/database"
	"go-rad/model"
	"go-rad/radius"
	"net/http"
	"time"
)

// -------------------------- user start -----------------------------
func listUser(c *gin.Context) {
	var params model.RadUser
	err := c.ShouldBindJSON(&params)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}

	whereSql := "1=1 "
	whereArgs := make([]interface{}, 0)
	if params.UserName != "" {
		whereSql += "and ru.username like ? "
		whereArgs = append(whereArgs, "%"+params.UserName+"%")
	}

	if params.RealName != "" {
		whereSql += "and ru.real_name like ? "
		whereArgs = append(whereArgs, "%"+params.RealName+"%")
	}

	if params.Status != 0 {
		whereSql += "and ru.status = ?"
		whereArgs = append(whereArgs, params.Status)
	}

	var users []model.RadUserProduct
	totalCount, _ := database.DataBaseEngine.Table("rad_user").
		Alias("ru").Select(`ru.id,ru.username,ru.real_name,ru.product_id,
			ru.status,ru.available_time,ru.available_flow,ru.expire_time,
			ru.concurrent_count,ru.should_bind_mac_addr,ru.should_bind_vlan,ru.mac_addr,ru.vlan_id,
			ru.vlan_id2,ru.framed_ip_addr,ru.installed_addr,ru.mobile,ru.email,
			ru.pause_time,ru.create_time,ru.update_time,ru.description, sp.*`).
		Where(whereSql, whereArgs...).
		Limit(params.PageSize, (params.Page-1)*params.PageSize).
		Join("INNER", []string{"rad_product", "sp"}, "ru.product_id = sp.id").
		FindAndCount(&users)

	pagination := model.NewPagination(users, totalCount, params.Page, params.PageSize)
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "success", Data: pagination})
}

func fetchUserOrderRecord(c *gin.Context) {
	var user model.RadUser
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	var records []model.UserOrderRecordProduct
	database.DataBaseEngine.Join("INNER", "rad_product", "rad_product.id = user_order_record.product_id").Where("user_id = ?", user.Id).Asc("user_order_record.status").Find(&records)
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "success", Data: records})
}

func updateUser(c *gin.Context) {
	var user model.RadUser
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	session := database.DataBaseEngine.NewSession()
	defer session.Close()
	session.Begin()
	count, _ := session.Table("rad_user").Where("username = ? and id != ?", user.UserName, user.Id).Count()
	if count > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "用户名重复",
		})
		session.Rollback()
		return
	}

	var oldUser model.RadUser
	session.ID(user.Id).Get(&oldUser)
	// 停机用户重新使用需要顺延过期时间
	if oldUser.Status == radius.UserPauseStatus && user.Status == radius.UserAvailableStatus {
		hours := time.Now().Sub(time.Time(oldUser.PauseTime)).Hours()
		user.ExpireTime = model.Time(time.Time(user.ExpireTime).AddDate(0, 0, int(hours)/24))
	}
	session.ID(user.Id).Update(&user)
	session.Commit()
	c.JSON(http.StatusOK, common.NewSuccessJsonResult("success", nil))
}

func addUser(c *gin.Context) {
	var user model.RadUser
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	user.Status = radius.UserAvailableStatus
	user.CreateTime = model.NowTime()
	fmt.Printf("%#v", user)
	session := database.DataBaseEngine.NewSession()
	defer session.Close()
	session.Begin()
	count, _ := session.Table("rad_user").Where("username = ?", user.UserName).Count()
	if count > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "用户名重复",
		})
		session.Rollback()
		return
	}

	var product model.RadProduct
	session.ID(user.ProductId).Get(&product)
	if product.Id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "产品不存在",
		})
		session.Rollback()
		return
	}
	user.Password = common.Encrypt(user.Password)
	PurchaseProduct(&user, &product, &model.RadUser{})
	session.InsertOne(&user)
	// 订购信息
	webSession := GlobalSessionManager.GetSessionByGinContext(c)
	manager := webSession.GetAttr("manager").(model.SysUser)
	orderRecord := model.UserOrderRecord{
		UserId:    user.Id,
		ProductId: product.Id,
		Price:     user.Price,
		SysUserId: manager.Id,
		OrderTime: model.NowTime(),
		Status:    radius.OrderUsingStatus,
		EndDate:   user.ExpireTime,
	}
	session.InsertOne(&orderRecord)

	session.Commit()
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "用户添加成功!"})
}

func deleteUser(c *gin.Context) {
	var user model.RadUser
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	user.Status = radius.UserDeletedStatus
	database.DataBaseEngine.Id(user.Id).Update(&user)
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "删除成功!"})
}

// get user info by id
func fetchUser(c *gin.Context) {
	var user model.RadUser
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	database.DataBaseEngine.Id(user.Id).Get(&user)
	user.Password = ""
	c.JSON(http.StatusOK, common.DefaultSuccessJsonResult(user))
}

func continueProduct(c *gin.Context) {
	var user model.RadUser
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	session := database.DataBaseEngine.NewSession()
	defer session.Close()
	session.Begin()
	bookOrderCount, e := session.Table("user_order_record").Where("user_id = ? and status = ?", user.Id, radius.OrderBookStatus).Count()

	if e != nil {
		session.Rollback()
		c.JSON(http.StatusOK, common.NewErrorJsonResult("用户预定了套餐失败"+err.Error()))
		return
	}

	if bookOrderCount > 0 {
		session.Rollback()
		c.JSON(http.StatusOK, common.NewErrorJsonResult("用户已经预定了套餐暂未生效，不允许再次预定"))
		return
	}

	var oldUser model.RadUser
	session.ID(user.Id).Get(&oldUser)
	var newProduct model.RadProduct
	session.ID(user.ProductId).Get(&newProduct)

	var oldProduct model.RadProduct
	session.ID(oldUser.ProductId).Get(&oldProduct)

	webSession := GlobalSessionManager.GetSessionByGinContext(c)
	manager := webSession.GetAttr("manager").(model.SysUser)

	if newProduct.Id == 0 {
		session.Rollback()
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "产品不存在",
		})
		return
	}

	var expireDate model.Time
	var orderStatus int
	if radius.IsExpire(&oldUser, &oldProduct) { // 产品到期, 直接更新产品信息
		PurchaseProduct(&oldUser, &newProduct, &user)
		expireDate = oldUser.ExpireTime
		orderStatus = radius.OrderUsingStatus
	} else {
		// 产品未到期续订同一产品，修改过期时间
		if oldUser.ProductId == user.ProductId {
			oldUser.ExpireTime = model.Time(time.Time(oldUser.ExpireTime).AddDate(0, newProduct.ServiceMonth*user.Count, 0))
			expireDate = model.Time(oldUser.ExpireTime)
			orderStatus = radius.OrderUsingStatus
		} else {
			// 产品未到期续订不同产品，作为预定订单，当产品到期定时任务更换为预定产品
			expire, _ := common.GetStdTimeFromString("2099-12-31 23:59:59")
			expireDate = model.Time(expire)
			orderStatus = radius.OrderBookStatus
		}
	}

	orderRecord := model.UserOrderRecord{
		UserId:    user.Id,
		ProductId: newProduct.Id,
		Price:     user.Price,
		SysUserId: manager.Id,
		OrderTime: model.NowTime(),
		Status:    orderStatus,
		Count:     user.Count,
		EndDate:   expireDate,
	}

	session.InsertOne(&orderRecord)
	session.AllCols().ID(oldUser.Id).Update(&oldUser)
	session.Commit()
	c.JSON(http.StatusOK, common.NewSuccessJsonResult("续订成功!", nil))
}

func PurchaseProduct(user *model.RadUser, product *model.RadProduct, continueUser *model.RadUser) {
	user.ShouldBindMacAddr = product.ShouldBindMacAddr
	user.ProductId = product.Id
	user.ShouldBindVlan = product.ShouldBindVlan
	user.ConcurrentCount = product.ConcurrentCount
	user.AvailableTime = product.ProductDuration
	user.AvailableFlow = product.ProductFlow
	if product.Type == common.MonthlyProduct {
		expire := time.Time(user.ExpireTime)
		if time.Time(expire).IsZero() || continueUser.BeContinue {
			expire = time.Now()
		}

		month := product.ServiceMonth
		if continueUser.BeContinue {
			month *= continueUser.Count
		}

		expire = time.Time(time.Date(expire.Year(), expire.Month()+time.Month(month), expire.Day(), 23, 59, 59, 0, expire.Location()))
		user.ExpireTime = model.Time(expire)
	} else if product.Type == common.TimeProduct {
		if time.Time(user.ExpireTime).IsZero() || continueUser.BeContinue {
			expireTime, _ := common.GetStdTimeFromString("2099-12-31 23:59:59")
			user.ExpireTime = model.Time(expireTime)
		}
	} else if product.Type == common.FlowProduct {
		if product.FlowClearCycle == common.DefaultFlowClearCycle {
			expireTime, _ := common.GetStdTimeFromString("2099-12-31 23:59:59")
			user.ExpireTime = model.Time(expireTime)
		} else if product.FlowClearCycle == common.DayFlowClearCycle {
			user.ExpireTime = model.Time(common.GetNextDayLastTime())
		} else if product.FlowClearCycle == common.MonthFlowClearCycle {
			user.ExpireTime = model.Time(common.GetMonthLastTime())
		} else if product.FlowClearCycle == common.FixedPeriodFlowClearCycle {
			if time.Time(user.ExpireTime).IsZero() || continueUser.BeContinue {
				user.ExpireTime = model.Time(common.GetDayLastTimeAfterAYear())
			}
		}
	}
}

// -------------------------- user end ----------------------------------
