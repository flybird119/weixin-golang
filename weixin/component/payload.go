package component

type GetAcessTokenResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
}
