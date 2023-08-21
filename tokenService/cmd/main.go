package main

import (
	tokenservice "tokenService/internal/tokenService"
	"tokenService/protos/token"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
)

func main() {
	etcdReg := etcd.NewRegistry(
		registry.Addrs("172.19.0.10:2379"),
	)
	microService := micro.NewService(
		micro.Name("tokenService"),
		micro.Address("0.0.0.0:8085"),
		micro.Registry(etcdReg),
		micro.Metadata(map[string]string{"protocol": "http"}),
	)
	microService.Init()
	token.RegisterTokenServiceHandler(microService.Server(), new(tokenservice.TokenService))
	microService.Run()
}
