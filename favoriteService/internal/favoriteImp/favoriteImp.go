package favoriteImp

import (
	"context"
	"favoriteService/models"
	"favoriteService/rpcClients"
	"favoriteService/services/favoriteService"
	"favoriteService/services/publishService"
	"fmt"
	"sync"
)

// 请在此完成服务的业务逻辑
type FavoriteService struct {
}

func (*FavoriteService) FavoriteAction(ctx context.Context, request *favoriteService.DouyinFavoriteActionRequest, response *favoriteService.DouyinFavoriteActionResponse) error {
	token := request.Token
	videoId := request.VideoId
	actionType := request.ActionType

	//从 token 中解析出登录用户的 id
	loginUserId, err := rpcClients.GetIdByToken(token)
	if err != nil {
		response.StatusCode = -1
		response.StatusMsg = "登录失效，请重新登录"
		return nil
	}

	favorite := &models.Favorite{
		UserId:  loginUserId,
		VideoId: videoId,
	}

	authorId, err := rpcClients.GetAuthorIdByVideoId(videoId)
	if err != nil {
		fmt.Println("获取视频作者id失败")
	}
	if actionType == 1 {
		_ = models.NewFavoriteDaoInstance().CreateFavorite(favorite)
		rpcClients.UpdateFavoriteCount(loginUserId, 1, int64(actionType))
		rpcClients.UpdateFavoritedCount(authorId, 1, int64(actionType))
		response.StatusCode = 0
		response.StatusMsg = "点赞成功"
		return nil
	} else if actionType == 2 {
		_ = models.NewFavoriteDaoInstance().DeleteFavorite(favorite)
		rpcClients.UpdateFavoriteCount(loginUserId, 1, int64(actionType))
		rpcClients.UpdateFavoritedCount(authorId, 1, int64(actionType))
		response.StatusCode = 0
		response.StatusMsg = "取消点赞成功"
		return nil
	} else {
		response.StatusCode = -1
		response.StatusMsg = "actionType出错"
		return nil
	}
}

func (*FavoriteService) FavoriteList(ctx context.Context, request *favoriteService.DouyinFavoriteListRequest, response *favoriteService.DouyinFavoriteListResponse) error {
	userId := request.UserId
	videoIdList := models.NewFavoriteDaoInstance().GetFavoriteVideoIdList(userId)

	// 使用协程来逐个获取关注对象的用户信息
	// 使用通道来传递用户信息
	videoInfoChan := make(chan *publishService.Video)
	// var userList []*relationService.User
	var wg sync.WaitGroup
	for _, videoId := range videoIdList {
		wg.Add(1)
		go func(id int64) {
			defer wg.Done()
			videoInfo, err := rpcClients.GetVideoInfoById(id)
			if err != nil {
				fmt.Println("[Error : 查询视频列表, 向 videoService 请求视频信息失败，视频id：]", id)
				return
			}
			videoInfoChan <- videoInfo
		}(videoId)
	}

	// 等待所有协程完成
	go func() {
		wg.Wait()
		close(videoInfoChan) // 关闭数据，表示不再写入数据
	}()

	// 收集用户信息
	var videoList []*favoriteService.Video
	for videoInfo := range videoInfoChan {
		videoList = append(videoList, BuildVideo(videoInfo))
	}

	response.StatusCode = 0
	response.StatusMsg = "查询关注列表成功"
	response.VideoList = videoList
	return nil
}
func (*FavoriteService) IsFavorite(ctx context.Context, request *favoriteService.DouyinIsFavoriteRequest, response *favoriteService.DouyinIsFavoriteResponse) error {
	userId := request.UserId
	videoId := request.VideoId
	isFav := models.NewFavoriteDaoInstance().IsFavorite(userId, videoId)
	response.IsFavorite = isFav
	return nil
}

func BuildVideo(video *publishService.Video) *favoriteService.Video {
	return &favoriteService.Video{
		Id:            video.Id,
		Author:        BuildUser(video.Author),
		PlayUrl:       video.PlayUrl,
		CoverUrl:      video.CoverUrl,
		FavoriteCount: video.FavoriteCount,
		CommentCount:  video.CommentCount,
		IsFavorite:    video.IsFavorite,
		Title:         video.Title,
	}
}

func BuildUser(user *publishService.User) *favoriteService.User {
	return &favoriteService.User{
		Id:            user.Id,
		Name:          user.Name,
		FollowCount:   user.FollowCount,
		FavoriteCount: user.FollowerCount,
		IsFollow:      user.IsFollow,
	}
}
