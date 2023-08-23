package models

import (
	"fmt"
	"sync"
	"time"
)

type Comment struct {
	CommentId int64 `gorm:"primary_key"`
	UserId    int64
	VideoId   int64
	Content   string
	CreateAt  time.Time
	DeletedAt time.Time
}

func (Comment) TableName() string {
	return "comment"
}

type CommentDao struct {
}

var commentDao *CommentDao
var commentOnce sync.Once

func NewCommentDaoInstance() *CommentDao {
	commentOnce.Do(
		func() {
			commentDao = &CommentDao{}
		})
	return commentDao
}

/*
*
创建一个新的 Comment，返回 Content
*/
func (*CommentDao) CreateComment(comment *Comment) (*Comment, error) {
	/*comment := Comment{UserId: userid, VideoId: videoId, Content: content, CreatedAt: time.Now()}*/
	result := SqlSession.Create(&comment)

	if result.Error != nil {
		return nil, result.Error
	}

	return comment, nil
}

/*
*
根据 CommentId 删除 Comment
*/
func (*CommentDao) DeleteComment(commentId int64) error {
	err := SqlSession.Where("comment_id = ?", commentId).Delete(&Comment{}).Error

	if err != nil {
		fmt.Printf("删除失败")
	}

	return err
}

/*
*
根据 VideoId 获取 Comment 列表
*/
func (*CommentDao) GetCommentListByVideoId(videoId int64) ([]*Comment, error) {
	var comment []*Comment

	err := SqlSession.Where("video_id = ?", videoId).Find(&comment).Error

	if err != nil {
		fmt.Printf("查询失败")
		return nil, err
	}

	return comment, err
}

/*
*
根据 CommentId 获取 UserId
*/
func (*CommentDao) GetUserIdByCommentId(commentId int64) (int64, error) {
	comment := Comment{CommentId: commentId}
	err := SqlSession.Where("comment_id = ?", commentId).First(&comment).Error
	return comment.UserId, err
}
