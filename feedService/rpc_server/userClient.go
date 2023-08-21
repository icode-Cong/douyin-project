package rpc_server

import (
	"context"
	etcdInit "feedService/rpc_server/etcd"
	"feedService/services/userService"
	"fmt"

	"github.com/micro/go-micro/v2"
)

/*
*
调用user查询用户信息
*/
func GetUserInfo(userId int64, token string) (*userService.User, error) {
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

/**
输入userId列表，查询User实体列表
*/
// func GetUsersInfo(userId []int64, token string) ([]*usersproto.User, error) {
// 	userMicroService := micro.NewService(micro.Registry(etcdInit.EtcdReg))
// 	usersService := usersproto.NewToRelationService("userService", userMicroService.Client())

// 	var req usersproto.GetUsersByIdsRequest

// 	req.UserId = userId
// 	req.Token = token

// 	resp, err := usersService.GetUsersByIds(context.TODO(), &req)
// 	if err != nil {
// 		fmt.Println("调用远程UserInfo服务失败，具体错误如下")
// 		fmt.Println(err)
// 	}

// 	return resp.UserList, err
// }
