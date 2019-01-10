package main

import (
	"go.etcd.io/etcd/clientv3"
	"time"
	"fmt"
	"encoding/json"
	"context"
)
type SecInfoConf struct{
	ProductId int
	StartTime int
	EndTime int
	Status int
	Total int
	Left int
}

const(
	EtcdKey="/liudz/secskill/product"
)
func SetLogConfToEtcd (){
cli,err := clientv3.New(clientv3.Config{
	Endpoints:[]string{"192.168.37.104:2379"},
	DialTimeout:5*time.Second,
	})
	if err != nil{
		fmt.Println("connect failed,err ",err)
			}
			fmt.Println("connect succ")
defer cli.Close()
var SecInfoConfArr []SecInfoConf
	SecInfoConfArr=append(SecInfoConfArr,
		SecInfoConf{
			ProductId:1022,
			StartTime:1517495965,
			EndTime:1517495666,
			Status:0,
			Total:10000,
			Left:10000},
		SecInfoConf{
			ProductId:1032,
			StartTime:1517495965,
			EndTime:1517495666,
			Status:0,
			Total:9000,
			Left:9000},
		)

	data,err:=json.Marshal(SecInfoConfArr)
	if err != nil{
		fmt.Println("json failed",err)
		return
	}
	ctx,cancel:=context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	_,err=cli.Put(ctx,EtcdKey,string(data))
	//cancel()
	if err != nil{
		fmt.Println("put failed:",err)
		return
	}

	ctx,_=context.WithTimeout(context.Background(),time.Second)
resp,err :=	cli.Get(ctx,EtcdKey)
cancel()
	if err != nil{
		fmt.Println("get failed",err)
		return
	}
for _,env := range resp.Kvs{
	fmt.Printf("%s:%s\n",env.Key,env.Value)
}

}

func main(){
	SetLogConfToEtcd()
}