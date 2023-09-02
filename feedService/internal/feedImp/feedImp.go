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
	if latestTimeStamp == 0 {
		//若没有时间戳，默认为当前时间
		latestTimeStamp = time.Now().Unix()
	}
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
	isLogined := false
	var loginUserId int64
	if req.Token != "" {
		var err error
		loginUserId, err = rpcClients.GetIdByToken(req.Token)
		if err == nil {
			isLogined = true
		} else {
			fmt.Printf("获取登录用户Id失败，userToken:%v", req.Token)
		}
	}

	// 从数据库中按时间查询视频列表
	videos := models.NewVideoDaoInstance().GetVideosByTime(&latestTime, 5)

	//返回结果
	var videoResult []*feedService.Video
	for _, video := range videos {
		authorId := video.UserId
		//查询作者信息
		author, err := rpcClients.GetUserInfo(authorId, req.Token)
		if err != nil {
			fmt.Printf("查询作者信息失败\nauthorId:%v\ntoken:%v", authorId, req.Token)
		}
		isFavorite := false
		if isLogined {
			isFavorite, err = rpcClients.IsFavorite(loginUserId, video.VideoId)
			if err != nil {
				fmt.Printf("查询点赞信息失败\nloginUserId:%v\nvideoId:%v", loginUserId, video.VideoId)
			}
		}
		videoResult = append(videoResult, BuildProtoVideo(video, author, isFavorite))
	}
	if len(videos) > 0 {
		resp.NextTime = videos[len(videos)-1].CreateAt.Unix()
	} else {
		//修改为app传来的参数，防止获取视频重复
		resp.NextTime = t.Unix()
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
