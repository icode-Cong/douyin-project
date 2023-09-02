# 各个服务

## gateService

通过接受app传来的请求调用各个服务，将结果返回给app

## feedService

- 表结构

```go
type Video struct {
    VideoId       int64 `gorm:"primaryKey"`
    UserId        int64
    PlayUrl       string
    CoverUrl      string
    FavoriteCount int64
    CommentCount  int64
    Title         string
    CreateAt      time.Time
    UpdateAt      *time.Time
    DeleteAt      *time.Time
}
```

- api
  
```proto
service FeedService{
    //推流
  rpc Feed(DouyinFeedRequest) returns (DouyinFeedResponse);
}
```

- 未实现

1. 返回用户信息：avatar，background_image，signature

## favoriteService

- 表结构

```go
type Favorite struct {
    UserId  int64 `gorm:"primaryKey"`
    VideoId int64 `gorm:"primaryKey"`
}
```

- api
  
```proto
service FavoriteService {
    //点赞操作
    rpc FavoriteAction(DouyinFavoriteActionRequest) returns(DouyinFavoriteActionResponse);
    //获取点赞视频列表
    rpc FavoriteList(DouyinFavoriteListRequest) returns(DouyinFavoriteListResponse);
}
```

## publishService

- 表结构

```go
type Video struct {
    VideoId       int64 `gorm:"primaryKey"`
    UserId        int64
    PlayUrl       string
    CoverUrl      string
    FavoriteCount int64
    CommentCount  int64
    Title         string
    CreateAt      time.Time
    UpdateAt      *time.Time
    DeleteAt      *time.Time
}
```

- api
  
```proto
service PublishService{
    //获取发布列表
  rpc PublishList(DouyinPublishListRequest) returns (DouyinPublishListResponse);
    //发布视频
  rpc Publish(DouyinPublishActionRequest) returns(DouyinPublishActionResponse);
    //获取视频作者Id
  rpc GetAuthorId(GetAuthorIdRequest) returns(GetAuthorIdResponse);
    //获取视频信息
  rpc GetVideoInfo(GetVideoInfoRequest) returns(GetVideoInfoResponse);
}
```

## relationService

- 表结构
  
```go
type Relation struct {
    FromUserId int64
    ToUserId   int64
    CreateAt   time.Time
}
```

- api

```proto
service RelationService{
    //关注或者取消关注
  rpc RelationAction(DouyinRelationActionRequest) returns (DouyinRelationActionResponse); 
    //获取关注列表
  rpc FollowList(DouyinRelationFollowListRequest) returns(DouyinRelationFollowListResponse);
    //获取粉丝列表
  rpc FollowerList(DouyinRelationFollowerListRequest) returns(DouyinRelationFollowerListResponse);
    //获取好友列表
  rpc FriendList(DouyinRelationFriendListRequest) returns(DouyinRelationFriendListResponse); 
    //询问是否关注
  rpc IsFollowed(IsFollowedRequest) returns(IsFollowedResponse); 
}
```

## tokenService

- 表结构（无）

- api

```proto
service TokenService {
    //根据token解析出Id
    rpc ParseTokenToId(ParseTokenToIdRequest) returns(ParseTokenToIdResponse);
}
```

## userService

- 表结构

```go
type User struct {
    UserId          int64 `gorm:"primary_key"`
    Name            string
    FollowingCount  int64
    FollowerCount   int64
    Password        string
    Avatar          string
    BackgroundImage string
    Signature       string
    TotalFavorited  int64
    WorkCount       int64
    FavoriteCount   int64
    // CreateAt        time.Time
    // DeleteAt        time.Time
}
```

- api

```proto
service UserService {
    //登录
    rpc Login(DouyinUserLoginRequest) returns(DouyinUserLoginResponse);
    //注册
    rpc Register(DouyinUserRegisterRequest) returns(DouyinUserRegisterResponse);
    //查询用户信息
    rpc UserInfo(DouyinUserRequest) returns(DouyinUserResponse);
    //更新关注数
    rpc UpdateFollowCount(DouyinUpdateFollowCountRequest) returns(DouyinUpdateFollowCountResponse);
    //更新粉丝数
    rpc UpdateFollowerCount(DouyinUpdateFollowerCountRequest) returns(DouyinUpdateFollowerCountResponse);
    //更新总点赞数
    rpc UpdateTotalFavorited(DouyinUpdateTotalFavoritedRequest) returns(DouyinUpdateTotalFavoritedResponse);
    //更新点赞数
    rpc UpdateFavoriteCount(DouyinUpdateFavoriteCountRequest) returns(DouyinUpdateFavoriteCountResponse);
    //更新作品数
    rpc UpdateWorkCount(DouyinUpdateWorkCountRequest) returns(DouyinUpdateWorkCountResponse);
    //批量查询用户
    rpc MultiUserInfo(DouyinMultiUserInfoRequest) returns(DouyinMultiUserInfoResponse);
}

```

## messageService

不知道干啥的

- task

1. 各个微服务最后需要同步到servicesRepository
2. 关注的表重复定义？
