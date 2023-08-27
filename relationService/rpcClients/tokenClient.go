package rpcClients

import (
	"context"
	etcdInit "relationService/rpcClients/etcd"
	"relationService/services/tokenService"

	"github.com/micro/go-micro/v2"
)

func GetIdByToken(token string) (int64, error) {
	// 构造一个 tokenService 的客户端,要构造别的客户端，只要将这里所有的token改成对应的user,publish,.....就可以了
	tokenMicroService := micro.NewService(micro.Registry(etcdInit.EtcdReg))
	tokenServiceInstance := tokenService.NewTokenService("tokenService", tokenMicroService.Client())

	// 构造请求体
	var request tokenService.ParseTokenToIdRequest
	request.Token = token

	// 发送请求，调用接口
	response, err := tokenServiceInstance.ParseTokenToId(context.TODO(), &request)

	// 返回结果
	return int64(response.UserId), err
}
