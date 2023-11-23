package feishu

import (
	"errors"
	"github.com/AY7295/notifer/shared"
	"testing"
)

var (
	config = &Config{
		Lark: Lark{
			AppId:     "",
			AppSecret: "",
		},
		GroupNotify: false,
	}
	app = &shared.App{
		Name:    "TestApp",
		Mobiles: []string{""},
	}
)

func TestMain(m *testing.M) {
	err := config.Init()
	if err != nil {
		panic(err)
	}

	m.Run()
}

func Test_api_GetToken(t *testing.T) {
	err := config.api.refreshToken(config.Lark)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(config.api.auth)
	t.Log(config.api.expire)

}

func Test_api_GetChatIds(t *testing.T) {
	chatIds, err := config.api.getChatIds()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(chatIds)
}

func Test_api_GetOpenIds(t *testing.T) {
	openIds, err := config.api.getOpenIds(app.Mobiles)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(openIds)
}

func Test_api_SendCard(t *testing.T) {
	var (
		userIds []string
		chatIds []string
		err     error
		info    = shared.NewInformation("TestError", shared.WithErrors(errors.New("error1"), errors.New("error2")))
	)

	userIds, err = config.api.getOpenIds(app.Mobiles)
	if err != nil {
		t.Error(err)
		return
	}
	chatIds, err = config.api.getChatIds()
	if err != nil {
		t.Error(err)
		return
	}
	ss := info.Format()
	t.Log(ss)
	err = config.api.sendCard(
		newCard(getTemplate(shared.Error), app.Name, shared.Error.String(), ss, builderAt(userIds...)),
		userIds, chatIds...,
	)
	if err != nil {
		t.Error(err)
		return
	}
}
