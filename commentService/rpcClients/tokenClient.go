package rpcClients

import (
	etcdInit "commentService/rpcClients/etcd"
	tokenproto "commentService/services/tokenService"
	"context"
	"fmt"

	"github.com/micro/go-micro/v2"
)

/*
*
判断 token 是否有效，若有效则返回 userId，否则 err
*/
func GetIdByToken(token string) (int64, error) {
	tokenMicroService := micro.NewService(micro.Registry(etcdInit.EtcdReg))
	tokenService := tokenproto.NewTokenService("tokenService", tokenMicroService.Client())

	var req tokenproto.ParseTokenToIdRequest
	req.Token = token

	resp, err := tokenService.ParseTokenToId(context.TODO(), &req)
	if err != nil {
		fmt.Println(err)
	}

	return int64(resp.UserId), err
}
