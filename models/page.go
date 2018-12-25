package models

import (
	"fmt"
	"time"

	"github.com/go-xorm/xorm"
	"wozaizhao.com/book/common"
)

type Page struct {
	Id          int `json:"id" xorm:"pk autoincr unique"`
	BookId      int
	ContentId   int
	Sn          int
	Title       string
	MdUrl       string
	Created     time.Time `xorm:"-" json:"-"`
	CreatedUnix int64
	Updated     time.Time `xorm:"-" json:"-"`
	UpdatedUnix int64
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

func InsertPage(bookid int, contentid int, title string, mdurl string) (p *Page) {
	p = new(Page)
	p.BookId = bookid
	p.ContentId = contentid
	p.Sn = GetPageSn(contentid)
	p.Title = title
	p.MdUrl = mdurl

	// var engine *xorm.Engine
	// var err error
	// engine,err = GetEngine()

	// if err !=nil {
	// 	fmt.Println("数据库初始化失败")
	// }

	//如果表不存在就创建表
	var tablePage = &Page{}

	errc := engine.CreateTables(tablePage)

	if errc != nil {
		common.Log("InsertPage CreateTables Error:", errc)
	}

	affected, err := engine.Insert(p)
	// fmt.Println(affected)
	if err != nil {
		common.Log("InsertPage Insert Error:", err)
	} else {
		common.Log("InsertPage Insert Successfully:", affected)
	}

	return
}

func GetPages(id int) (pages []Page) {
	// engine,err := GetEngine()
	// if err !=nil {
	// 	fmt.Println("数据库初始化失败")
	// }
	errc := engine.Table("page").Cols("content_id", "book_id", "sn", "title").Where("content_id = ?", id).Find(&pages)
	if errc != nil {
		fmt.Println(errc)
	}
	return

}

func GetPageById(id string) (page *Page) {
	page = new(Page)
	// engine,err := GetEngine()
	// if err !=nil {
	// 	fmt.Println("数据库初始化失败")
	// }
	has, errc := engine.Id(id).Get(page)
	if errc != nil {
		common.Log("GetPageById Get Error:", errc)
	}
	common.Log("GetPageById Get:", has)

	return

}

func GetPage(bookid string, contentid string, pageid string) (page *Page) {
	page = new(Page)
	// engine,err := GetEngine()
	// if err !=nil {
	// 	fmt.Println("数据库初始化失败")
	// }
	has, errc := engine.Where("book_id = ? AND content_id = ? AND sn = ?", bookid, contentid, pageid).Get(page)
	if errc != nil {
		common.Log("GetPage Get Error:", errc)
	}

	common.Log("GetPage Get:", has)
	return
}

//查找当前目录下的页面index
func GetPageSn(contentid int) int {
	var pages []Page
	// engine,err := GetEngine()
	// if err !=nil {
	// 	fmt.Println("数据库初始化失败")
	// }
	errc := engine.Table("page").Desc("sn").Cols("sn").Where("content_id = ?", contentid).Find(&pages)
	if errc != nil {
		fmt.Println(errc)
	}
	if len(pages) == 0 {
		return 1
	}
	return pages[0].Sn + 1
}

func PageExist(bookid int, contentid int, pageid int) bool {
	// engine, err := GetEngine()
	// if err != nil {
	// 	fmt.Println("数据库初始化失败")
	// }
	has, errc := engine.Table("page").Where("book_id = ? AND content_id = ? AND sn = ?", bookid, contentid, pageid).Exist()
	if errc != nil {
		common.Log("PageExist Exist Error:", errc)
	}
	return has
}
