package main

import (
	"fmt"
	"gateway/services/feedService"
	"gateway/services/publishService"
	"gateway/services/user"
	"gateway/weblib"
	"gateway/wrappers"
	"time"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-micro/v2/web"
)

func main() {
	etcdReg := etcd.NewRegistry(
		registry.Addrs("172.19.0.10:2379"),
	)

	// userService 用户微服务
	userMicroService := micro.NewService(
		micro.Name("userService.client"),
		micro.WrapClient(wrappers.NewUserWrapper),
	)
	userService := user.NewUserService("userService", userMicroService.Client())

	// publishService 投稿微服务
	publishMicroService := micro.NewService(
		micro.Name("publishService.client"),
		micro.WrapClient(wrappers.NewPublishWrapper),
	)
	publishService := publishService.NewPublishService("publishService", publishMicroService.Client())
	fmt.Printf("ohh %v\n", publishService)

	// feed 视频流微服务
	feedMicroService := micro.NewService(
		micro.Name("feedService.client"),
		micro.WrapClient(wrappers.NewFeedWrapper),
	)
	feedService := feedService.NewFeedService("feedService", feedMicroService.Client())

	serviceMap := make(map[string]interface{})
	serviceMap["userService"] = userService
	serviceMap["publishService"] = publishService
	serviceMap["feedService"] = feedService

	server := web.NewService(
		web.Name("httpService"),
		web.Address("0.0.0.0:4000"),
		//将服务调用实例使用gin处理
		web.Handler(weblib.NewRouter(serviceMap)),
		web.Registry(etcdReg),
		web.RegisterTTL(time.Second*300),
		web.RegisterInterval(time.Second*150),
		web.Metadata(map[string]string{"protocol": "http"}),
	)
	server.Init()
	server.Run()
}