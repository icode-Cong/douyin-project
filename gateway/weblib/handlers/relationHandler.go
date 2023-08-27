package handlers

import (
	"context"
	"fmt"
	"gateway/services/relationService"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RelationAction no practical effect, just check if token is valid
func RelationAction(ginCtx *gin.Context) {
	var relationReq relationService.DouyinRelationActionRequest
	//获取request的信息
	relationReq.ToUserId, _ = strconv.ParseInt(ginCtx.Query("to_user_id"), 10, 64)
	actionType, _ := strconv.Atoi(ginCtx.Query("action_type"))
	relationReq.ActionType = int32(actionType)
	relationReq.Token = ginCtx.Query("token")
	fmt.Println("关注：", relationReq.ToUserId)

	// 从gin.Key中取出服务实例
	relationServiceInstance := ginCtx.Keys["relationService"].(relationService.RelationService)
	fmt.Println("得到服务实例：", relationServiceInstance)
	//调用comment微服务，将context的上下文传入
	relationResp, err := relationServiceInstance.RelationAction(context.Background(), &relationReq)
	PanicIfPublishError(err)
	//返回
	ginCtx.JSON(http.StatusOK, relationService.DouyinRelationActionResponse{
		StatusCode: relationResp.StatusCode,
		StatusMsg:  relationResp.StatusMsg,
	})
}

// FollowList all users have same follow list
func FollowList(ginCtx *gin.Context) {
	var relationReq relationService.DouyinRelationFollowerListRequest

	relationReq.Token = ginCtx.Query("token")
	relationReq.UserId, _ = strconv.ParseInt(ginCtx.Query("user_id"), 10, 64)

	relationServiceInstance := ginCtx.Keys["relationService"].(relationService.RelationService)
	resp, err := relationServiceInstance.FollowerList(context.Background(), &relationReq)
	PanicIfPublishError(err)

	ginCtx.JSON(http.StatusOK, relationService.DouyinRelationFollowerListResponse{
		StatusCode: resp.StatusCode,
		StatusMsg:  resp.StatusMsg,
		UserList:   resp.UserList,
	})
}

// FollowerList all users have same follower list
func FollowerList(ginCtx *gin.Context) {
	var relationReq relationService.DouyinRelationFollowListRequest

	relationReq.Token = ginCtx.Query("token")
	relationReq.UserId, _ = strconv.ParseInt(ginCtx.Query("user_id"), 10, 64)

	relationServiceInstance := ginCtx.Keys["relationService"].(relationService.RelationService)
	resp, err := relationServiceInstance.FollowList(context.Background(), &relationReq)
	PanicIfPublishError(err)

	ginCtx.JSON(http.StatusOK, relationService.DouyinRelationFollowListResponse{
		StatusCode: resp.StatusCode,
		StatusMsg:  resp.StatusMsg,
		UserList:   resp.UserList,
	})
}

// FriendList all users have same friend list
func FriendList(ginCtx *gin.Context) {
	var relationReq relationService.DouyinRelationFriendListRequest

	relationReq.Token = ginCtx.Query("token")
	relationReq.UserId, _ = strconv.ParseInt(ginCtx.Query("user_id"), 10, 64)

	relationServiceInstance := ginCtx.Keys["relationService"].(relationService.RelationService)
	resp, err := relationServiceInstance.FriendList(context.Background(), &relationReq)
	PanicIfPublishError(err)

	ginCtx.JSON(http.StatusOK, relationService.DouyinRelationFriendListResponse{
		StatusCode: resp.StatusCode,
		StatusMsg:  resp.StatusMsg,
		UserList:   resp.UserList,
	})
}
