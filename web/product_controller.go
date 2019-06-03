package web

import (
	"github.com/gin-gonic/gin"
	"go-rad/common"
	"go-rad/database"
	"go-rad/model"
	"net/http"
)

// -------------------------- product start -----------------------------

func addProduct(c *gin.Context) {
	var product model.RadProduct
	err := c.ShouldBindJSON(&product)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	product.Status = 1
	product.CreateTime = model.NowTime()
	database.DataBaseEngine.InsertOne(&product)
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "套餐添加成功!"})
}

func updateProduct(c *gin.Context) {
	var product model.RadProduct
	err := c.ShouldBindJSON(&product)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	product.UpdateTime = model.NowTime()
	session := database.DataBaseEngine.NewSession()
	session.Begin()
	database.DataBaseEngine.Id(product.Id).Update(&product)
	database.DataBaseEngine.Cols("concurrent_count").ID(product.Id).Update(&product)
	session.Commit()
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "修改成功!"})
}

func fetchProductList(c *gin.Context) {
	var products []model.RadProduct
	database.DataBaseEngine.Where("status = ?", 1).Find(&products)
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "success", Data: products})
}

func listProduct(c *gin.Context) {
	var product model.RadProduct
	err := c.ShouldBindJSON(&product)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	whereSql := "1=1 "
	whereArgs := make([]interface{}, 0)
	if product.Name != "" {
		whereSql += "and name like ? "
		whereArgs = append(whereArgs, "%"+product.Name+"%")
	}

	if product.Status != 0 {
		whereSql += "and status = ? "
		whereArgs = append(whereArgs, product.Status)
	}

	if product.Type != 0 {
		whereSql += "and type = ?"
		whereArgs = append(whereArgs, product.Type)
	}

	var products []model.RadProduct
	totalCount, _ := database.DataBaseEngine.Table("rad_product").Where(whereSql, whereArgs...).
		Limit(product.PageSize, (product.Page-1)*product.PageSize).
		FindAndCount(&products)

	pagination := model.NewPagination(products, totalCount, product.Page, product.PageSize)
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "success", Data: pagination})
}

func deleteProduct(c *gin.Context) {
	var product model.RadProduct
	err := c.ShouldBindJSON(&product)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	product.UpdateTime = model.NowTime()
	product.Status = 2
	database.DataBaseEngine.Cols("status, update_time").ID(product.Id).Update(&product)
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "停用成功!"})
}

func getProductInfo(c *gin.Context) {
	var product model.RadProduct
	err := c.ShouldBindJSON(&product)
	if err != nil {
		c.JSON(http.StatusOK, common.JsonResult{Code: 1, Message: err.Error()})
		return
	}
	database.DataBaseEngine.Id(product.Id).Get(&product)
	c.JSON(http.StatusOK, common.DefaultSuccessJsonResult(product))
}

// -------------------------- product end -----------------------------
