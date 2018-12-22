package models

import (
	"fmt"
	"time"

	"github.com/go-xorm/xorm"
	"wozaizhao.com/book/common"
)

type Favorite struct {
	Id          int `json:"id" xorm:"pk autoincr unique"`
	BookId      int
	OpenID      string    `json:"openId"`
	Status      int       //1关注 2不关注
	Created     time.Time `xorm:"-" json:"-"`
	CreatedUnix int64
	Updated     time.Time `xorm:"-" json:"-"`
	UpdatedUnix int64
}

func (f *Favorite) BeforeInsert() {
	f.CreatedUnix = time.Now().Unix()
	f.UpdatedUnix = f.CreatedUnix
}

func (f *Favorite) BeforeUpdate() {
	f.UpdatedUnix = time.Now().Unix()
}

func (f *Favorite) AfterSet(colName string, _ xorm.Cell) {
	switch colName {
	case "created_unix":
		f.Created = time.Unix(f.CreatedUnix, 0).Local()
	case "updated_unix":
		f.Updated = time.Unix(f.UpdatedUnix, 0).Local()
	}
}

func InsertFavorite(openid string, bookid int, status int) {
	f := new(Favorite)
	f.BookId = bookid
	f.OpenID = openid
	f.Status = status

	engine, err := GetEngine()

	if err != nil {
		fmt.Println("数据库初始化失败")
	}

	//如果表不存在就创建表
	var tableFavorite = &Favorite{}

	errc := engine.CreateTables(tableFavorite)

	if errc != nil {
		fmt.Println(errc)
	}

	affected, err := engine.Insert(f)
	fmt.Println(affected)
	if err != nil {
		fmt.Println(err) // 如果不是如预期的那么就报错
	} else {
		fmt.Println("数据插入成功") //记录一些你期望记录的信息
	}
}

func FavoriteExsit(openid string, bookid int) bool {
	engine, err := GetEngine()
	if err != nil {
		fmt.Println("数据库初始化失败")
	}

	has, err := engine.Table("favorite").Where("open_i_d = ? and book_id = ?", openid, bookid).Exist()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return has
}

func UpdateFavorite(openid string, bookid, status int) {
	f := new(Favorite)
	// f.BookId = bookid
	// f.OpenID = openid
	f.Status = status
	// fmt.Println(f)
	affected, err := engine.Where("open_i_d = ? and book_id = ?", openid, bookid).Update(f)
	fmt.Println(affected)
	if err != nil {
		fmt.Println(err) // 如果不是如预期的那么就报错
	} else {
		fmt.Println("数据更新成功") //记录一些你期望记录的信息
	}
}

//用户是否已订阅本书
func HasUserSubscribe(openid string, bookid int) bool {
	engine, err := GetEngine()
	if err != nil {
		fmt.Println("数据库初始化失败")
	}

	has, err := engine.Table("favorite").Where("open_i_d = ? and book_id = ? and status = ?", openid, bookid, common.SUBSCRIBED).Exist()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return has
}

//本书总订阅人数
func Subscription(bookid int) int64 {
	f := new(Favorite)
	total, err := engine.Where("book_id = ? and status = ?", bookid, common.SUBSCRIBED).Count(f)
	if err != nil {
		fmt.Println(err)
	}
	return total
}
