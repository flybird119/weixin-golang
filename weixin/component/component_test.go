package component

import (
	"testing"
	com "wechat_component/component"
)

func TestAccessToken(t *testing.T) {
	component := com.New("wxb39264d954771192", "01e4ae1dd265db7cd94611c2b9bb7ff5", "goushuyun", "goushuyungoushuyungoushuyungoushuyungoushuy")

	token, expire := component.GetRegularApi().GetAccessToken("ticket@@@tj2xWy0ZEsS4rY9XKF1TJdo68lJWpfO723JLilh_QkceXd_fA6NxoB5gWAInsXiMw75JoRDoTs8K9sDYRz0Fkg")

	t.Log(token, expire)

}
