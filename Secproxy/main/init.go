package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"time"
	etcd_client "go.etcd.io/etcd/clientv3"
	"encoding/json"
	"fmt"
	"context"
)


func initEtcd()(err error){
	fmt.Println(b)
	cli,err:=etcd_client.New(etcd_client.Config{
		Endpoints:[]string{	secKillconf.EtcdConf.etcdAddr},
		DialTimeout:  time.Duration(secKillconf.EtcdConf.timeout) * time.Second,

	})
	fmt.Println("etcd_client.New success")
	if err != nil {
		// handle error!
		logs.Error("error",err)
		return
	}
	//这句很重要
	etcdclient=cli
	return

}


var (
	redisPool *redis.Pool
	etcdclient *etcd_client.Client

)

func initRedis()(err error){
	redisPool=&redis.Pool{
	MaxIdle:secKillconf.RedisConf.redisIdleTimeout,
	MaxActive:secKillconf.RedisConf.redisMaxActive,
	IdleTimeout:time.Duration(secKillconf.RedisConf.redisIdleTimeout)*time.Second,
	Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp",secKillconf.RedisConf.redisAddr)
	},
}

conn:=redisPool.Get()
defer conn.Close()

_,err=conn.Do("ping")
if err!=nil{
logs.Error("ping redis failed % v",err)

}
	return

}

//将文件中的日志级别转换成标准的,级别是整型变量
func convert(level string) int{
	//switch {}  switch后面不写，默认及时true
	switch(level){
	case "debug":
		return logs.LevelDebug
	case "warn":
		return logs.LevelWarning
	case "info":
		return logs.LevelInfo
	}

	return logs.LevelDebug
}

//初始化日志配置的函数
func initLog()(err error){
	config:=make(map[string]interface{})
	config["filename"]=secKillconf.logpath
	config["level"]=convert(secKillconf.loglevel)
configStr,err:=json.Marshal(config)
	if err!=nil{
		logs.Error("configStr failed,err %v",err)
		return
	}

	logs.SetLogger(logs.AdapterFile,string(configStr))
return
}



//读取etcd中的秒杀产品信息
func loadSecConf()(err error) {
	//secKillconf.EtcdConf.etcd_sec_key
	//字符串拼接
/*	key := fmt.Sprintf("%s/product", secKillconf.EtcdConf.etcdseckey)
	fmt.Println("key is", key)
*/

	resp, err := etcdclient.Get(context.Background(), secKillconf.EtcdConf.etcdSecProductkey)
	if err != nil {
		logs.Error("get %v from etcd failed ,err %v", secKillconf.EtcdConf.etcdSecProductkey, err)
		return
	}
	//把返回的json字符串转换成对象
	//这个例子，只有一个key
	//从etcd取出的东西，要全局保存成结构体
	var secProductInfo []SecProductInfoConf

	for k, v := range resp.Kvs {
		logs.Debug("key[%s] values[%s] 	", k, v)
		err := json.Unmarshal(v.Value, &secProductInfo)
		if err != nil {
		}
		logs.Error("json Unmarshal result: %v", secProductInfo)
	}
//保存在secKillconf
secKillconf.secProduction=secProductInfo
	return

}

func initSec()(err error){

	//log 在本文件中定义
	err=initLog()
	if err!=nil{
		logs.Error("init log failed,err %v",err)
		return
	}


	//etcd
	err=initEtcd()
	if err!=nil{
		logs.Error("init etcd failed,err %v",err)
		return
	}

	//redis
	err=initRedis()
	if err!=nil{
		logs.Error("init redis failed,err %v",err)
		return
	}

	//读取etcd中的内容
	err=loadSecConf()
    logs.Info("init sec succ")
	return
}

