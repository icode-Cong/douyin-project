package main

import (
	"favoriteService/configs"
	"favoriteService/internal/favoriteImp"
	"favoriteService/services/favoriteService"

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
		micro.Name("favoriteService"), // 微服务名字
		micro.Address("0.0.0.0:8085"),
		micro.Registry(etcdReg), // etcd注册件
		micro.Metadata(map[string]string{"protocol": "http"}),
	)
	microService.Init()
	favoriteService.RegisterFavoriteServiceHandler(microService.Server(), new(favoriteImp.FavoriteService))
	microService.Run()
}
