package models

import (
	"fmt"
	"time"

	"github.com/go-xorm/xorm"
)

type User struct {
	Id          int       `json:"id" xorm:"pk autoincr unique"`
	OpenID      string    `json:"openId"`
	NickName    string    `json:"nickName"`
	Gender      int       `json:"gender"` //性别，0-未知，1-男，2-女
	City        string    `json:"city"`
	Province    string    `json:"province"`
	Country     string    `json:"country"`
	AvatarURL   string    `json:"avatarUrl"`
	UnionID     string    `json:"unionId"`
	SessionKey  string    `json:"session_key"`
	Skey        string    `json:"skey"`
	Created     time.Time `xorm:"-" json:"-"`
	CreatedUnix int64
	Updated     time.Time `xorm:"-" json:"-"`
	UpdatedUnix int64
}

func (u *User) BeforeInsert() {
	u.CreatedUnix = time.Now().Unix()
	u.UpdatedUnix = u.CreatedUnix
}

func (u *User) BeforeUpdate() {
	u.UpdatedUnix = time.Now().Unix()
}

func (u *User) AfterSet(colName string, _ xorm.Cell) {
	switch colName {
	case "created_unix":
		u.Created = time.Unix(u.CreatedUnix, 0).Local()
	case "updated_unix":
		u.Updated = time.Unix(u.UpdatedUnix, 0).Local()
	}
}

//用户是否存在？
func UserExist(openid string) bool {
	engine, err := GetEngine()
	if err != nil {
		fmt.Println("数据库初始化失败")
	}

	has, err := engine.Table("user").Where("open_i_d = ?", openid).Exist()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return has
}

//更新用户SessionKey和token
func UpdateUserToken(openid, session_key, skey string) {
	user := new(User)
	user.SessionKey = session_key
	user.Skey = skey
	affected, err := engine.Where("open_i_d = ?", openid).Update(user)
	fmt.Println(affected)
	if err != nil {
		fmt.Println(err) // 如果不是如预期的那么就报错
	} else {
		fmt.Println("数据更新成功") //记录一些你期望记录的信息
	}
}

//保存用户信息
func SaveUser(openId, nickName string, gender int, city, province, country, avatarUrl, unionId, session_key, skey string) {
	u := new(User)
	u.OpenID = openId
	u.NickName = nickName
	u.Gender = gender //性别，0-未知，1-男，2-女
	u.City = city
	u.Province = province
	u.Country = country
	u.AvatarURL = avatarUrl
	u.UnionID = unionId
	u.SessionKey = session_key
	u.Skey = skey

	engine, err := GetEngine()

	if err != nil {
		fmt.Println("数据库初始化失败")
	}

	//如果表不存在就创建表
	var tableUser = &User{}

	errc := engine.CreateTables(tableUser)

	if errc != nil {
		fmt.Println(errc)
	}

	affected, err := engine.Insert(u)
	fmt.Println(affected)
	if err != nil {
		fmt.Println(err) // 如果不是如预期的那么就报错
	} else {
		fmt.Println("数据插入成功") //记录一些你期望记录的信息
	}
}

//根据token查找openid
func Skey2OpenId(skey string) string {
	u := new(User)
	engine, err := GetEngine()
	if err != nil {
		fmt.Println("数据库初始化失败")
	}
	has, errc := engine.Where("skey = ?", skey).Get(u)
	if errc != nil {
		fmt.Println(errc)
	}
	if has {
		return u.OpenID
	}
	return ""
}
