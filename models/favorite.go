package models

import (
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

	// engine, err := GetEngine()

	// if err != nil {
	// 	fmt.Println("数据库初始化失败")
	// }

	//如果表不存在就创建表
	var tableFavorite = &Favorite{}

	errc := engine.CreateTables(tableFavorite)

	if errc != nil {
		common.Log("CreateTables Error:", errc)
	}
	affected, err := engine.Insert(f)
	if err != nil {
		common.Log("InsertFavorite Insert Error:", err)
	} else {
		common.Log("Insert Successfully:", affected)
	}
}

func FavoriteExsit(openid string, bookid int) bool {
	// engine, err := GetEngine()
	// if err != nil {
	// 	fmt.Println("数据库初始化失败")
	// }

	has, err := engine.Table("favorite").Where("open_i_d = ? and book_id = ?", openid, bookid).Exist()
	if err != nil {
		common.Log("FavoriteExsit Exist Error:", err)
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
	if err != nil {
		common.Log("UpdateFavorite Update Error:", err)
	} else {
		common.Log("Update Successfully:", affected)
	}
}

//用户是否已订阅本书
func HasUserSubscribe(openid string, bookid int) bool {
	// engine, err := GetEngine()
	// if err != nil {
	// 	fmt.Println("数据库初始化失败")
	// }

	has, err := engine.Table("favorite").Where("open_i_d = ? and book_id = ? and status = ?", openid, bookid, common.SUBSCRIBED).Exist()
	if err != nil {
		common.Log("HasUserSubscribe Exist Error:", err)
	}
	common.Log("HasUserSubscribe:", has)
	return has
}

//本书总订阅人数
func Subscription(bookid int) int64 {
	f := new(Favorite)
	total, err := engine.Where("book_id = ? and status = ?", bookid, common.SUBSCRIBED).Count(f)
	if err != nil {
		common.Log("Subscription Count Error:", err)
	}
	return total
}
