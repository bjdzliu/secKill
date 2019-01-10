package controller

import "github.com/astaxie/beego"

type SkillController struct{
	beego.Controller
}

func(c *SkillController) Seckill(){
	c.Data["json"] = "sec kill"
	// send a json
	c.ServeJSON()
}

func(c *SkillController) Secinfo(){
	c.Data["json"] = "sec info"
	c.ServeJSON()
}