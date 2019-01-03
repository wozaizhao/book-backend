package models

import (
	"time"

	"github.com/go-xorm/xorm"
	"wozaizhao.com/book/common"
)

type Book struct {
	Id          int `json:"id" xorm:"pk autoincr unique"`
	Priority    int
	Name        string
	Cate        string
	Cover       string
	Slogan      string
	Bg          string
	Color       string
	Tag         string
	Intro       string
	Path        string
	Url         string
	Count       int
	Status      int       `xorm:"default 0"`
	Created     time.Time `xorm:"-" json:"-"`
	CreatedUnix int64
	Updated     time.Time `xorm:"-" json:"-"`
	UpdatedUnix int64
}

func (b *Book) BeforeInsert() {
	b.CreatedUnix = time.Now().Unix()
	b.UpdatedUnix = b.CreatedUnix
}

func (b *Book) BeforeUpdate() {
	b.UpdatedUnix = time.Now().Unix()
}

func (b *Book) AfterSet(colName string, _ xorm.Cell) {
	switch colName {
	case "created_unix":
		b.Created = time.Unix(b.CreatedUnix, 0).Local()
	case "updated_unix":
		b.Updated = time.Unix(b.UpdatedUnix, 0).Local()
	}
}

func GetBooks() (books []Book) {
	// engine, err := GetEngine()
	// if err != nil {
	// 	fmt.Println("数据库初始化失败")
	// }
	errc := engine.Table("book").Where("status = ?", 1).Cols("id", "name", "cate", "cover", "slogan", "bg", "color", "tag").Desc("priority").Find(&books)
	if errc != nil {
		common.Log("GetBooks Find Error:", errc)
	}
	return

}

func GetSelfBooks(openid string) (books []Book) {
	// engine, err := GetEngine()
	// if err != nil {
	// 	fmt.Println("数据库初始化失败")
	// }
	sql := `select * from book, favorite where book.id  = favorite.book_id and favorite.status = 1 and favorite.open_i_d = "` + openid + `"`
	errc := engine.Sql(sql).Find(&books)
	if errc != nil {
		common.Log("GetSelfBooks Find Error:", errc)
	}
	return

}

func GetBook(id int) (b *Book) {
	b = new(Book)
	// engine, err := GetEngine()
	// if err != nil {
	// 	fmt.Println("数据库初始化失败")
	// }
	has, errc := engine.Id(id).Get(b)
	if errc != nil {
		common.Log("GetBook Get Error:", errc)
	}
	common.Log("GetBook Get:", has)
	return
}

func InsertBook(priority int, name string, cate string, cover string, slogan string, bg string, color string, tag string, intro string, path string, url string) (b *Book) {
	b = new(Book)
	b.Priority = priority
	b.Name = name
	b.Cate = cate
	b.Cover = cover
	b.Slogan = slogan
	b.Bg = bg
	b.Color = color
	b.Tag = tag
	b.Intro = intro
	b.Path = path
	b.Url = url

	// var engine *xorm.Engine
	// var err error
	// engine, err = GetEngine()

	// if err != nil {
	// 	fmt.Println("数据库初始化失败")
	// }

	//如果表不存在就创建表
	var tableBook = &Book{}

	errc := engine.CreateTables(tableBook)

	if errc != nil {
		common.Log("InsertBook CreateTables Error:", errc)
	}

	affected, err := engine.Insert(b)
	if err != nil {
		common.Log("InsertBook Insert Error:", err)
	} else {
		common.Log("InsertBook Insert Successfully:", affected)
	}
	return
}

func ReadBookRecord(bookid int) {
	b := new(Book)
	has, errc := engine.Where("id = ?", bookid).Get(b)
	if has {
		b.Count = b.Count + 1
		affected, err := engine.Where("id = ?", bookid).Update(b)
		if err != nil {
			common.Log("ReadBookRecord Update Error:", err)
		} else {
			common.Log("Update Successfully:", affected)
		}
	}
	if errc != nil {
		common.Log("ReadBookRecord Get Error:", errc)
	}

}

func GetBookCount() int64 {
	b := new(Book)
	total, err := engine.Count(b)
	if err != nil {
		common.Log("GetBookCount Count Error:", err)
	}
	return total
}
