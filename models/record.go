package models

import (
	"time"

	"github.com/go-xorm/xorm"
	"wozaizhao.com/book/common"
)

type Record struct {
	Id          int    `json:"id" xorm:"pk autoincr unique"`
	OpenID      string `json:"openId"`
	PageId      int
	Created     time.Time `xorm:"-" json:"-"`
	CreatedUnix int64
	Updated     time.Time `xorm:"-" json:"-"`
	UpdatedUnix int64
}

func (r *Record) BeforeInsert() {
	r.CreatedUnix = time.Now().Unix()
	r.UpdatedUnix = r.CreatedUnix
}

func (r *Record) BeforeUpdate() {
	r.UpdatedUnix = time.Now().Unix()
}

func (r *Record) AfterSet(colName string, _ xorm.Cell) {
	switch colName {
	case "created_unix":
		r.Created = time.Unix(r.CreatedUnix, 0).Local()
	case "updated_unix":
		r.Updated = time.Unix(r.UpdatedUnix, 0).Local()
	}
}

func Recording(openid string, pageid int) {
	r := new(Record)
	r.OpenID = openid
	r.PageId = pageid

	//如果表不存在就创建表
	var tableRecord = &Record{}

	errc := engine.CreateTables(tableRecord)

	if errc != nil {
		common.Log("CreateTables Error:", errc)
	}
	affected, err := engine.Insert(r)
	if err != nil {
		common.Log("InsertFavorite Insert Error:", err)
	} else {
		common.Log("Insert Successfully:", affected)
	}
}

//总的阅读数
func AllRead() int64 {
	r := new(Record)
	total, err := engine.Count(r)
	if err != nil {
		common.Log("Subscription Count Error:", err)
	}
	return total
}

//我的阅读数
func MyRead(openid string) int64 {
	r := new(Record)
	total, err := engine.Where("open_i_d = ?", openid).Count(r)
	if err != nil {
		common.Log("Subscription Count Error:", err)
	}
	return total
}
