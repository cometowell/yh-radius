package task

import (
	"go-rad/common"
	"go-rad/database"
	"go-rad/logger"
	"go-rad/model"
	"go-rad/radius"
	"go-rad/web"
	"time"
)

// user expire task
// if user has continue book order then change the product for user.
func UserExpireTask() {
	logger.Logger.Info("用户到期定时任务执行...")
	session := database.DataBaseEngine.NewSession()
	defer session.Close()
	today := time.Now().Format(common.DateFormat)
	var users []model.RadUser
	err := session.Table("rad_user").Alias("ru").
		Join("INNER", []string{"rad_product", "rp"}, "ru.product_id = rp.id").
		Where("(ru.expire_time < ?) or (ru.available_time <= 0 and rp.type = 2) or (ru.available_flow <= 0 and rp.type = 3)", today).Find(&users)
	if err != nil {
		logger.Logger.Warn("user expire task occur error: " + err.Error())
		return
	}

	if len(users) == 0 {
		return
	}

	for _, user := range users {
		var record model.UserOrderRecord
		_, err = session.Where("user_id = ? and status = 1", user.Id).Get(&record)
		if err != nil {
			logger.Logger.Warnf("user:%s find order record, %s%s", user.UserName, "user expire task occur error: ", err.Error())
			continue
		}
		if record.Id == 0 {
			continue
		}
		var product model.RadProduct
		_, err = session.Where("id = ?", record.ProductId).Get(&product)
		if err != nil {
			logger.Logger.Warnf("user :%s find product, %s%s", user.UserName, "user expire task occur error: ", err.Error())
			continue
		}
		web.PurchaseProduct(&user, &product, &model.RadUser{Count: record.Count, BeContinue: true})
		user.ProductId = product.Id
		_, err := session.AllCols().ID(user.Id).Update(&user)
		if err != nil {
			logger.Logger.Warnf("user:%s update to product: %s, %s%s", user.UserName, product.Name, "user expire task occur error: ", err.Error())
			continue
		}
		record.Status = radius.OrderUsingStatus
		session.Cols("status").ID(record.Id).Update(&record)
	}

	session.Commit()
}
