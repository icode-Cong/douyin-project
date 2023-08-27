package handlers

import (
	"context"
	"gateway/services/commentService"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CommentAction(ginCtx *gin.Context) {
	var commentReq commentService.DouyinCommentActionRequest
	//获取request的信息
	commentReq.CommentId, _ = strconv.ParseInt(ginCtx.Query("comment_id"), 10, 64)
	commentReq.CommentText = ginCtx.Query("comment_text")
	actionType, _ := strconv.Atoi(ginCtx.Query("action_type"))
	commentReq.ActionType = int32(actionType)
	commentReq.Token = ginCtx.Query("token")
	commentReq.VideoId, _ = strconv.ParseInt(ginCtx.Query("video_id"), 10, 64)

	// 从gin.Key中取出服务实例
	commentServiceInstance := ginCtx.Keys["commentService"].(commentService.CommentService)
	//调用comment微服务，将context的上下文传入
	commentResp, _ := commentServiceInstance.CommentAction(context.Background(), &commentReq)

	//返回
	ginCtx.JSON(http.StatusOK, commentService.DouyinCommentActionResponse{
		StatusCode: commentResp.StatusCode,
		StatusMsg:  commentResp.StatusMsg,
		Comment:    commentResp.Comment,
	})

}
func CommentList(ginCtx *gin.Context) {
	var commentReq commentService.DouyinCommentListRequest

	commentReq.Token = ginCtx.Query("token")
	commentReq.VideoId, _ = strconv.ParseInt(ginCtx.Query("video_id"), 10, 64)

	commentServiceInstance := ginCtx.Keys["commentService"].(commentService.CommentService)
	commentResp, _ := commentServiceInstance.CommentList(context.Background(), &commentReq)

	ginCtx.JSON(http.StatusOK, commentService.DouyinCommentListResponse{
		StatusCode:  commentResp.StatusCode,
		StatusMsg:   commentResp.StatusMsg,
		CommentList: commentResp.CommentList,
	})

}
