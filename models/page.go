package models

import (
	"time"
	"fmt"
	"github.com/go-xorm/xorm"
)

type Page struct {
	Id int `json:"id" xorm:"pk autoincr unique"`
	BookId int
	ContentId int
	Sn int
	Title string
	MdUrl string
	Created              time.Time   `xorm:"-" json:"-"`
	CreatedUnix          int64
	Updated              time.Time   `xorm:"-" json:"-"`
	UpdatedUnix          int64
}

func (p *Page) BeforeInsert() {
	p.CreatedUnix = time.Now().Unix()
	p.UpdatedUnix = p.CreatedUnix
}

func (p *Page) BeforeUpdate() {
	p.UpdatedUnix = time.Now().Unix()
}

func (p *Page) AfterSet(colName string, _ xorm.Cell) {
	switch colName {
	case "created_unix":
		p.Created = time.Unix(p.CreatedUnix, 0).Local()
	case "updated_unix":
		p.Updated = time.Unix(p.UpdatedUnix, 0).Local()
	}
}

func InsertPage(bookid int,contentid int,sn int, title string, mdurl string){
	p := new(Page)
	p.BookId = bookid
	p.ContentId= contentid
	p.Sn = sn
	p.Title = title
	p.MdUrl = mdurl

	var engine *xorm.Engine
	var err error
	engine,err = GetEngine()

	if err !=nil {
		fmt.Println("数据库初始化失败")
	}

	//如果表不存在就创建表
	var tablePage= &Page{}

	errc := engine.CreateTables(tablePage)

	if errc != nil {
		fmt.Println(errc)
	}

	affected, err := engine.Insert(p)
	fmt.Println(affected)
	if err != nil {
		fmt.Println(err) // 如果不是如预期的那么就报错
	} else {
		fmt.Println("数据插入成功") //记录一些你期望记录的信息
	}
}

func GetPages(id int) []Page {
	var pages []Page
	engine,err := GetEngine()
	if err !=nil {
		fmt.Println("数据库初始化失败")
	}
	errc := engine.Table("page").Cols("content_id","book_id","sn","title").Where("content_id = ?",id).Find(&pages)
	if  errc != nil {
		fmt.Println(errc)
	}
	return pages

}