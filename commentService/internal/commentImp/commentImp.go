package commentImp

import (
	"commentService/models"
	"commentService/rpcClients"
	"commentService/services/commentService"
	"commentService/services/userService"
	"context"
	"fmt"
	"sync"
	"time"
)

type CommentService struct {
}

// 创建删除评论
func (*CommentService) CommentAction(ctx context.Context, request *commentService.DouyinCommentActionRequest, response *commentService.DouyinCommentActionResponse) error {
	// 调用 tokenService 微服务, 通过 token 判断是否登录,如果查不到对应的 userId ,如果查不到则返回提示,要求重新登录
	userId, err := rpcClients.GetIdByToken(request.Token)
	if err != nil {
		response.StatusCode = 1
		response.StatusMsg = "登陆失败，请重新登陆"
		return nil
	}

	// 根据 userId 获取 user 实体
	user, _ := rpcClients.GetUserInfo(userId, request.Token)

	// 创建 comment 实体
	comment := &models.Comment{
		UserId:   userId,
		VideoId:  request.VideoId,
		Content:  request.CommentText,
		CreateAt: time.Now(),
	}

	// 获取 actionType 判断是创建还是删除 Comment
	actionType := request.ActionType

	if actionType == 1 {
		// 创建 Comment
		comment, err := models.NewCommentDaoInstance().CreateComment(comment)
		if err != nil {
			response.StatusCode = 1
			response.StatusMsg = "创建评论失败"
			return nil
		}

		response.StatusCode = 0
		response.StatusMsg = "创建评论成功"
		response.Comment = BuildComment(comment, user)

	} else if actionType == 2 {
		// 删除 Comment
		commentId := request.CommentId

		// 判断是否删除自己的评论
		commentUserId, _ := models.NewCommentDaoInstance().GetUserIdByCommentId(commentId)
		if userId != commentUserId {
			response.StatusCode = 1
			response.StatusMsg = "仅能删除自己的评论"
			return nil
		}

		err := models.NewCommentDaoInstance().DeleteComment(commentId)
		if err != nil {
			response.StatusCode = 1
			response.StatusMsg = "删除时异常"
			return nil
		}

		response.StatusCode = 0
		response.StatusMsg = "删除评论成功"
		response.Comment = &commentService.Comment{}

	} else {
		response.StatusCode = 1
		response.StatusMsg = "actionType 错误"
	}

	return nil
}

func (*CommentService) CommentList(ctx context.Context, request *commentService.DouyinCommentListRequest, response *commentService.DouyinCommentListResponse) error {
	// 调用 tokenService 微服务, 通过 token 判断是否登录,如果查不到对应的 userId ,如果查不到则返回提示,要求重新登录
	_, err := rpcClients.GetIdByToken(request.Token)
	if err != nil {
		response.StatusCode = 1
		response.StatusMsg = "登陆失败，请重新登陆"
		return nil
	}

	// 根据 videoId 获取 commentList
	var commentList []*models.Comment
	commentList, err = models.NewCommentDaoInstance().GetCommentListByVideoId(request.VideoId)
	if err != nil {
		response.StatusCode = 1
		response.StatusMsg = "获取评论列表失败"
		return nil
	}

	// 使用协程来逐个获取关注对象的用户信息
	// 使用通道来传递用户信息
	commentChan := make(chan *commentService.Comment)
	// var userList []*relationService.User
	var wg sync.WaitGroup
	for _, comment := range commentList {
		wg.Add(1)
		go func(id int64) {
			defer wg.Done()
			userInfo, err := rpcClients.GetUserInfo(id, request.Token)
			if err != nil {
				fmt.Println("[Error : 查询关注列表, 向 userService 请求用户信息失败，用户id：]", id)
				return
			}
			commentChan <- BuildComment(comment, userInfo)
		}(comment.UserId)
	}

	// 等待所有协程完成
	go func() {
		wg.Wait()
		close(commentChan) // 关闭数据，表示不再写入数据
	}()

	// 收集用户信息
	var commentRes []*commentService.Comment
	for comment := range commentChan {
		commentRes = append(commentRes, comment)
	}

	return nil
}

func BuildComment(comment *models.Comment, user *userService.User) *commentService.Comment {
	return &commentService.Comment{
		Id:         comment.CommentId,
		User:       BuildUser(user),
		Content:    comment.Content,
		CreateDate: comment.CreateAt.Format("2006-01-02 15:04:05"),
	}
}

func BuildUser(user *userService.User) *commentService.User {
	return &commentService.User{
		Id:            user.Id,
		Name:          user.Name,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		IsFollow:      user.IsFollow,
	}
}
