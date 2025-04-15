package contracts

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"solution/internal/domain/dto"
)

type AuthService interface {
	GenerateToken(ctx context.Context, subject string) (string, *dto.HttpErr)
	GetSubject(ctx context.Context, tokenString *jwt.Token) (string, *dto.HttpErr)

	ExchangeCodeForToken(ctx context.Context, code string) (*dto.YandexTokenInfo, error)
	ExchangeTokenForUserInfo(ctx context.Context, token string) (*dto.YandexUserInfo, error)
}
