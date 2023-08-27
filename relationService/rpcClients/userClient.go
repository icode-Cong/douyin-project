package rpcClients

import (
	"context"
	"fmt"
	etcdInit "relationService/rpcClients/etcd"
	"relationService/services/userService"

	"github.com/micro/go-micro/v2"
)

func UpdateFollowCount(userId int64, count int64, actionType int64) error {
	userMicroService := micro.NewService(micro.Registry(etcdInit.EtcdReg))
	userServiceInstance := userService.NewUserService("userService", userMicroService.Client())

	var req userService.DouyinUpdateFollowCountRequest

	req.UserId = userId
	req.Count = count
	req.ActionType = actionType

	resp, err := userServiceInstance.UpdateFollowCount(context.TODO(), &req)
	if err != nil || resp.StatusCode != 0 {
		fmt.Println("调用远程UpdateFollowCount服务失败")
		return err
	}
	return nil
}

func UpdateFollowerCount(userId int64, count int64, actionType int64) error {
	userMicroService := micro.NewService(micro.Registry(etcdInit.EtcdReg))
	userServiceInstance := userService.NewUserService("userService", userMicroService.Client())

	var req userService.DouyinUpdateFollowerCountRequest

	req.UserId = userId
	req.Count = count
	req.ActionType = actionType

	resp, err := userServiceInstance.UpdateFollowerCount(context.TODO(), &req)
	if err != nil || resp.StatusCode != 0 {
		fmt.Println("调用远程UpdateFollowerCount服务失败")
		return err
	}
	return nil
}

func GetUsersInfoByIds(userIds []int64) ([]*userService.User, error) {
	userMicroService := micro.NewService(micro.Registry(etcdInit.EtcdReg))
	userServiceInstance := userService.NewUserService("userService", userMicroService.Client())

	var req userService.DouyinMultiUserInfoRequest

	req.UserId = userIds

	resp, err := userServiceInstance.MultiUserInfo(context.TODO(), &req)
	if err != nil || resp.StatusCode != 0 {
		fmt.Println("调用远程MultiUserInfo服务失败")
	}
	return resp.Users, err
}

func GetUserInfoById(userId int64, token string) (*userService.User, error) {
	userMicroService := micro.NewService(micro.Registry(etcdInit.EtcdReg))
	userServiceInstance := userService.NewUserService("userService", userMicroService.Client())

	var req userService.DouyinUserRequest

	req.UserId = userId
	req.Token = token

	resp, err := userServiceInstance.UserInfo(context.TODO(), &req)
	if err != nil {
		fmt.Println("调用远程UserInfo服务失败，具体错误如下")
		fmt.Println(err)
	}

	user := &userService.User{
		Id:            resp.User.Id,
		Name:          resp.User.Name,
		FollowCount:   resp.User.FollowCount,
		FollowerCount: resp.User.FollowerCount,
		IsFollow:      resp.User.IsFollow,
	}
	return user, err
}
