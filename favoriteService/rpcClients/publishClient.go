package rpcClients

import (
	"context"
	etcdInit "favoriteService/rpcClients/etcd"
	"favoriteService/services/publishService"
	"fmt"

	"github.com/micro/go-micro/v2"
)

/*
*
查询作者id
*/
func GetAuthorIdByVideoId(videoId int64) (int64, error) {
	publishMicroService := micro.NewService(micro.Registry(etcdInit.EtcdReg))
	publishServiceInstance := publishService.NewPublishService("publishService", publishMicroService.Client())

	var req publishService.GetAuthorIdRequest

	req.VideoId = videoId

	resp, err := publishServiceInstance.GetAuthorId(context.TODO(), &req)
	if err != nil {
		fmt.Println("调用远程GetAutorhId服务失败，具体错误如下")
		fmt.Println(err)
	}

	return resp.AuthorId, err
}

func GetVideoInfoById(videoId int64) (*publishService.Video, error) {
	publishMicroService := micro.NewService(micro.Registry(etcdInit.EtcdReg))
	publishServiceInstance := publishService.NewPublishService("publishService", publishMicroService.Client())

	var req publishService.GetVideoInfoRequest

	req.VideoId = videoId

	resp, err := publishServiceInstance.GetVideoInfo(context.TODO(), &req)
	if err != nil {
		fmt.Println("调用远程GetVideoInfo服务失败，具体错误如下")
		fmt.Println(err)
	}

	return resp.VideoInfo, err
}
