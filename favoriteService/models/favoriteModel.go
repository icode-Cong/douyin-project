package models

import (
	"sync"
)

// 请在此完成对 favarite 服务所用表的定义
type Favorite struct {
	UserId  int64 `gorm:"primaryKey"`
	VideoId int64 `gorm:"primaryKey"`
}
type FavoriteDao struct {
}

var favDao *FavoriteDao
var favOnce sync.Once //单例模式，只生成一个VideoDao实例，提高性能

func NewFavDaoInstance() *FavoriteDao {
	favOnce.Do(
		func() {
			favDao = &FavoriteDao{}
		})
	return favDao
}

// 查询是否存在user点赞video的记录
func (f *FavoriteDao) IsFavorite(userId int64, VideoId int64) (bool, error) {
	ret := false
	var rec Favorite
	result := SqlSession.Where("User_id = ? and Video_id = ?", userId, VideoId).Find(rec)
	err := result.Error
	if err != nil {
		return ret, err
	}
	ret = true
	return ret, err
}

// 根据UserId获取关注视频Id列表
func (f *FavoriteDao) GetFavVideoIdList(userId int64) ([]int64, error) {
	var rec []int64
	result := SqlSession.Where("User_id = ? ", userId).Find(&rec)
	err := result.Error
	if err != nil {
		return rec, err
	}
	return rec, err
}

func (f *FavoriteDao) AddFavRec(userId int64, videoId int64) error {
	fav := Favorite{
		UserId:  userId,
		VideoId: videoId,
	}
	result := SqlSession.Create(&fav)
	return result.Error
}

func (f *FavoriteDao) DelFavRec(userId int64, videoId int64) error {
	fav := Favorite{
		UserId:  userId,
		VideoId: videoId,
	}
	result := SqlSession.Delete(&fav)
	return result.Error
}
