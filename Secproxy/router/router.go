package router

import (
	"github.com/astaxie/beego"
	"secKill/Secproxy/controller"
)

func init(){
	beego.Router("/seckill",&controller.SkillController{},"*:Seckill")
	beego.Router("/secinfo",&controller.SkillController{},"*:Secinfo")
}
