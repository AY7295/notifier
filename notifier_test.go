package notifier

import (
	"github.com/AY7295/notifer/pkg/feishu"
	"github.com/AY7295/notifer/shared"
	"testing"
)

var (
	app = shared.App{
		Name:    "TestAppName",
		Mobiles: []string{},
	}
	config = feishu.Config{
		Lark: feishu.Lark{
			AppId:     "",
			AppSecret: "",
		},
		GroupNotify: true,
	}
)

func TestInit(t *testing.T) {
	fsBuilder, err := NewFeishuBuilder(config)
	if err != nil {
		t.Error(err)
		return
	}

	Global.Init(app, WithNotifier(shared.Error, fsBuilder.Build(shared.Error)))
	Global.Notify(shared.Error, shared.NewInformation("TestError"))
}
