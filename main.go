package main

import (
	"wozaizhao.com/book/server"
)

func main(){
	r := server.SetupRouter()
	r.Run(":8080")
}