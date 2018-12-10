package models

import (
	"time"
	"fmt"
	"github.com/go-xorm/xorm"
)

type Content struct {
	Id int `json:"id" xorm:"pk autoincr unique"`
	BookId int
	Sn int
	Title string
	Created              time.Time   `xorm:"-" json:"-"`
	CreatedUnix          int64
	Updated              time.Time   `xorm:"-" json:"-"`
	UpdatedUnix          int64
}

func (c *Content) BeforeInsert() {
	c.CreatedUnix = time.Now().Unix()
	c.UpdatedUnix = c.CreatedUnix
}

func (c *Content) BeforeUpdate() {
	c.UpdatedUnix = time.Now().Unix()
}

func (c *Content) AfterSet(colName string, _ xorm.Cell) {
	switch colName {
	case "created_unix":
		c.Created = time.Unix(c.CreatedUnix, 0).Local()
	case "updated_unix":
		c.Updated = time.Unix(c.UpdatedUnix, 0).Local()
	}
}

func InsertContent(bookid int,sn int, title string){
	c := new(Content)
	c.BookId = bookid
	c.Sn = sn
	c.Title = title

	var engine *xorm.Engine
	var err error
	engine,err = GetEngine()

	if err !=nil {
		fmt.Println("数据库初始化失败")
	}

	//如果表不存在就创建表
	var tableContent= &Content{}

	errc := engine.CreateTables(tableContent)

	if errc != nil {
		fmt.Println(errc)
	}

	affected, err := engine.Insert(c)
	fmt.Println(affected)
	if err != nil {
		fmt.Println(err) // 如果不是如预期的那么就报错
	} else {
		fmt.Println("数据插入成功") //记录一些你期望记录的信息
	}
}

func GetContents(id int) []Content {
	var contents []Content
	engine,err := GetEngine()
	if err !=nil {
		fmt.Println("数据库初始化失败")
	}
	errc := engine.Asc("sn").Where("book_id = ?",id).Find(&contents)
	if  errc != nil {
		fmt.Println(errc)
	}
	return contents

}