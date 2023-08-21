package userImp

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"userService/models"
	"userService/services/userService"
	"userService/utils/sha256"
)

type UserService struct {
}

func (*UserService) Login(ctx context.Context, request *userService.DouyinUserLoginRequest, response *userService.DouyinUserLoginResponse) error {
	username := request.Username
	password := request.Password

	// 判断用户名和密码是否为空
	if username == "" || password == "" {
		response.StatusCode = -1
		response.StatusMsg = "用户名或密码不能为空"
		response.UserId = -1
		response.Token = ""
		return nil
	}

	// 在数据库中查找用户是否存在
	user, err := models.NewUserDaoInstance().GetUserByName(username)
	if err != nil {
		return err
	}

	// 比较用户密码是否一致
	if sha256.Sha256(password) != user.Password {
		response.StatusCode = -1
		response.StatusMsg = "用户名或密码错误"
		response.UserId = -1
		response.Token = ""
		return nil
	}

	// 成功登录,此处不生成 token, token由gateway统一颁发
	response.StatusCode = 0
	response.StatusMsg = "登录成功"
	response.UserId = user.UserId
	return nil

}

func (*UserService) Register(ctx context.Context, request *userService.DouyinUserRegisterRequest, response *userService.DouyinUserRegisterResponse) error {
	username := request.Username
	password := request.Password

	// 判断用户名和密码是否为空
	if username == "" || password == "" {
		response.StatusCode = -1
		response.StatusMsg = "用户名或密码不能为空"
		response.UserId = -1
		response.Token = ""
		return nil
	}

	// 在数据库中查找用户是否存在
	if _, err := models.NewUserDaoInstance().GetUserByName(username); err == nil {
		response.StatusCode = -1
		response.StatusMsg = "用户名已经存在"
		response.UserId = -1
		response.Token = ""
		return nil
	}

	// 构造 user 行记录并存入数据库中
	user := &models.User{
		Name:           username,
		Password:       sha256.Sha256(password),
		FollowingCount: 0,
		FollowerCount:  0,
	}
	_, err := models.NewUserDaoInstance().CreateUser(user)
	if err != nil {
		response.StatusCode = -1
		response.StatusMsg = "创建用户失败"
		response.UserId = -1
		response.Token = ""
		return err
	}

	// 从数据库中读取该行,获得用户id(目前id是自增主键),成功注册
	user, _ = models.NewUserDaoInstance().GetUserByName(username)
	response.StatusCode = 0
	response.StatusMsg = "注册成功"
	response.UserId = user.UserId
	response.Token = ""
	return nil
}

func (*UserService) UserInfo(ctx context.Context, request *userService.DouyinUserRequest, response *userService.DouyinUserResponse) error {
	// 判断 token 是否为空
	// if request.Token == "" {
	// 	response.StatusCode = -1
	// 	response.StatusMsg = "登录失效，请重新登录"
	// 	response.User = &userService.User{}
	// 	return nil
	// }

	// 从 token 中解析出登录用户的 id
	// _, err := rpcClients.GetIdByToken(request.Token)
	// if err != nil {
	// 	response.StatusCode = -1
	// 	response.StatusMsg = "登录失效，请重新登录"
	// 	response.User = &userService.User{}
	// 	return nil
	// }

	userIdReq := request.UserId
	var user *models.User
	var userString string

	// 从 redis 缓存中查找是否缓存了被请求的用户信息
	count, err := models.RedisSession.Exists(models.Ctx, strconv.FormatInt(userIdReq, 10)).Result()
	if err != nil {
		log.Println(err)
	}

	if count > 0 {
		// 缓存命中,在 redis 缓存中通过 user_id 来查找用户信息
		userString, err = models.RedisSession.Get(models.Ctx, strconv.FormatInt(userIdReq, 10)).Result()
		if err != nil {
			// 查找出错
			log.Println("[redis] 查询用户信息出错:", err)
		}
		json.Unmarshal([]byte(userString), &user)
		log.Println("[redis] 得到:", &user)
	} else {
		// 缓存失效,从数据库中通过 user_id 来查找用户信息
		user, err = models.NewUserDaoInstance().GetUserById(userIdReq)
		if err != nil {
			// 查找出错
			response.StatusCode = 1
			response.StatusMsg = "查询不到该用户信息"
			return err
		}
		// 缓存查找到的数据
		userValue, _ := json.Marshal(&user)
		_ = models.RedisSession.Set(models.Ctx, strconv.FormatInt(userIdReq, 10), userValue, 0).Err()
	}

	// 查找双方关注关系
	// isFollow, err := models.NewUserDaoInstance().FindRelationById(userIdReq, userIdLogined)
	// if err != nil {
	// 	response.StatusCode = -1
	// 	response.StatusMsg = "查询relation数据库，两人是否有关注关系的时候失败"
	// 	return err
	// }

	// 查询成功
	response.StatusCode = 0
	response.StatusMsg = "查询用户信息成功"
	response.User = BuildUser(user, false)
	return nil
}

func BuildUser(item *models.User, isFollow bool) *userService.User {
	user := userService.User{
		Id:             item.UserId,
		Name:           item.Name,
		FollowCount:    item.FollowingCount,
		FollowerCount:  item.FollowerCount,
		IsFollow:       isFollow,
		TotalFavorited: item.TotalFavorited,
		WorkCount:      item.WorkCount,
		FavoriteCount:  item.FavoriteCount,
	}
	return &user
}
