package database

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"go-rad/model"
	"strings"
	"testing"
)

func TestJoin(t *testing.T) {
	tsEngine, _ := xorm.NewEngine("mysql",
		"root:root@tcp(127.0.0.1:3306)/radius?charset=utf8")

	tsEngine.ShowSQL(true)

	mr := make([]model.SysResource, 0)
	tsEngine.Table("sys_user").Alias("sm").
		Join("INNER", []string{"sys_user_role_rel", "smr"}, "sm.id = smr.sys_user_id").
		Join("INNER", []string{"sys_role", "sr"}, "smr.role_id = sr.id").
		Join("INNER", []string{"sys_role_resource_rel", "srr"}, "srr.role_id = sr.id").
		Join("INNER", []string{"sys_resource", "r"}, "srr.role_id = r.id").
		Where("sm.id = ?", 1).
		Find(&mr)

	fmt.Println(mr)

	fmt.Println(strings.Repeat("-", 50))

	smr := []model.SysUserRole{}
	err := tsEngine.Table("sys_user").Alias("sm").
		Join("INNER", []string{"sys_user_role_rel", "smr"}, "sm.id = smr.sys_user_id").
		Join("INNER", []string{"sys_role", "sr"}, "smr.role_id = sr.id").
		Join("INNER", []string{"sys_role_resource_rel", "srr"}, "sr.id = srr.role_id").
		Join("INNER", []string{"sys_resource", "r"}, "srr.resource_id = r.id").
		Find(&smr)

	if err != nil {
		panic(err)
	}

	for _, v := range smr {
		fmt.Printf("%+v", v)
	}
}

func TestJoin1(t *testing.T) {
	tsEngine, _ := xorm.NewEngine("mysql",
		"root:root@tcp(127.0.0.1:3306)/radius?charset=utf8")

	tsEngine.ShowSQL(true)
	//var users = make([]RadUserProduct, 0)
	//total, err := tsEngine.Table("rad_user").Alias("r").
	//	Limit(10, 0).Join("INNER", []string{"rad_product", "sp"}, "r.product_id = sp.id").
	//	FindAndCount(&users)
	//if err != nil {
	//	fmt.Println("你妹的异常了", err.Error())
	//}
	//users := make([]RadUser, 0)
	var users model.RadUser
	tsEngine.Get(&users)

	fmt.Printf("%#v", users) //, total, err)
}

type Manager struct {
	model.SysUser `xorm:"extends"`
	Roles         []model.SysRole `xorm:"extends"`
}

func (Manager) TableName() string {
	return "sys_user"
}

func TestCollection(t *testing.T) {
	tsEngine, _ := xorm.NewEngine("mysql",
		"root:root@tcp(127.0.0.1:3306)/radius?charset=utf8")

	tsEngine.ShowSQL(true)
	var managers []Manager
	tsEngine.Table("sys_user").Alias("sm").
		Join("INNER", []string{"sys_user_role_rel", "smr"}, "sm.id = smr.sys_user_id").
		Join("INNER", []string{"sys_role", "sr"}, "smr.role_id = sr.id").
		Find(&managers)

	fmt.Printf("%#v", managers)
}

func TestDepartment(t *testing.T) {
	tsEngine, _ := xorm.NewEngine("mysql",
		"root:root@tcp(127.0.0.1:3306)/radius?charset=utf8")

	tsEngine.ShowSQL(true)

	var departments []model.Department
	count, _ := tsEngine.Cols("sd.*, d.name").Table("sys_department").Alias("sd").
		Join("LEFT", []string{"sys_department", "d"}, "sd.parent_id = d.id").
		Where("sd.status = 1").
		Limit(10, 0).
		FindAndCount(&departments)

	fmt.Printf("%d, %#v", count, departments)
}
