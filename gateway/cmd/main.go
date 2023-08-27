package main

import (
	"fmt"
	"gateway/services/commentService"
	"gateway/services/favoriteService"
	"gateway/services/feedService"
	"gateway/services/messageService"
	"gateway/services/publishService"
	"gateway/services/relationService"
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

	// message 消息微服务
	messageMicroService := micro.NewService(
		micro.Name("messageService.client"),
		micro.WrapClient(wrappers.NewMessageWrapper),
	)
	messageService := messageService.NewMessageService("messageService", messageMicroService.Client())

	// relation 关注微服务
	relationMicroService := micro.NewService(
		micro.Name("relationService.client"),
		micro.WrapClient(wrappers.NewRelationWrapper),
	)
	relationService := relationService.NewRelationService("relationService", relationMicroService.Client())

	// favorite 点赞微服务
	favoriteMicroService := micro.NewService(
		micro.Name("favoriteService.client"),
		micro.WrapClient(wrappers.NewFavoriteWrapper),
	)
	favoriteService := favoriteService.NewFavoriteService("favoriteService", favoriteMicroService.Client())

	// comment 评论微服务
	commentMicroService := micro.NewService(
		micro.Name("commentService.client"),
		micro.WrapClient(wrappers.NewCommentWrapper),
	)
	commentService := commentService.NewCommentService("commentService", commentMicroService.Client())

	serviceMap := make(map[string]interface{})
	serviceMap["userService"] = userService
	serviceMap["publishService"] = publishService
	serviceMap["feedService"] = feedService
	serviceMap["messageService"] = messageService
	serviceMap["relationService"] = relationService
	serviceMap["favoriteService"] = favoriteService
	serviceMap["commentService"] = commentService

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
