package main

import (
	"time"
)

// user expire task
// if user has continue book order then change the product for user.
func userExpireTask() {
	logger.Info("用户到期定时任务执行...")
	session := engine.NewSession()
	defer session.Close()
	today := time.Now().Format(DateFormat)
	var users []RadUser
	err := session.Where("expire_time < ?", today).Find(&users)
	if err != nil || len(users) == 0 {
		logger.Warn("user expire task occur error: " + err.Error())
		return
	}
	for _, user := range users {
		var record UserOrderRecord
		_, err = session.Where("user_id = ? and status = 1", user.Id).Get(&record)
		if err != nil {
			logger.Warnf("user:%s find order record, %s%s", user.UserName, "user expire task occur error: ", err.Error())
			continue
		}
		if record.Id == 0 {
			continue
		}
		var product RadProduct
		_, err = session.Where("product_id = ?", record.ProductId).Get(&product)
		if err != nil {
			logger.Warnf("user :%s find product, %s%s", user.UserName, "user expire task occur error: ", err.Error())
			continue
		}
		purchaseProduct(&user, &product)
		i, err := session.ID(user.Id).Update(&user)
		if i == 0 || err != nil {
			logger.Warnf("user:%s update to product: %s, %s%s", user.UserName, product.Name, "user expire task occur error: ", err.Error())
			continue
		}
	}

	session.Commit()
}
