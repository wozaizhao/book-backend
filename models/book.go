package models

import (
	"time"
	"github.com/go-xorm/xorm"
	"fmt"
)

type Book struct {
	Id int `json:"id" xorm:"pk autoincr unique"`
	Priority int
	Name string
	Cate string
	Cover string
	Slogan string
	Bg string
	Color string
	Tag string
	Intro string
	Path string
	Created              time.Time   `xorm:"-" json:"-"`
	CreatedUnix          int64
	Updated              time.Time   `xorm:"-" json:"-"`
	UpdatedUnix          int64
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

func GetBooks() []Book {
	var books []Book
	engine,err := GetEngine()
	if err !=nil {
		fmt.Println("数据库初始化失败")
	}
	errc := engine.Table("book").Cols("name","cate","cover","slogan","bg","color","tag").Find(&books);
	if  errc != nil {
		fmt.Println(errc)
	}
	return books

}

func GetBook(id int) *Book {
	book := new(Book)
	engine,err := GetEngine()
	if err !=nil {
		fmt.Println("数据库初始化失败")
	}
	has, errc := engine.Id(id).Get(book)
	if  errc != nil {
		fmt.Println(errc)
	}
	if has {
		return book
	}
	return nil

}

func InsertBook(priority int,name string,cate string, cover string,slogan string,bg string,color string,tag string,intro string,path string){
	b := new(Book)
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

	var engine *xorm.Engine
	var err error
	engine,err = GetEngine()

	if err !=nil {
		fmt.Println("数据库初始化失败")
	}

	//如果表不存在就创建表
	var tableBook = &Book{}

	errc := engine.CreateTables(tableBook)

	if errc != nil {
		fmt.Println(errc)
	}

	affected, err := engine.Insert(b)
	fmt.Println(affected)
	if err != nil {
		fmt.Println(err) // 如果不是如预期的那么就报错
	} else {
		fmt.Println("数据插入成功") //记录一些你期望记录的信息
	}
}