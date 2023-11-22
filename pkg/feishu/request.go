package feishu

import (
	"errors"
	"fmt"
	"github.com/AY7295/notifer/shared"
	"github.com/Mmx233/tool"
	"time"
)

const (
	timeout          = 10 * time.Second
	getChatIdsUrl    = "https://open.feishu.cn/open-apis/im/v1/chats"
	getTokenUrl      = "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal"
	getUserOpenIdUrl = "https://open.feishu.cn/open-apis/contact/v3/users/batch_get_id?user_id_type=open_id"
	sendCardUrl      = ""
)

func newApiUrlError(url string) error {
	return errors.New(fmt.Sprintf("api [%s] may has changed", url))
}

func wrongPhoneNumber(phoneNumber string) error {
	return errors.New(fmt.Sprintf("wrong phone number: %s", phoneNumber))
}

func newApi() *api {
	return &api{
		client: tool.NewHttpTool(tool.GenHttpClient(&tool.HttpClientOptions{
			Transport: tool.GenHttpTransport(&tool.HttpTransportOptions{
				Timeout:           timeout,
				IdleConnTimeout:   3 * timeout,
				SkipSslCertVerify: true,
			}),
			Timeout: timeout,
		})),
	}
}

type api struct {
	client *tool.Http
}

func (a *api) GetOpenIds(phones shared.Phones, tenantAccessToken string) ([]string, error) {
	_, body, err := a.client.Post(&tool.DoHttpReq{
		Url:    getUserOpenIdUrl,
		Header: map[string]any{"Authorization": tenantAccessToken},
		Body:   phones,
	})
	if err != nil {
		return nil, err
	}
	fmt.Println(body)
	ids := make([]string, 0, len(phones.Mobiles))
	userList, ok := body["data"].(map[string]any)["user_list"].([]any)
	if !ok {
		return nil, newApiUrlError(getUserOpenIdUrl)
	}

	for _, user := range userList {
		userInfo, ok := user.(map[string]any)
		if !ok {
			return nil, newApiUrlError(getUserOpenIdUrl)
		}

		id, ok := userInfo["user_id"].(string)
		if !ok {
			phoneNumber, ok := userInfo["mobile"].(string)
			if !ok {
				return nil, newApiUrlError(getUserOpenIdUrl)
			}

			return nil, errors.Join(newApiUrlError(getUserOpenIdUrl), wrongPhoneNumber(phoneNumber))
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (a *api) GetChatIds(tenantAccessToken string) ([]string, error) {
	_, body, err := a.client.Get(&tool.DoHttpReq{
		Url:    getChatIdsUrl,
		Header: map[string]any{"Authorization": tenantAccessToken},
		Query:  map[string]any{"page_size": 100},
	})
	if err != nil {
		return nil, err
	}

	groups, ok := body["data"].(map[string]any)["items"].([]any)
	if !ok {
		return nil, newApiUrlError(getChatIdsUrl)
	}

	ids := make([]string, 0, len(groups))
	for _, group := range groups {
		id, ok := group.(map[string]any)["chat_id"].(string)
		if !ok {
			return nil, newApiUrlError(getChatIdsUrl)
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (a *api) GetToken(app Lark) (string, time.Duration, error) {
	_, body, err := a.client.Post(&tool.DoHttpReq{
		Url:  getTokenUrl,
		Body: app,
	})
	if err != nil {
		return "", 0, err
	}

	token, ok := body["tenant_access_token"].(string)
	if !ok {
		return "", 0, newApiUrlError(getTokenUrl)
	}

	expire, ok := body["expire"].(float64)
	if !ok {
		return "", 0, newApiUrlError(getTokenUrl)
	}

	return token, time.Duration(expire) * time.Second, nil
}

func (a *api) SendCard() error {
	panic("implement me")
}
