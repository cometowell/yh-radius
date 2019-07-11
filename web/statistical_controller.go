package web

import (
	"github.com/gin-gonic/gin"
	"go-rad/database"
	"net/http"
	"time"
)

type PieStruct struct {
	Value    interface{} `json:"value"`
	Name     string      `json:"name"`
	Selected bool        `json:"selected"`
}

// 新用户发展统计
type NewUserStatistic struct {
	CreateDate string `xorm:"create_date" json:"createDate"`
	Total      int    `xorm:"total" json:"total"`
}

// 套餐订购统计
type ProductOrderStatistic struct {
	ProductName string `xorm:"product_name" json:"productName"`
	Total       int    `xorm:"total" json:"total"`
}

type AreaUserStatistic struct {
	AreaName string `xorm:"area_name" json:"areaName"`
	Total    int    `xorm:"total" json:"total"`
}

type OnlineAndFlowTrendStatistic struct {
	StartHour       int     `xorm:"start_hour" json:"startHour"`
	Total           int     `xorm:"total" json:"total"`
	TotalDownStream float64 `xorm:"total_down_stream" json:"totalDownStream"`
	TotalUpStream   float64 `xorm:"total_up_stream" json:"totalUpStream"`
}

// 统计新用户发展
// 日为单位 1周内
func statisticNewUser(c *gin.Context) {
	sql := `SELECT
			create_date,
			count(*) AS total
		FROM
			(
				SELECT
					*, SUBSTR(create_time FROM 1 FOR 10) AS create_date
				FROM
					rad_user
			) c
		where DATE_SUB(NOW(),INTERVAL 7 DAY) <= create_date
		GROUP BY
			c.create_date
		ORDER BY c.create_date`
	var newUserStatistics []NewUserStatistic
	database.DataBaseEngine.SQL(sql).Find(&newUserStatistics)
	now := time.Now()
	xAxis := [7]string{}
	data := [7]interface{}{}

	for i := 0; i < 7; i++ {
		var find bool = false
		date := now.AddDate(0, 0, -i).Format("2006-01-02")
		xAxis[i] = date
		for _, item := range newUserStatistics {
			if item.CreateDate == date {
				data[i] = item.Total
				find = true
				break
			}
		}
		if !find {
			data[i] = 0
		}
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"xAxis": xAxis,
		"data":  data,
	})
}

// 8小时内上下行流量统计
func statisticOnlineAndFlowTrend(c *gin.Context) {
	sql := `SELECT
			ANY_VALUE(start_hour) as start_hour,
			sum(total_down_stream) / 1024 / 1024 as total_down_stream,
			sum(total_up_stream) / 1024 / 1024 as total_up_stream,
			count(*) AS total
			FROM
				(
					SELECT
						start_time,
						total_down_stream,
						total_up_stream,
						SUBSTR(start_time FROM 12 FOR 2) AS start_hour
					FROM
						rad_user_online_log
					UNION
						SELECT
							start_time,
							total_down_stream,
							total_up_stream,
							SUBSTR(start_time FROM 12 FOR 2) AS start_hour
						FROM
							rad_online_user
						WHERE
							DATE_SUB(NOW(), INTERVAL 8 HOUR) <= start_time
				) t
			GROUP BY
				t.start_hour
			ORDER BY t.start_hour`
	var onlineAndFlowTrendStatistics []OnlineAndFlowTrendStatistic
	database.DataBaseEngine.SQL(sql).Find(&onlineAndFlowTrendStatistics)
	now := time.Now()
	xAxis := [7]int{}
	totalFlow := [7]interface{}{}
	input := [7]interface{}{}
	output := [7]interface{}{}
	total := [7]interface{}{}
	for i := 0; i < 7; i++ {
		var find = false
		date := now.Add(-time.Duration(i) * time.Hour)
		hour := date.Hour()
		xAxis[i] = hour
		for _, item := range onlineAndFlowTrendStatistics {
			if item.StartHour == hour {
				input[i] = item.TotalUpStream
				output[i] = item.TotalDownStream
				total[i] = item.Total
				totalFlow[i] = item.TotalUpStream + item.TotalDownStream
				find = true
				break
			}
		}
		if !find {
			input[i] = 0
			output[i] = 0
			total[i] = 0
			totalFlow[i] = 0
		}
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"xAxis":     xAxis,
		"total":     total,
		"totalFlow": totalFlow,
		"input":     input,
		"output":    output,
	})
}

// 区域用户统计
// 饼图
func statisticAreaUser(c *gin.Context) {
	sql := `SELECT
			ANY_VALUE(area_name) as area_name,
			count(*) AS total
		FROM
			(
				SELECT
					ra.status, ra.name as area_name, ra.id as area_id
				FROM
					rad_user ru, rad_area ra, rad_town rt where ru.town_id = rt.id and rt.area_id = ra.id
			) c
		where c.status = 1
		GROUP BY
			c.area_id`
	var areaUserStatistics []AreaUserStatistic
	database.DataBaseEngine.SQL(sql).Find(&areaUserStatistics)
	areaNames := make([]string, len(areaUserStatistics))
	total := make([]PieStruct, len(areaUserStatistics))

	for index, item := range areaUserStatistics {
		areaNames[index] = item.AreaName
		total[index] = PieStruct{item.Total, item.AreaName, index == 0}
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"areaNames": areaNames,
		"total":     total,
	})

}

// 套餐订购统计，饼图
func statisticProductOrderTrend(c *gin.Context) {
	sql := `SELECT
			any_value(product_name) as product_name,
			count(*) AS total
		FROM
			(
				SELECT
					ru.*, SUBSTR(ru.create_time FROM 1 FOR 10) AS create_date, rp.name as product_name
				FROM
					rad_user ru, rad_product rp where ru.product_id = rp.id
			) c
		where DATE_SUB(NOW(),INTERVAL 7 DAY) <= create_date
		GROUP BY
			c.product_id`
	var productOrderStatistics []ProductOrderStatistic
	database.DataBaseEngine.SQL(sql).Find(&productOrderStatistics)
	productNames := make([]string, len(productOrderStatistics))
	total := make([]PieStruct, len(productOrderStatistics))

	for index, item := range productOrderStatistics {
		productNames[index] = item.ProductName
		total[index] = PieStruct{item.Total, item.ProductName, index == 0}
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"productNames": productNames,
		"total":        total,
	})
}
