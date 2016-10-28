package wxpaltutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/toolbox"
	"hello/models"
	"io/ioutil"
	"net/http"
	"strings"
)

type AccessTokenResponse struct {
	AccessToken string  `json:"access_token"`
	ExpiresIn   float64 `json:"expires_in"`
}

type AccessTokenErrorResponse struct {
	ErrCode string  `json:"errcode"`
	ErrMsg  float64 `json:"errmsg"`
}

func FetchAccessToken(appID, appSecret, accessTokenUrl string) (string, error) {
	request := strings.Join([]string{accessTokenUrl, "?grant_type=client_credential&appid=", appID, "&secret=", appSecret}, "")
	resp, err := http.Get(request)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("request url error!")
		return "", err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("request body error!")
		return "", err
	}
	if bytes.Contains(data, []byte("access_token")) {
		atr := AccessTokenResponse{}
		err = json.Unmarshal(data, &atr)
		if err != nil {
			fmt.Println("request json error!")
			return "", nil
		}
		return atr.AccessToken, nil
	} else {
		fmt.Println("send get request, get msg error!")
		ater := AccessTokenErrorResponse{}
		err = json.Unmarshal(data, &ater)
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf("%s", ater.ErrMsg)
	}
}

func GetAndUpdateDBWxAtoken(o orm.Ormer) error {
	at := models.WxAccessToken{Id: 1}
	o.ReadOrCreate(&at, "id")
	accessToken, err := FetchAccessToken("", "", "https://api.weixin.qq.com/cgi-bin/token")

	if err != nil {
		fmt.Println("get accesstoken error")
		return err
	}

	at.AccessToken = accessToken

	o.Update(&at, "access_token")

	return nil
}

func WxTimeTask(o orm.Ormer) {
	timeStr := "0 */1 * * * *"
	t2 := toolbox.NewTask("getAtoken", timeStr, func() error {
		err := GetAndUpdateDBWxAtoken(o)
		if err != nil {
			//todo
		}
		return nil
	})

	toolbox.AddTask("tk2", t2)
	toolbox.StartTask()
}
