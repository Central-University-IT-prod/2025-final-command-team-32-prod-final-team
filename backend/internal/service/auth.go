package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"solution/internal/config"
	"solution/internal/domain/contracts"
	"solution/internal/domain/dto"
	"solution/internal/wrapper"
	"solution/pkg/logger"

	"github.com/golang-jwt/jwt/v5"
	"github.com/valyala/fasthttp"
)

var _ contracts.AuthService = (*authService)(nil)

var (
	YANDEX_PROVIDER = "yandex"
)

type authService struct {
	secretKey     string
	yaClientId    string
	yaSecret      string
	yaRedirectUrl string
}

func NewAuthService(config *config.Config) *authService {
	return &authService{
		secretKey:     config.SecretKey,
		yaClientId:    config.YaClientId,
		yaSecret:      config.YaSecret,
		yaRedirectUrl: config.YaRedirectUrl,
	}
}

func (s *authService) GenerateToken(ctx context.Context, username string) (string, *dto.HttpErr) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(3 * 24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		logger.FromCtx(ctx).Error(ctx, fmt.Sprintf("failed to generate token with ERR: %s", err.Error()))
		return "", wrapper.InternalServerErr("Failed to create subject")
	}
	return tokenString, nil
}

func (s *authService) GetSubject(ctx context.Context, token *jwt.Token) (string, *dto.HttpErr) {
	sub, err := token.Claims.GetSubject()
	if err != nil {
		logger.FromCtx(ctx).Error(ctx, fmt.Sprintf("failed to read token with ERR: %s", err.Error()))
		return "", wrapper.InternalServerErr("Failed to retrieve auth subject")
	}
	return sub, nil
}

func (s *authService) ExchangeCodeForToken(ctx context.Context, code string) (*dto.YandexTokenInfo, error) {
	clientID := s.yaClientId
	clientSecret := s.yaSecret
	redirectURI := s.yaRedirectUrl

	// Prepare the request to Yandex's token endpoint
	formData := url.Values{}
	formData.Set("grant_type", "authorization_code")
	formData.Set("code", code)
	formData.Set("client_id", clientID)
	formData.Set("client_secret", clientSecret)
	formData.Set("redirect_uri", redirectURI)

	// Send a POST request to Yandex
	resp, err := http.PostForm("https://oauth.yandex.com/token", formData)
	if err != nil {
		logger.FromCtx(ctx).Info(ctx, "ERROR ON oauth.yandex.com/token")
		logger.FromCtx(ctx).Info(ctx, fmt.Sprintf("ERROR: %s | BODY: %s", err, resp.Body))
		return nil, err
	}
	respBody := dto.YandexTokenInfo{}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		logger.FromCtx(ctx).Info(ctx, "ERROR ON decoding")
		logger.FromCtx(ctx).Info(ctx, fmt.Sprintf("ERROR: %s ", err))
		return nil, err
	}
	return &respBody, nil
}

func (s *authService) ExchangeTokenForUserInfo(ctx context.Context, token string) (*dto.YandexUserInfo, error) {
	client := fasthttp.Client{}
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	req.Header.Add("Authorization", fmt.Sprintf("OAuth %s", token))
	req.Header.SetRequestURI(fmt.Sprintf("https://login.yandex.ru/info?jwt_secret=%s&format=json", s.yaSecret))
	req.Header.SetMethod("GET")
	err := client.Do(req, resp)
	if err != nil {
		logger.FromCtx(ctx).Info(ctx, "ERROR ON login.yandex.ru/info")
		return nil, err
	}

	data := dto.YandexUserInfo{}
	logger.FromCtx(ctx).Info(ctx, fmt.Sprintf("LOGIN YANDEX INFO BODY: %s", string(resp.Body())))
	logger.FromCtx(ctx).Info(ctx, fmt.Sprintf("STATUS: %d", resp.StatusCode()))
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		logger.FromCtx(ctx).Info(ctx, "ERROR ON UNMARSHAL")
		logger.FromCtx(ctx).Info(ctx, fmt.Sprintf("ERR: %s", err.Error()))
		return nil, err
	}
	return &data, nil
}
