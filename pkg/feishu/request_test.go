package feishu

import (
	"github.com/AY7295/notifer/shared"
	"testing"
)

var (
	config = &Config{
		Lark: Lark{
			ID:     "cli_a381acfe9ff99013",
			Secret: "ey3WcQNlnTi3YOr2b0cifdow3AiHOZU8",
		},
	}
	phones = shared.Phones{Mobiles: []string{"18980710863", "13535814223"}}
)

func TestMain(m *testing.M) {
	err := config.Init()
	if err != nil {
		panic(err)
	}

	m.Run()
}

func Test_api_GetToken(t *testing.T) {
	api := newApi()

	token, expire, err := api.GetToken(config.Lark)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(token)
	t.Log(expire)

}

func Test_api_GetChatIds(t *testing.T) {
	api := newApi()
	chatIds, err := api.GetChatIds(config.tenantAccessToken)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(chatIds)
}

func Test_api_GetOpenIds(t *testing.T) {
	api := newApi()
	openIds, err := api.GetOpenIds(phones, config.tenantAccessToken)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(openIds)
}
