package models

import (
	"log"
	"github.com/joho/godotenv"
	"os"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var engine *xorm.Engine

//获取连接字符串
func getDataSource() (string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	
	usrname := os.Getenv("USRNAME")
	password := os.Getenv("PASSWORD")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	database := os.Getenv("DATABASE")

    datasource := usrname + ":" + password  + "@tcp(" + host + ":" + port + ")/" + database
	 
    return datasource
}

//初始化数据库引擎
func initDB() error {
	var ds = getDataSource()
	fmt.Println(ds)
	var err error
	engine, err = xorm.NewEngine("mysql",ds)
	if err!= nil{
		return err
	}
	return nil
}

//获取数据库引擎
func GetEngine() (*xorm.Engine, error){
   if (engine == nil){
	   err := initDB()
	   if err!= nil {
		return nil,err
	}
   }
   return engine,nil
}