package models

import (
	"sync"
	"time"
)

// 请在此完成对 relation 服务所用表的定义
type Relation struct {
	FromUserId int64
	ToUserId   int64
	CreateAt   time.Time
}

func (Relation) TableName() string {
	return "relation"
}

type RelationDao struct {
}

var relationDao *RelationDao
var relationOnce sync.Once //单例模式，只生成一个messageDao实例，提高性能

func NewRelationDaoInstance() *RelationDao {
	relationOnce.Do(
		func() {
			relationDao = &RelationDao{}
		})
	return relationDao
}

func (*RelationDao) CreateRelation(relation *Relation) error {
	result := SqlSession.Create(&relation)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (*RelationDao) DeleteRelation(relation *Relation) error {
	err := SqlSession.Where("from_user_id=? and to_user_id=?", relation.FromUserId, relation.ToUserId).Delete(relation).Error
	if err != nil {
		return err
	}
	return nil
}

func (*RelationDao) GetFollowIdList(userId int64) ([]int64, error) {
	var userIds []int64
	var relations []*Relation
	err := SqlSession.Where("from_user_id=?", userId).Find(&relations).Error
	if err != nil {
		return userIds, err
	}
	for i := 0; i < len(relations); i++ {
		userIds = append(userIds, relations[i].ToUserId)
	}
	return userIds, nil
}

func (*RelationDao) GetFollowerIdList(userId int64) ([]int64, error) {
	var userIds []int64
	var relations []*Relation
	err := SqlSession.Where("to_user_id=?", userId).Find(&relations).Error
	if err != nil {
		return userIds, err
	}
	for i := 0; i < len(relations); i++ {
		userIds = append(userIds, relations[i].ToUserId)
	}
	return userIds, nil
}

func (*RelationDao) GetFriendIdList(userId int64) []int64 {
	var userIds []int64
	var relations []*Relation
	SqlSession.Exec("SELECT * FROM relation WHERE from_user_id = ? AND to_user_id IN (SELECT from_user_id FROM relation WHERE to_user_id = ?)", 7, 7).Find(&relations)
	for i := 0; i < len(relations); i++ {
		userIds = append(userIds, relations[i].ToUserId)
	}
	return userIds
}
