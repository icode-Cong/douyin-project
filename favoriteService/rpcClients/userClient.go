package rpcClients

import (
	"context"
	etcdInit "favoriteService/rpcClients/etcd"
	"favoriteService/services/userService"
	"fmt"

	"github.com/micro/go-micro/v2"
)

func UpdateFavoriteCount(userId int64, count int64, actionType int64) error {
	userMicroService := micro.NewService(micro.Registry(etcdInit.EtcdReg))
	userServiceInstance := userService.NewUserService("userService", userMicroService.Client())

	var req userService.DouyinUpdateFavoriteCountRequest

	req.UserId = userId
	req.Count = count
	req.ActionType = actionType

	resp, err := userServiceInstance.UpdateFavoriteCount(context.TODO(), &req)
	if err != nil || resp.StatusCode != 0 {
		fmt.Println("调用远程UpdateFavoriteCount服务失败")
		return err
	}
	return nil
}

func UpdateFavoritedCount(userId int64, count int64, actionType int64) error {
	userMicroService := micro.NewService(micro.Registry(etcdInit.EtcdReg))
	userServiceInstance := userService.NewUserService("userService", userMicroService.Client())

	var req userService.DouyinUpdateTotalFavoritedRequest

	req.UserId = userId
	req.Count = count
	req.ActionType = actionType

	resp, err := userServiceInstance.UpdateTotalFavorited(context.TODO(), &req)
	if err != nil || resp.StatusCode != 0 {
		fmt.Println("调用远程UpdateTotalFavorited服务失败")
		return err
	}
	return nil
}
