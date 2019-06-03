package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-rad/common"
	"go-rad/database"
	"go-rad/model"
	"net/http"
	"reflect"
)

// pagination with where condition string
func PageByWhereSql(c *gin.Context, result interface{}, whereSql string, whereArgs []interface{}) {
	pageSize, _ := c.Get("pageSize")
	current, _ := c.Get("current")
	totalCount, _ := database.DataBaseEngine.Omit("password").Where(whereSql, whereArgs...).Limit(pageSize.(int), (current.(int)-1)*pageSize.(int)).FindAndCount(result)
	pagination := model.NewPagination(result, totalCount, current.(int), pageSize.(int))
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "success", Data: pagination})
}

// pagination with conditions
// conditions should be struct or map
func PageByConditions(c *gin.Context, result interface{}, conditions interface{}) {
	pageSize, _ := c.Get("pageSize")
	current, _ := c.Get("current")
	totalCount, err := database.DataBaseEngine.Limit(pageSize.(int), (current.(int)-1)*pageSize.(int)).FindAndCount(result, conditions)
	if err != nil {
		panic(err)
	}
	pagination := model.NewPagination(result, totalCount, current.(int), pageSize.(int))
	c.JSON(http.StatusOK, common.JsonResult{Code: 0, Message: "success", Data: pagination})
}

// struct to map
func structToMap(data interface{}) (dst map[string]interface{}) {
	dst = make(map[string]interface{})
	dataType := reflect.TypeOf(data)
	dataValue := reflect.ValueOf(data)
	for i := 0; i < dataType.NumField(); i++ {
		field := dataType.Field(i)
		val := dataValue.FieldByName(field.Name)
		if val.IsValid() { // filter zero value
			dst[field.Name] = val.Interface()
		}
	}
	return dst
}

// build xorm where sql and where args
func buildWhereSql(params map[string]interface{}, limitConditions map[string]string) (whereSql string, whereArgs []interface{}) {
	whereSql += "1=1"
	template := "and %s %s ? "
	whereArgs = make([]interface{}, 0)
	for key, value := range params {
		condition, ok := limitConditions[key]
		if !ok {
			condition = "="
		}
		whereSql += fmt.Sprintf(template, key, condition)
		whereArgs = append(whereArgs, value)
	}
	return
}
