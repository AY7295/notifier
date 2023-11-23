package feishu

import (
	"errors"
	"fmt"
	"github.com/Mmx233/tool"
	"log"
	"time"
)

/*
* [https://open.feishu.cn/document/server-docs/api-call-guide/calling-process/overview] here contains all the apis of lark
**/
const (
	timeout          = 10 * time.Second
	getChatIdsUrl    = "https://open.feishu.cn/open-apis/im/v1/chats"
	getTokenUrl      = "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal"
	getUserOpenIdUrl = "https://open.feishu.cn/open-apis/contact/v3/users/batch_get_id?user_id_type=open_id"
	sendMessageUrl   = "https://open.feishu.cn/open-apis/im/v1/messages"
)

func newApiUrlError(url string) error {
	return errors.New(fmt.Sprintf("api [%s] may has changed", url))
}

func wrongPhoneNumber(phoneNumber string) error {
	return errors.New(fmt.Sprintf("wrong phone number: %s", phoneNumber))
}

func failedRequest(response map[string]any) error {
	return errors.New(fmt.Sprintf("failed request: %v", response))
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

func (a *api) getOpenIds(phones []string) ([]string, error) {
	_, response, err := a.client.Post(&tool.DoHttpReq{
		Url:    getUserOpenIdUrl,
		Header: a.auth,
		Body: struct {
			Mobiles []string `json:"mobiles"`
		}{phones},
	})
	if err != nil {
		return nil, err
	}
	if code, ok := response["code"].(float64); !ok || code != 0 {
		return nil, failedRequest(response)
	}

	ids := make([]string, 0, len(phones))
	userList, ok := response["data"].(map[string]any)["user_list"].([]any)
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

func (a *api) getChatIds() ([]string, error) {
	_, response, err := a.client.Get(&tool.DoHttpReq{
		Url:    getChatIdsUrl,
		Header: a.auth,
		Query:  map[string]any{"page_size": 100},
	})
	if err != nil {
		return nil, err
	}
	if code, ok := response["code"].(float64); !ok || code != 0 {
		return nil, failedRequest(response)
	}

	groups, ok := response["data"].(map[string]any)["items"].([]any)
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

func (a *api) refreshToken(app Lark) error {
	_, response, err := a.client.Post(&tool.DoHttpReq{
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
		return newApiUrlError(getTokenUrl)
	}

	a.expire = time.Now().Add(time.Duration(expire) * time.Second)
	return nil
}

var (
	cardMessageBody = struct {
		ReceiveId string `json:"receive_id"`
		Content   string `json:"content"`
		MsgType   string `json:"msg_type"`
	}{
		MsgType: "interactive",
	}
	userTarget = map[string]any{
		"receive_id_type": "open_id",
	}
	groupTarget = map[string]any{
		"receive_id_type": "chat_id",
	}
)

func (a *api) sendCard(card string, userIds []string, chatIds ...string) error {
	log.Println(card)
	cardMessageBody.Content = card
	var errs []error
	for _, id := range userIds {
		cardMessageBody.ReceiveId = id
		_, err := a.sendMessage(cardMessageBody, userTarget)
		if err != nil {
			errs = append(errs, err)
		}
	}
	for _, id := range chatIds {
		cardMessageBody.ReceiveId = id
		_, err := a.sendMessage(cardMessageBody, groupTarget)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (a *api) sendMessage(body any, target map[string]any) (map[string]any, error) {
	_, response, err := a.client.Post(&tool.DoHttpReq{
		Url:    sendMessageUrl,
		Header: a.auth,
		Query:  target,
		Body:   body,
	})
	if err != nil {
		return nil, err
	}

	if code, ok := response["code"].(float64); !ok || code != 0 {
		return nil, failedRequest(response)
	}
	return response["data"].(map[string]any), nil
}
