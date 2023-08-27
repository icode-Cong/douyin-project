package main

import (
	"commentService/configs"
	"commentService/internal/commentImp"
	"commentService/services/commentService"

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
		micro.Name("commentService"), // 微服务名字
		micro.Address("0.0.0.0:8091"),
		micro.Registry(etcdReg), // etcd注册件
		micro.Metadata(map[string]string{"protocol": "http"}),
	)
	microService.Init()
	commentService.RegisterCommentServiceHandler(microService.Server(), new(commentImp.CommentService))
	microService.Run()
}
