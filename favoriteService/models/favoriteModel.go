package models

import (
	"sync"
)

// 请在此完成对 favarite 服务所用表的定义
type Favorite struct {
	UserId  int64 `gorm:"primaryKey"`
	VideoId int64 `gorm:"primaryKey"`
}

func (Favorite) TableName() string {
	return "favorite"
}

type FavoriteDao struct {
}

var favoriteDao *FavoriteDao
var favoriteOnce sync.Once //单例模式，只生成一个VideoDao实例，提高性能

func NewFavoriteDaoInstance() *FavoriteDao {
	favoriteOnce.Do(
		func() {
			favoriteDao = &FavoriteDao{}
		})
	return favoriteDao
}

// 查询是否存在user点赞video的记录
func (f *FavoriteDao) IsFavorite(userId int64, videoId int64) bool {
	var rec Favorite
	result := SqlSession.Where("User_id = ? and Video_id = ?", userId, videoId).Limit(1).Find(&rec)
	err := result.Error
	if err != nil {
		return false
	}
	return true
}

// 根据UserId获取关注视频Id列表
func (*FavoriteDao) GetFavoriteVideoIdList(userId int64) []int64 {
	var videoIds []int64
	var favorites []*Favorite
	SqlSession.Where("user_id = ? ", userId).Find(&favorites)
	for i := 0; i < len(favorites); i++ {
		videoIds = append(videoIds, favorites[i].VideoId)
	}
	return videoIds
}

func (*FavoriteDao) CreateFavorite(favorite *Favorite) error {
	result := SqlSession.Create(&favorite)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (*FavoriteDao) DeleteFavorite(favorite *Favorite) error {
	err := SqlSession.Where("user_id=? and video_id=?", favorite.UserId, favorite.VideoId).Delete(favorite).Error
	if err != nil {
		return err
	}
	return nil
}
