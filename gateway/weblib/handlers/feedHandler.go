package handlers

import (
	"context"
	"fmt"
	feedService "gateway/services/feedService"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 视频流
func Feed(ginCtx *gin.Context) {
	//数据绑定
	var feedReq feedService.DouyinFeedRequest
	lastTime, err := strconv.ParseInt(ginCtx.Query("latest_time"), 10, 64)
	PanicIfFeedError(err)
	feedReq.LatestTime = lastTime
	feedReq.Token = ginCtx.Query("token")

	//设置最长响应时长
	ctx, _ := context.WithTimeout(ginCtx, time.Minute*1)

	// 从gin.Key中取出服务实例
	feedServiceInstance := ginCtx.Keys["feedService"].(feedService.FeedService)
	// 调用
	fmt.Printf("feedServiceInstance %v\n", feedServiceInstance)
	feedResp, err := feedServiceInstance.Feed(ctx, &feedReq)
	PanicIfFeedError(err)
	//返回
	ginCtx.JSON(http.StatusOK, feedService.DouyinFeedResponse{
		StatusCode: feedResp.StatusCode,
		StatusMsg:  feedResp.StatusMsg,
		VideoList:  feedResp.VideoList,
		NextTime:   feedResp.NextTime,
	})
}
