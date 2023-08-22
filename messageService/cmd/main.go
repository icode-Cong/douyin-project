package main

import (
	"messageService/configs"
	"messageService/internal/messageImp"
	"messageService/services/messageService"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
)

func main() {
	configs.Init()
	etcdReg := etcd.NewRegistry(
		registry.Addrs("172.19.0.10:2379"),
	)
	microService := micro.NewService(
		micro.Name("messageService"), // 微服务名字
		micro.Address("0.0.0.0:8089"),
		micro.Registry(etcdReg), // etcd注册件
		micro.Metadata(map[string]string{"protocol": "http"}),
	)
	microService.Init()
	messageService.RegisterMessageServiceHandler(microService.Server(), new(messageImp.MessageService))
	microService.Run()
}
