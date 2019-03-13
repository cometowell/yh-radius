package main

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"strings"
	"testing"
)

func TestJoin(t *testing.T) {
	tsEngine, _ := xorm.NewEngine("mysql",
		"root:root@tcp(127.0.0.1:3306)/radius?charset=utf8")

	tsEngine.ShowSQL(true)

	mr := make([]SysResource, 0)
	tsEngine.Table("sys_manager").Alias("sm").
		Join("INNER", []string{"sys_manager_role_rel", "smr"}, "sm.id = smr.manager_id").
		Join("INNER", []string{"sys_role", "sr"}, "smr.role_id = sr.id").
		Join("INNER", []string{"sys_role_resource_rel", "srr"}, "srr.role_id = sr.id").
		Join("INNER", []string{"sys_resource", "r"}, "srr.role_id = r.id").
		Where("sm.id = ?", 1).
		Find(&mr)

	fmt.Println(mr)

	fmt.Println(strings.Repeat("-", 50))

	smr := []SysManagerRole{}
	err := tsEngine.Table("sys_manager").Alias("sm").
		Join("INNER", []string{"sys_manager_role_rel", "smr"}, "sm.id = smr.manager_id").
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
