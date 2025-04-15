package dto

const YandexProvider = "yandex"

type YandexTokenInfo struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_i"`
	RefreshToken string `json:"refresh_token"`
}

type YandexUserInfo struct {
	Login    string `json:"login"`
	Id       string `json:"id"`
	ClientId string `json:"client_id"`
	Uid      string `json:"ui"`
	Psuid    string `json:"psuid"`
}
