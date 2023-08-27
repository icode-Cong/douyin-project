package handlers

import (
	"context"
	"fmt"
	messageService "gateway/services/messageService"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

/*
*
聊天记录
*/
func MessageChat(ginCtx *gin.Context) {
	var messageReq messageService.DouyinMessageChatRequest

	messageReq.Token = ginCtx.Query("token")
	messageReq.ToUserId, _ = strconv.ParseInt(ginCtx.Query("to_user_id"), 10, 64)
	messageReq.PreMsgTime, _ = strconv.ParseInt(ginCtx.Query("pre_msg_time"), 10, 64)
	fmt.Println("查询to_user_id: ", messageReq.ToUserId)

	messageServiceInstance := ginCtx.Keys["messageService"].(messageService.MessageService)
	resp, _ := messageServiceInstance.MessageChat(context.Background(), &messageReq)

	ginCtx.JSON(http.StatusOK, messageService.DouyinMessageChatResponse{
		StatusCode:  resp.StatusCode,
		StatusMsg:   resp.StatusMsg,
		MessageList: resp.MessageList,
	})
}

/*
*
发送消息
*/
func MessageAction(ginCtx *gin.Context) {
	var messageReq messageService.DouyinMessageActionRequest

	messageReq.Token = ginCtx.Query("token")
	messageReq.ToUserId, _ = strconv.ParseInt(ginCtx.Query("to_user_id"), 10, 64)
	actionType, _ := strconv.Atoi(ginCtx.Query("action_type"))
	messageReq.ActionType = int32(actionType)
	messageReq.Content = ginCtx.Query("content")

	messageServiceInstance := ginCtx.Keys["messageService"].(messageService.MessageService)
	resp, _ := messageServiceInstance.MessageAction(context.Background(), &messageReq)

	ginCtx.JSON(http.StatusOK, messageService.DouyinMessageActionResponse{
		StatusCode: resp.StatusCode,
		StatusMsg:  resp.StatusMsg,
	})
}
