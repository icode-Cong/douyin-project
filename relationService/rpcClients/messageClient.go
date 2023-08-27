package rpcClients

import (
	"context"
	"fmt"
	etcdInit "relationService/rpcClients/etcd"
	"relationService/services/messageService"

	"github.com/micro/go-micro/v2"
)

func GetMessagesByIds(userId int64, friendIds []int64) ([]*messageService.Message, error) {
	messageMicroService := micro.NewService(micro.Registry(etcdInit.EtcdReg))
	messageServiceInstance := messageService.NewMessageService("messageService", messageMicroService.Client())
	fmt.Println("建立请求消息客户端：", messageServiceInstance)
	var req messageService.DouyinLatestMessagesRequest
	var idPairs []*messageService.UserIdPair
	for _, friednId := range friendIds {
		idPairs = append(idPairs, &messageService.UserIdPair{
			UserIdOne: userId,
			UserIdTwo: friednId,
		})
	}
	fmt.Printf("查找最新消息：%v", idPairs)
	req.IdPairs = idPairs

	resp, err := messageServiceInstance.LatestMessages(context.TODO(), &req)
	if err != nil || resp.StatusCode != 0 {
		fmt.Println("调用远程LatestMessages服务失败")
		fmt.Println(err)
	}
	return resp.MessageList, err
}
