package relationImp

import (
	"context"
	"relationService/services/relationService"
)

type RelationService struct {
}

// 请在此完成服务的业务逻辑,四个接口，上面这个结构体，不需要增加任何内容
func (*RelationService) RelationAction(ctx context.Context, request *relationService.DouyinRelationActionRequest, response *relationService.DouyinRelationActionResponse) error {
	return nil
}
func (*RelationService) FollowList(ctx context.Context, request *relationService.DouyinRelationFollowListRequest, response *relationService.DouyinRelationFollowListResponse) error {
	return nil
}
func (*RelationService) FollowerList(ctx context.Context, request *relationService.DouyinRelationFollowerListRequest, response *relationService.DouyinRelationFollowerListResponse) error {
	return nil
}
func (*RelationService) FriendList(ctx context.Context, request *relationService.DouyinRelationFriendListRequest, response *relationService.DouyinRelationFriendListResponse) error {
	return nil
}
