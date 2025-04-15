package wrapper

import "solution/internal/domain/dto"

func NewYandexUser(tokenInfo *dto.YandexTokenInfo, userInfo *dto.YandexUserInfo) *dto.UserAuth {
	return &dto.UserAuth{
		Login:       userInfo.Login,
		Provider:    dto.YandexProvider,
		AccessToken: tokenInfo.AccessToken,
		YandexId:    userInfo.Id,
	}
}
