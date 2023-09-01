package rpc_server

import (
	"context"
	etcdInit "feedService/rpcClients/etcd"
	"feedService/services/favoriteService"
	"fmt"

	"github.com/micro/go-micro/v2"
)

/*
*
调用user查询用户信息
*/
func IsFavorite(userId int64, videoId int64) (bool, error) {
	favoriteMicroService := micro.NewService(micro.Registry(etcdInit.EtcdReg))
	favoriteServiceInstance := favoriteService.NewFavoriteService("favoriteService", favoriteMicroService.Client())

	var req favoriteService.DouyinIsFavoriteRequest

	req.UserId = userId
	req.VideoId = videoId

	resp, err := favoriteServiceInstance.IsFavorite(context.TODO(), &req)
	if err != nil {
		fmt.Println("调用远程favoriteInfo服务失败，具体错误如下")
		fmt.Println(err)
	}

	return resp.IsFavorite, err
}
