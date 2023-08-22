package messageImp

import (
	"context"
	"messageService/services/messageService"
)

type MessageService struct {
}

// 请在此完成服务的业务逻辑,四个接口，上面这个结构体，不需要增加任何内容
func (*MessageService) MessageList(ctx context.Context, request *messageService.DouyinMessageChatRequest, response *messageService.DouyinMessageChatResponse) error {
	return nil
}
func (*MessageService) MessageAction(ctx context.Context, request *messageService.DouyinMessageActionRequest, response *messageService.DouyinMessageActionResponse) error {
	return nil
}
