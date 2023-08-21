package favoriteImp

import (
	"context"
	"favoriteService/services/favoriteService"
)

// 请在此完成服务的业务逻辑
type FavoriteService struct {
}

func (*FavoriteService) FavoriteAction(ctx context.Context, request *favoriteService.DouyinFavoriteActionRequest, response *favoriteService.DouyinFavoriteActionResponse) error {
	return nil
}

func (*FavoriteService) FavoriteList(ctx context.Context, request *favoriteService.DouyinFavoriteListRequest, response *favoriteService.DouyinFavoriteListResponse) error {
	return nil
}
