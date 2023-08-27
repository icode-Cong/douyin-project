package relationImp

import (
	"context"
	"fmt"
	"math/rand"
	"relationService/models"
	"relationService/rpcClients"
	"relationService/services/messageService"
	"relationService/services/relationService"
	"relationService/services/userService"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type RelationService struct {
}

func (*RelationService) RelationAction(ctx context.Context, request *relationService.DouyinRelationActionRequest, response *relationService.DouyinRelationActionResponse) error {
	token := request.Token
	toUserId := request.ToUserId
	actionType := request.ActionType

	loginUserId, err := rpcClients.GetIdByToken(token)
	if err != nil {
		response.StatusCode = -1
		response.StatusMsg = "登录失效，请重新登录"
		return nil
	}

	updateRelationAndCache := func() {
		relation := &models.Relation{
			FromUserId: loginUserId,
			ToUserId:   toUserId,
			CreateAt:   time.Now(),
		}
		if actionType == 1 {
			_ = models.NewRelationDaoInstance().CreateRelation(relation)
			rpcClients.UpdateFollowCount(loginUserId, 1, int64(actionType))
			rpcClients.UpdateFollowerCount(toUserId, 1, int64(actionType))
			err := models.RedisSession.SAdd(ctx, fmt.Sprintf("follow:user:%d", loginUserId), toUserId).Err()
			if err != nil {
				fmt.Println("[警告：更新关注关系缓存出错：]", err)
			}
			err = models.RedisSession.SAdd(ctx, fmt.Sprintf("follower:user:%d", toUserId), loginUserId).Err()
			if err != nil {
				fmt.Println("[警告：更新粉丝关系缓存出错：]", err)
			}
		} else if actionType == 2 {
			_ = models.NewRelationDaoInstance().DeleteRelation(relation)
			rpcClients.UpdateFollowCount(loginUserId, 1, int64(actionType))
			rpcClients.UpdateFollowerCount(toUserId, 1, int64(actionType))
			// 在取消关注时，从关注列表缓存中移除取消关注的用户ID
			err := models.RedisSession.SRem(ctx, fmt.Sprintf("follow:user:%d", loginUserId), toUserId).Err()
			if err != nil {
				fmt.Println("[警告：更新关注关系缓存出错：]", err)
			}
			// 同时从粉丝列表缓存中移除取消关注的用户ID
			err = models.RedisSession.SRem(ctx, fmt.Sprintf("follower:user:%d", toUserId), loginUserId).Err()
			if err != nil {
				fmt.Println("[警告：更新粉丝关系缓存出错：]", err)
			}
		}
	}
	// 根据 actionType 执行不同操作
	switch actionType {
	case 1: // 关注
		updateRelationAndCache()
		response.StatusCode = 0
		response.StatusMsg = "关注成功"
	case 2: // 取消关注
		updateRelationAndCache()
		response.StatusCode = 0
		response.StatusMsg = "取消关注成功"
	default:
		response.StatusCode = -1
		response.StatusMsg = "actionType 出错"
	}
	return nil
}

func (*RelationService) FollowList(ctx context.Context, request *relationService.DouyinRelationFollowListRequest, response *relationService.DouyinRelationFollowListResponse) error {
	// 由网关负责登录状态的校验，服务内不再校验
	userId := request.UserId
	token := request.Token

	// 获取关注对象的id数组
	// ？？？？？这里需要用 redis 优化吗
	// 尝试从 Redis 中获取关注 ID 列表
	cachedFollowIds, err := models.RedisSession.SMembers(ctx, fmt.Sprintf("follow:user:%d", userId)).Result()
	if err != nil && err != redis.Nil {
		fmt.Println("[Warning: 从缓存中获取关注列表失败，错误：]", err)
	}

	var followIds []int64
	for _, idStr := range cachedFollowIds {
		id, convErr := strconv.ParseInt(idStr, 10, 64)
		if convErr != nil {
			fmt.Println("[警告：无法转换关注 ID 到 int64：]", convErr)
			continue
		}
		followIds = append(followIds, id)
	}

	// 如果缓存中没有，从数据库获取并存入 Redis 缓存
	if len(followIds) == 0 {
		followIds, err := models.NewRelationDaoInstance().GetFollowIdList(userId)
		if err != nil {
			response.StatusCode = -1
			response.StatusMsg = "查询关注列表失败"
			return nil
		}

		// 存入 Redis 缓存，设置随机的过期时间
		expiration := time.Duration(rand.Intn(600)+300) * time.Second // 随机设置过期时间在 300 到 900 秒之间（5 到 15 分钟)
		err = models.RedisSession.SAdd(ctx, fmt.Sprintf("follow:user:%d", userId), followIds).Err()
		if err != nil {
			fmt.Println("[警告：缓存关注 ID 列表出错：]", err)
		}
		err = models.RedisSession.Expire(ctx, fmt.Sprintf("follow:user:%d", userId), expiration).Err()
		if err != nil {
			fmt.Println("[警告：设置缓存过期时间出错：]", err)
		}
	}

	// 使用协程来逐个获取关注对象的用户信息
	// 使用通道来传递用户信息
	userInfoChan := make(chan *userService.User)
	// var userList []*relationService.User
	var wg sync.WaitGroup
	for _, followId := range followIds {
		wg.Add(1)
		go func(id int64) {
			defer wg.Done()
			userInfo, err := rpcClients.GetUserInfoById(id, token)
			if err != nil {
				fmt.Println("[Error : 查询关注列表, 向 userService 请求用户信息失败，用户id：]", id)
				return
			}
			userInfoChan <- userInfo
		}(followId)
	}

	// 等待所有协程完成
	go func() {
		wg.Wait()
		close(userInfoChan) // 关闭数据，表示不再写入数据
	}()

	// 收集用户信息
	var userList []*relationService.User
	for userInfo := range userInfoChan {
		userList = append(userList, BuildUser(userInfo))
	}

	response.StatusCode = 0
	response.StatusMsg = "查询关注列表成功"
	response.UserList = userList
	return nil
}

// func (*RelationService) FollowerList(ctx context.Context, request *relationService.DouyinRelationFollowerListRequest, response *relationService.DouyinRelationFollowerListResponse) error {
// 	userId := request.UserId
// 	token := request.Token

// 	// ？？？？？这里需要用 redis 优化吗
// 	followerIds, err := models.NewRelationDaoInstance().GetFollowerIdList(userId)
// 	if err != nil {
// 		response.StatusCode = -1
// 		response.StatusMsg = "查询粉丝列表失败"
// 		return nil
// 	}

// 	userInfoChan := make(chan *userService.User)
// 	var wg sync.WaitGroup
// 	for _, followerId := range followerIds {
// 		wg.Add(1)
// 		go func(id int64) {
// 			defer wg.Done()
// 			userInfo, err := rpcClients.GetUserInfoById(id, token)
// 			if err != nil {
// 				fmt.Println("[Error : 查询粉丝列表, 向 userService 请求用户信息失败，用户id：]", id)
// 				return
// 			}
// 			userInfoChan <- userInfo
// 		}(followerId)
// 	}

// 	go func() {
// 		wg.Wait()
// 		close(userInfoChan)
// 	}()

// 	var userList []*relationService.User
// 	for userInfo := range userInfoChan {
// 		userList = append(userList, BuildUser(userInfo))
// 	}

// 	response.StatusCode = 0
// 	response.StatusMsg = "查询粉丝列表成功"
// 	response.UserList = userList
// 	return nil
// }

func (*RelationService) FollowerList(ctx context.Context, request *relationService.DouyinRelationFollowerListRequest, response *relationService.DouyinRelationFollowerListResponse) error {
	// 由网关负责登录状态的校验，服务内不再校验
	userId := request.UserId
	token := request.Token

	// 获取关注对象的id数组
	// ？？？？？这里需要用 redis 优化吗
	// 尝试从 Redis 中获取关注 ID 列表
	cachedFollowerIds, err := models.RedisSession.SMembers(ctx, fmt.Sprintf("follower:user:%d", userId)).Result()
	if err != nil && err != redis.Nil {
		fmt.Println("[Warning: 从缓存中获取关注列表失败，错误：]", err)
	}

	var followerIds []int64
	for _, idStr := range cachedFollowerIds {
		id, convErr := strconv.ParseInt(idStr, 10, 64)
		if convErr != nil {
			fmt.Println("[警告：无法转换关注 ID 到 int64：]", convErr)
			continue
		}
		followerIds = append(followerIds, id)
	}

	// 如果缓存中没有，从数据库获取并存入 Redis 缓存
	if len(followerIds) == 0 {
		followerIds, err := models.NewRelationDaoInstance().GetFollowerIdList(userId)
		if err != nil {
			response.StatusCode = -1
			response.StatusMsg = "查询关注列表失败"
			return nil
		}

		// 存入 Redis 缓存，设置随机的过期时间
		expiration := time.Duration(rand.Intn(600)+300) * time.Second // 随机设置过期时间在 300 到 900 秒之间（5 到 15 分钟)
		err = models.RedisSession.SAdd(ctx, fmt.Sprintf("follower:user:%d", userId), followerIds).Err()
		if err != nil {
			fmt.Println("[警告：缓存关注 ID 列表出错：]", err)
		}
		err = models.RedisSession.Expire(ctx, fmt.Sprintf("follow:user:%d", userId), expiration).Err()
		if err != nil {
			fmt.Println("[警告：设置缓存过期时间出错：]", err)
		}
	}

	// 使用协程来逐个获取关注对象的用户信息
	// 使用通道来传递用户信息
	userInfoChan := make(chan *userService.User)
	// var userList []*relationService.User
	var wg sync.WaitGroup
	for _, followerId := range followerIds {
		wg.Add(1)
		go func(id int64) {
			defer wg.Done()
			userInfo, err := rpcClients.GetUserInfoById(id, token)
			if err != nil {
				fmt.Println("[Error : 查询关注列表, 向 userService 请求用户信息失败，用户id：]", id)
				return
			}
			userInfoChan <- userInfo
		}(followerId)
	}

	// 等待所有协程完成
	go func() {
		wg.Wait()
		close(userInfoChan) // 关闭数据，表示不再写入数据
	}()

	// 收集用户信息
	var userList []*relationService.User
	for userInfo := range userInfoChan {
		userList = append(userList, BuildUser(userInfo))
	}

	response.StatusCode = 0
	response.StatusMsg = "查询关注列表成功"
	response.UserList = userList
	return nil
}

func (*RelationService) IsFollowed(ctx context.Context, request *relationService.IsFollowedRequest, response *relationService.IsFollowedResponse) error {
	userId := request.UserId
	token := request.Token

	loginUserId, err := rpcClients.GetIdByToken(token)
	if err != nil {
		response.StatusCode = -1
		response.StatusMsg = "登录失效，请重新登录"
		return nil
	}

	isFollowed, err := models.RedisSession.SIsMember(ctx, fmt.Sprintf("follow:user:%d", loginUserId), userId).Result()

	if err != nil {
		followIds, err := models.NewRelationDaoInstance().GetFollowIdList(loginUserId)
		if err != nil {
			response.StatusCode = -1
			response.StatusMsg = "获取关注列表失败"
			return nil
		}

		err = models.RedisSession.SAdd(ctx, fmt.Sprintf("follow:user:%d", loginUserId), followIds).Err()
		if err != nil {
			fmt.Println("[警告：更新关注关系缓存出错：]", err)
		}
		isFollowed = models.RedisSession.SIsMember(ctx, fmt.Sprintf("follow:user:%d", loginUserId), userId).Val()
	}

	response.StatusCode = 0
	response.StatusMsg = "查询关注状态成功"
	response.IsFollow = isFollowed
	return nil
}

func (*RelationService) FriendList(ctx context.Context, request *relationService.DouyinRelationFriendListRequest, response *relationService.DouyinRelationFriendListResponse) error {
	token := request.Token

	// 解析 token 得到登录用户id
	loginUserId, err := rpcClients.GetIdByToken(token)
	if err != nil {
		response.StatusCode = -1
		response.StatusMsg = "登录失效，请重新登录"
		return nil
	}

	friendIds := models.NewRelationDaoInstance().GetFriendIdList(loginUserId)
	friends, err1 := rpcClients.GetUsersInfoByIds(friendIds)
	fmt.Println("err1 :", err1)
	messages, err2 := rpcClients.GetMessagesByIds(loginUserId, friendIds)
	fmt.Println("err2 :", err2)
	// var userList []*relationService.FriendUser
	// for _, friend := range friends {
	// 	for _, message := range messages {
	// 		if friend.Id == message.ToUserId {
	// 			userList = append(userList, BuildFriendUser(friend, message, 1)) // 1 表示发送
	// 			break
	// 		} else if friend.Id == message.FromUserId {
	// 			userList = append(userList, BuildFriendUser(friend, message, 0)) // 0 表示接收
	// 			break
	// 		}
	// 	}
	// }
	// 构建消息映射，将每个朋友的消息按照发送者分类

	messageMap := make(map[int64]*messageService.Message) // 假设 Message 为消息结构体
	for _, message := range messages {
		messageFrom := message.FromUserId
		messageTo := message.ToUserId
		if messageFrom != loginUserId {
			messageMap[messageFrom] = message
		} else {
			messageMap[messageTo] = message
		}
	}
	fmt.Println(messageMap)
	var userList []*relationService.FriendUser
	for _, friend := range friends {
		friendMessage, exists := messageMap[friend.Id]
		var messageType int64
		if exists {
			if friendMessage.FromUserId == friend.Id {
				messageType = 0 // 表示接收
			} else {
				messageType = 1 // 表示发送
			}
		}
		fmt.Println(friend.Id)
		fmt.Println(friendMessage)
		userList = append(userList, BuildFriendUser(friend, friendMessage, messageType))
	}
	response.StatusCode = 0
	response.StatusMsg = "查询朋友列表成功"
	response.UserList = userList
	return nil
}

func BuildFriendUser(user *userService.User, message *messageService.Message, msgType int64) *relationService.FriendUser {
	return &relationService.FriendUser{
		Id:              user.Id,
		Name:            user.Name,
		FollowCount:     user.FollowCount,
		FollowerCount:   user.FollowerCount,
		IsFollow:        user.IsFollow,
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		Signature:       user.Signature,
		TotalFavorited:  user.TotalFavorited,
		WorkCount:       user.WorkCount,
		FavoriteCount:   user.FavoriteCount,
		Message:         message.Content,
		MsgType:         msgType,
	}
}

func BuildUser(user *userService.User) *relationService.User {
	return &relationService.User{
		Id:              user.Id,
		Name:            user.Name,
		FollowCount:     user.FollowCount,
		FollowerCount:   user.FollowerCount,
		IsFollow:        user.IsFollow,
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		Signature:       user.Signature,
		TotalFavorited:  user.TotalFavorited,
		WorkCount:       user.WorkCount,
		FavoriteCount:   user.FavoriteCount,
	}
}
