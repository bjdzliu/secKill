package main

import "github.com/astaxie/beego"
import (
	_ "secKill/Secproxy/router"
)

func main(){
	err:=InitConfig()
	if err!=nil{
		panic(err)
		return
	}
	err2:=initSec()
	if err2!=nil{
		panic(err2)
	return
	}
	beego.Run()
}
