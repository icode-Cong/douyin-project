package main

import (
	"userService/configs"
	"userService/internal/userImp"
	"userService/services/userService"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
)

func main() {
	// 初始化 mysql, redis 的连接
	configs.Init()
	// 注册 etcd
	etcdReg := etcd.NewRegistry(
		registry.Addrs("172.19.0.10:2379"),
	)
	// 注册微服务
	microService := micro.NewService(
		micro.Name("userService"),
		micro.Address("0.0.0.0:8082"),
		micro.Registry(etcdReg),
		micro.Metadata(map[string]string{"protocol": "http"}),
	)
	microService.Init()
	// 服务注册
	_ = userService.RegisterUserServiceHandler(microService.Server(), new(userImp.UserService))
	// _ = to_relation_proto.RegisterToRelationServiceHandler(microService.Server(), new(to_relation.ToRelationService))
	// _ = to_publish_proto.RegisterToPublishServiceHandler(microService.Server(), new(to_publish.ToPublishService))
	// _ = to_favorite_proto.RegisterToFavoriteServiceHandler(microService.Server(), new(to_favorite.ToFavoriteService))
	// 启动微服务
	_ = microService.Run()
}
