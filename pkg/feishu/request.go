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
		auth: make(map[string]any),
	}
}

type api struct {
	// client: the http-client of the api
	client *tool.Http
	// auth: the tenantAccessToken of the lark-bot, need to be refreshed every 2 hours
	auth map[string]any
	// expire: the expiring time of the auth
	expire time.Time
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
		return err
	}
	if code, ok := response["code"].(float64); !ok || code != 0 {
		return failedRequest(response)
	}

	token, ok := response["tenant_access_token"].(string)
	if !ok {
		return newApiUrlError(getTokenUrl)
	}
	a.auth["Authorization"] = "Bearer " + token

	expire, ok := response["expire"].(float64)
	if !ok {
		return "", 0, newApiUrlError(getTokenUrl)
		return newApiUrlError(getTokenUrl)
	}

	a.expire = time.Now().Add(time.Duration(expire) * time.Second)
	return nil
}

	}

	return token, time.Duration(expire) * time.Second, nil
}

func (a *api) SendCard() error {
	panic("implement me")
}
