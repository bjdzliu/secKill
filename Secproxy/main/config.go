package main

import (
	"github.com/astaxie/beego"
	"fmt"
	"github.com/astaxie/beego/logs"
)

/*
为什么定义这个文件：
放在main中，initConfig 后变量就消失了。无法为其他包使用

 */


 //声明全局变量 secKillconf,保存redis的配置信息
 var(
 	secKillconf=&SecKillConf{}
 	b="test" //全局变量，同一个包内直接引用
 )

 type RedisConf struct{
 	 redisAddr string
	 redisMaxidle int
	 redisMaxActive int
	 redisIdleTimeout int
 }

 type EtcdConf struct{
 	etcdAddr string
 	timeout int
 	etcdSecKeyPrefix string
 	etcdSecProductkey string
}

//SecKillConf 配置由redisconf和etcdconf组成。
type SecKillConf struct{
    RedisConf
	EtcdConf
	logpath string
	loglevel string
	secProduction	[]SecProductInfoConf
}

type SecProductInfoConf struct{
ProductId int
StartTime int
EndTime int
Status int
Total int
Left int
}

func InitConfig()(err error){
	fmt.Print(beego.BConfig.RunMode)

	redisAddr := beego.AppConfig.String("redis_addr")
	etcdAddr := beego.AppConfig.String("etcd_addr")
	etcdTimeout,_:=beego.AppConfig.Int("etcd_timeout")


	secKillconf.RedisConf.redisAddr=redisAddr
	//
	secKillconf.EtcdConf.etcdAddr=etcdAddr
	secKillconf.EtcdConf.timeout=etcdTimeout
	secKillconf.EtcdConf.etcdSecKeyPrefix=beego.AppConfig.String("etcd_sec_key_prefix")
	if len(secKillconf.EtcdConf.etcdSecKeyPrefix) == 0 {
		err=fmt.Errorf("init config from etcd key error%v",err)
		return
	}
	productkey:=beego.AppConfig.String("etcd_product_key")

	//拼接出 /liudz/seckill/product
secKillconf.EtcdConf.etcdSecProductkey=fmt.Sprintf("%s/%s",secKillconf.EtcdConf.etcdSecKeyPrefix,productkey)


	fmt.Printf("runmode is %s\n",beego.AppConfig.String("runmode"))
	logs.Debug("read addr %+v",redisAddr)
	logs.Debug("etcd addr %+v",etcdAddr)

	if len(redisAddr) ==0 || len(etcdAddr) ==0{
		err=fmt.Errorf("init config failed redis[%s],etcd[%s]",redisAddr,etcdAddr)
		return
	}

//从conf中获取3个变量的值，赋值给RedisConf struct
	redisMaxidle,err:=beego.AppConfig.Int("redis_max_idle")
	if err!=nil{
		err=fmt.Errorf("redis_max_idle err %v",err)
		return
	}

	redisMaxActive,err:=beego.AppConfig.Int("redis_max_active")
	if err!=nil{
		err=fmt.Errorf("redis_max_active err %v",err)
		return
	}

	redisIdleTimeout,err:=beego.AppConfig.Int("redis_idle_time")
	if err!=nil{
		err=fmt.Errorf("redis_idle_time err %v",err)
		return
	}
	secKillconf.RedisConf.redisMaxidle=redisMaxidle
	secKillconf.RedisConf.redisMaxActive=redisMaxActive
	secKillconf.RedisConf.redisIdleTimeout=redisIdleTimeout


//获取日志配置
secKillconf.loglevel=beego.AppConfig.String("log_level")
secKillconf.logpath=beego.AppConfig.String("log_path")
//获得日志的配置文件后，做一个初始化配置； initSec()
//整体思路：从配置文件读配置，赋值给struct，编一个初始化函数，实现初始化
	return
}