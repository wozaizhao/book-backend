package models

import (
	"fmt"
	"testing"
	"github.com/go-xorm/xorm"
)

func Test_AddBook(t *testing.T) {
    
    b := new(Book)
	b.Name = "Vuejs Learn"
	b.Cover = "/vue/logo.png"
	b.Slogan = "vuejs"
	// b.Bg = ""
	// b.Color = ""
	b.Tag = "EN"
	b.Intro ="Vue.js features an incrementally adoptable architecture that focuses on declarative rendering and component composition. Advanced features required for complex applications such as routing, state management and build tooling are offered via officially maintained supporting libraries and packages."


	var engine *xorm.Engine
	var err error
	engine,err = GetEngine()

	if err !=nil {
		t.Log("数据库初始化失败")
	}

	//如果表不存在就创建表
	var tableBook = &Book{}

	errc := engine.CreateTables(tableBook)

	if errc != nil {
		t.Error(errc)
	}

	affected, err := engine.Insert(b)
	fmt.Println(affected)
	if err != nil {
		t.Error(err) // 如果不是如预期的那么就报错
	} else {
		t.Log("数据插入成功") //记录一些你期望记录的信息
	}
}