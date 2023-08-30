package publishImp

import (
	"context"
	"fmt"
	"publishService/configs"
	"publishService/models"
	"publishService/rpcClients"
	"publishService/services/publishService"
	"publishService/services/userService"
	utils "publishService/utils"
	"time"

	"github.com/gofrs/uuid"
)

type PublishService struct {
}

// 发布视频
func (*PublishService) Publish(ctx context.Context, request *publishService.DouyinPublishActionRequest, response *publishService.DouyinPublishActionResponse) error {
	fmt.Println("被调用")
	// 调用 tokenService 微服务, 通过 token 判断是否登录,如果查不到对应的 userId ,如果查不到则返回提示,要求重新登录
	userId, err := rpcClients.GetIdByToken(request.Token)
	if err != nil {
		response.StatusCode = -1
		response.StatusMsg = "登录失效，请重新登录"
		return nil
	}
	title := request.Title

	// 使用 uuid 为视频生成一个随机的文件名,并构造其在 OSS 上的 url
	videoUUID, _ := uuid.NewV4()
	videoPath := time.Now().Format("2006-01-02") + "/" + videoUUID.String() + ".mp4"
	videoUrl := "https://" + configs.Conf.OssConf.BucketName + "." + configs.Conf.OssConf.Endpoint + "/" + videoPath

	// 使用 uuid 为视频封面生成一个随机的文件名,并构造其在 OSS 上的 url
	coverUUID, _ := uuid.NewV4()
	coverPath := time.Now().Format("2006-01-02") + "/" + coverUUID.String() + ".jpg"
	coverUrl := "https://" + configs.Conf.OssConf.BucketName + "." + configs.Conf.OssConf.Endpoint + "/" + coverPath

	// 开一个协程来负责上传视频和封面,避免等待过久
	go func() {
		utils.UploadVideo(videoPath, request.Data)
		// 使用 ffmpeg 截取一帧来作为视频封面
		// time.Sleep(2 * time.Second)
		coverFrame, _ := utils.CutFrameAsCover(videoUrl)
		_ = utils.UploadCover(coverPath, coverFrame)
	}()
	// 调用 userService 微服务, 更新用户的作品数
	rpcClients.UpdateWorkCount(userId, 1, 1)

	// 构造 video 行记录,并添加到数据库中
	video := &models.Video{
		UserId:        userId,
		Title:         title,
		CoverUrl:      coverUrl,
		PlayUrl:       videoUrl,
		FavoriteCount: 0,
		CommentCount:  0,
		CreateAt:      time.Now(),
		UpdateAt:      time.Now(),
	}
	if _, err := models.NewVideoDaoInstance().CreateVideo(video); err != nil {
		response.StatusCode = -1
		response.StatusMsg = "上传视频失败"
		return err
	}

	response.StatusCode = 0
	response.StatusMsg = "上传视频成功"
	return nil
}

// 获取发布视频列表
func (*PublishService) PublishList(ctx context.Context, request *publishService.DouyinPublishListRequest, response *publishService.DouyinPublishListResponse) error {
	// 调用 tokenService 微服务, 只有登录用户才能查看视频列表
	if request.Token == "" {
		response.StatusCode = -1
		response.StatusMsg = "登录失效，请重新登录"
		return nil
	}
	_, err := rpcClients.GetIdByToken(request.Token)
	if err != nil {
		response.StatusCode = -1
		response.StatusMsg = "登录失效，请重新登录"
		return nil
	}

	// 根据请求中的 UserId 从数据库中获取对应的视频列表
	var videoResult []*publishService.Video
	videos, err := models.NewVideoDaoInstance().GetVideosByUserId(request.UserId)
	if err != nil {
		response.StatusCode = -1
		response.StatusMsg = "查询视频列表失败"
		return err
	}

	// 调用 userService 微服务,请求对应的视频作者信息
	userInfo, _ := rpcClients.GetUserInfo(request.UserId, request.Token)
	for _, video := range videos {
		videoResult = append(videoResult, BuildVideo(video, userInfo, false))
	}

	response.StatusCode = 0
	response.StatusMsg = "查询视频列表成功"
	response.VideoList = videoResult

	return nil
}

func (*PublishService) GetAuthorId(ctx context.Context, request *publishService.GetAuthorIdRequest, response *publishService.GetAuthorIdResponse) error {
	videoId := request.VideoId

	video, err := models.NewVideoDaoInstance().FindVideoById(videoId)
	if err != nil {
		response.StatusCode = -1
		response.StatusMsg = "查询作者Id失败"
	}
	response.StatusCode = 0
	response.StatusMsg = "查询作者Id成功"
	response.AuthorId = video.UserId
	return nil
}

func (*PublishService) GetVideoInfo(ctx context.Context, request *publishService.GetVideoInfoRequest, response *publishService.GetVideoInfoResponse) error {
	videoId := request.VideoId

	video, err := models.NewVideoDaoInstance().FindVideoById(videoId)
	if err != nil {
		response.StatusCode = -1
		response.StatusMsg = "查询视频失败"
	}
	userInfo, _ := rpcClients.GetUserInfo(video.UserId, "")
	response.StatusCode = 0
	response.StatusMsg = "查询视频成功"
	response.VideoInfo = BuildVideo(video, userInfo, true)
	return nil
}

func BuildVideo(video *models.Video, user *userService.User, isFavorite bool) *publishService.Video {
	return &publishService.Video{
		Id:            video.VideoId,
		Author:        BuildUser(user),
		PlayUrl:       video.PlayUrl,
		CoverUrl:      video.CoverUrl,
		FavoriteCount: video.FavoriteCount,
		CommentCount:  video.CommentCount,
		IsFavorite:    isFavorite,
		Title:         video.Title,
	}
}

func BuildUser(user *userService.User) *publishService.User {
	return &publishService.User{
		Id:            user.Id,
		Name:          user.Name,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		IsFollow:      user.IsFollow,
	}
}
