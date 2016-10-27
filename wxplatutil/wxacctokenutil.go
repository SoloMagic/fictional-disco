package wxplatutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/toolbox"
	"io/ioutil"
	"net/http"
	"strings"
)

//sql replace
var accMap map[int]string

type AccessTokenResponse struct {
	AccessToken string  `json:"access_token"`
	ExpiresIn   float64 `json:"expires_in"`
}

type AccessTokenErrorResponse struct {
	ErrCode string  `json:"errcode"`
	ErrMsg  float64 `json:"errmsg"`
}

func FetchAccessToken(appID, appSecret, accessTokenFetchUrl string) (string, error) {
	requestLine := strings.Join([]string{accessTokenFetchUrl, "?grant_type=client_credential&appid=", appID, "&secret=", appSecret}, "")
	resp, err := http.Get(requestLine)
	if err != nil {
		fmt.Println("send get request atoken is error!", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("send get request atoken return body error!")
		return "", err
	}

	if bytes.Contains(body, []byte("access_token")) {
		at := AccessTokenResponse{}
		err = json.Unmarshal(body, &at)
		if err != nil {
			fmt.Println("send get request atoken return json error!")
			return "", err
		}
		return at.AccessToken, nil
	} else {
		fmt.Println("send get request, get msg error!")
		ater := AccessTokenErrorResponse{}
		err = json.Unmarshal(body, &ater)
		fmt.Println("send get request return wechat error message %+v\n!", ater)
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf("%s", ater.ErrMsg)
	}
}

func GetAndUpdateMapToken(appID, appSecret, accessTokenFetchUrl string) error {
	accMap = make(map[int]string, 100)

	accessToken, err := FetchAccessToken(appID, appSecret, accessTokenFetchUrl)
	fmt.Println(accessToken)
	if err != nil {
		fmt.Println("send request return error !!")
		return err
	}

	accMap[len(accMap)+1] = accessToken
	fmt.Println("one ", accMap)

	return nil
}

func TestTimeTask(appID, appSecret, accessTokenFetchUrl string) map[int]string {
	fmt.Println(accMap)
	timeStr := "0 */60 * * * *"
	t2 := toolbox.NewTask("getAtoken", timeStr, func() error {
		err := GetAndUpdateMapToken(appID, appSecret, accessTokenFetchUrl)
		if err != nil {
			//todo 向微信请求access_token失败 结合业务逻辑处理
		}
		return nil
	})

	toolbox.AddTask("tk2", t2)
	toolbox.StartTask()

	return accMap
}
