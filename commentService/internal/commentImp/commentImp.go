package commentImp

import (
	"commentService/services/commentService"
	"context"
)

// 请在此完成服务的业务逻辑
type CommentService struct {
}

func (*CommentService) CommentAction(ctx context.Context, request *commentService.DouyinCommentActionRequest, reponse *commentService.DouyinCommentActionResponse) error {
	// .. = request.Token
	// .. = request.ActionType
	// .. = request.CommentId
	// .. = request.CommentText
	// .. = request.VideoId

	// response.StatusCode = ...
	// reponse.StatusMsg = ...
	// reponse.Comment = ...
	return nil
}

func (*CommentService) CommentList(ctx context.Context, request *commentService.DouyinCommentListRequest, reponse *commentService.DouyinCommentListResponse) error {
	return nil
}
