package handlers

import (
	"bytes"
	"context"
	"fmt"
	publishService "gateway/services/publishService"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 上传视频
func Publish(ginCtx *gin.Context) {
	var publishReq publishService.DouyinPublishActionRequest

	publishReq.Title = ginCtx.PostForm("title")
	publishReq.Token = ginCtx.PostForm("token")
	fileHeader, _ := ginCtx.FormFile("data")
	fmt.Println(publishReq.Title)

	file, err := fileHeader.Open()
	if err != nil {
		PanicIfPublishError(err)
		return
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	if _, err = io.Copy(buf, file); err != nil {
		PanicIfPublishError(err)
		return
	}

	publishReq.Data = buf.Bytes()

	ctx, _ := context.WithTimeout(ginCtx, time.Minute*1)
	// 从gin.Key中取出服务实例
	fmt.Printf("here %v", ginCtx.Keys["publishService"])
	publishServiceInstance := ginCtx.Keys["publishService"].(publishService.PublishService)
	publishResp, err := publishServiceInstance.Publish(ctx, &publishReq)
	PanicIfPublishError(err)

	ginCtx.JSON(http.StatusOK, publishService.DouyinPublishActionResponse{
		StatusCode: publishResp.StatusCode,
		StatusMsg:  publishResp.StatusMsg,
	})
}

// 发布列表
func PublishList(ginCtx *gin.Context) {
	var publishReq publishService.DouyinPublishListRequest
	publishReq.Token = ginCtx.Query("token")
	ctx, _ := context.WithTimeout(ginCtx, time.Minute*1)

	//user_id绑定req.userId
	userId, err := strconv.ParseInt(ginCtx.Query("user_id"), 10, 64)
	PanicIfPublishError(err)
	publishReq.UserId = userId

	// 从gin.Key中取出服务实例
	publishServiceInstance := ginCtx.Keys["publishService"].(publishService.PublishService)
	publishResp, err := publishServiceInstance.PublishList(ctx, &publishReq)
	PanicIfPublishError(err)

	ginCtx.JSON(http.StatusOK, publishService.DouyinPublishListResponse{
		StatusCode: publishResp.StatusCode,
		StatusMsg:  publishResp.StatusMsg,
		VideoList:  publishResp.VideoList,
	})
}
