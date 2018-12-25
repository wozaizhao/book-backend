package main

import (
	"wozaizhao.com/book/common"
	"wozaizhao.com/book/models"
	"wozaizhao.com/book/server"
)

func main() {
	common.Loginit()
	common.Log("main", "Server Start")
	models.DBinit()
	r := server.SetupRouter()
	r.Run(":8080")
}
