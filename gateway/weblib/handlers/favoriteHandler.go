package handlers

import (
	"gateway/services/favoriteService"
	"net/http"

	"github.com/gin-gonic/gin"

	"strconv"
)

func FavoriteAction(ginCtx *gin.Context) {
	actionType, _ := strconv.Atoi(ginCtx.Query("action_type"))
	token := ginCtx.Query("token")
	vid, _ := strconv.ParseInt(ginCtx.Query("video_id"), 10, 64)

	favServiceInstance := ginCtx.Keys["favoriteService"].(favoriteService.FavoriteService)
	var req favoriteService.DouyinFavoriteActionRequest
	req.ActionType = int32(actionType)
	req.Token = token
	req.VideoId = vid
	action, _ := favServiceInstance.FavoriteAction(ginCtx, &req)

	ginCtx.JSON(http.StatusOK, favoriteService.DouyinFavoriteActionResponse{
		StatusCode: action.StatusCode,
		StatusMsg:  action.StatusMsg,
	})

}
func FavoriteList(ginCtx *gin.Context) {

	token := ginCtx.Query("token")

	uid, _ := strconv.ParseInt(ginCtx.Query("user_id"), 10, 64)

	favServiceInstance := ginCtx.Keys["favoriteService"].(favoriteService.FavoriteService)

	var req favoriteService.DouyinFavoriteListRequest

	req.Token = token
	req.UserId = uid
	action, _ := favServiceInstance.FavoriteList(ginCtx, &req)

	ginCtx.JSON(http.StatusOK, favoriteService.DouyinFavoriteListResponse{
		StatusCode: action.StatusCode,
		StatusMsg:  action.StatusMsg,
		VideoList:  action.VideoList,
	})

}
