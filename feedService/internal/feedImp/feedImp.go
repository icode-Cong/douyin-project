package feedImp

import (
	"context"
	"feedService/models"
	rpcClients "feedService/rpcClients"
	"feedService/services/feedService"
	"fmt"

	// "feedService/services/favorite_to_video_proto"
	// usersproto "feedService/services/to_relation"
	"feedService/services/userService"
	"time"
)

type FeedService struct {
}

// 数据流
func (*FeedService) Feed(ctx context.Context, req *feedService.DouyinFeedRequest, resp *feedService.DouyinFeedResponse) error {

	// 读取时间戳并转换格式，获得查询时间
	latestTimeStamp := req.LatestTime
	fmt.Println("得到请求时间戳", latestTimeStamp)
	format := "2006-01-02 15:04:05"
	var t time.Time
	if latestTimeStamp < 9999999999 {
		t = time.Unix(latestTimeStamp, 0)
	} else {
		t = time.Unix(latestTimeStamp/1000, 0)
	}
	latestTime := t.Format(format)
	fmt.Printf("查询时间 %v\n", latestTime)
	//获得登录用户Id
	loginUserId, err := rpcClients.GetIdByToken(req.Token)
	isLogin := true
	if err != nil {
		isLogin = false
		fmt.Printf("获取登录用户Id失败，userToken:%v", req.Token)
	}

	// 从数据库中按时间查询视频列表
	videos := models.NewVideoDaoInstance().GetVideosByTime(&latestTime, 5)

	// 对获得的视频列表中的每个视频，查询其作者信息
	var videoResult []*feedService.Video
	var authors []*userService.User

	if true {
		for _, video := range videos {
			authorId := video.UserId
			author, _ := rpcClients.GetUserInfo(authorId, req.Token)
			authors = append(authors, author)
		}
		for _, video := range videos {
			for _, user := range authors {
				if video.UserId == user.Id {
					isFavorite := false
					if isLogin {
						isFavorite, _ = rpcClients.IsFavorite(loginUserId, video.VideoId)
					}
					videoResult = append(videoResult, BuildProtoVideo(video, user, isFavorite))
					break
				}
			}
		}
	}
	if len(videos) > 0 {
		resp.NextTime = videos[len(videos)-1].CreateAt.Unix()
	} else {
		resp.NextTime = time.Now().Unix()
	}
	resp.StatusCode = 0
	resp.StatusMsg = "查询视频成功"
	resp.VideoList = videoResult

	return nil
}

func BuildProtoVideo(video *models.Video, user *userService.User, isFavorite bool) *feedService.Video {
	return &feedService.Video{
		Id:            video.VideoId,
		Author:        BuildProtoUser(user),
		PlayUrl:       video.PlayUrl,
		CoverUrl:      video.CoverUrl,
		FavoriteCount: video.FavoriteCount,
		CommentCount:  video.CommentCount,
		IsFavorite:    isFavorite,
		Title:         video.Title,
	}
}

func BuildProtoUser(user *userService.User) *feedService.User {
	return &feedService.User{
		Id:             user.Id,
		Name:           user.Name,
		FollowCount:    user.FollowCount,
		FollowerCount:  user.FollowerCount,
		IsFollow:       user.IsFollow,
		TotalFavorited: user.TotalFavorited,
		WorkCount:      user.WorkCount,
		FavoriteCount:  user.FavoriteCount,
	}
}
