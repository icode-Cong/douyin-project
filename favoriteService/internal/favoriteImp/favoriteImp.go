package favoriteImp

import (
	"context"
	"favoriteService/models"
	"favoriteService/services/favoriteService"
	"userService/rpcClients"
)

// 请在此完成服务的业务逻辑
type FavoriteService struct {
}

func (*FavoriteService) FavoriteAction(ctx context.Context, request *favoriteService.DouyinFavoriteActionRequest, response *favoriteService.DouyinFavoriteActionResponse) error {
	videoId := request.VideoId
	actionType := request.ActionType
	if request.Token == "" {
		//response todo
		return nil
	}
	//从 token 中解析出登录用户的 id
	userId, err := rpcClients.GetIdByToken(request.Token)
	if err != nil {
		// response.StatusCode = -1
		// response.StatusMsg = "登录失效，请重新登录"
		// response.User = &userService.User{}
		// return nil
	}
	video, err := models.NewVideoDaoInstance().FindVideoById(videoId)
	if err != nil {
		return err
	}
	if actionType == 1 {
		//视频获赞数+1
		models.NewVideoDaoInstance().AddFavoriteCount(videoId, 1)
		//用户喜欢数+1
		models.NewUserDaoInstance().AddFavoriteCount(userId, 1)
		//插入一条点赞记录
		models.NewFavDaoInstance().AddFavRec(userId, videoId)
		//视频作者被点赞数+1
		models.NewUserDaoInstance().AddTotalFavorited(video.UserId, 1)
	} else if actionType == 2 {
		models.NewVideoDaoInstance().AddFavoriteCount(videoId, -1)
		models.NewUserDaoInstance().AddFavoriteCount(userId, -1)
		models.NewFavDaoInstance().DelFavRec(userId, videoId)
		models.NewUserDaoInstance().AddTotalFavorited(video.UserId, -1)
	}
	return nil
}

func (*FavoriteService) FavoriteList(ctx context.Context, request *favoriteService.DouyinFavoriteListRequest, response *favoriteService.DouyinFavoriteListResponse) error {
	if request.Token == "" {
		//response todo
		return nil
	}
	//从 token 中解析出登录用户的 id
	userId, err := rpcClients.GetIdByToken(request.Token)
	if err != nil {
		return err
	}
	videoIdList, err := models.NewFavDaoInstance().GetFavVideoIdList(userId)
	if err != nil {
		return err
	}
	videoList := make([]*favoriteService.Video, 0)
	for _, id := range videoIdList {
		video, err := models.NewVideoDaoInstance().FindVideoById(id)
		if err != nil {
			return err
		}
		authorId := video.UserId
		author, err := models.NewUserDaoInstance().GetUserById(authorId)
		if err != nil {
			return err
		}
		videoList = append(videoList, &favoriteService.Video{
			Id:      video.VideoId,
			PlayUrl: video.PlayUrl,
			Author: &favoriteService.User{
				Id:              author.UserId,
				Name:            author.Name,
				FollowCount:     author.FollowingCount,
				FollowerCount:   author.FollowerCount,
				IsFollow:        false, //unknow
				Avatar:          author.Avatar,
				BackgroundImage: author.BackgroundImage,
				Signature:       author.Signature,
				TotalFavorited:  author.TotalFavorited,
				WorkCount:       author.WorkCount,
				FavoriteCount:   author.WorkCount,
			},
			CoverUrl:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    true,
			Title:         video.Title,
		})
	}
	response = &favoriteService.DouyinFavoriteListResponse{
		StatusCode: 0,
		StatusMsg:  "获取点赞视频列表成功",
		VideoList:  videoList,
	}
	return nil
}
