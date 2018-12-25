package models

import (
	"time"

	"github.com/go-xorm/xorm"
	"wozaizhao.com/book/common"
)

type Content struct {
	Id          int `json:"id" xorm:"pk autoincr unique"`
	BookId      int
	Sn          int
	Title       string
	Created     time.Time `xorm:"-" json:"-"`
	CreatedUnix int64
	Updated     time.Time `xorm:"-" json:"-"`
	UpdatedUnix int64
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

func InsertContent(bookid int, sn int, title string) (c *Content) {
	c = new(Content)
	c.BookId = bookid
	c.Sn = sn
	c.Title = title

	// var engine *xorm.Engine
	// var err error
	// engine,err = GetEngine()

	// if err !=nil {
	// 	fmt.Println("数据库初始化失败")
	// }

	//如果表不存在就创建表
	var tableContent = &Content{}

	errc := engine.CreateTables(tableContent)

	if errc != nil {
		common.Log("InsertContent CreateTables Error:", errc)
	}

	affected, err := engine.Insert(c)
	// fmt.Println(affected)
	if err != nil {
		common.Log("InsertContent Insert Error:", err)
	} else {
		common.Log("InsertContent Insert Successfully:", affected)
	}
	return
}

func GetContents(id int) (c []Content) {
	// engine,err := GetEngine()
	// if err !=nil {
	// 	fmt.Println("数据库初始化失败")
	// }
	errc := engine.Asc("sn").Where("book_id = ?", id).Find(&c)
	if errc != nil {
		common.Log("GetContents Find Error:", errc)
	}
	return

}
